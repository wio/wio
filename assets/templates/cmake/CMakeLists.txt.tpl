set(WCOSA_PATH wcosa)
set(VER {{cmake-version}})
set(NAME {{project-name}})

# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "${WCOSA_PATH}/toolchain/cmake/CosaToolchain.cmake")

cmake_minimum_required(VERSION ${VER})

project(${NAME} C CXX ASM)

# add search paths for all the user libraries and build them
% lib-search
{{include_directories("{{lib-path}}")}}
{{generate_arduino_library({{name}}\n\tSRCS {{srcs}}\n\tHDRS {{hdrs}}\n\tBOARD {{board}})}}
{{target_compile_definitions({{name}} PRIVATE __AVR_Cosa__ {{custom-definitions}})}}
% end

file(GLOB_RECURSE SRC_FILES "../src/*.cpp" "../src/*.cc" "../src/*.c")

# create the firmware
% firmware-gen
{{generate_arduino_firmware({{name}}\n\tSRCS ${SRC_FILES}\n\tARDLIBS {{cosa-libraries}}\n\tLIBS {{libs}}\n\tPORT {{port}}\n\tBOARD {{board}})}}
{{target_compile_definitions({{name}} PRIVATE __AVR_Cosa__ {{custom-definitions}})}}
% end
