package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

func LoadPath(path string) (string, error) {
	var b strings.Builder

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	root := path
	if !info.IsDir() {
		root = filepath.Dir(path)
	}

	var gitIgnore *ignore.GitIgnore
	gitignorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		gitIgnore, _ = ignore.CompileIgnoreFile(gitignorePath)
	}

	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(root, p)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if gitIgnore != nil && gitIgnore.MatchesPath(relPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(p)
		if err != nil {
			return nil
		}

		if !isText(content) {
			return nil
		}

		b.WriteString(formatFile(p, content))
		return nil
	})

	return b.String(), err
}

func isText(content []byte) bool {
	checkLen := len(content)
	if checkLen > 1024 {
		checkLen = 1024
	}

	count := 0
	for _, c := range content[:checkLen] {
		if c == 0 {
			return false
		}
		if c < 0x09 || (c > 0x0D && c < 0x20) {
			count++
		}
	}

	return float64(count)/float64(checkLen) < 0.05
}

func formatFile(path string, content []byte) string {
	ext := filepath.Ext(path)
	var b strings.Builder
	b.WriteString("\n# FILE: " + path + "\n")
	b.WriteString("```" + strings.TrimPrefix(ext, ".") + "\n")
	b.Write(content)
	b.WriteString("\n```\n")
	return b.String()
}

func (w *WorkspaceMenager) Add(path string, once bool) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	filesContent, err := LoadPath(abs)
	if err != nil {
		return err
	}
	w.Entries = append(w.Entries, Entry{Path: abs, Content: filesContent, addOnce: once})
	return nil
}

func (w *WorkspaceMenager) WorkspaceToString(triggerClearing bool) string {
	var b strings.Builder
	if len(w.Entries) == 0 {
		return ""
	}

	b.WriteString("This is files content provided by user: ")
	for num, e := range w.Entries {
		part := fmt.Sprintf("%d file: %s ", num, e.Content)
		b.WriteString(part)
	}

	if triggerClearing {
		w.Entries = slices.DeleteFunc(w.Entries, func(e Entry) bool {
			return e.addOnce
		})
	}

	return b.String()
}

func (w *WorkspaceMenager) Clear(path string) {
	if path == "all" {
		w.Entries = nil
		return
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		w.Entries = slices.DeleteFunc(w.Entries, func(e Entry) bool {
			return e.Path == path
		})
		return
	}

	cleanPath := filepath.Clean(path)

	if !info.IsDir() {
		w.Entries = slices.DeleteFunc(w.Entries, func(e Entry) bool {
			return filepath.Clean(e.Path) == cleanPath
		})
		return
	}

	w.Entries = slices.DeleteFunc(w.Entries, func(e Entry) bool {
		entryPath := filepath.Clean(e.Path)

		if entryPath == cleanPath {
			return true
		}

		return strings.HasPrefix(entryPath, cleanPath+string(os.PathSeparator))
	})
}
