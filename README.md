# lil-git

A minimal version control system implementation in Go, inspired by Git. This project was created as a learning exercise based on a challenge from [quera.org](https://quera.org).

## Features

- Initialize version control repositories
- Stage files individually or all at once
- Commit changes with messages
- Check repository status (modified and staged files)
- View commit history
- Checkout previous commits using relative references (~1, ~2, ^, ^^)

## Requirements

Go 1.19 or later

## Installation

```bash
go get vc
```

## Usage

```go
package main

import (
    "vc/commands"
    "vc/workdir"
)

func main() {
    // Initialize a work directory
    wd := workdir.InitEmptyWorkDir()
    wd.CreateFile("README.md")
    wd.WriteToFile("README.md", "# My Project")
    
    // Initialize version control
    vc := commands.Init(wd)
    
    // Stage and commit files
    vc.AddAll()
    vc.Commit("initial commit")
    
    // Make changes
    wd.AppendToFile("README.md", "\nNew content")
    
    // Check status
    status := vc.Status()
    // status.ModifiedFiles contains ["README.md"]
    
    // Stage specific files
    vc.Add("README.md")
    vc.Commit("update README")
    
    // View commit history
    log := vc.Log()
    // log contains ["update README", "initial commit"]
    
    // Checkout previous commit
    wd, _ := vc.Checkout("~1")
    content, _ := wd.CatFile("README.md")
    // content is the state from "initial commit"
}
```

## Project Structure

```
.
├── commands/
│   └── commands.go    # Version control operations (Init, Add, Commit, Status, Log, Checkout)
├── workdir/
│   └── workdir.go     # File system abstraction for managing files and directories
└── test/
    ├── init_test.go   # Test setup and initialization
    ├── vc_test.go     # Version control command tests
    └── workdir_test.go # Work directory operation tests
```

## Testing

```bash
go test ./...
```