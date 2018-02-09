"""
Main Script that calls other scripts to make wcosa work
"""

from __future__ import absolute_import

import argparse

from wcosa.command import handle, monitor, package_manager, use
from wcosa.objects.objects import Board, Fore, Generator, IDE, Path, Port
from wcosa.parsers import board_parser
from wcosa.utils import helper, output


def parse():
    """Adds command line arguments and returns the options"""

    parser = argparse.ArgumentParser(description='WCosa create, build and upload Cosa AVR projects')
    parser.add_argument(
        '--path',
        default=helper.get_working_directory(),
        help='path to run action on (default: current directory)',
        type=str)
    subparsers = parser.add_subparsers(dest='action', metavar='action')
    subparsers.required = True
    create_parser = subparsers.add_parser(
            'create',
            help='create project')
    create_parser.add_argument(
        '--board',
        help='board to use for wcosa project',
        required=True,
        type=str)
    create_parser.add_argument(
        '--ide',
        help='create project structure for specific ide (default: none)',
        type=str)
    update_parser = subparsers.add_parser(
        'update',
        help='update project')
    update_parser.add_argument(
        '--board',
        help='board to use for wcosa project',
        type=str)
    update_parser.add_argument(
        '--ide',
        help='update project structure for specific ide (default: none)',
        type=str)
    build_parser = subparsers.add_parser(
        'build',
        help='build project')
    build_parser.add_argument(
        '--generator',
        help='makefile generator to use for build (default: Unix Makefiles)',
        type=str)
    build_parser.add_argument(
        '--make',
        help='path to make binary',
        type=str)
    build_parser.add_argument(
        '--cmake',
        help='path to cmake binary',
        type=str)
    upload_parser = subparsers.add_parser(
        'upload',
        help='upload project')
    upload_parser.add_argument(
        '--port',
        help='port to upload the AVR traget to (default: automatic)',
        type=str)
    monitor_parser = subparsers.add_parser(
        'monitor',
        help='monitor AVR device')
    monitor_parser.add_argument(
        '--port',
        help='port to monitor the AVR traget at (default: automatic)',
        type=str)
    monitor_parser.add_argument(
        '--baud',
        help='buad rate for serial (default: 9600)',
        type=int)
    subparsers.add_parser('boards', help='print supported boards')
    subparsers.add_parser('clean', help='clean build files')

    package_parser = subparsers.add_parser(
        'package',
        help='manipulate packages')
    package_subparsers = package_parser.add_subparsers(
        dest='package_command',
        metavar='command')
    package_subparsers.required = True
    install_parser = package_subparsers.add_parser(
        'install',
        help='install package(s)')
    install_parser.add_argument(
        'package',
        nargs='*',
        type=str)
    remove_parser = package_subparsers.add_parser(
        'remove',
        help='remove package(s)')
    remove_parser.add_argument(
        'package',
        nargs='*',
        type=str)
    update_parser = package_subparsers.add_parser(
        'update',
        help='update all packages')

    return parser.parse_args()


def print_boards():
    """Print all the available boards and their name"""

    boards = board_parser.get_all_board(helper.get_wcosa_path() + '/wcosa/boards.json')

    output.writeln('Boards compatible with this project are: ', Fore.CYAN)

    for curr_board in boards:
        name = board_parser.get_board_properties(curr_board, helper.get_wcosa_path() + '/wcosa/boards.json')['name']
        output.writeln('{:15s} --->\t{}'.format(curr_board, name))


def main():
    options = parse()

    path = Path(options.path)

    # based on the action call scripts
    if options.action == 'boards':
        print_boards()
    elif options.action == 'create':
        handle.create_wcosa(path, Board(options.board), IDE(options.ide))
    elif options.action == 'update':
        handle.update_wcosa(path, Board(options.board), IDE(options.ide))
    elif options.action == 'build':
        use.build_wcosa(path, Generator(options.generator),
                        options.make, options.cmake)
    elif options.action == 'upload':
        use.upload_wcosa(path, Port(options.port))
    elif options.action == 'clean':
        use.clean_wcosa(path)
    elif options.action == 'monitor':
        monitor.serial_monitor(options.port, options.baud)
    elif options.action == 'package':
        if options.package_command == 'install':
            package_manager.package_install_many(
                    options.path,
                    ' '.join(options.package).split(', '))
        elif options.package_command == 'update':
            package_manager.package_update_all(options.path)
        elif options.package_command == 'remove':
            package_manager.package_uninstall_many(
                    options.path,
                    ' '.join(options.package).split(', '))


if __name__ == '__main__':
    main()
