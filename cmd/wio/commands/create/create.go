// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Creates and initializes a wio project. It also works as an updater when called on already created projects.
package create

import (
    . "wio/cmd/wio/types"
)

// Executes the create command provided configuration packet
func Execute(config ConfigCreate) {

    // create a directory
    // create common project structure
    // move common template files

    if config.AppType == "app" {
        executeApp(config)
    } else {
        executeLib(config)
    }
}
