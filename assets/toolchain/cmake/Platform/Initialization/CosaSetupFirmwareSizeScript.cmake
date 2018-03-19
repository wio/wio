# Find the script provided by `arduino-cmake`
find_file(firmware_size_script_path
        NAMES CalculateFirmwareSize.cmake
        PATHS ${ARDUINO_CMAKE_PATH}/Platform/Extras)

if (NOT firmware_size_script_path OR NOT EXISTS ${firmware_size_script_path})
    fatal("Unable to find path to template FIRMWARE_SIZE_SCRIPT")
endif ()

if (NOT COSA_AVRSIZE_PROGRAM OR NOT EXISTS ${COSA_AVRSIZE_PROGRAM})
    fatal("COSA_AVRSIZE_PROGRAM must be specified")
endif ()

# Replace the placeholder with the determined `avrsize` path
file(READ ${firmware_size_script_path} firmware_size_script_template)
string(REGEX REPLACE "PLACEHOLDER_1" "${COSA_AVRSIZE_PROGRAM}" firmware_size_script ${firmware_size_script_template})
set(COSA_FIRMWARE_SCRIPT_PATH ${CMAKE_BINARY_DIR}/CMakeFiles/FirmwareSize.cmake)
file(WRITE ${COSA_FIRMWARE_SCRIPT_PATH} "${firmware_size_script}")

write_sep()
info("Template: ${firmware_size_script_path}")
info("Cached:   ${COSA_FIRMWARE_SCRIPT_PATH}")

unset(firmware_size_script_path)
unset(firmware_size_script_template)
unset(firmware_size_script)
