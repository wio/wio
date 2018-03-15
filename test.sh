#!/usr/bin/env bash
set -xe # Be verbose and exit immediately if any command fails

# install wcosa
python setup.py -q install

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
mkdir src/module
cp ../test/main.cpp src/main.cpp
cp ../test/config.json config.json
cp ../test/module/* src/module/
wcosa package install waterloop/wlib@1.0.0
wcosa update
wcosa build

# Test Case 3
rm -rf test-wcosa
mkdir test-wcosa
cd test-wcosa
wcosa package install waterloop/wlib
test -d .pkg/wlib && test -d pkg/wlib # Exists and linked at pkg/wlib
wcosa package install waterloop/wlib at wlib-alt
test -d wlib-alt # Now also linked at wlib-alt
test $(readlink wlib-alt) = $(readlink pkg/wlib) # Point to same location
wcosa package remove waterloop/wlib at wlib-alt
test -d .pkg/wlib && test ! -d wlib-alt # Still exists, but not at wlib-alt
wcosa package remove waterloop/wlib
test ! -d .pkg/wlib && test ! -L pkg/wlib # Does not exist and no dangling link
