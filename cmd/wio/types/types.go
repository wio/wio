// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

// DConfig contains configurations for default commandline arguments
type DConfig struct {
    Ide string
    Framework string
    Platform string
    File string
    Port string
    Version string
    Board string
    Btarget string
    Utarget string
}

type ConfigCreate struct {
    AppType     string
    Directory   string
    Board       string
    Framework   string
    Platform    string
    Ide         string
    Tests       bool
}
