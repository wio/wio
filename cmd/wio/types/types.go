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

// Structure to handle targets created
type TargetsStruct map[string]*TargetStruct
type TargetStruct struct {
    Board string
    Compile_flags []string
}

// Structures to handle Libraries imported
type LibrariesStruct map[string]*LibraryStruct
type LibraryStruct struct {
    Url           string            // local libraries will have project relative path
    Version       string
    Compile_flags []string
}

// Structure to hold information about application
type AppStruct struct {
    Name string
    Platform string
    Framework string
    Default_target string
    Ide string
    Targets TargetsStruct
}

// Structure to hold information about library
type LibStruct struct {
    Name string
    Version string
    Authors []string
    License []string
    Platform string
    Framework []string
    Board []string
    Compile_flags []string
    Ide string
    Default_target string
    Targets TargetsStruct
}
