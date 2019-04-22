package resources

import (
	"errors"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"os"
)

var box packd.Box = nil

func Initialize(path string) {
	box = packr.NewBox(path)
}

// CopyFile copies file from embedded resources to system fs. If destination file is present,
// it overrides the file content based on if override is specified
func CopyFile(source string, dest string, override bool) error {
	if _, err := os.Stat(dest); err == nil && !override {
		return nil
	}

	destFile, err := os.Create(dest) // creates if file doesn't exist
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcData, err := ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, srcData, os.ModePerm)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	return err
}

// MustCopy copies file from embedded resources to destination and override if destination file exists
func MustCopyFile(source string, dest string) error {
	return CopyFile(source, dest, true)
}

// CopyMultipleFiles copies multiple files from embedded resources to destination
func CopyMultipleFiles(sources []string, dests []string, overrides []bool) error {
	if len(sources) != len(dests) || len(dests) != len(overrides) {
		return errors.New("length of sources, destinations and overrides is not equal")
	}

	for i := 0; i < len(sources); i++ {
		if err := CopyFile(sources[i], dests[i], overrides[i]); err != nil {
			return err
		}
	}

	return nil
}

// ReadFile reads the embedded resource file and provides it's content
func ReadFile(path string) ([]byte, error) {
	return box.Find(path)
}
