"""
Handle handles creating and updating WCosa projects
"""

import os

from shutil import copyfile
from colorama import Fore
from wcosa.others.output import write, writeln
from wcosa.others import helper
from wcosa.templates import config
from wcosa.templates import cmake


def create_wcosa(path, board, ide):
    """Creates WCosa project from scratch"""

    project_path = path

    if path is None:
        project_path = helper.get_working_directory()

    if ide is None:
        ide = ""
    else:
        ide = ide.strip(" ")

    templates_path = helper.linux_path(helper.get_wcosa_path() + "/templates")
    user_config_path = helper.linux_path(project_path + "/config.json")
    internal_config_path = helper.linux_path(project_path + "/wcosa/internal-config.json")
    general_cmake_path = helper.linux_path(project_path + "/wcosa/CMakeLists.txt")

    write("Creating work environment - ", color=Fore.CYAN)

    # check if path exists
    if not os.path.exists(path) or not os.path.isdir(path):
        writeln("aborted")
        write("Path specified for project creation does not exist or is not a directory", color=Fore.RED)
        quit(2)

    # check if there are already files/directories in that folder
    if len(helper.get_dirs(path)) > 0 or len(helper.get_files(path)) > 0:
        writeln("aborted")
        write("The directory should be empty where the project should be created. Use wcosa update instead",
              color=Fore.RED)
        quit(2)

    # create src, lib, and wcosa folders
    helper.create_folder(project_path + "/src", True)
    helper.create_folder(project_path + "/lib", True)
    helper.create_folder(project_path + "/wcosa", True)
    helper.create_folder(project_path + "/wcosa/bin", True)

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + "/cmake/CMakeLists.txt.tpl", general_cmake_path)
    copyfile(templates_path + "/config/internal-config.json.tpl", internal_config_path)
    copyfile(templates_path + "/config/config.json.tpl", user_config_path)

    if ide == "clion":
        copyfile(templates_path + "/ide/clion/CMakeLists.txt.tpl", project_path + "/CMakeLists.txt")
        copyfile(templates_path + "/ide/clion/CMakeListsPrivate.txt.tpl", project_path + "/CMakeListsPrivate.txt")

    writeln("done")
    write("Updating configurations based on the system - ", color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, "None", ide)  # give a dummy port right now
    project_data = config.fill_internal_config(internal_config_path, path, user_data)

    cmake.parse_update(general_cmake_path, project_data)

    if ide != "":
        cmake.parse_update(project_path + "/CMakeLists.txt", project_data)
        cmake.parse_update(project_path + "/CMakeListsPrivate.txt", project_data)

    writeln("done")
    writeln("Project Created and structure:", color=Fore.YELLOW)
    writeln("src    ->    All source files go here:", color=Fore.YELLOW)
    writeln("lib    ->    All custom libraries go here", color=Fore.YELLOW)
    writeln("wcosa  ->    All the build files are here (do no modify)", color=Fore.YELLOW)


def update_wcosa(path, board):
    """Updates existing WCosa project"""

    write("Updating work environment - ", color=Fore.CYAN)

    project_path = path

    if path is None:
        project_path = helper.get_working_directory()

    templates_path = helper.get_wcosa_path() + "/templates"
    user_config_path = project_path + "/config.json"
    internal_config_path = project_path + "/wcosa/internal-config.json"
    general_cmake_path = project_path + "/wcosa/CMakeLists.txt"

    # create src, lib, and wcosa folders
    helper.create_folder(project_path + "/src")
    helper.create_folder(project_path + "/lib")
    helper.create_folder(project_path + "/wcosa")
    helper.create_folder(project_path + "/wcosa/bin")

    # copy all then CMakeLists templates and configuration templates
    copyfile(templates_path + "/cmake/CMakeLists.txt.tpl", general_cmake_path)
    copyfile(templates_path + "/config/internal-config.json.tpl", internal_config_path)

    writeln("done")
    write("Updating configurations with new changes - ", color=Fore.CYAN)

    user_data = config.fill_user_config(user_config_path, board, "None")  # give a dummy port right now
    ide = user_data["ide"]

    if ide == "clion":
        copyfile(templates_path + "/ide/clion/CMakeLists.txt.tpl", project_path + "/CMakeLists.txt")

    project_data = config.fill_internal_config(internal_config_path, path, user_data)

    cmake.parse_update(general_cmake_path, project_data)

    if ide != "":
        cmake.parse_update(project_path + "/CMakeLists.txt", project_data)

    writeln("done")
