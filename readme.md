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
   - 5.6 [remove](#remove)
   - 5.7 [config](#config)  
   - 5.8 [exit / help](#exit--help)
6. [Memory System](#memory-system)  
7. [Workspace Management](#workspace-management)  
8. [Token Tracking](#token-tracking)  
9. [Running the Application](#running-the-application)  
10. [Extending / Contributing](#contributing)  
11. [License](#license)  

---

## Architecture Overview
```
github.com/pfczx/goagentai
│
├─ agent/          # Core *Agent* struct, factories and core logic
│   ├─ agent.go        ← Agent definition and Generate wrapper
│   ├─ agent_runner.go ← Execution flow, context building & token handling
│   ├─ config.go       ← YAML schema, validation, comment injection
│   ├─ list_command.go ← Directory provider and model/provider lists
│   ├─ profile.go      ← Profile lifecycle (creation, loading, default)
│   ├─ remove_command.go ← Deletion of profiles, history, used‑tokens
│   ├─ switch_command.go ← Switching provider / model / internal provider
│   └─ workspace_commands.go ← Workspace add / remove / list operations
│
├─ cli/            # REPL & command dispatch utilities
│   ├─ commands.go   ← Mapping of CLI keywords to callbacks
│   ├─ repl.go       ← Interactive prompt loop
│   ├─ singlerun.go  ← One‑shot handling for `goagentai <cmd>`
│   └─ utils.go      ← Helper rendering, env loading, exit handling
│
├─ llm/            # Provider abstractions (Groq, HuggingFace, OpenRouter)
│   ├─ groq.go       ← Direct HTTP interaction with Groq
│   ├─ huggingface.go← HTTP interaction with HF inference API
│   ├─ llm.go        ← Provider interface and factory
│   └─ open_router.go← HTTP interaction with OpenRouter
│
├─ memory/         # Persistent storage (SQLite) and semantic operations
│   ├─ db_handler.go ← Migration handling with goose
│   ├─ generated/*   ← SQLC‑generated types & query functions
│   ├─ long_term.go  ← Summary generation and IDF weighting
│   ├─ menager.go    ← Public façade for Short‑/Long‑Term Memory
│   ├─ short_term.go ← JSON storage of recent turns
│   └─ text_analyzer.go ← TF‑IDF analyzer for relevance selection
│
├─ workspace/      # OPTIONAL persistent workspace of file “content” snippets
│   ├─ menager.go   ← Simple slice of *Entry* structs
│   └─ workspace.go ← Load, add, clear, render utilities
│
├─ token/          # Token consumption accounting
│    ├─ token.go     ← Simple in‑memory read/write of usage counters
│    └─ menager.go   ← Public façade for token package, save/load utilities
│    
│
├─ go.mod / go.sum # Dependency list (includes https://github.com/kbinani/screenshot)
└─ main.go         # Profile bootstrap and entry point
```

---

## Prerequisites for building from source

| Item | Reason |
|------|--------|
| **Go 1.25+** | Required by `go.mod`. |
| **Git** | To clone the repository. |
| **Terminal with ANSI colour support** | For nicer REPL prompting (charmbracelet packages rely on it). |
|** XGB for linux/BSD or cgo for OSX**| screenshot package dependancies github.com/kbinani/screenshot |

XGB or cgo will be installed automatically, but your system should support their requirements

---

## Installation from source

```bash
# Clone the repository
git clone https://github.com/pfczx/goagentai.git
cd goagentai

# System wide installation
go install
```

Invoke it via `go run ./...`, or build binary with `go build -o goagentai` or `go build -o goagentai.exe` for windows 

---

## Configuration & Profiles


### Profile Directory Layout (auto generated)
```
~/.config/goagent/
├─ latestProfile                     ← Latest used profile name 
│
├─long_term.db                       ← SQLite long term memory database 
│
├─ profiles/
│   ├─ <profile_name>/               ← Individual folders per profile
│   │   ├─ config.yaml               ← Human‑editable configuration witg config command
│   │   └─ shortTermMemory.json      ← Auto‑created memory file
│   └─ default/                      ← Auto‑generated fallback profile
```

---

### Each profile config contains:

| Field | Description |
|-------|-------------|
| `name` | Human readable identifier (also the folder name). |
| `path` | Absolute path to the profile folder. |
| `provider` | One of `groq`, `openrouter`, `huggingface` |
| `model` | Model identifier accepted by the selected provider (e.g. `gpt-4-turbo`). |
| `memory_on` | Toggle memory system. |
| … | A total of ~15 flags controlling memory limits, summarisation settings, token budget, output format, etc. Check agent/config.go for more information|

---

### Init command
```bash
goagentai init profile1
```

The command creates `~/.config/goagent/profiles/profile1/config.yaml` populated with the default template from `config.go`.  
To create additional profiles, repeat the command with a new name (e.g. `goagentai init myprofile`).

---

### Editing Configuration

```bash
goagentai config
```

The command launches `$EDITOR` (defaults to `nano` on Unix, `notepad` on Windows) on the selected profile’s `config.yaml`. Once the editor exits, the CLI reloads the profile automatically.

---

## CLI Reference

All commands are invoked as

```
goagentai <subcommand> [args...]
```

or interactively via REPL (`goagentai`).  
Aliases are documented in the help output (`goagentai help`).

---

### 1. `init` / `i`

Creates a new profile.

```
goagentai init <profile_name>
```

*If a profile with that name already exists the command aborts.*

---

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


Output is rendered:

* With Markdown if `profile.config.MdFormat` is `true` (requires `github.com/charmbracelet/glamour`).  
* Otherwise plain text is printed.

The command returns token usage statistics when available.

---

### 3. `switch` / `s`

Switches various runtime settings:

| Alias | Usage |
|-------|-------|
| `profile` | `goagentai switch profile <new_name>` |
| `provider` | `goagentai switch provider <provider_name>` |
| `internal-provider` | `goagentai switch internal-provider <number or name>` |
| `model` | `goagentai switch model <number or model_name>` |

Numbers refer to the order shown by `list`/`list providers`.

---

### 4. `list` / `l`

Prints catalogues of **profiles**, **providers**, **internal providers**, or **models**.

Examples:

```
goagentai list profiles
goagentai list providers
goagentai list models -img   # (only applies to providers supporting multimodal models)
```

---

### 5. `workspace` / `w`

Manipulates a *file workspace* attached to the current profile.

| Sub‑command | Description |
|-------------|-------------|
| `add <path>` (`-p` to not be erased after next ask) | Adds a file/directory (auto‑clear after next `ask`). |
| `remove <path or all>` | Deletes a specific entry (`all` clears the whole workspace). |
| `list` (`-all` for listing with content )  | Displays stored paths. If `all` is supplied the full file content is rendered. |

The workspace content is injected into LLM prompts when `ask` command is running.

---

### 6. `remove` / `r`

Deletes profile‑related data or usage statistics.

| Argument | Effect |
|----------|--------|
| `profile <name>` | Removes profile folder **only after switching away** from it. |
| `history <st or lt>` | Clears short‑term or long‑term stored memories of the current profile. |
| `used-tokens` | Resets the token usage counter stored on disk. |

---

### 7. `config` / `c`

Opens the profile’s `config.yaml` for manual editing; afterwards the agent is re‑initialised with the updated configuration.

---

### 8. `exit` / `e` and `help` / `h`

 `exit` – Cleanly closes the DB and terminates.  
 `help` – Prints a formatted list of all commands (generated from `cli/commands.go`).


---

## Memory System

1. **Short‑Term Memory** – JSON file (`shortTermMemory.json`) appended after each user turn.  
   * Limited in length by `ShortTermMemoryLimit`.  
   * Evaluation flag (`ShortTermMemoryEvaluation`) is queried at the end of a response; a “y” can be supplied by the user to mark a turn as useful, influencing future context selection.

2. **Long‑Term Memory** – Persistent SQLite tables (`long_term_memory`, `short_term_memory`).  
   * Periodically summarised when the buffer size exceeds `LongTermMemoryBufferSize`.  
   * Summaries are generated by a provider dedicated to summarisation (configurable).  
   * An IDF‑weighted TF‑IDF classifier (`text_analyzer.go`) selects the most relevant chunks for a given query.

3. **Semantic Search**  
   * Compute TF‑IDF vectors for each stored chunk.  
   * `SelectRelevantChunks` matches query tokens against these vectors and returns up to `LongTermMemoryChunksToAdd` of the highest‑scoring summaries.

4. **Safety** 
   * The summarisation step runs in a background goroutine to avoid blocking the CLI.  
   * Panics are recovered and logged, preventing the application from terminating abruptly.

---

## Workspace Management

The workspace is a *content‑only* annex that can hold the textual representation of files you wish to keep visible to the model.  
It is deliberately simple:

* Files are read (excluding binary or `.gitignore`‑matched content).  
* The rendered markdown (`# FILE: <path>` … ```filetype```) is stored in a slice.  
* Only plain‑text file types are kept; binary or large binary blobs are ignored.  

The workspace can be cleared with the `workspace remove <path or "all">` command, all for clearing all of workspace entries.

---

## Token Tracking

A tiny `TokenMenager` keep track of used tokens durning session.  
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
GoAgent@default (groq::groq::openai/gpt-oss-20b) 36036 / 1000000 (3.60%)
```

Execute any of the commands above. The REPL supports line editing, history, and automatic colourisation.

---

