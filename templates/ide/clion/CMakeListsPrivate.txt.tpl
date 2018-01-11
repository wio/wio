set(WCOSA_CMD python "{{wcosa-path}}/wcosa/wcosa.py")

SET(CMAKE_C_COMPILER avr-gcc)
SET(CMAKE_CXX_COMPILER avr-g++)
SET(CMAKE_CXX_FLAGS_DISTRIBUTION "{{cmake-cxx-flags}}")
SET(CMAKE_C_FLAGS_DISTRIBUTION "{{cmake-c-flags}}")
set(CMAKE_CXX_STANDARD {{cmake-cxx-standard}})

% def-search
{{add_definitions({{user-definition}})}}
% end

# add search paths for cosa core
% cosa-search
{{include_directories("{{wcosa-core}}")}}
{{include_directories("{{wcosa-board}}")}}
% end

# add search paths for all the user libraries
% lib-search

{{include_directories("{{lib-path}}")}}
% end

file(GLOB_RECURSE SRC_FILES "src/*.cpp" "src/*.cc" "src/*.c")
