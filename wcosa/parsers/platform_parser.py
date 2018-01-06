"""@package parsers
Parses the platform.txt file and gathers information about the current platform
"""

from wcosa.objects import settings
from wcosa.utils import helper


def get_raw_flags(lines, identifier, include_extra):
    """Gathers raw flags before updating templates"""

    raw_flags = ''

    for line in lines:
        if 'compiler.' + identifier + '.flags=' in line:
            raw_flags += line[line.find('=') + 1:].strip(' ').strip('\n')
        elif include_extra and 'compiler.' + identifier + '.extra_flags=' in line:
            raw_flags += ' ' + line[line.find('=') + 1:].strip(' ').strip('\n')

    return raw_flags


def get_c_compiler_flags(board_properties, platform_path, include_extra=True):
    """Get template filled c compiler flags"""

    with open(helper.linux_path(platform_path)) as f:
        raw_flags = get_raw_flags(f.readlines(), 'c', include_extra)

    processed_flags = ''

    for flag in raw_flags.split(' '):
        data = {'build.mcu': board_properties['mcu'], 'build.f_cpu': board_properties['f_cpu'],
                'runtime.ide.version': settings.get_settings_value('arduino-version')}
        processed_flags += helper.fill_template(flag, data) + ' '

    return processed_flags.strip(' ')


def get_cxx_compiler_flags(board_properties, platform_path, include_extra=True):
    """Get template filled cxx compiler flags"""

    with open(helper.linux_path(platform_path)) as f:
        raw_flags = get_raw_flags(f.readlines(), 'cpp', include_extra)

    processed_flags = ''

    for flag in raw_flags.split(' '):
        data = {'build.mcu': board_properties['mcu'], 'build.f_cpu': board_properties['f_cpu'],
                'runtime.ide.version': settings.get_settings_value('arduino-version')}
        processed_flags += helper.fill_template(flag, data) + ' '

    return processed_flags.strip(' ')
