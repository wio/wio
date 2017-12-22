set(CMAKE_TOOLCHAIN_FILE cmake/CosaToolchain.cmake) # Cosa Toolchain

cmake_minimum_required(VERSION 3.1.0)

project(CosaExample C CXX ASM)

add_subdirectory(example)
