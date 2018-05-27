// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
package main

import (
    "wio/cmd/wio/utils/io"
    "github.com/urfave/cli"
    "time"
    "os"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/pac"
    "wio/cmd/wio/utils/io/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/commands/build"
    "wio/cmd/wio/commands/clean"
    "wio/cmd/wio/commands/run"
)

func main() {
    // override help template
    cli.AppHelpTemplate =
        `Wio a simplified development process for embedded applications.
Create, Build, Test, and Upload AVR projects from Commandline.

Common Commands:
    
    wio create <project type> [options] <output directory>
        Create a new Wio project in the specified directory.
    
    wio build [options]
        Build the Wio project based on all the configurations defined
    
    wio upload [options]
        Upload the Wio project to an attached embedded device
    
    wio run [options]
        Builds, Tests, and Uploads the Wio projects

Usage: {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
Global options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
Available commands:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
Global options:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

Copyright:
   {{.Copyright}}
   {{end}}{{if .Version}}
Vesrion:
   {{.Version}}
   {{end}}
Run "wio command <help>" for more information about a command.
`

    cli.CommandHelpTemplate =
        `{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

Category:
   {{.Category}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .VisibleFlags}}

Available commands:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`

    cli.SubcommandHelpTemplate =
        `{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

Available commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`
    // get default configuration values
    defaults := types.DConfig{}
    err := io.AssetIO.ParseYml("config/defaults.yml", &defaults)
    if err != nil {
        log.Error(false, err.Error())
    }

    // command that will be executed
    var command commands.Command

    app := cli.NewApp()
    app.Name = "wio"
    app.Version = defaults.Version
    app.EnableBashCompletion = true
    app.Compiled = time.Now()
    app.Copyright = "Copyright (c) 2018 Waterloop"
    app.Usage = "Create, Build and Upload AVR projects"

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create pkg <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: false}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: false}
                    },
                },
            },
        },
        {
            Name:  "update",
            Usage: "Updates the current project and fixes any issues.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Updates a wio package, intended to be used by other people",
                    UsageText: "wio update pkg <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: true}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Updates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio update app <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: true}
                    },
                },
            },
        },
        {
            Name:      "build",
            Usage:     "Builds the wio project.",
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Build a specified target instead of building the default",
                    Value: defaults.Btarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = build.Build{Context: c}
            },
        },
        {
            Name:      "clean",
            Usage:     "Cleans all the build files for the project.",
            UsageText: "wio clean",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Cleans build files for a specified target instead of cleaning all the targets",
                    Value: defaults.Btarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = clean.Clean{Context: c}
            },
        },
        {
            Name:      "run",
            Usage:     "Builds and Uploads the project to a device (provide port flag to trigger upload)",
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = run.Run{Context: c}
            },
        },
        /*
        {
            Name:      "test",
            Usage:     "Runs unit tests available in the project.",
            UsageText: "wio test",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },

        {
            Name:      "monitor",
            Usage:     "Runs the serial monitor.",
            UsageText: "wio monitor [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "gui",
                    Usage: "Runs the GUI version of the serial monitor tool",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doctor",
            Usage:     "Guide development tooling and system configurations.",
            UsageText: "wio doctor",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "analyze",
            Usage:     "Analyzes C/C++ code statically.",
            UsageText: "wio analyze",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doxygen",
            Usage:     "Runs doxygen tool to create documentation for the code.",
            UsageText: "wio doxygen",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        */
        {
            Name:  "pac",
            Usage: "Package manager for Wio projects.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "publish",
                    Usage:     "Publish the wio package to the package manager site (npm site)",
                    UsageText: "wio pac publish [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.PUBLISH}
                    },
                },
                cli.Command{
                    Name:      "get",
                    Usage:     "Gets all the packages mentioned in wio.yml file and vendor folder.",
                    UsageText: "wio pac get [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "clean",
                            Usage: "Cleans all the current packages and re get all of them.",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.GET}
                    },
                },
                cli.Command{
                    Name:      "update",
                    Usage:     "Updates all the packages mentioned in wio.yml file and vendor folder.",
                    UsageText: "wio pac update [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.UPDATE}
                    },
                },
                cli.Command{
                    Name:      "collect",
                    Usage:     "Creates vendor folder and puts all the packages in that folder.",
                    UsageText: "wio pac collect [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "dir",
                            Usage: "Directory for the project (default: current working directory)",
                            Value: getCurrDir(),
                        },
                        cli.StringFlag{Name: "pkg",
                            Usage: "Packages to collect instead of collecting all of the packages.",
                            Value: "none",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed.",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.COLLECT}
                    },
                },
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    if err = app.Run(os.Args); err != nil {
        panic(err)
    }

    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }

        command.Execute()
    }
}

// returns the current directory from where wio is being called
func getCurrDir() (string) {
    directory, err := os.Getwd()
    commands.RecordError(err, "")
    return directory
}
