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
    config.json
    CMakeLists.txt
    CMakeListsPrivate.txt
```

Sources files should be placed in the `src` directory, and libraries should be placed in folders inside
the `lib` directory. Build files are contained in `wcosa` and needs to be generated for each environment
that is running the WCosa project.

## Commands and Usage
Using commands has the format `wcosa [action]` where `[action]` is one of `create`, `update`, `build`, `clean`, `upload`.

Here is an example of creating a project:
```bash
mkdir wcosa-project
cd wcosa-project
wcosa create --board uno --ide clion
```
The `--board` flag specifies the desired target board for project build, and `--ide clion` tells WCosa to generate files to enable 
project import into `CLion` for code suggestion and completion. Use `wcosa boards` to see the list of available boards, supported Cosa.
This command will generate environment-specific files such as `CMakeListsPrivate.txt` and the `wcosa` folder.

To build the project, creating uploadable binaries, use `wcosa build`, and to clean the project, run `wcosa clean`. If a microcontroller
is plugged into your PC, running `wcosa upload` will attempt to autodetect the port and upload the program.
If you have multiple microcontrollers or if WCosa cannot detect the port, use `wcosa upload --port [port_name]` to specify
the desired upload port.

The command `wcosa update` can be used to update the target board with `wcosa update --board [new_board]` or if `config.json` is 
modified. To add build definitions to the main project or to submodules, modify `build-flags` to something as
```json
{
    "build-flags": "MAX_ALLOCATORS=4u STATIC_MEMORY"
}
```

To specify definitions for modules added under `lib`, 
```json
{
    "build-flags": "...",
    "module-flags": {
        "wlib": "BLOCK_SIZE=64u NUM_BLOCKS=400u"
    }
}
```
Then run `wcosa update` to update the internal config.

The command `wcosa update` should also be called when checking out a project from source control to create the environment
files. For example,
```bash
git clone --recursive https://github.com/teamwaterloop/goose-sensors.git
cd goose-sensors
wcosa update
```

## Installation
```bash
pip install wcosa
```

WCosa requires either `gcc-avr` or the Arduino SDK to be installed. __CMake__ is also required to build projects.

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
brew install avrdude
```

## Committers
Deep Dhillon (@dhillondeep)
Jeff Niu (@mogball)

## Contributors
Dmytro Shynkevych (@hermord)

