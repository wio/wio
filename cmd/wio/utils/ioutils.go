package utils

import (
    "io/ioutil"
    "runtime"
    "path/filepath"

    "gopkg.in/yaml.v2"
    "os"
    "io"
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

// Copies file from src to dist and if dest file exists, it overrides the file
// content based on if override is specified
func Copy(src string, dest string, override bool) (error) {
    if _, err := os.Stat(dest); err == nil  && !override {
        return nil
    }

    srcFile, err := os.Open(src)

    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(dest) // creates if file doesn't exist
    if err != nil {
        return err
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
    if err != nil {
        return err
    }

    err = destFile.Sync()
    if err != nil {
        return err
    }

    return nil
}

// Get's the path to the root directory of this project
func GetExecutableRootPath() (string, error) {
    _, configFileName, _, _ := runtime.Caller(0)
    return filepath.Abs(configFileName + "/../../../../")
}

