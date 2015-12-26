// Package exp12 provides caching of vcs per directory.
package exp12

import (
	"sync"

	"github.com/shurcooL/go/exp/13"
	"github.com/shurcooL/go/gists/gist7802150"
	"github.com/shurcooL/go/vcs"
)

// rootPath -> *VcsState
var repos = make(map[string]*exp13.VcsState)
var reposLock sync.Mutex

// path -> *Directory
var directories = make(map[string]*Directory)
var directoriesLock sync.Mutex

type Directory struct {
	path string

	Repo *exp13.VcsState // nil if no Vcs.

	gist7802150.DepNode2
}

func (this *Directory) Update() {
	if vcs := vcs.New(this.path); vcs != nil {
		reposLock.Lock()
		if repo, ok := repos[vcs.RootPath()]; ok {
			this.Repo = repo
		} else {
			this.Repo = exp13.NewVcsState(vcs)
			repos[vcs.RootPath()] = this.Repo
		}
		reposLock.Unlock()
	}
}

func newDirectory(path string) *Directory {
	this := &Directory{path: path}
	// No DepNode2I sources, so each instance can only be updated (i.e. initialized) once.
	return this
}

// LookupDirectory is safe to call concurrently.
func LookupDirectory(path string) *Directory {
	directoriesLock.Lock()
	defer directoriesLock.Unlock()
	if dir := directories[path]; dir != nil {
		return dir
	} else {
		dir = newDirectory(path)
		directories[path] = dir
		return dir
	}
}
