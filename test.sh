#!/usr/bin/env bash

python setup.py install
rm -rf test-wcosa
mkdir test-wcosa
wcosa create --board mega2560
wcosa build
