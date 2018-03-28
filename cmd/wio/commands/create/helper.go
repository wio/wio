// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// This contains helper function for generic template parsing and creating the project
package create

import (
    "os"
    "bufio"
    "strings"
    "fmt"
    "reflect"
    "errors"
    . "wio/cmd/wio/types"
    "gopkg.in/yaml.v2"
)

// Converts map interface read from YML to Libraries structure
func getLibrariesStruct(data map[string]interface{}) (LibrariesStruct) {
    libraries := LibrariesStruct{}
    for k := range data {
        dependencyMap := data[k].(map[string]interface{})
        library := &LibraryStruct{}

        library.Url = dependencyMap["url"].(string)
        library.Version = dependencyMap["version"].(string)
        var compileF []string
        compileFlags := reflect.ValueOf(dependencyMap["compile_flags"].(interface{}))
        for i := 0; i < compileFlags.Len(); i++ {
            compileF = append(compileF, compileFlags.Index(i).Interface().(string))
        }

        library.Compile_flags = compileF
        libraries[string(k)] = library
    }

    return libraries
}

// Uses reflection to set the field of structure based on map key
func setField(obj interface{}, name string, value interface{}) error {
    name = strings.Title(name)
    structValue := reflect.ValueOf(obj).Elem()
    structFieldValue := structValue.FieldByName(name)

    if !structFieldValue.IsValid() {
        return fmt.Errorf("no such field: %s in obj", name)
    }

    if !structFieldValue.CanSet() {
        return fmt.Errorf("cannot set %s field value", name)
    }

    val := reflect.ValueOf(value)
    structFieldType := structFieldValue.Type()

    if structFieldType != val.Type() {
        return errors.New("provided value type didn't match obj field type")
    }

    structFieldValue.Set(val)

    return nil
}

// Converts map interface read from YML to App structure
func getAppStruct(data map[string]interface{}) (AppStruct, error) {
    wio := AppStruct{}
    wioMap := data
    var err error = nil
    for k := range wioMap {
        if k == "targets" {
            wio.Targets = getTargetsStruct(wioMap[k].(map[string]interface{}))
        } else{
            err = setField(&wio, k, wioMap[k])
        }
    }
    return wio, err
}

// Parses slice interface and copies values over to an actual slice
func parseSlice(wioMap map[string]interface{}, tag string) ([]string) {
    var dataF []string
    data := reflect.ValueOf(wioMap[tag].(interface{}))
    for i := 0; i < data.Len(); i++ {
        dataF = append(dataF, data.Index(i).Interface().(string))
    }
    return dataF
}

// Converts map interface read from YML to Lib structure
func getLibStruct(data map[string]interface{}) (LibStruct, error) {
    wio := LibStruct{}
    wioMap := data
    var err error = nil
    for k := range wioMap {
        if k == "authors" {
            wio.Authors = parseSlice(wioMap, "authors")
        } else if k == "license" {
            wio.License = parseSlice(wioMap, "license")
        } else if k == "framework" {
            wio.Framework = parseSlice(wioMap, "framework")
        } else if k == "board" {
            wio.Board = parseSlice(wioMap, "board")
        } else if k == "compile_flags" {
            wio.Compile_flags = parseSlice(wioMap, "compile_flags")
        } else if k == "targets" {
            wio.Targets = getTargetsStruct(wioMap[k].(map[string]interface{}))
        } else{
            err = setField(&wio, k, wioMap[k])
        }
    }
    return wio, err
}

// Converts map interface read from YML to Targets structure
func getTargetsStruct(data map[string]interface{}) (TargetsStruct) {
    targets := TargetsStruct{}
    for k := range data {
        targetMap := data[k].(map[string]interface{})
        target := &TargetStruct{}
        target.Board = targetMap["board"].(string)
        var compileF []string
        compileFlags := reflect.ValueOf(targetMap["compile_flags"].(interface{}))
        for i := 0; i < compileFlags.Len(); i++ {
            compileF = append(compileF, compileFlags.Index(i).Interface().(string))
        }
        target.Compile_flags = compileF
        targets[string(k)] = target
    }

    return targets
}

// Reads lines from a path
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// WriteLines writes the lines to the given file.
func writeAppConfig(lines []string, path string) error {
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

// Prints Config file with nice spacing and info at the top
func PrettyWriteConfig(infoPath string, configData interface{}, configPath string) (error) {
    ymlData, err := yaml.Marshal(configData)
    if err != nil {
        return err
    }

    infoData, err := readLines(infoPath)
    if err != nil {
        return err
    }

    totalConfig := make([]string, 0)
    totalConfig = append(totalConfig, infoData...)
    totalConfig = append(totalConfig, string(ymlData))

    err = os.Remove(configPath)
    if err != nil {
        return err
    }

    err = writeAppConfig(totalConfig, configPath)
    return err
}
