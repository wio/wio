// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates a library to be published
package create

import (
    "fmt"
)

// Executes lib sub command of create command
func executeLib(config ConfigCreate) {
    fmt.Println(config)
}

