set(CMAKE_VER 3.0.0)
set(PROJECT_NAME {{PROJECT_NAME}})
set(PROJECT_PATH "{{PROJECT_PATH}}")
set(CMAKE_TOOLCHAIN_FILE "{{TOOLCHAIN_FILE}}")
set (CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}")
set(DEPENDENCY_FILE dependencies)

# properties
set(TARGET_NAME {{TARGET_NAME}})
set(BOARD {{BOARD}})
set(FRAMEWORK {{FRAMEWORK}})

cmake_minimum_required(VERSION ${CMAKE_VERSION})
project(${PROJECT_NAME} C CXX ASM)
cmake_policy(SET CMP0023 OLD)

include(${DEPENDENCY_FILE})

file(GLOB_RECURSE LIB_FILES "${PROJECT_PATH}/src/*.cpp" "${PROJECT_PATH}/src/*.cc" "${PROJECT_PATH}/src/*.c")
generate_arduino_library(${PROJECT_NAME}
	SRCS ${LIB_FILES}
    BOARD ${BOARD})
target_compile_definitions(${PROJECT_NAME} PRIVATE __AVR_${FRAMEWORK}__ {{PKG_COMPILE_FLAGS}})
target_include_directories(${PROJECT_NAME} PUBLIC "${PROJECT_PATH}/include")
target_include_directories(${PROJECT_NAME} PRIVATE "${PROJECT_PATH}/src")

file(GLOB_RECURSE SRC_FILES "${PROJECT_PATH}/tests/*.cpp" "${PROJECT_PATH}/tests/*.cc" "${PROJECT_PATH}/tests/*.c")
generate_arduino_firmware(${TARGET_NAME}
    SRCS ${SRC_FILES}
    LIBS ${PROJECT_NAME}
    BOARD ${BOARD}
    PORT {{PORT}})
target_compile_definitions(${TARGET_NAME} PRIVATE __AVR_${FRAMEWORK}__ {{TARGET_COMPILE_FLAGS}})