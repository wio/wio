"""@package module
Output is a wrapper over print to provide coloring and verbose mode
"""

from __future__ import print_function

import os
import sys

from colorama import (
    init,
    Style,
)


class Scope:
    """holds the verbose flag and output status flag"""

    verbose_flag = False
    output_status_flag = True
    original_stdout = sys.stdout


init()
scope = Scope()


def set_verbose(status):
    """turns verbose flag on and off"""

    Scope.verbose_flag = status


def output_status(status):
    """turn output on and off"""

    if status:
        sys.stdout = scope.original_stdout
    else:
        scope.original_stdout = sys.stdout
        sys.stdout = open(os.devnull, 'w')


def verbose(string, newline, color=Style.RESET_ALL):
    """if verbose option is one, then only write on the output"""

    if Scope.verbose_flag is True:
        write(string, color)

        if newline is True:
            write('\n', color)


def write(string='', color=Style.RESET_ALL):
    """write a string with color specified. No new line"""

    print(color, end='')
    print(string, end='')
    print(Style.RESET_ALL, end='')


def writeln(string='', color=Style.RESET_ALL):
    """write a string with color specified. New line"""

    print(color, end='')
    print(string, end='')
    print(Style.RESET_ALL)
