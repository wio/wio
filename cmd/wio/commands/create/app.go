// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates an executable application
package create

import (
    "fmt"
    . "wio/cmd/wio/types"
)

// Executes app sub command of create command
func executeApp(config ConfigCreate) {
    fmt.Println(config)
}

