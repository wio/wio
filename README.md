# Waterloop Cosa

Waterloop Cosa is a commandline template project for AVR based projects. The aim of this project is to make creation of AVR projects easier. The idea is that a user will use `wcosa` to create their project structure, work in that structure and use it's functionality to update, build and upload the project. This tool is cross platform and hence call be used on `windows`, `linux` and `mac os` . User can use any text editor they want and then use this tool from command line to build and upload their project. This project integrates Cosa build platform and has no support for Arduino standard libraries. This is hence an alterantive to Arduino standard libraries. Cosa is an object-oriented platform primarily written in C++, whereas Arduino standard libraries are written in ANSI C.

## Cosa

[Cosa](https://github.com/mikaelpatel/Cosa) boasts better performance and lower power consumption, while more powerful but complex than standard Arduino libraries. Cosa, being object oriented integrates well with other C++ program written in OOP style.

## WCosa

The template project ships with `wcosa`, a generic build script. All the scripts are written in python, making them cross platform. Python 3 is recommended for running these scripts. User should only call `wcosa.py` script and that in return will allow for create, update, build, upload and serial functionality for AVR projects. This tool uses CMake files behind the scenes to do everything. When this tool is run it creates `src`, `lib`, `wcosa` folders and `.gitignore`. Files and folders that need to be ignored are added to `.gitignore` file. `wcosa` will create the project in the folder the script was called from. This means a user can create a folder for their project, call this script and they have a project created and initialized. This however can be changed by using command line option `--path` or `-p`. This allows user to specify the path where the project should be created.  The project structure in bit more detail:

### scr

Source files go here. This includes all the header files as well

### lib

All 3rd party libraries go here. They must be in their seperate folder. It supports libraries that have no **src** folder and all files are in the root folder as well as libraries that have **src** folder to contains all the library files

### wcosa

All necessary build files used internally. This folder should not be modified as this may cause the build and update process to break. This folder should also be included in `.gitignore` as it contains machine specific details.

## Installation

### Linux

**Ubuntu**

```bash
sudo apt-get install gcc-avr avr-libc avrdude
```

**Arch**

```bash
sudo pacman -S avr-gcc avr-libc avrdue
```

### MacOS

Installation may take a while because MacOS has to build `gcc-avr`.

```bash
xcode-select --install
brew tap osx-cross/avr
brew install avr-gcc
```

This project is under construction and no stable version is available yet. When this project is closer to a stable version, more information will be provided on how to install and use this project. Installation procedures will be provided for more operating systems in future and more detail will be provided for what exactly does this tool do.
