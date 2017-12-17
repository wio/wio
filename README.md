# Waterloop Cosa Template
This is a template project for Arduino projects with the Cosa build platform. The aim is to provide an alternative to the standard Arduino libraries with Cosa, an object-oriented platform primarily in C++, whereas Arduino is ANSI C.

## Cosa
[Cosa](https://github.com/mikaelpatel/Cosa) boasts better performance and lower power consumption, while more powerful but complex than standard Arduino.

## WMake
The template project ships with `wmake`, a generic build script.

## Installation

### Linux Ubuntu
```bash
sudo apt-get install gcc-avr avr-libc avrdude
```

### MacOS
Installation may take a while because MacOS has to build `gcc-avr`.
```bash
xcode-select --install
brew tap osx-cross/avr
brew install avr-gcc
```

