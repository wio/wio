// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// This contains helper function for parsing paths and copying files based on paths
package create

import (
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "regexp"
    "strings"

    "wio/cmd/wio/utils"
    . "wio/cmd/wio/types"
)

var sep = string(filepath.Separator) // separator based on the OS

// Gives the root path to this execurable
func GetRoot() (string, error){
    rootPath, err := utils.GetExecutableRootPath()

    return rootPath, err
}

// Returns a path to template files. User needs to provide a relative path to get the full path
func GetAssetsPathRelative(fileName string) (string, error) {
    rootPath, err := GetRoot()

    return rootPath + string(filepath.Separator) + "assets" + string(filepath.Separator) + fileName, err
}

// Parses the paths.json file and uses that to get paths to copy files to and from
// It also stores all the paths as a map to be used later on
func ParsePathsAndCopy(path string, config ConfigCreate, tags []string) (map[string]string, error) {
    text, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var paths = Paths{}
    err = json.Unmarshal([]byte(text), &paths)
    if err != nil {
        return nil, err
    }

    var re, e = regexp.Compile(`{{.+}}`)
    if e != nil {
        return nil, e
    }

    rootPath, err := GetRoot()
    if err != nil {
        return nil, err
    }

    pathsMap := make(map[string]string)

    for i := 0; i < len(paths.Paths); i++ {
        for t := 0; t < len(tags); t++ {
            if paths.Paths[i].Id == tags[t] {
                src, e := filepath.Abs(re.ReplaceAllString(paths.Paths[i].Src, rootPath))
                if e != nil {
                    return nil, e
                }

                des, e := filepath.Abs(re.ReplaceAllString(paths.Paths[i].Des, config.Directory))
                if e != nil {
                    return nil, e
                }

                desBaseName := strings.Trim(strings.Replace(des, config.Directory, "", -1), sep)
                pathsMap[desBaseName] = des

                override := paths.Paths[i].Override

                err = utils.Copy(src, des, override)
                if err != nil {
                    return nil, err
                }
            }
        }
    }

    return pathsMap, err
}
