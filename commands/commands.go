package commands

import (
	"strings"
	"vc/workdir"
)

type VC struct {
	wd          *workdir.WorkDir
	commits     []*Commit
	status      *Status
	stagingArea map[string]string
}

type Commit struct {
	message string
	files   *workdir.WorkDir
}

type Status struct {
	ModifiedFiles []string
	StagedFiles   []string
}

func Init(wd *workdir.WorkDir) *VC {
	return &VC{
		wd:          wd,
		status:      &Status{},
		stagingArea: make(map[string]string),
	}
}

func (vc *VC) Add(files ...string) {
	for _, file := range files {
		content, err := vc.wd.CatFile(file)
		if err == nil {
			vc.stagingArea[file] = content
			alreadyStaged := false
			for _, stagedFile := range vc.status.StagedFiles {
				if stagedFile == file {
					alreadyStaged = true
					break
				}
			}
			if !alreadyStaged {
				vc.status.StagedFiles = append(vc.status.StagedFiles, file)
			}
		}
	}
}

func (vc *VC) AddAll() {
	files := vc.wd.ListFilesRoot()
	for _, file := range files {
		vc.Add(file)
	}
}

func (vc *VC) Status() *Status {
	files := vc.wd.ListFilesRoot()
	modifiedFiles := []string{}
	stagedFiles := vc.status.StagedFiles

	if len(vc.commits) == 0 {
		vc.status.ModifiedFiles = []string{}
		vc.status.StagedFiles = stagedFiles
		return vc.status
	}

	for _, file := range files {
		workingContent, err := vc.wd.CatFile(file)
		if err != nil {
			continue
		}

		stagedContent, isStaged := vc.stagingArea[file]

		if isStaged {
			if workingContent != stagedContent {
				modifiedFiles = append(modifiedFiles, file)
			}
		} else {
			lastCommit := vc.commits[len(vc.commits)-1]
			commitContent, err := lastCommit.files.CatFile(file)
			if err != nil {
				modifiedFiles = append(modifiedFiles, file)
			} else if workingContent != commitContent {
				modifiedFiles = append(modifiedFiles, file)
			}
		}
	}

	vc.status.ModifiedFiles = modifiedFiles
	vc.status.StagedFiles = stagedFiles
	return vc.status
}

func (vc *VC) Commit(message string) {
	newCommit := &Commit{
		message: message,
		files:   vc.wd.Clone(),
	}
	vc.commits = append(vc.commits, newCommit)
	vc.status = &Status{}
	vc.stagingArea = make(map[string]string)
}

func (vc *VC) GetWorkDir() *workdir.WorkDir {
	return vc.wd
}

func (vc *VC) GetCommit() *Commit {
	if len(vc.commits) == 0 {
		return nil
	}
	return vc.commits[len(vc.commits)-1]
}

func (vc *VC) Log() []string {
	if len(vc.commits) == 0 {
		return []string{}
	}

	var messages []string
	for i := len(vc.commits) - 1; i >= 0; i-- {
		messages = append(messages, vc.commits[i].message)
	}
	return messages
}

func (vc *VC) Checkout(commit string) (*workdir.WorkDir, error) {
	if len(vc.commits) == 0 {
		return vc.wd, nil
	}

	index := len(vc.commits) - 1

	if strings.HasPrefix(commit, "~") {
		steps := len(commit) - 1
		if commit == "~1" {
			steps = 1
		} else if commit == "~2" {
			steps = 2
		} else if commit == "~3" {
			steps = 3
		}
		index = len(vc.commits) - 1 - steps
	} else if strings.HasPrefix(commit, "^") {
		steps := len(commit)
		index = len(vc.commits) - 1 - steps
	}
	if index < 0 {
		index = 0
	}
	if index >= len(vc.commits) {
		index = len(vc.commits) - 1
	}

	vc.wd = vc.commits[index].files.Clone()
	return vc.wd, nil
}
