"""@package parsers
Wrappers around various flags, these handles default cases and whether to use them or not
"""

from colorama import Fore
import serial.tools.list_ports

from wcosa.parsers import board_parser
from wcosa.utils import (
    helper,
    output,
)


class Board:
    """Wrapper for the board flag"""

    def __init__(self, name):
        if name is None:
            self.same = True
            self.name = None
        else:
            self.name = name

            # verify the board
            boards = board_parser.get_all_board(helper.get_wcosa_path() + '/wcosa/boards.json')

            if name not in boards:
                output.writeln('Board Invalid. Run wcosa script with boards option to see all the valid boards',
                               Fore.RED)
                quit(2)

            self.same = False

    def use_same(self):
        """Returns true if the value of board stayed same (None is provided)"""

        return self.same

    def __str__(self):
        return self.name or ''


class IDE:
    """Wrapper for the ide flag"""

    def __init__(self, name):
        if name is None:
            self.same = True
            self.name = None
        else:
            self.name = name
            self.same = False

    def __str__(self):
        return self.name or ''

    def use_same(self):
        """Returns true if the value of ide stayed same (None is provided)"""

        return self.same


class Port:
    """Wrapper for the port flag"""

    def __init__(self, name):
        if name is None:
            self.same = True
            self.name = None
        else:
            self.name = name
            self.same = False

    def verify(self):
        # verify the port
        ports = list(serial.tools.list_ports.comports())
        devices = []

        for p in ports:
            devices.append(p.device)

        if self.name not in devices:
            output.writeln('There is no device connected to this port', Fore.RED)
            quit(2)

    def __str__(self):
        return self.name or ''

    def use_same(self):
        """Returns true if the value of port stayed same (None is provided)"""

        return self.same


class Path:
    """Wrapper for the path flag"""

    def __init__(self, path):
        if path is None:
            self.path = helper.get_working_directory()
        else:
            self.path = path

    def __str__(self):
        return self.path


class Generator:
    """Wrapper for the generator flag"""

    def __init__(self, name):
        if not name:
            self.name = None
            self.same = True
        else:
            self.name = name
            self.same = False

    def __str__(self):
        return self.name or ''

    def use_same(self):
        """Returns true if the value of generator stayed same (None is provided)"""

        return self.same
