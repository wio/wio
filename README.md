[![Build Status](https://travis-ci.org/waterloop/wcosa.svg?branch=master)](https://travis-ci.org/waterloop/wcosa)

# Waterloop Cosa

Waterloop Cosa is a commandline tool used to help create, build, and upload AVR projects built using Cosa, an object-oriented platform for Arduino. This project is built on top of 
[arduino-cmake](https://github.com/arduino-cmake/arduino-cmake) to provide a CMake toolchain for Cosa,
and wrapped in a Python script.

## Cosa

[Cosa](https://github.com/mikaelpatel/Cosa) boasts better performance and lower power consumption, while more powerful but complex than standard Arduino libraries. Cosa, being object oriented integrates well with other C++ program written in OOP style.

## Wcosa -> Wio
This project was used to known as `wcosa` but, our team had decided to completely rewamp the project
and hence change the name to `wio`. `wio` is under development right now now and when it is ready, it
will replace the old tool. Please stay tuned as development progresses

For `wcosa` tool, check it out here: [wcosa](https://github.com/waterloop/wcosa)


## Development
If anyone is interested in working on this projects, follow the following:

**Getting Setup**
* Clone this repo:
```
git clone https://www.github.com/dhillondeep/wio
```
* Create a feature branch and checkout that branch
* Install dependencies:
```
make deps
```
* Build project:
```
make build
```
* Run project:
```
// this without any arguments
make run

// with arguments
make run ARGS="create -h"
```

**Commands**

All the commands are split into seperate files inside `internal/commands` package.
You can write code for each command there and the main file in `cmd/wio` will pick up the changes.
Every command has an **execute** function. This function is an entry point and is provided will all the information.



## Committers
Deep Dhillon (@dhillondeep)

Jeff Niu (@mogball)

## Contributors
Dmytro Shynkevych (@hermord)
