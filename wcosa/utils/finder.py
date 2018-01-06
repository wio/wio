import os

from .memoize import memoized

make_programs = [
    'make',
    'make.exe',
    'mingw32-make',
    'mingw32-make.exe',
]

cmake_programs = [
    'cmake',
    'cmake.exe',
]


@memoized
def cmd_exists(cmd):
    """checks if the program exists in the PATH"""

    return any(
        os.access(os.path.join(path, cmd), os.X_OK)
        for path in os.environ['PATH'].split(os.pathsep)
    )


def get_program(programs):
    for program in programs:
        if cmd_exists(program):
            return program
    return None


@memoized
def get_generator_for(make_program):
    if 'mingw' in make_program:
        return 'MinGW Makefiles'
    return 'Unix Makefiles'


@memoized
def get_make_program():
    return get_program(make_programs)


@memoized
def get_cmake_program():
    return get_program(cmake_programs)
