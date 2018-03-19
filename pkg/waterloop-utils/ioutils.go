package waterloop_utils

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

// Reads the file and provides it's content as a string
func FileToString(fileName string) (string, error) {
    buff, err := ioutil.ReadFile(fileName)
    str := string(buff)

    return str, err
}

// Converts a String to Yml struct
func ToYmlStruct(data string, out interface{}) (error) {
    e := yaml.Unmarshal([]byte(data), out)

    return e
}
