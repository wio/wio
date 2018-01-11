#!/usr/bin/env bash

# install wcosa
python setup.py install

# Test Case 1
rm -rf test-wcosa
mkdir test-wcosa
cd test-wcosa
wcosa create --board mega2560
wcosa build
wcosa clean
wcosa update --board nano
wcosa build
cd ..

# Test Case 2
rm -rf test-wcosa
mkdir test-wcosa
cd test-wcosa
wcosa create --ide clion --board mega2560
git init
rm src/main.cpp
rm config.json
cp ../test/main.cpp src/main.cpp
cp ../test/config.json config.json
git submodule add https://github.com/teamwaterloop/wlib.git lib/wlib
wcosa update
wcosa build

