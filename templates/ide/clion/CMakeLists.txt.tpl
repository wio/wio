set(VER {{cmake-version}})
set(NAME {{project-name}})

cmake_minimum_required(VERSION ${VER})

project(${NAME} C CXX ASM)

include("CMakeListsPrivate.txt")

add_custom_target(
    WCOSA_BUILD ALL
    COMMAND ${WCOSA_CMD} build
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_CLEAN ALL
    COMMAND ${WCOSA_CMD} clean
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPDATE_ALL ALL
    COMMAND ${WCOSA_CMD} update
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPLOAD ALL
    COMMAND ${WCOSA_CMD} upload
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_executable(${PROJECT_NAME} ${SRC_FILES})
