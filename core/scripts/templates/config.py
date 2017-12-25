"""@package templates
Parses and completes the config templates
"""

import json
import os

from core.scripts.others import helper
from core.scripts.parsers import platform_parser, board_parser


def fill_internal_config(path, curr_path, ide, user_config_data):
    """fills the internal config file that will be used for internal build"""

    internal_config_file = open(helper.linux_path(path))
    internal_config_data = json.load(internal_config_file)
    internal_config_file.close()

    settings_file = open(helper.linux_path(os.path.dirname(__file__) + "/../settings.json"))
    settings_data = json.load(settings_file)
    settings_file.close()

    internal_config_data["project-name"] = os.path.basename(curr_path)
    internal_config_data["ide"] = ide
    internal_config_data["board"] = user_config_data["board"]
    internal_config_data["port"] = user_config_data["port"]
    internal_config_data["wcosa-path"] = helper.linux_path(os.path.abspath(os.path.dirname(__file__) + "/../../"))
    internal_config_data["current-path"] = helper.linux_path(curr_path)
    internal_config_data["cmake-version"] = settings_data["cmake-version"]

    # get c and cxx flags
    board_properties = board_parser.get_board_properties(user_config_data["board"],
                                                         internal_config_data["wcosa-path"] + "/core/boards.txt")
    internal_config_data["cmake-c-flags"] = platform_parser.get_c_compiler_flags(board_properties,
                                                                                 internal_config_data[
                                                                                     "wcosa-path"] +
                                                                                 "/toolchain/cosa/platform.txt",
                                                                                 settings_data["include-extra-flags"])
    internal_config_data["cmake-cxx-flags"] = platform_parser.get_cxx_compiler_flags(board_properties,
                                                                                     internal_config_data[
                                                                                         "wcosa-path"] +
                                                                                     "/toolchain/cosa/platform.txt",
                                                                                     settings_data[
                                                                                         "include-extra-flags"])
    internal_config_data["cmake-cxx-standard"] = settings_data["cmake-cxx-standard"]
    internal_config_data["custom-definitions"] = user_config_data
    internal_config_data["custom-definitions"] = " -D" + board_properties["id"]  # board ID
    internal_config_data["custom-definitions"] = internal_config_data["custom-definitions"].strip(" ")
    internal_config_data["cosa-libraries"] = user_config_data["cosa-libraries"]

    internal_config_file = open(helper.linux_path(path), "w")
    json.dump(internal_config_data, internal_config_file, indent=settings_data["json-indent"])
    internal_config_file.close()

    return internal_config_data


def fill_user_config(path, board, port):
    """fills the user config file that will be used for internal build"""

    user_config_file = open(helper.linux_path(path))
    user_config_data = json.load(user_config_file)
    user_config_file.close()

    settings_file = open(helper.linux_path(os.path.dirname(__file__) + "/../settings.json"))
    settings_data = json.load(settings_file)
    settings_file.close()

    user_config_data["board"] = board
    user_config_data["framework"] = settings_data["framework"]
    user_config_data["port"] = port

    user_config_file = open(helper.linux_path(path), "w")
    json.dump(user_config_data, user_config_file, indent=settings_data["json-indent"])
    user_config_file.close()

    return user_config_data
