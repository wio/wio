package ioutils

import (
    "io/ioutil"
    "runtime"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

const (
    WINDOWS string = "windows"
    DARWIN  string = "darwin"
    LINUX   string = "linux"
)

// Reads the file and provides it's content as a string
func FileToString(fileName string) (string, error) {
    fileName, _ = GetPath(fileName)
    buff, err := ioutil.ReadFile(fileName)
    str := string(buff)

    return str, err
}

// Writes string data to a file
func StringToFile(fileName string, data string) {
    ioutil.WriteFile(fileName, []byte(data), 0064)
}

// Returns operating system from three types (windows, darwin, and linux)
func GetOS() (string) {
    goos := runtime.GOOS

    if goos == "windows" {
        return WINDOWS
    } else if goos == "darwin" {
        return DARWIN
    } else {
        return LINUX
    }
}

// Converts the path provided into operating system preferred path
func GetPath(currPath string) (string, error) {
    return filepath.Abs(currPath)
}

// Converts a String to Yml struct
func ToYmlStruct(data string, out interface{}) (error) {
    e := yaml.Unmarshal([]byte(data), out)

    return e
}

