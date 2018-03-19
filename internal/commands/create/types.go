// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// A struct to hold all the data for create command
package create

type ConfigCreate struct {
    AppType     string
    Directory   string
    Board       string
    Framework   string
    Platform    string
    Ide         string
    Tests       bool
}

