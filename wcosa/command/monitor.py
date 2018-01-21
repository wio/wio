"""
Handles basic serial monitor capabilities
"""

from __future__ import print_function

import sys

from serial.tools import miniterm

from wcosa.command.use import serial_ports
from wcosa.objects.objects import Fore
from wcosa.utils.output import write, writeln


def serial_monitor(port, baud):
    """
    Open serial monitor to the specified port with the given baud rate.
    """
    write('Serial Monitor ', Fore.CYAN)
    write(port, Fore.YELLOW)
    write(' @ ', Fore.CYAN)
    writeln(baud, Fore.YELLOW)

    sys.argv = ['monitor', '--exit-char', '3']
    try:
        miniterm.main(
            default_port=port or serial_ports(),
            default_baudrate=baud or 9600)
    except KeyboardInterrupt:
        writeln('\nExiting serial monitor', Fore.CYAN)
