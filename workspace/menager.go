package workspace

import ()

type WorkspaceMenager struct {
	Entries []Entry
}

type Entry struct {
	Path    string
	Content string
	addOnce bool
}

func NewWorkspaceMenager() *WorkspaceMenager {
	var e []Entry
	return &WorkspaceMenager{
		Entries: e,
	}
}
