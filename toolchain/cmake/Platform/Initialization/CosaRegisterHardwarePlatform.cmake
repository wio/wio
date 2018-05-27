write_sep()
info("Obtaining hardware information from ${COSA_SDK_PATH}")

# Find paths to directories and files
# Cosa does not provide a `programmers.txt`
set(COSA_CORES_PATH ${COSA_SDK_PATH}/cores)
set(COSA_VARIANTS_PATH ${COSA_SDK_PATH}/variants/arduino)
set(COSA_BOOTLOADERS_PATH ${COSA_SDK_PATH}/bootloaders)
set(COSA_BOARDS_PATH ${COSA_SDK_PATH}/boards.txt)

info("Founds paths")
info("Cores:       ${COSA_CORES_PATH}")
info("Variants:    ${COSA_VARIANTS_PATH}")
info("Bootloaders: ${COSA_BOOTLOADERS_PATH}")
info("Boards:      ${COSA_BOARDS_PATH}")

if (NOT COSA_CORES_PATH OR NOT EXISTS ${COSA_CORES_PATH})
    fatal("Failed to find COSA_CORES_PATH to `cores`")
endif ()
if (NOT COSA_VARIANTS_PATH OR NOT EXISTS ${COSA_VARIANTS_PATH})
    fatal("Failed to find COSA_VARIANTS_PATH to `variants/arduino`")
endif ()
if (NOT COSA_BOOTLOADERS_PATH OR NOT EXISTS ${COSA_BOOTLOADERS_PATH})
    fatal("Failed to find COSA_BOOTLOADERS_PATH to `bootloaders`")
endif ()
if (NOT COSA_BOARDS_PATH OR NOT EXISTS ${COSA_BOARDS_PATH})
    fatal("Failed to find COSA_BOARDS_PATH to `boards.txt`")
endif ()

# Read in `boards.txt`
set(SETTINGS_LIST COSA_BOARDS)
set(SETTINGS_PATH ${COSA_BOARDS_PATH})
include(LoadArduinoPlatformSettings)

# Display example boards read
list(GET COSA_BOARDS 0 example_board_0)
list(GET COSA_BOARDS 1 example_board_1)
list(GET COSA_BOARDS 2 example_board_2)
info("Parsed `boards.txt` (e.g. ${example_board_0}, ${example_board_1}, ${example_board_2})")
unset(example_board_0)
unset(example_board_1)
unset(example_board_2)

# Have to manually set
# ${cosa_board}.build.core=cosa
# ${cosa_board}.build.variant=${cosa_board}
foreach (cosa_board ${COSA_BOARDS})
    set(${cosa_board}.build.core cosa)
    set(${cosa_board}.build.variant ${cosa_board})
endforeach ()

# Read in variant boards
info("Reading variants from directory")
file(GLOB variant_sub_dir ${COSA_VARIANTS_PATH}/*)
unset(COSA_VARIANTS CACHE)
foreach (variant_dir ${variant_sub_dir})
    if (IS_DIRECTORY ${variant_dir})
        get_filename_component(variant ${variant_dir} NAME)
        set(COSA_VARIANTS ${COSA_VARIANTS} ${variant} CACHE INTERNAL "A list of registered variants")
        set(${variant}.path ${variant_dir} CACHE INTERNAL "The path to variant ${variant}")
        if (WCOSA_DEBUG)
            info("Variant [${variant}]: ${${variant}.path}")
        endif ()
    endif ()
endforeach ()
list(LENGTH COSA_VARIANTS length_cosa_variants)
info("Found and cached ${length_cosa_variants} variants")
unset(variant_sub_dir)
unset(length_cosa_variants)

# Read in cores
info("Reading cores from directory")
file(GLOB cores_sub_dir ${COSA_CORES_PATH}/*)
unset(COSA_CORES CACHE)
foreach (core_dir ${cores_sub_dir})
    if (IS_DIRECTORY ${core_dir})
        get_filename_component(core ${core_dir} NAME)
        set(COSA_CORES ${COSA_CORES} ${core} CACHE INTERNAL "A list of registered cores")
        set(${core}.path ${core_dir} CACHE INTERNAL "The path to core ${core}")
        if (WCOSA_DEBUG)
            info("Core [${core}]: ${${core}.path}")
        endif ()
    endif ()
endforeach ()
list(LENGTH COSA_CORES length_cosa_cores)
info("Found and cached ${length_cosa_cores} cores")
unset(cores_sub_dir)
unset(length_cosa_cores)
