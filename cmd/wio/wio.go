//+build !test

package main

import (
	"fmt"
	"github.com/dhillondeep/afero"
	"log"
	"wio/internal/config"
	"wio/internal/constants"
	"wio/internal/evaluators/hillang"
	"wio/pkg/sys"
	"wio/templates"
)

func main() {
	sys.SetFileSystem(afero.NewMemMapFs())

	if err := config.CreateConfig(templates.ProjectCreation{
		Type:        constants.APP,
		ProjectName: "Wio",
		ProjectPath: "/projects",
		Platform:    "native",
	}); err != nil {
		log.Fatalf(err.Error())
	}

	configFile, warn, err := config.ReadConfig("/projects")
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(warn) > 0 {
		fmt.Println("Warnings...")
		fmt.Println(warn)
	}

	if err := hillang.Initialize(configFile.GetVariables(), configFile.GetArguments()); err != nil {
		log.Fatal(err.Error())
	}

	projectName, err := configFile.GetProject().GetName(hillang.GetDefaultEvalConfig())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Project Name: " + projectName)
}
