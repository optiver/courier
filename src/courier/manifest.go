package main

type Manifest map[string]Dependency

type Dependency interface {
	IgnoreDir() string
	DirToCopy() string
}

type GitDependency struct {
	VCS string `json:"vcs"` // XXX it doesn't make sense from a data model point of view to have this field. It's to help converting to JSON.
	URL string `json:"url"`
	Ref string `json:"ref"`
	Dir string `json:"dir"`
}

func (d GitDependency) IgnoreDir() string { return ".git" }
func (d GitDependency) DirToCopy() string { return d.Dir }

type SVNDependency struct {
	VCS string  `json:"vcs"` // XXX it doesn't make sense from a data model point of view to have this field. It's to help converting to JSON.
	URL string  `json:"url"`
	Rev *string `json:"rev,omitempty"`
}

func (d SVNDependency) IgnoreDir() string { return ".svn" }
func (d SVNDependency) DirToCopy() string { return "" } // Copy the whole thing.
