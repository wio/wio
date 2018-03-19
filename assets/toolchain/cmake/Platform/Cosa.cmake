cmake_minimum_required(VERSION 3.1.0)

if (NOT COSA_SDK_PATH)
    message(FATAL_ERROR "Error: COSA_SDK_PATH is not defined")
endif ()
if (NOT ARDUINO_CMAKE_PATH)
    message(FATAL_ERROR "Error: ARDUINO_CMAKE_PATH is not defined")
endif ()

if (COSA_SCRIPT_EXECUTED)
    return()
endif ()

# Set module paths to include `arduino-cmake` components
set(ARDUINO_CMAKE_PLATFORM_PATH ${ARDUINO_CMAKE_PATH}/Platform)
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH}
        ${ARDUINO_CMAKE_PLATFORM_PATH}
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Initialization
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/BoardFlags
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Libraries
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Targets
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Sketch
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Core/Examples
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Extras
        ${ARDUINO_CMAKE_PLATFORM_PATH}/Generation)

# Set module paths to include `cosa-cmake` components
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH}
        ${CMAKE_CURRENT_LIST_DIR}/Initialization
        ${CMAKE_CURRENT_LIST_DIR}/Pollute
        ${CMAKE_CURRENT_LIST_DIR}/Utils
        ${CMAKE_CURRENT_LIST_DIR}/Vendor)

# Include vendored files
include(JsonParser)

# Include utilities
include(CosaOutput)

# If on Windows, ensure that ARDUINO_SDK_PATH is specified
if (WIN32 OR CYGWIN OR MINGW)
    if (NOT ARDUINO_SDK_PATH OR NOT EXISTS ${ARDUINO_SDK_PATH})
        fatal("Windows systems must specify ARDUINO_SDK_PATH")
    endif ()
endif ()

# Include external utilities
include(CMakeParseArguments)
include(VariableValidator)

# Initialization scripts
include(CosaInitializer)

# Include all scripts from `arduino-cmake`
include(ArduinoCmakeScripts)

write_sep()

# Mark configuration as complete
set(COSA_SCRIPT_EXECUTED True)
