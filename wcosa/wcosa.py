"""
Main Script that calls other scripts to make wcosa work
"""

from __future__ import absolute_import

import argparse

from wcosa.command import handle
from wcosa.command import use
from wcosa.objects.objects import *


def parse():
    """Adds command line arguments and returns the options"""

    parser = argparse.ArgumentParser(description="WCosa create, build and upload Cosa AVR projects")

    parser.add_argument(
        'action',
        help='action to perform (create, update, build, upload, serial and boards')
    parser.add_argument(
        '--board',
         help='board to use for wcosa project',
        type=str)
    parser.add_argument(
        '--port',
        help='port to upload the AVR traget to (default: automatic)',
        type=str)
    parser.add_argument(
        '--baud',
        help='buad rate for serial (default: 9600)',
        type=int)
    parser.add_argument(
        '--ide',
        help='create specific project structure for specific ide (default: none)',
        type=str)
    parser.add_argument(
        '--path',
        help='path to create the project at (default: curr dir)',
        type=str)
    parser.add_argument(
        '--generator',
        help='makefile generator to use for build (default: Unix Makefiles)',
        type=str)
    parser.add_argument(
        '--make',
        help='path to make binary',
        type=str)
    parser.add_argument(
        '--cmake',
        help='path to cmake binary',
        type=str)

    return parser.parse_args()


def print_boards():
    """Print all the available boards and their name"""

    boards = board_parser.get_all_board(helper.get_wcosa_path() + "/wcosa/boards.json")

    output.writeln("Boards compatible with this project are: ", Fore.CYAN)

    for curr_board in boards:
        name = board_parser.get_board_properties(curr_board, helper.get_wcosa_path() + "/wcosa/boards.json")["name"]
        output.writeln('{:15s} --->\t{}'.format(curr_board, name))


def provided(*args):
    """checks if give flags are specified during command line usage"""

    if any(flag is not None for flag in args):
        return True


def main():
    options = parse()

    board = Board(options.board)
    ide = IDE(options.ide)
    port = Port(options.port)
    path = Path(options.path)
    generator = Generator(options.generator)
    cmake = options.cmake
    make = options.make

    # based on the action call scripts
    if options.action == "boards":
        print_boards()
    elif options.action == "create":
        if provided(options.port, options.generator, options.baud):
            output.writeln("Create only requires path, board and ide, other flags are ignored", Fore.YELLOW)

        if options.board is not None:
            handle.create_wcosa(path, board, ide)
        else:
            output.writeln("Board is needed for creating wcosa project", Fore.RED)
    elif options.action == "update":
        if provided(options.port, options.generator, options.baud):
            output.writeln("Update only requires path, board and ide, other flags are ignored", Fore.YELLOW)

        handle.update_wcosa(path, board, ide)
    elif options.action == "build":
        if provided(options.port, options.ide, options.board, options.baud):
            output.writeln("Build only requires path and generator, other flags are ignored", Fore.YELLOW)

        use.build_wcosa(path, generator, make, cmake)
    elif options.action == "upload":
        if provided(options.ide, options.board, options.generator, options.baud):
            output.writeln("Upload only requires path and port, other flags are ignored", Fore.YELLOW)

        use.upload_wcosa(path, port)
    elif options.action == "clean":
        if provided(options.ide, options.board, options.generator, options.port, options.baud):
            output.writeln("Clean only requires path, other flags are ignored", Fore.YELLOW)

        use.clean_wcosa(path)

if __name__ == "__main__":
    main()
