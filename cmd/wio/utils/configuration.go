// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package utils contains utilities/files useful throughout the app
// This file contains all the function to manipulate project configuration file

package utils

import (
    "bufio"
    "strings"
    "fmt"
    "gopkg.in/yaml.v2"
    "os"

    . "wio/cmd/wio/utils/io"
)

// Adds spacing and other formatting for project configuration
func writeProjectConfig(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        tokens := strings.Split(line, "\n")
        for _, token := range tokens {
            if strings.Contains(token, "targets:") ||
                (strings.Contains(token, "libraries:") && !strings.Contains(token, "#   libraries:")) {
                fmt.Fprint(w, "\n")
            }
            fmt.Fprintln(w, token)
        }
    }

    return w.Flush()
}

// Write configuration for the project with information on top and nice spacing
func PrettyPrintConfig(projectType string, projectConfig interface{}, filePath string) (error) {
    infoPath := "templates" + Sep + "config" + Sep + "project-" + projectType + "-help"

    var ymlData []byte
    var infoData []byte
    var err error

    // get data
    if ymlData, err = yaml.Marshal(projectConfig); err != nil { return err }
    if infoData, err = AssetIO.ReadFile(infoPath); err != nil { return err }

    infoDataSlice :=  strings.Split(string(infoData), "\n")
    totalConfig := make([]string, 0)
    totalConfig = append(totalConfig, infoDataSlice...)
    totalConfig = append(totalConfig, string(ymlData))

    if err = os.Remove(filePath); err != nil { return err }

    err = writeProjectConfig(totalConfig, filePath)
    return err
}
