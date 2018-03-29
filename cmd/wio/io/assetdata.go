// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package contains interfaces to access data like assets and stuff related to io
// This file contains an interface to access data stored in binary
package io

import (
    "os"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/yaml.v2"
)

// interface to access Asset Data
type AssetData struct {}

// Copies asset from binary to actual location in file system
func (assetData AssetData) CopyAsset(assetPath string, destinationPath  string, override bool) (error) {
    if _, err := os.Stat(destinationPath); err == nil  && !override {
        return nil
    }

    dest, err := os.Create(destinationPath) // creates if file doesn't exist
    if err != nil {
        return err
    }
    defer dest.Close()

    srcData, err := Asset(assetPath)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(destinationPath, srcData, 0644)
    if err != nil {
        return err
    }

    err = dest.Sync()
    return err
}

// Reads data in bytes from the asset stored in binary
func (assetData AssetData) Read(fileName string) ([]byte, error) {
    return Asset(fileName)
}

// Parses JSON from the data in asset specified
func (assetData AssetData) ParseJson(fileName string, out interface{}) (err error) {
    text, err := Asset(fileName)
    if err != nil {
        return err
    }

    err = json.Unmarshal([]byte(text), out)
    return err
}

// Parses YML from the data in asset specified
func (assetData AssetData) ParseYml(fileName string, out interface{}) (error) {
    text, err := Asset(fileName)
    if err != nil {
        return err
    }

    return yaml.Unmarshal(text, out)
}
