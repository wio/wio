// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for each of them

package cmake

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/parsers"
    "wio/cmd/wio/utils"
    "io/ioutil"
    "path/filepath"
    "strings"
)

var configFileName = "wio.yml"
var vendorPkgFolder = "vendor"
var remotePkgFolder = "remote"
var lockFileName = "pkg.lock"
var executionWayFile = "executionWay.txt"

// Parses Package based on the path provides and its dependencies and it modifies dependency tree to store the data
func parsePackage(packagePath string, depTree *parsers.DependencyTree, dependencies types.DependenciesTag,
    hash string) (error) {

    configFile := packagePath + io.Sep + configFileName
    dependenciesPath := packagePath

    var config types.PkgConfig
    var pkgName string

    dependencyTag := parsers.DependencyTag{}

    // check if wio.yml file exists for the dependency and based on that gather the name
    // of the package
    if utils.PathExists(configFile) {
        if err := io.NormalIO.ParseYml(configFile, &config); err != nil {
            return err
        }

        pkgName = config.MainTag.Name
        dependencyTag.Compile_flags = config.MainTag.Compile_flags
    } else {
        pkgName = filepath.Base(packagePath)
    }

    if len(hash) <= 0 {
        hash += pkgName
    } else {
        hash += "__" + pkgName
    }

    dependencyTag.Name = pkgName
    dependencyTag.Hash = hash
    dependencyTag.Path = packagePath

    depTree.Config = dependencyTag
    depTree.Child = make([]*parsers.DependencyTree, 0)

    // parse vendor dependencies
    vendorPath := dependenciesPath + io.Sep + "vendor"
    if utils.PathExists(vendorPath) {
        if dirs, err := ioutil.ReadDir(vendorPath); err != nil {
            return err
        } else if len(dirs) > 0 {
            for dir := range dirs {
                if dirs[dir].Name()[0] == '.' {
                    continue
                }

                currTree := parsers.DependencyTree{}

                if err := parsePackage(vendorPath+io.Sep+dirs[dir].Name(), &currTree, config.DependenciesTag, hash);
                    err != nil {
                    return err
                }
                depTree.Child = append(depTree.Child, &currTree)
            }
        }
    }

    // parse remote dependencies
    remotePath := dependenciesPath + io.Sep + "remote"
    if utils.PathExists(remotePath) {
        if dirs, err := ioutil.ReadDir(remotePath); err != nil {
            return err
        } else if len(dirs) > 0 {
            for dir := range dirs {
                if dirs[dir].Name()[0] == '.' {
                    continue
                }

                currTree := parsers.DependencyTree{}

                if err := parsePackage(remotePath+io.Sep+dirs[dir].Name(), &currTree, config.DependenciesTag, hash);
                    err != nil {
                    return err
                }
                depTree.Child = append(depTree.Child, &currTree)
            }
        }
    }

    if val, ok := dependencies[depTree.Config.Name]; ok {
        allFlags := utils.AppendIfMissing(val.Compile_flags, depTree.Config.Compile_flags)

        depTree.Config.Compile_flags = allFlags
    }

    return nil
}

// Parses all the packages recursively with the help of ParsePackage function
func parsePackages(packagesPath string, depTrees []*parsers.DependencyTree, dependencies types.DependenciesTag,
    hash string) ([]*parsers.DependencyTree, error) {

    if !utils.PathExists(packagesPath) {
        return depTrees, nil
    }

    if dirs, err := ioutil.ReadDir(packagesPath); err != nil {
        return depTrees, err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            currTree := parsers.DependencyTree{}

            // ignore hidden files
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            if err := parsePackage(packagesPath+io.Sep+dirs[dir].Name(), &currTree, dependencies, hash);
                err != nil {
                return depTrees, err
            }

            depTrees = append(depTrees, &currTree)
        }
    }

    return depTrees, nil
}

// Parses all the packages and their dependencies and creates pkg.lock file
func createDependencyTree(projectPath string, dependencies types.DependenciesTag) ([]*parsers.DependencyTree, error) {
    wioPath := projectPath + io.Sep + ".wio"
    librariesLocalPath := wioPath + io.Sep + "pkg" + io.Sep + vendorPkgFolder
    librariesRemotePath := wioPath + io.Sep + "pkg" + io.Sep + remotePkgFolder

    dependencyTrees := make([]*parsers.DependencyTree, 0)

    // parse all the libraries and their dependencies in vendor folder
    dependencyTrees, err := parsePackages(librariesLocalPath, dependencyTrees, dependencies, "")
    if err != nil {
        return nil, err
    }

    // parse all the libraries and their dependencies in remote folder
    dependencyTrees, err = parsePackages(librariesRemotePath, dependencyTrees, dependencies, "")
    if err != nil {
        return nil, err
    }

    // return Dependency Tree
    return dependencyTrees, nil
}

func buildCMakeString(templateString string, finalString string, depTree *parsers.DependencyTree) (string) {
    currString := templateString

    currString = strings.Replace(currString, "{{DEPENDENCY_PATH}}", filepath.ToSlash(depTree.Config.Path), -1)
    currString = strings.Replace(currString, "{{DEPENDENCY_NAME}}", depTree.Config.Hash, -1)
    currString = strings.Replace(currString, "{{DEPENDENCY_COMPILE_FLAGS}}",
        strings.Join(depTree.Config.Compile_flags, " "), -1)
    currString += "\n"

    for dep := range depTree.Child {
        currLinkString := "target_link_libraries({{DEPENDENCY_NAME}} PUBLIC {{DEP_NAME}})"
        currLinkString = strings.Replace(currLinkString, "{{DEPENDENCY_NAME}}", depTree.Config.Hash, -1)
        currLinkString = strings.Replace(currLinkString, "{{DEP_NAME}}", depTree.Child[dep].Config.Hash, -1)

        currString += currLinkString + "\n"
    }

    finalString += currString + "\n"

    for dep := range depTree.Child {
        finalString = buildCMakeString(templateString, finalString, depTree.Child[dep])
    }

    return finalString
}

// Given all the dependency tree data, it will create cmake files for each dependency
func createDependencyCMakeString(depTree *parsers.DependencyTree) (string, error) {

    dependenciesTemplate, err := io.AssetIO.ReadFile("templates" + io.Sep + "cmake" + io.Sep +
        "dependencies.txt")
    if err != nil {
        return "", err
    }

    dependencyBuildString := buildCMakeString(string(dependenciesTemplate), "", depTree)

    return dependencyBuildString, nil
}

// This parses dependencies and creates cmake files for all of these. It does this for a target that is provided
// to it.
func ParseDepsAndCreateCMake(projectPath string, dependencies types.DependenciesTag) ([]*parsers.DependencyTree, error) {

    // create dependency tree
    dependencyTree, err := createDependencyTree(projectPath, dependencies)
    if err != nil {
        return nil, err
    }

    finalBuildString := ""

    for tree := range dependencyTree {
        if str, err := createDependencyCMakeString(dependencyTree[tree]);
            err != nil {
            return nil, err
        } else {
            finalBuildString += str
        }
    }

    if err := io.NormalIO.WriteFile(projectPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"dependencies.cmake",
        []byte(finalBuildString)); err != nil {
        return nil, err
    }

    return dependencyTree, nil
}
