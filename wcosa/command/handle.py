"""
Handle handles creating and updating WCosa projects
"""

import os
from shutil import copyfile

from wcosa.objects.objects import Fore, Port
from wcosa.templates import cmake, config
from wcosa.utils import helper
from wcosa.utils.output import write, writeln


def create_folders(project_path, override=False):
    """Creates required folders (src, bin and wcosa) in the project directory"""

    helper.create_folder(project_path + '/src', override)
    helper.create_folder(project_path + '/lib', override)
    helper.create_folder(project_path + '/wcosa', override)
    helper.create_folder(project_path + '/wcosa/bin', override)


def verify_path(path):
    """check if the project path is correct"""

    if not os.path.exists(path) or not os.path.isdir(path):
        writeln('\nPath specified for project creation does not exist or is not a directory', color=Fore.RED)
        quit(2)


def create_wcosa(path, board, ide):
    """Creates WCosa project from scratch"""

    path = str(path)

    write('Creating work environment - ', color=Fore.CYAN)
    verify_path(path)

    templates_path = helper.linux_path(helper.get_wcosa_path() + '/templates')
    user_config_path = helper.linux_path(path + '/config.json')
    internal_config_path = helper.linux_path(path + '/wcosa/internal-config.json')
    general_cmake_path = helper.linux_path(path + '/wcosa/CMakeLists.txt')
    src_path = helper.linux_path(path + '/src/main.cpp')

    # check if there are already src and lib folders. We do not want to delete those folders
    if len(helper.get_dirs(path)) > 0 and (os.path.exists(path + '/src') or os.path.exists(path + '/lib')):
        writeln('\nThere is already a src and/or lib folder in this directory. Use wcosa update instead',
                color=Fore.RED)
        quit(2)

    # create src, lib, and wcosa folders
    create_folders(path, True)

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + '/cmake/CMakeLists.txt.tpl', general_cmake_path)
    copyfile(templates_path + '/config/internal-config.json.tpl', internal_config_path)
    copyfile(templates_path + '/config/config.json.tpl', user_config_path)
    copyfile(templates_path + '/examples/main.cpp', src_path)

    writeln('done')
    write('Updating configurations based on the system - ', color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, Port(None), ide)
    project_data = config.fill_internal_config(internal_config_path, path, user_data)

    cmake.parse_update(general_cmake_path, project_data)

    if user_data['ide'] == 'clion':
        copyfile(templates_path + '/ide/clion/CMakeLists.txt.tpl', path + '/CMakeLists.txt')
        copyfile(templates_path + '/ide/clion/CMakeListsPrivate.txt.tpl', path + '/CMakeListsPrivate.txt')
        copyfile(templates_path + '/gitignore-files/.gitignore-clion', path + '/.gitignore')

        cmake.parse_update(path + '/CMakeLists.txt', project_data)
        cmake.parse_update(path + '/CMakeListsPrivate.txt', project_data)
    else:
        copyfile(templates_path + '/gitignore-files/.gitignore-general', path + '/.gitignore')

    writeln('done')
    writeln('Project Created and structure:', color=Fore.YELLOW)
    writeln('src    ->    All source files go here:', color=Fore.YELLOW)
    writeln('lib    ->    All custom libraries go here', color=Fore.YELLOW)
    writeln('wcosa  ->    All the build files are here (do no modify)', color=Fore.YELLOW)


def update_wcosa(path, board, ide):
    """Updates existing WCosa project"""

    path = str(path)

    write('Updating work environment - ', color=Fore.CYAN)
    verify_path(path)

    templates_path = helper.linux_path(helper.get_wcosa_path() + '/templates')
    user_config_path = path + '/config.json'
    internal_config_path = path + '/wcosa/internal-config.json'
    general_cmake_path = path + '/wcosa/CMakeLists.txt'

    # create src, lib, and wcosa folders
    create_folders(path)

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + '/cmake/CMakeLists.txt.tpl', general_cmake_path)
    copyfile(templates_path + '/config/internal-config.json.tpl', internal_config_path)

    # recopy any missing files
    if not os.path.exists(path + '/config.json'):
        copyfile(templates_path + '/config/config.json.tpl', user_config_path)

    writeln('done')
    write('Updating configurations with new changes - ', color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, Port(None), ide)
    project_data = config.fill_internal_config(internal_config_path, path, user_data)
    cmake.parse_update(general_cmake_path, project_data)

    if user_data['ide'] == 'clion':
        copyfile(templates_path + '/ide/clion/CMakeListsPrivate.txt.tpl', path + '/CMakeListsPrivate.txt')

        # recopy any missing files
        if not os.path.exists(path + '/CMakeLists.txt'):
            copyfile(templates_path + '/ide/clion/CMakeLists.txt.tpl', path + '/CMakeLists.txt')

        if not os.path.exists(path + '/.gitignore-files'):
            copyfile(templates_path + '/gitignore-files/.gitignore-clion', path + '/.gitignore')

        cmake.parse_update(path + '/CMakeLists.txt', project_data)
        cmake.parse_update(path + '/CMakeListsPrivate.txt', project_data)
    elif not os.path.exists(path + '/.gitignore'):
        copyfile(templates_path + '/gitignore-files/.gitignore-general', path + '/.gitignore')

    writeln('done')
