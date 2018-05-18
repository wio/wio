// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands package, which contains all the commands provided by the tool.
// Uploads the project to a device
package upload

import (
    "github.com/urfave/cli"
    "github.com/GuoXi/go-serial/enumerator"
    "fmt"
    "github.com/google/gousb/usbid"
    "github.com/google/gousb"
)

type Upload struct {
    Context *cli.Context
    error
}

// get context for the command
func (upload Upload) GetContext() (*cli.Context) {
    return upload.Context
}

// Runs the build command when cli build option is provided
func (upload Upload) Execute() {
    RunUpload(upload.Context.String("dir"), upload.Context.String("target"),
        upload.Context.String("port"))
}

func getDeviceDescription (pid string) (string, error) {
    ctx := gousb.NewContext()
    defer ctx.Close()

    description := ""

    _, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
        if pid == desc.Product.String() {
            description = usbid.Describe(desc)
        }

        return false
    })

    return description, err
}


// This function allows other packages to call build as well. This is also used when cli build is executed
func RunUpload(directoryCli string, targetCli string, port string) {
    /*
    ctx := gousb.NewContext()
    defer ctx.Close()

    // OpenDevices is used to find the devices to open.
    devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
        // The usbid package can be used to print out human readable information.
        fmt.Printf("%03d.%03d %s:%s %s\n", desc.Bus, desc.Address, desc.Vendor, desc.Product, usbid.Describe(desc))
        fmt.Printf("  Protocol: %s\n", usbid.Classify(desc))

        fmt.Println(desc.Port)

        // The configurations can be examined from the DeviceDesc, though they can only
        // be set once the device is opened.  All configuration references must be closed,
        // to free up the memory in libusb.
        for _, cfg := range desc.Configs {
            // This loop just uses more of the built-in and usbid pretty printing to list
            // the USB devices.
            fmt.Printf("  %s:\n", cfg)
            for _, intf := range cfg.Interfaces {
                fmt.Printf("    --------------\n")
                for _, ifSetting := range intf.AltSettings {
                    fmt.Printf("    %s\n", ifSetting)
                    fmt.Printf("      %s\n", usbid.Classify(ifSetting))
                    for _, end := range ifSetting.Endpoints {
                        fmt.Printf("      %s\n", end)
                    }
                }
            }
            fmt.Printf("    --------------\n")
        }

        // After inspecting the descriptor, return true or false depending on whether
        // the device is "interesting" or not.  Any descriptor for which true is returned
        // opens a Device which is retuned in a slice (and must be subsequently closed).
        return false
    })

    // All Devices returned from OpenDevices must be closed.
    defer func() {
        for _, d := range devs {
            d.Close()
        }
    }()

    // OpenDevices can occaionally fail, so be sure to check its return value.
    if err != nil {
        log.Fatalf("list: %s", err)
    }

    for _, dev := range devs {
        // Once the device has been selected from OpenDevices, it is opened
        // and can be interacted with.
        _ = dev
    }
    */


    ports, _ := enumerator.GetDetailedPortsList()

    for _, port := range ports {
        if port.PID != "" {
            desc, err := getDeviceDescription(port.PID)
            if err != nil {
                panic(err)
            } else {
                fmt.Println(desc)
            }
        }
    }



    /*
    directory, err := filepath.Abs(directoryCli)
    commands.RecordError(err, "")

    // find a right port

    // run build
    build.RunBuild(directoryCli, targetCli, false, port)

    // then run the upload target from the make
    targetsDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"
    targetPath := targetsDirectory + io.Sep + targetCli

    cmakeCommand := exec.Command("make", "upload")
    cmakeCommand.Dir = targetPath
    cmakeCommand.Stdout = os.Stdout
    if err := cmakeCommand.Run(); err != nil {
        os.Stderr.WriteString(err.Error())
    }

    commands.RecordError(err, "")
    */
}
