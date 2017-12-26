"""@package module
Output is a wrapper over print to provide coloring and verbose mode
"""

from __future__ import print_function
from colorama import init
from colorama import Style


class VerboseScope:
    """holds the verbose flag"""
    verbose_flag = False


init()
scope = VerboseScope()


def set_verbose(status):
    """turns verbose flag on and off"""
    VerboseScope.verbose_flag = status


def verbose(string, newline, color=Style.RESET_ALL):
    """if verbose option is one, then only write on the output"""
    if VerboseScope.verbose_flag is True:
        write(string, color)

        if newline is True:
            write("\n", color)


def write(string="", color=Style.RESET_ALL):
    """write a string with color specified. No new line"""

    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL, end="")


def writeln(string="", color=Style.RESET_ALL):
    """write a string with color specified. New line"""

    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL)
