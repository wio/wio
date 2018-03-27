package create

import (
    "path/filepath"
)

var sep string = string(filepath.Separator)
var rootPath string = ""

func GetRoot() (string){
    return rootPath
}

// Returns a path to template files. User needs to provide a relative path to get the full path
func GetTemplatesRelativeFile(fileName string) (string) {
    root := GetRoot()

    return root + string(filepath.Separator) + "assets" + string(filepath.Separator) + "templates" +
        string(filepath.Separator) + fileName
}
