"""@package parsers
Parses the boards.txt file and gathers information about the current board
"""

import os
import json

from core.scripts.others import helper


def create_boards_tree(board_file_path, new_board_path):
    """Create a json version of boards file from cosa"""

    board_file = open(helper.linux_path(board_file_path))
    board_str = board_file.readlines()

    tree = {}
    curr_board = ""

    for line in board_str:
        if "name=" in line:
            curr_board = line[:line.find(".")]
            tree[curr_board] = {}
        elif "mcu=" in line:
            tree[curr_board]["mcu"] = line[line.find('='):].strip("=").strip("\n").strip(" ")
        elif "f_cpu=" in line:
            tree[curr_board]["f_cpu"] = line[line.find('='):].strip("=").strip("\n").strip(" ")
        elif "board=" in line:
            tree[curr_board]["id"] = line[line.find('='):].strip("=").strip("\n").strip(" ")

    file = open(new_board_path, "w")
    json.dump(tree, file, indent=4)
    file.close()


def get_board_properties(board, board_path):
    """parses the board file returns the properties of the board specified"""

    board_file = open(helper.linux_path(os.path.abspath(board_path)))
    board_data = json.load(board_file)
    board_file.close()

    return board_data[board]


def get_all_board(board_path):
    """parses the board file returns the properties of the board specified"""

    board_file = open(helper.linux_path(os.path.abspath(board_path)))
    board_data = json.load(board_file)
    board_file.close()

    keys = []

    for key in board_data:
        keys.append(key)

    return keys
