// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package type contains types for use by other packages
// This file contains all the types that are used throughout the application

package types

type CliArgs struct {
    AppType   string
    Directory string
    Board     string
    Framework string
    Platform  string
    Ide       string
    Tests     bool
    Update    bool
}

// type for the targets tag in the configuration file
type TargetsTag struct {
    Default_target string                `yaml:"default"`
    Targets        map[string]*TargetTag `yaml:"create"`
}

// Structure to handle individual target inside targets
type TargetTag struct {
    Board         string
    Compile_flags []string
}

// Structure to hold information about project type: app
type AppTag struct {
    Name      string
    Platform  string
    Framework string
    Ide       string
}

// Structure to hold information about project type: lib
type PkgTag struct {
    Name          string
    Version       string
    Authors       []string
    License       []string
    Platform      string
    Framework     []string
    Board         []string
    Compile_flags []string
    Ide           string
}

type AppConfig struct {
    MainTag         AppTag          `yaml:"app"`
    TargetsTag      TargetsTag      `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

type PkgConfig struct {
    MainTag         PkgTag          `yaml:"pkg"`
    TargetsTag      TargetsTag      `yaml:"targets"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}

// Structure to handle individual library inside libraries
type DependencyTag struct {
    Url           string
    Ref           string
    Compile_flags []string
}

// type for the libraries tag in the main wio.yml file
type DependenciesTag map[string]*DependencyTag

// DConfig contains configurations for default commandline arguments
type DConfig struct {
    Ide       string
    Framework string
    Platform  string
    Port      string
    Version   string
    Board     string
    Btarget   string
    Utarget   string
}
