# Waterloop Cosa

Waterloop Cosa is a commandline tool used to help create, build, and upload AVR projects built using Cosa, an object-oriented platform for Arduino. This project is built on top of 
[arduino-cmake](https://github.com/arduino-cmake/arduino-cmake) to provide a CMake toolchain for Cosa,
and wrapped in a Python script.

## Cosa

[Cosa](https://github.com/mikaelpatel/Cosa) boasts better performance and lower power consumption, while more powerful but complex than standard Arduino libraries. Cosa, being object oriented integrates well with other C++ program written in OOP style.

## WCosa

This project provides `wcosa`, a build script written in Python. The script allows user to `create`, 
`build`, `upload`, and `monitor` AVR projects. This tool uses the `cmake` toolchain behind the scenes.
Running the creation scripts generates a project with the structure

```
project/
    lib/
    src/
    wcosa/
        bin/
        CMakeLists.txt
    .gitignore
    CMakeLists.txt
```

Sources files should be placed in the `src` directory, and libraries should be placed in folders inside
the `lib` directory. Build files are contained in `wcosa` and needs to be generated for each environment
that is running the WCosa project.

## Installation
WCosa requires either `gcc-avr` or the Arduino SDK to be installed.i

### Windows
We recommend installing the Arduino SDK.
1. Download and install the [Arduino IDE](https://www.arduino.cc/en/Main/Software)
2. Add the Arduino installation directory and the subdirectory `\hardware\tools\avr\bin`
to your System PATH; these may look like
    * `C:\Program Files (x86)\Arduino`
    * `C:\Program Files (x86)\Arduino\hardware\tools\bin`

There are some [avr-gcc](http://blog.zakkemble.co.uk/avr-gcc-builds/) builds available for Windows
but these are untested.

### Linux
You may choose to install the Arduino SDK or the required tools and binaries
from the commandline.

**Ubuntu**

```bash
sudo apt-get install gcc-avr avr-libc avrdude
```

**Arch**

```bash
sudo pacman -S avr-gcc avr-libc avrdude
```

### MacOS
You may install the Arduino SDK or build `avr-gcc` using `brew`. Keep in mind that
building `avr-gcc` may take some time.

```bash
xcode-select --install
brew tap osx-cross/avr
brew install avr-gcc
```

