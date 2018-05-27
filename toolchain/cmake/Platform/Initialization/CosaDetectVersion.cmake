#=============================================================================#
# Detects the Cosa SDK Version based on its contained `package_cosa_index.json`
# file. The results are stored in the following variables
#
#    ${COSA_SDK_VERSION}         -> the full version (major.minor.patch)
#    ${COSA_SDK_VERSION}_MAJOR   -> the major version
#    ${COSA_SDK_VERSION}_MINOR   -> the minor version
#    ${COSA_SDK_VERSION}_PATCH   -> the patch version
#
#=============================================================================#

info("Determining `cosa` version")

find_file(PACKAGE_COSA_INDEX_PATH
        NAMES package_cosa_index.json
        PATHS ${COSA_SDK_PATH})

if (NOT PACKAGE_COSA_INDEX_PATH)
    fatal("Failed to find `package_cosa_index.json` in COSA_SDK_PATH")
endif ()

info("Found package file: ${PACKAGE_COSA_INDEX_PATH}")

# Read in the file and parse JSON
file(READ ${PACKAGE_COSA_INDEX_PATH} package_cosa_index_raw)
sbeParseJson(package_cosa_index_json package_cosa_index_raw)

# Lower index in the index file represent newer versions
set(version_access_string package_cosa_index_json.packages_0.platforms_0.version)
set(COSA_SDK_VERSION ${${version_access_string}} CACHE STRING "")

# Get version parts
string(REPLACE "." ";" cosa_sdk_version_split ${COSA_SDK_VERSION})
list(GET cosa_sdk_version_split 0 cosa_sdk_version_major)
list(GET cosa_sdk_version_split 1 cosa_sdk_version_minor)
list(GET cosa_sdk_version_split 2 cosa_sdk_version_patch)
set(COSA_SDK_VERSION_MAJOR ${cosa_sdk_version_major} CACHE STRING "")
set(COSA_SDK_VERSION_MINOR ${cosa_sdk_version_minor} CACHE STRING "")
set(COSA_SDK_VERSION_PATCH ${cosa_sdk_version_patch} CACHE STRING "")

info("Identified `cosa` version: ${COSA_SDK_VERSION_MAJOR}.${COSA_SDK_VERSION_MINOR}.${COSA_SDK_VERSION_PATCH}")

# Clean up variables
sbeClearJson(package_cosa_index_json)
unset(version_access_string)
unset(package_cosa_index_raw)
unset(cosa_sdk_version_split)
unset(cosa_sdk_version_major)
unset(cosa_sdk_version_minor)
unset(cosa_sdk_version_patch)
