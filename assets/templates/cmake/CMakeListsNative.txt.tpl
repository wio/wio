######################################################################
# This is auto-generated by wio
######################################################################

set(CMAKE_VER 3.0.0)
set(PROJECT_NAME {{PROJECT_NAME}})
set(PROJECT_PATH "{{PROJECT_PATH}}")
set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR})
set(DEPENDENCY_FILE dependencies)

# C++ standard
set(CMAKE_CXX_STANDARD {{CPP_STANDARD}})
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# C standard
set(CMAKE_C_STANDARD {{C_STANDARD}})
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)

# Properties
set(TARGET_NAME {{TARGET_NAME}})
set(OS {{OS}})
set(PLATFORM {{PLATFORM}})
set(FRAMEWORK {{FRAMEWORK}})
set(ENTRY {{ENTRY}})

# CMAKE
cmake_minimum_required(VERSION ${CMAKE_VERSION})
project(${PROJECT_NAME} C CXX ASM)

# Variables
set(PLATFORM {{PLATFORM}})
set(FRAMEWORK {{FRAMEWORK}})
set(BOARD {{BOARD}})

file(GLOB_RECURSE ${TARGET_NAME}_files
    ${PROJECT_PATH}/${ENTRY}/*.cpp
    ${PROJECT_PATH}/${ENTRY}/*.cc
    ${PROJECT_PATH}/${ENTRY}/*.c)

add_executable(${TARGET_NAME} ${${TARGET_NAME}_files})

target_compile_definitions(
    ${TARGET_NAME}
    PRIVATE
    WIO_PLATFORM_${PLATFORM}
    WIO_FRAMEWORK_${FRAMEWORK}
    WIO_OS_${OS}
    {{TARGET_COMPILE_DEFINITIONS}})

target_compile_options(
    ${TARGET_NAME}
    PRIVATE
    {{TARGET_COMPILE_FLAGS}})

include(${DEPENDENCY_FILE})
