// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio
package pac

import (
    "github.com/urfave/cli"
    "github.com/go-errors/errors"
    "fmt"
)

const (
    GET     = "get"
    UPDATE  = "update"
    COLLECT = "collect"
)

type Pac struct {
    Context *cli.Context
    Type    string
    error
}

// Get context for the command
func (pac Pac) GetContext() (*cli.Context) {
    return pac.Context
}

// Executes the libraries command
func (pac Pac) Execute() {

    switch pac.Type {
    case GET:
        pac.error = errors.New("GG")
        pac.handleGet(pac.Context)
        break
    case UPDATE:
        break
    case COLLECT:
        break
    }
}
func (pac Pac) handleGet(context *cli.Context) {
    if pac.error != nil {
        return
    }

    fmt.Println("GGGGG")
}
