# WCosa `cmake` Toolchain
WCosa ships with a `cmake` toolchain built on top of [arduino-cmake](https://github.com/arduino-cmake/arduino-cmake)
that does Cosa-specific configuration. On the surface, the usage is similar.

## Creating Firmware
An example `CMakeLists.txt` looks like this

```cmake
set(CMAKE_TOOLCHAIN_FILE cmake/CosaToolchain.cmake)

cmake_minimum_required(VERSION 3.1.0)
project(CosaExample C CXX ASM)

generate_arduino_firmware(blink
        SRCS main.cpp
        ARDLIBS FAT16 ManchesterCodec exlib
        PORT /dev/ttyACM0
        BOARD mega2560)
```

The first line allows access to the `cosa-cmake` toolchain, and should point `CMAKE_TOOLCHAIN_FILE`
to the location of `CosaToolchain.cmake`.

The command `generate_arduino_firmware` creates a target `blink` with one source `main.cpp`
and links the libraries `FAT16`, `ManchesterCodec` and `exlib`. Arguments for `ARDLIBS` are
searched for in the Cosa libraries directory, in the current working directory, and in a `libraries`
subdirectory.

`BOARD` specifies the desired target board. If a board has multiple CPUs, the desired one
must be specified with `BOARD_CPU`. If `PORT` is specified, an upload target will be created.

In the example above, the folder structure is

```
example/
    libraries/
        exlib/
            exlib.h
            exlib.cpp
    CMakeLists.txt
    main.cpp
```
