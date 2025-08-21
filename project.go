package wig

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProjectManager struct {
	root string
}

func NewProjectManager() ProjectManager {
	root, _ := os.Getwd()
	return ProjectManager{
		root: root,
	}
}

// TODO: this must be also base on git dir
func (p ProjectManager) GetRoot() (root string) {
	return p.root
}

// Find project root by file path. Project root must have .git directory in it.
// otherwise "working directory" will be returned.
func (p ProjectManager) FindRoot(buf *Buffer) (root string, err error) {
	root = p.root
	fp := filepath.Dir(buf.FilePath)

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = fp
	r, err := cmd.Output()
	if err != nil {
		return root, nil
	}

	return strings.TrimSpace(string(r)), nil
}

// Returns Dir of current buffer or working directory if buffer has no valid file path.
func (p ProjectManager) Dir(buf *Buffer) (dir string) {
	if len(buf.FilePath) > 0 {
		r := filepath.Dir(buf.FilePath)
		if r != "." {
			return r
		}
	}
	return p.root
}

