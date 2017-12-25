set(WCOSA_CMD "python {{wcosa-path}}/core/wcosa.py")
set(VER {{cmake-version}})
set(NAME {{project-name}})

cmake_minimum_required(VERSION ${VER})

project(${NAME} C CXX ASM)

SET(CMAKE_C_COMPILER avr-gcc)
SET(CMAKE_CXX_COMPILER avr-g++)
SET(CMAKE_CXX_FLAGS_DISTRIBUTION "{{cmake-cxx-flags}}")
SET(CMAKE_C_FLAGS_DISTRIBUTION "{{cmake-c-flags}}")
set(CMAKE_CXX_STANDARD {{cmake-cxx-standard}})

% def-search
{{add_definitions({{user-definition}})}}
% end

# add search paths for all the user libraries
% lib-search
{{include_directories({{lib-path}})}}
% end

add_custom_target(
    WCOSA_BUILD ALL
    COMMAND ${WCOSA_CMD} build --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_CLEAN ALL
    COMMAND ${WCOSA_CMD} clean --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPDATE_ALL ALL
    COMMAND ${WCOSA_CMD} update --ide clion
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(
    WCOSA_UPLOAD ALL
    COMMAND ${WCOSA_CMD} upload
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

file(GLOB RECURSE SRC_FILES "src/*.cpp", "src/*.cc", "src/*.c")

add_executable(${SRC_FILES})
