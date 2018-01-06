"""@package parsers
Parses the boards.txt file and gathers information about the current board
"""

from collections import OrderedDict
import json

from wcosa.utils import helper


def create_boards_tree(board_file_path, new_board_path):
    """Create a json version of boards file from cosa"""

    with open(helper.linux_path(board_file_path)) as f:
        board_str = f.readlines()

    tree = {}
    curr_board = ''

    for line in board_str:
        if 'name=' in line:
            curr_board = line[:line.find('.')]
            tree[curr_board] = {}
            tree[curr_board]['name'] = line[line.find('='):].strip('=').strip('\n').strip(' ')
        elif 'mcu=' in line:
            tree[curr_board]['mcu'] = line[line.find('='):].strip('=').strip('\n').strip(' ')
        elif 'f_cpu=' in line:
            tree[curr_board]['f_cpu'] = line[line.find('='):].strip('=').strip('\n').strip(' ')
        elif 'board=' in line:
            tree[curr_board]['id'] = line[line.find('='):].strip('=').strip('\n').strip(' ')

    with open(helper.linux_path(new_board_path), 'w') as f:
        json.dump(tree, f, indent=4)


def get_board_properties(board, board_path):
    """parses the board file returns the properties of the board specified"""

    with open(helper.linux_path(board_path)) as f:
        board_data = json.load(f, object_pairs_hook=OrderedDict)

    return board_data[board]


def get_all_board(board_path):
    """parses the board file returns the properties of the board specified"""

    with open(helper.linux_path(board_path)) as f:
        board_data = json.load(f, object_pairs_hook=OrderedDict)

    keys = []

    for key in board_data:
        keys.append(key)

    return keys
