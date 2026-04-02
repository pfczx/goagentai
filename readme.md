# goagentai

**A CLI for LLM‑driven agents with integrated memory, workspace, and token accounting.**  
The project is written in Go and bundles a lightweight file‑based SQLite database, a modular configuration system, and out‑of‑the‑box support for several LLM providers (Groq, HuggingFace, OpenRouter).

---

## Table of Contents
1. [Architecture Overview](#architecture-overview)  
2. [Prerequisites](#prerequisites)  
3. [Installation](#installation)  
4. [Configuration & Profiles](#configuration--profiles)  
5. [Command‑Line Interface (CLI) Reference](#cli-reference)  
   - 5.1 [init](#init)  
   - 5.2 [ask/#40263](#ask)  
   - 5.3 [switch](#switch)  
   - 5.4 [list](#list)  
   - 5.5 [workspace](#workspace)  
   - 5.6 [remove](#remove)     - 5.7 [config](#config)  
   - 5.8 [exit / help](#exit--help)  6. [Memory System](#memory-system)  
7. [Workspace Management](#workspace-management)  
8. [Token Tracking](#token-tracking)  
9. [Running the Application](#running-the-application)  
10. [Extending / Contributing](#contributing)  
11. [License](#license)  

---

## Architecture Overview```
github.com/pfczx/goagentai
│
├─ agent/          # Core *Agent* struct, factories and core logic│   ├─ agent.go        ← Agent definition and Generate wrapper
│   ├─ agent_runner.go ← Execution flow, context building & token handling
│   ├─ config.go       ← YAML schema, validation, comment injection
│   ├─ list_command.go ← Directory provider and model/provider lists
│   ├─ profile.go      ← Profile lifecycle (creation, loading, default)
│   ├─ remove_command.go ← Deletion of profiles, history, used‑tokens│   ├─ switch_command.go ← Switching provider / model / internal provider
│   └─ workspace_commands.go ← Workspace add / remove / list operations
│
├─ cli/            # REPL & command dispatch utilities│   ├─ commands.go   ← Mapping of CLI keywords to callbacks
│   ├─ repl.go       ← Interactive prompt loop
│   ├─ singlerun.go  ← One‑shot handling for `goagentai <cmd>`
│   └─ utils.go      ← Helper rendering, env loading, exit handling
│
├─ llm/            # Provider abstractions (Groq, HuggingFace, OpenRouter)
│   ├─ groq.go       ← Direct HTTP interaction with Groq
│   ├─ huggingface.go← HTTP interaction with HF inference API
│   ├─ llm.go        ← Provider interface and factory
│   └─ open_router.go← HTTP interaction with OpenRouter│
├─ memory/         # Persistent storage (SQLite) and semantic operations
│   ├─ db_handler.go ← Migration handling with goose
│   ├─ generated/*   ← SQLC‑generated types & query functions│   ├─ long_term.go  ← Summary generation and IDF weighting
│   ├─ menager.go    ← Public façade for Short‑/Long‑Term Memory
│   ├─ short_term.go ← JSON storage of recent turns│   └─ text_analyzer.go ← TF‑IDF analyzer for relevance selection
│
├─ workspace/      # OPTIONAL persistent workspace of file “content” snippets│   ├─ menager.go   ← Simple slice of *Entry* structs
│   └─ workspace.go ← Load, add, clear, render utilities
│
├─ token/          # Token consumption accounting
│   └─ menager.go   ← Simple in‑memory read/write of usage counters
│
├─ go.mod / go.sum # Dependency list (includes https://github.com/kbinani/screenshot)
└─ main.go         # Profile bootstrap and entry point
```

---

## Prerequisites

| Item | Reason |
|------|--------|
| **Go 1.25+** | Required by `go.mod`. |
| **Git** | To clone the repository. |
| **SQLite library** (optional) | Used for internal DB; Go driver `github.com/mattn/go-sqlite3` has CGO dependencies. |
| **Environment variables** | Create a `.env` file in `~/.config/goagent/` containing one of the following keys: `GROQ`, `HUGGING_FACE`, `OPEN_ROUTER`.  The value is the provider‑specific API key. |
| **Terminal with ANSI colour support** | For nicer REPL prompting (charmbracelet packages rely on it). |

---

## Installation

```bash
# Clone the repository
git clone https://github.com/pfczx/goagentai.git
cd goagentai

# Build a single binary (optional)
go build -o goagentai ./main.go
```

Move the binary to a location on your `$PATH` (e.g. `~/.local/bin/goagentai`) or invoke it via `go run ./...`.

---

## Configuration & Profiles### Profile Directory Layout

```
~/.config/goagent/
│
├─ profiles/
│   ├─ <profile_name>/               ← Individual folders per profile
│   │   ├─ config.yaml               ← Human‑editable configuration
│   │   └─ shortTermMemory.json      ← Auto‑created memory file
│   └─ default/                      ← Auto‑generated fallback profile```

Each profile contains:

| Field | Description |
|-------|-------------|
| `name` | Human readable identifier (also the folder name). |
| `path` | Absolute path to the profile folder. |
| `provider` | One of `groq`, `openai`, `openrouter`, `huggingface` (as defined in `.env`). |
| `model` | Model identifier accepted by the selected provider (e.g. `openai/gpt-4-turbo`). |
| `memory_on` | Toggle persistent memory usage. |
| … | A total of ~15 flags controlling memory limits, summarisation settings, token budget, output format, etc. |

### Default Configuration

Run **once** to create the default profile:

```bash
goagentai init default
```

The command creates `~/.config/goagent/profiles/default/config.yaml` populated with the template from `config.go`.  
To create additional profiles, repeat the command with a new name (e.g. `goagentai init myprofile`).

### Editing Configuration

```bash
goagentai config```

The command launches `$EDITOR` (defaults to `nano` on Unix, `notepad` on Windows) on the selected profile’s `config.yaml`. Once the editor exits, the CLI reloads the profile automatically.

---

## CLI Reference

All commands are invoked as

```
goagentai <subcommand> [args...]
```

or interactively via REPL (`goagentai`).  
Aliases are documented in the help output (`goagentai help`).

### 1. `init` / `i`

Creates a new profile.

```
goagentai init <profile_name>
```

*If a profile with that name already exists the command aborts.*

### 2. `ask` / `a`

**Ask the LLM** using the currently selected profile.

#### Basic Syntax
```
goagentai ask <prompt>
```

#### Optional Flags
| Flag | Effect |
|------|--------|
| `-s` | Capture the current display as a screenshot, encode as base64 and send to the model. |
| `-p` (workspace) | Add current workspace content to the LLM context (see *workspace*). |

Output is rendered:

* With Markdown if `profile.config.MdFormat` is `true` (requires `github.com/charmbracelet/glamour`).  
* Otherwise plain text is printed.

The command returns token usage statistics when available.

### 3. `switch` / `s`

Switches various runtime settings:

| Alias | Usage |
|-------|-------|
| `profile` | `goagentai switch profile <new_name>` |
| `provider` | `goagentai switch provider <provider_name>` |
| `internal-provider` | `goagentai switch internal-provider <0|name>` |
| `model` | `goagentai switch model <0|model_name>` |

Numbers refer to the order shown by `list`/`list providers`.

### 4. `list` / `l`

Prints catalogues of **profiles**, **providers**, **internal providers**, or **models**.

Examples:

```
goagentai list profiles
goagentai list providers
goagentai list models -img   # (only applies to providers supporting multimodal models)
```

### 5. `workspace` / `w`

Manipulates a *file workspace* attached to the current profile.

| Sub‑command | Description |
|-------------|-------------|
| `add <path>` (`-p` to persist) | Adds a file/directory, optionally `add‑once` (auto‑clear after next `ask`). |
| `remove <path|all>` | Deletes a specific entry (`all` clears the whole workspace). |
| `list [<all|...]` | Displays stored paths. If `all` is supplied the full file content is rendered. |

The workspace content is injected into LLM prompts when `triggerScreenshot` is enabled or when memory is to be enriched.

### 6. `remove` / `r`

Deletes profile‑related data or usage statistics.

| Argument | Effect |
|----------|--------|
| `profile <name>` | Removes profile folder **only after switching away** from it. |
| `history <st|lt>` | Clears short‑term or long‑term stored memories of the current profile. |
| `used-tokens` | Resets the token usage counter stored on disk. |

### 7. `config` / `c`

Opens the profile’s `config.yaml` for manual editing; afterwards the agent is re‑initialised with the updated configuration.

### 8. `exit` / `e` and `help` / `h`

* `exit` – Cleanly closes the DB and terminates.  * `help` – Prints a formatted list of all commands (generated from `cli/commands.go`).

---

## Memory System

1. **Short‑Term Memory** – JSON file (`shortTermMemory.json`) appended after each user turn.  
   * Limited in length by `ShortTermMemoryLimit`.  
   * Evaluation flag (`ShortTermMemoryEvaluation`) is queried at the end of a response; a “y” can be supplied by the user to mark a turn as useful, influencing future context selection.

2. **Long‑Term Memory** – Persistent SQLite tables (`long_term_memory`, `short_term_memory`).  
   * Periodically summarised when the buffer size exceeds `LongTermMemoryBufferSize`.  
   * Summaries are generated by a provider dedicated to summarisation (configurable).  
   * An IDF‑weighted TF‑IDF classifier (`text_analyzer.go`) selects the most relevant chunks for a given query.

3. **Semantic Search** –  
   * Compute TF‑IDF vectors for each stored chunk.  
   * `SelectRelevantChunks` matches query tokens against these vectors and returns up to `LongTermMemoryChunksToAdd` of the highest‑scoring summaries.

4. **Safety** –  
   * The summarisation step runs in a background goroutine to avoid blocking the CLI.  
   * Panics are recovered and logged, preventing the application from terminating abruptly.

---

## Workspace Management

The workspace is a *content‑only* annex that can hold the textual representation of files you wish to keep visible to the model.  
It is deliberately simple:

* Files are read (excluding binary or `.gitignore`‑matched content).  
* The rendered markdown (`# FILE: <path>` … ```filetype```) is stored verbatim in a slice.  
* Only plain‑text file types are kept; binary or large binary blobs are ignored.  

The workspace can be cleared or filtered per request; the current snapshot can be dumped with the `-img` flag in `list models`.

---

## Token Tracking

A tiny `TokenMenager` persists a usage counter under `~/.config/goagent/token/tokenBalanceUsed`.  
* `AddUsage(tokens)` increments the counter and warns once a `TokenBalanceLimit` (set per‑profile) is reached.  
* `Status()` prints a human‑readable “used / limit (percentage)” line in the REPL prompt.  
* `ResetUsage()` clears the counter, useful after a fresh conversation or before a batch of tests.

The usage information is displayed after each answer and added to the persisted value when summarisation occurs.

---

## Running the Application

### Non‑interactive (one‑shot) usage

```bash
# Switch to a profile (optional)
goagentai switch profile myprofile

# Ask a question, ask for markdown output
goagentai ask "Explain quantum entanglement in three sentences" -p
```

If `-s` is supplied, the CLI will pause and wait for you to specify a display number to capture a screenshot. The image is then Base64 encoded and sent as part of the prompt.

### Interactive REPL

```bash
goagentai
```

You will see a prompt of the form:

```
GoAgent@default (groq::groq::openai/gpt-oss-20b) 123 tokens
```

Execute any of the commands above. The REPL supports line editing, history, and automatic colourisation.

### Running Tests

The repository does not ship a test suite yet, but the generated `memory/db/generated` package can be exercised with `go test ./...` after you install the `sqlc` command and manually run `sqlc generate`.

---

## Extending / Contributing1. **Fork** the repository and clone your fork locally.  
2. **Create a branch** for your feature or bug‑fix.  
3. **Run** `go vet ./...` and `go test ./...` (once test infrastructure is added).  4. **Generate SQLC code** (if you modify queries):

   ```bash
   sqlc generate
   ```

5. **Submit a pull request** with a clear description of the change and any required updates to documentation (this README).  

The project follows the conventional Go module layout and respects semantic versioning. Keep the generated files (`memory/db/generated/*`) under version control; they are required for compilation.

---

## License

© 2025 pfczx. This project is released under the **MIT License**.  
See the `LICENSE` file for full terms.

--- 
