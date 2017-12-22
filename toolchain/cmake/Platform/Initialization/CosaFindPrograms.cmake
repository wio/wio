write_sep()

# Find examples and libraries
SET(COSA_EXAMPLES_PATH "${COSA_SDK_PATH}/examples")
SET(COSA_LIBRARIES_PATH "${COSA_SDK_PATH}/libraries")

if (NOT COSA_EXAMPLES_PATH OR NOT EXISTS ${COSA_EXAMPLES_PATH})
    fatal("Failed to find COSA_EXAMPLES_PATH to `examples`")
endif ()
if (NOT COSA_LIBRARIES_PATH OR NOT EXISTS ${COSA_LIBRARIES_PATH})
    fatal("Failed to find COSA_LIBRARIES_PATH to `libraries`")
endif ()

#===================================================#
# Search paths for `avrdude`
# Keep a log of lists here for posterity
#
# MacOS (commandline)
# - /usr/local/bin/avrdude
#
#===================================================#

# If ARDUINO_SDK_PATH is provided, search there first
if (EXISTS ${ARDUINO_SDK_PATH})
    find_program(COSA_AVRDUDE_PROGRAM
            NAMES avrdude
            PATHS ${ARDUINO_SDK_PATH}
            PATH_SUFFIXES hardware/tools hardware/tools/avr/bin
            NO_DEFAULT_PATH)
endif ()

# Search known paths first
set(cosa_avrdude_known_paths
        /usr/bin
        /usr/local/bin
        /usr/local/Cellar/avrdude/6.3/bin)
find_program(COSA_AVRDUDE_PROGRAM
        NAMES avrdude
        PATHS ${cosa_avrdude_known_paths}
        DOC "Path to avrdude programmer binary.")

# Search through environment PATH
find_program(COSA_AVRDUDE_PROGRAM
        NAMES avrdude
        DOC "Path to avrdude programmer binary.")

if (NOT COSA_AVRDUDE_PROGRAM OR NOT EXISTS ${COSA_AVRDUDE_PROGRAM})
    fatal("Unable to find path to `avrdude`")
endif ()

#===================================================#
# Search paths for `avr-size`
# Keep a log of lists here for posterity
#
# MacOS (commandline)
# - /usr/local/bin/avr-size
#
#===================================================#

if (EXISTS ${ARDUINO_SDK_PATH})
    find_program(COSA_AVRSIZE_PROGRAM
            names avr-size
            PATHS ${ARDUINO_SDK_PATH}
            PATH_SUFFIXES hardware/tools hardware/tools/avr/bin
            NO_DEFAULT_PATH)
endif ()

set(cosa_avr_binutils_paths
        /usr/bin
        /usr/local/bin
        /usr/local/Cellar/avr-binutils/2.29/bin)

find_program(COSA_AVRSIZE_PROGRAM
        names avr-size
        PATHS ${cosa_avr_binutils_paths}
        DOC "Path to avr-size program binary.")

find_program(COSA_AVRSIZE_PROGRAM
        names avr-size
        DOC "Path to avr-size program binary.")

if (NOT COSA_AVRSIZE_PROGRAM OR NOT EXISTS ${COSA_AVRSIZE_PROGRAM})
    fatal("Unable to find path to `avr-size`")
endif ()

# Cosa ships with its own `avrdude.conf`
set(COSA_AVRDUDE_CONFIG_PATH ${COSA_SDK_PATH}/build/avrdude.conf)
# Else try to find it elsewhere
if (NOT COSA_AVRDUDE_CONFIG_PATH OR NOT EXISTS ${COSA_AVRDUDE_CONFIG_PATH})
    warning("Unable to find `avrdude.conf` in `cosa`")
    warning("Searching in default locations instead")
    find_file(COSA_AVRDUDE_CONFIG_PATH
            NAMES avrdude.conf
            PATHS ${ARDUINO_SDK_PATH} /etc /etc/avrdude
            PATH_SUFFIXES hardware/tools hardware/tools/avr/etc
            DOC "Path to avrdude programmer configuration file.")
endif ()

# Cosa ships with boilerplate `Arduino.h` and `Cosa.h` headers
set(HEADER_COSA_H ${COSA_SDK_PATH}/cores/cosa/Cosa.h)
set(HEADER_ARDUINO_H ${COSA_SDK_PATH}/cores/cosa/Arduino.h)
if (NOT HEADER_COSA_H OR NOT EXISTS ${HEADER_COSA_H})
    warning("Unable to find HEADER_COSA_H")
endif ()
if (NOT HEADER_ARDUINO_H OR NOT EXISTS ${HEADER_ARDUINO_H})
    warning("Unable to find HEADER_ARDUINO_H")
endif ()

# Check for CMAKE_OBJCOPY
if (NOT CMAKE_OBJCOPY OR NOT EXISTS ${CMAKE_OBJCOPY})
    find_program(COSA_AVROBJCOPY_PROGRAM
            NAMES avr-objcopy
            PATHS cosa_avr_binutils_paths
            DOC "Path to avr-objcopy binary."
            NO_DEFAULT_PATH)
    find_program(COSA_AVROBJCOPY_PROGRAM
            NAMES avr-objcopy
            DOC "Path to avr-objcopy binary.")
    set(CMAKE_OBJCOPY ${COSA_AVROBJCOPY_PROGRAM})
else ()
    set(COSA_AVROBJCOPY_PROGRAM ${CMAKE_OBJCOPY})
endif ()

if (NOT CMAKE_OBJCOPY OR NOT EXISTS ${CMAKE_OBJCOPY})
    fatal("Failed to find `avr-objcopy`")
endif ()

info("Found paths and programs")
info("avrdude:      ${COSA_AVRDUDE_PROGRAM}")
info("avr-size:     ${COSA_AVRSIZE_PROGRAM}")
info("avr-objcopy:  ${COSA_AVROBJCOPY_PROGRAM}")
info("Examples:     ${COSA_EXAMPLES_PATH}")
info("Libraries:    ${COSA_LIBRARIES_PATH}")
info("avrdude.conf: ${COSA_AVRDUDE_CONFIG_PATH}")
info("Cosa.h:       ${HEADER_COSA_H}")
info("Arduino.h:    ${HEADER_ARDUINO_H}")
