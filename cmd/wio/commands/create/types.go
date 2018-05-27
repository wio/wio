// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Types for all the things being used in this package
package create

type Data struct {
    Id       string
    Src      string
    Des      string
    Override bool
}

type Paths struct {
    Paths []Data
}
