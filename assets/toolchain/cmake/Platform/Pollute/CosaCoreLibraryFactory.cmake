# We have to pollute the cmake namespace so that Board.h
# for the specified kind is pulled in correctly
function(make_core_library OUTPUT_VAR BOARD_ID)
    set(CORE_LIB_NAME ${BOARD_ID}_CORE)
    _get_board_property(${BOARD_ID} build.core BOARD_CORE)
    # Grab the variant type
    _get_board_property(${BOARD_ID} build.variant VARIANT)
    # Use the variant path to find its header file
    find_file(HEADER_BOARD_H
            NAMES Board.h
            PATHS ${${VARIANT}.path})
    info("Board.h: ${HEADER_BOARD_H}")
    # Ensure that the header file exists
    if (NOT HEADER_BOARD_H OR NOT EXISTS ${HEADER_BOARD_H})
        fatal("Failed to find `Board.h` for variant ${VARIANT}")
    endif ()
    # Include the directory so that the board defintions flow
    get_filename_component(HEADER_BOARD_H_DIRECTORY ${HEADER_BOARD_H} DIRECTORY)
    include_directories(${HEADER_BOARD_H_DIRECTORY})
    # Add the core library
    if (BOARD_CORE)
        if (NOT TARGET ${CORE_LIB_NAME})
            set(BOARD_CORE_PATH ${${BOARD_CORE}.path})
            find_sources(CORE_SRCS ${BOARD_CORE_PATH} True)
            # Debian/Ubuntu fix
            list(REMOVE_ITEM CORE_SRCS "${BOARD_CORE_PATH}/main.cxx")
            add_library(${CORE_LIB_NAME} ${CORE_SRCS})
            set_board_flags(ARDUINO_COMPILE_FLAGS ARDUINO_LINK_FLAGS ${BOARD_ID} FALSE)
            set_target_properties(${CORE_LIB_NAME} PROPERTIES
                    COMPILE_FLAGS "${ARDUINO_COMPILE_FLAGS}"
                    LINK_FLAGS "${ARDUINO_LINK_FLAGS}")
        endif ()
        set(${OUTPUT_VAR} ${CORE_LIB_NAME} PARENT_SCOPE)
    endif ()
endfunction()
