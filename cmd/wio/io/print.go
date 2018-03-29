// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package contains interfaces to access data like assets and stuff related to io
// This file contains an interface to print output to io in various colors and modes
package io

import (
    "fmt"
    "os"

    "github.com/shiena/ansicolor"
    "wio/cmd/wio/utils"
    "github.com/fatih/color"
)

var Norm = writer{verbose:true}
var Verb = writer{verbose:false}
var w = ansicolor.NewAnsiColorWriter(os.Stdout)

// user should not touch this
type writer struct {
    verbose bool
}

// Turns verbose mode on. This is the mode when Vprintf and Vprintln functions work
func SetVerbose() {
    Verb.verbose = true
}

// Black is a convenient helper function to print with black foreground.
func (writer writer) Black(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[30m", "\x1b[39m")
    } else {
        fmt.Print(color.BlackString(format, a...))
    }
}

// Red is a convenient helper function to print with red foreground.
func (writer writer) Red(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[31m", "\x1b[39m")
    } else {
        fmt.Print(color.RedString(format, a...))
    }
}

// Green is a convenient helper function to print with green foreground.
func (writer writer) Green(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[32m", "\x1b[39m")
    } else {
        fmt.Print(color.GreenString(format, a...))
    }
}

// Yellow is a convenient helper function to print with yellow foreground.
func (writer writer) Yellow(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[33m", "\x1b[39m")
    } else {
        fmt.Print(color.YellowString(format, a...))
    }
}

// Blue is a convenient helper function to print with blue foreground.
func (writer writer) Blue(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[34m", "\x1b[39m")
    } else {
        fmt.Print(color.BlueString(format, a...))
    }
}

// Magenta is a convenient helper function to print with magenta foreground.
func (writer writer) Magenta(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[35m", "\x1b[39m")
    } else {
        fmt.Print(color.MagentaString(format, a...))
    }
}

// Cyan is a convenient helper function to print with cyan foreground.
func (writer writer) Cyan(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[36m", "\x1b[39m")
    } else {
        fmt.Print(color.CyanString(format, a...))
    }
}

// White is a convenient helper function to print with white foreground.
func (writer writer) White(format string, a ...interface{}) {
    if !writer.verbose { return }

    if utils.GetOS() == utils.WINDOWS {
        text := "%s" + format + "%s"
        fmt.Fprintf(w, text, "\x1b[37m", "\x1b[39m")
    } else {
        fmt.Print(color.WhiteString(format, a...))
    }
}

// Normal is a convenient helper function to print with default/normal foreground.
func (writer writer) Normal(format string, a ...interface{}) {
    if !writer.verbose { return }

    text := "%s" + format + "%s"
    fmt.Fprintf(w, text, "\x1b[39m", "\x1b[39m")
}

// Special function to be used when printing error logs.
// It terminates the program after printing the logs
func (writer writer) Error(format string, a ...interface{})  {
    if !writer.verbose {
        Norm.Red("Turn Verbose mode to see the detailed error\n")
        os.Exit(2)
    }

    writer.Red("\nError Report: \n")
    writer.Normal(format, a...)
    os.Exit(2)
}

// Special function to be used when using Verbose mode.
// In this mode, color can be set and other verbose default things can be defined
func (writer writer) Verbose(format string, a ...interface{}) {
    if !writer.verbose { return }

    writer.Normal(format, a)
}
