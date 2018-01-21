"""
Handle handles building and uploading wcosa projects
"""

from collections import OrderedDict
import copy
import json
import os
import shutil
import subprocess

import serial.tools.list_ports

from wcosa.command import handle
from wcosa.objects import settings
from wcosa.objects.objects import (
    Board,
    Fore,
    Generator,
    IDE,
    Path,
)
from wcosa.utils import helper, output
from wcosa.utils.finder import (
    get_cmake_program,
    get_generator_for,
    get_make_program,
)


def build_wcosa(path, generator, make=None, cmake=None):
    """build wcosa project, cmake and make"""

    path = str(path)

    output.writeln('Wcosa project build started', Fore.GREEN)
    output.write('Verifying the project structure - ', Fore.GREEN)

    # check if path is valid
    if not os.path.exists(path):
        output.write('Project path is invalid: ' + path, Fore.RED)
        quit(2)

    # confirm if bin folder/path exists
    if os.path.exists(helper.linux_path(path + '/wcosa/bin')):
        os.chdir(helper.linux_path(path + '/wcosa/bin'))
    else:
        output.writeln('\nNot a valid WCosa project', Fore.RED)
        quit(2)

    output.writeln('done')

    cmake_program = cmake or get_cmake_program()
    make_program = make or get_make_program()

    output.write('Verifying cmake and make installs - ', Fore.GREEN)

    # check if cmake is in environment paths (unix/linux based systems)
    if not cmake_program:
        output.writeln(
            '\ncmake does not exist, please install it or make sure it is in your environment PATH',
            Fore.RED)
        quit(2)

    # check if make is in environment paths (unix/linux based systems)
    if not make_program:
        output.writeln(
            '\nmake does not exist, please install it or make sure it is in your environment PATH',
            Fore.RED)
        quit(2)

    output.writeln('done')

    if not str(generator):
        generator = Generator(get_generator_for(make_program))

    # check if path is valid and get build information from the user config
    if not os.path.exists(path + '/config.json'):
        output.write(
            'Project user configuration file does not exist, recreate or update the project',
            Fore.RED)
        quit(2)
    else:
        with open(path + '/config.json') as f:
            data = json.load(f, object_pairs_hook=OrderedDict)
            board = data['board']

    # check if the current build files are for the current board
    # clean and build if boards are different
    if os.path.exists(helper.get_working_directory() + '/Makefile'):
        with open(helper.get_working_directory() + '/Makefile') as f:
            makefile_str = ''.join(f.readlines())

        if '\n' + board not in makefile_str:
            output.writeln('Since a new board is detected, full build will be triggered', Fore.GREEN)
            clean_wcosa(path)

    output.writeln('Running the build using cmake and ' + str(generator), Fore.GREEN)
    cmake_code = subprocess.call(['cmake', '-G', str(generator), '..'])

    if cmake_code != 0:
        output.writeln('Project build unsuccessful, cmake exited with error code ' + str(cmake_code), Fore.RED)
        quit(2)

    make_code = subprocess.call([make_program])

    if make_code != 0:
        output.writeln('Project build unsuccessful, make exited with error code ' + str(make_code), Fore.RED)
        quit(2)

    output.writeln('Project successfully built', Fore.GREEN)


def get_serial_devices():
    """get all the valid serial devices/ports"""

    ports = list(serial.tools.list_ports.comports())
    devices = []

    for p in ports:
        devices.append(p.device)

    return devices


def serial_ports():
    """Returns the serial port by scanning all of them"""

    ports = list(serial.tools.list_ports.comports())

    if not ports:
        output.writeln('No device is connected at the moment', Fore.RED)
        quit(2)

    for p in ports:
        test_str = str(p.description).lower()
        if 'arduino' in test_str or 'Arduino' in test_str:
            return p.device

    # if no arduino port is found choose the first port
    output.writeln(
        'No Arduino port found, choosing the first available one. Specify the port if you want another port',
        Fore.YELLOW)

    return ports[0].device


def upload_wcosa(path, port):
    """upload wcosa project to the port specified or an automatically selected one"""

    path = str(path)

    output.writeln('Wcosa project upload started', Fore.GREEN)
    output.write('Verifying build files - ', Fore.GREEN)

    # confirm if bin folder/path exists
    if os.path.exists(helper.linux_path(path + '/wcosa/bin')):
        os.chdir(helper.linux_path(path + '/wcosa/bin'))
    else:
        output.writeln('\nNot a valid WCosa project', Fore.RED)

    # if project has not been built right
    if not os.path.exists(helper.get_working_directory() + '/Makefile'):
        output.writeln('\nNo Makefile, build the project first ', Fore.RED)
        quit(2)

    output.writeln('done')

    # check if path is valid and get the user config
    if not os.path.exists(path + '/config.json'):
        output.write('Project user configuration file does not exist, recreate or update the project', Fore.RED)
        quit(2)
        user_port = None
        data = None
    else:
        with open(path + '/config.json') as f:
            data = json.load(f, object_pairs_hook=OrderedDict)
            user_port = data['port'].strip(' ')

    if not port.use_same():
        output.write('Using the port provided - ' + str(port), Fore.GREEN)
    else:
        # if port is defined in the config file, then use that port
        if user_port.lower() != 'None'.lower() and user_port != '':
            port = user_port

            output.writeln('Using the port from the config file: ' + user_port, Fore.GREEN)
            # check if the port is valid
            if port not in get_serial_devices():
                output.writeln('\nPort provided does not have a valid device connected to it', Fore.RED)
                quit(2)
        else:
            output.writeln('Automatically selecting a port', Fore.GREEN)
            port = serial_ports()
            output.writeln('Port chosen: ' + port, Fore.GREEN)

    # save the port in the config file so that update can update the project based on that
    temp_data = copy.copy(data)
    temp_data['port'] = str(port)
    with open(path + '/config.json', 'w') as f:
        json.dump(temp_data, f, indent=settings.get_settings_value('json-indent'))

    # do not print the updating logs
    output.output_status(False)
    handle.update_wcosa(Path(path), Board(None), IDE(None))
    build_wcosa(path, Generator(None))
    output.output_status(True)

    # reset the port in the user config to the port that was there
    with open(path + '/config.json', 'w') as f:
        json.dump(data, f, indent=settings.get_settings_value('json-indent'))

    output.writeln('Upload triggered on ' + str(port), Fore.GREEN)

    upload_code = subprocess.call([get_make_program(), 'upload'])

    if upload_code != 0:
        output.writeln('Project upload unsuccessful, make exited with error code ' + str(upload_code), Fore.RED)
        quit(2)

    output.writeln('Project upload successful', Fore.GREEN)


def clean_wcosa(path):
    """cleans the bin folder, deleting all the build files"""

    path = str(path)

    # confirm if bin folder/path exists
    if os.path.exists(helper.linux_path(path + '/wcosa/bin')):
        os.chdir(helper.linux_path(path + '/wcosa/bin'))
    else:
        output.writeln('\nNot a valid WCosa project', Fore.RED)
        quit(2)

    # check if the current build files are for the current board
    # clean and build if boards are different
    try:
        output.write('Cleaning build files - ', Fore.GREEN)

        for folder in helper.get_dirs(helper.get_working_directory()):
            shutil.rmtree(folder)

        for file in helper.get_files(helper.get_working_directory()):
            os.unlink(file)
    except IOError:
        output.writeln('Error while cleaning build files', Fore.RED)

    output.writeln('done')
