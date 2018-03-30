// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package type contains types for use by other packages
// This file contains all the types that are used throughout the application

package types

type CliArgs struct {
    AppType     string
    Directory   string
    Board       string
    Framework   string
    Platform    string
    Ide         string
    Tests       bool
}

// type for the targets tag in the configuration file
type TargetsTag map[string]*TargetSubTags

// type for the libraries tag in the configuration file
type LibrariesTag map[string]*LibrarySubTags

// Structure to handle individual target inside targets
type TargetSubTags struct {
    Board string
    Compile_flags []string
}

// Structure to handle individual library inside libraries
type LibrarySubTags struct {
    Url           string            // local libraries will have project relative path
    Version       string
    Compile_flags []string
}

// Structure to hold information about project type: app
type AppTag struct {
    Name string
    Platform string
    Framework string
    Default_target string
    Ide string
    Targets TargetsTag
}

// Structure to hold information about project type: lib
type LibTag struct {
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
    Targets TargetsTag
}

type AppConfig struct {
    MainTag AppTag              `yaml:"app"`
    LibrariesTag LibrariesTag   `yaml:"libraries"`
}

type LibConfig struct {
    MainTag LibTag              `yaml:"lib"`
    LibrariesTag LibrariesTag   `yaml:"libraries"`
}
