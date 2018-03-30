// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// This contains helper function for generic template parsing and creating the project
package create

import (
    "regexp"
    "path/filepath"

    . "wio/cmd/wio/utils/io"
    "os"
    "wio/cmd/wio/utils/types"
)

// Parses the paths.json file and uses that to get paths to copy files to and from
// It also stores all the paths as a map to be used later on
func parsePathsAndCopy(jsonPath string, projectPath string, tags []string) (error) {
    var paths = Paths{}
    if err := AssetIO.ParseJson(jsonPath, &paths); err != nil { return err }

    var re, e = regexp.Compile(`{{.+}}`)
    if e != nil {
        return e
    }

    var sources []string
    var destinations []string
    var overrides []bool

    for i := 0; i < len(paths.Paths); i++ {
        for t := 0; t < len(tags); t++ {
            if paths.Paths[i].Id == tags[t] {
                sources = append(sources, paths.Paths[i].Src)

                destination, e := filepath.Abs(re.ReplaceAllString(paths.Paths[i].Des, projectPath))
                if e != nil {
                    return e
                }

                destinations = append(destinations, destination)
                overrides = append(overrides, paths.Paths[i].Override)
            }
        }
    }

    return AssetIO.CopyMultipleFiles(sources, destinations, overrides)
}

// generic method to create structure based on a list of paths provided
func createStructure(projectPath string, relPaths ...string) (error) {
    for path := 0; path < len(relPaths); path++ {
        fullPath := projectPath + Sep + relPaths[path]

        if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
            return err
        }
        Verb.Verbose(`* Created "` +  relPaths[path] + `" folder` + "\n")
    }

    return nil
}

// generic method to copy templates based on command line arguments
func copyTemplates(cliArgs *types.CliArgs) (error) {
    strArray := make([]string, 1)
    strArray[0] = cliArgs.AppType + "-gen"

    Verb.Verbose("\n")
    if cliArgs.Ide == "clion" {
        Verb.Verbose("* Clion Ide available so ide template set up will be used\n")
        strArray = append(strArray, cliArgs.AppType + "-clion")
    } else {
        Verb.Verbose("* General template setup will be used\n")
    }

    jsonPath := "config" + Sep + "paths.json"
    if err := parsePathsAndCopy(jsonPath, cliArgs.Directory, strArray); err != nil { return err }
    Verb.Verbose("* All Template files created in their right position\n")

    return nil
}

// Generic method to copy cmake files for each target
func copyTargetCMakes(projectPath string, projectType string, targets types.TargetsTag) {
    for target := range targets {
        srcPath := "templates/cmake/CMakeTarget.cmake-" + projectType + ".tpl"
        destPath := projectPath + Sep + ".wio" + Sep + "targets" + Sep + target + ".cmake"
        AssetIO.CopyFile(srcPath, destPath, true)
    }
}

// Appends a string to string slice only if it is missing
func AppendIfMissing(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}