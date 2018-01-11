"""@package templates
Parses and completes the cmake templates
"""

import os

from wcosa.utils import helper
from wcosa.utils.helper import (
    get_dirnames,
    get_dirs,
    get_files_recursively,
    quote_join,
)

DEF_SEARCH_TAG = '% def-search'
LIB_SEARCH_TAG = '% lib-search'
COSA_SEARCH_TAG = '% cosa-search'
FIRMWARE_GEN_TAG = '% firmware-gen'
END_TAG = '% end'
FILL_BLOCK_START = '{{'
FILL_BLOCK_END = '}}'


SRC_FILE_EXTS = ('.cpp', '.c', '.cc')
HDR_FILE_EXTS = ('.hh', '.h')


def lib_search(content, project_data):
    """searches for library paths and then completes the templates to include search paths and build library"""

    str_to_return = ''
    for lib in get_dirs(project_data['current-path'] + '/lib'):
        src_files = []
        hdr_files = []
        lib_paths = lib

        # check if there is a src folder
        src_found = False

        # handle src folder
        for sub_dir in get_dirs(lib):
            if os.path.basename(sub_dir) == 'src':
                lib_paths += '/src'

                # add all the src extensions in src folder
                src_files += get_files_recursively(sub_dir, SRC_FILE_EXTS)

                # add all the header extensions in src folder
                hdr_files += get_files_recursively(sub_dir, HDR_FILE_EXTS)
                src_found = True
                break

        if src_found is not True:
            # add all the src extensions
            src_files += get_files_recursively(lib, SRC_FILE_EXTS)

            # add all the header extensions
            hdr_files += get_files_recursively(lib, HDR_FILE_EXTS)

        # go through all files and generate cmake tags
        lib_name = os.path.basename(lib)
        lib_defs = ''
        if lib_name in project_data['module-definitions']:
            lib_defs = project_data['module-definitions'][lib_name]
        data = {
            'lib-path': [lib_paths],
            'name': lib_name,
            'wcosa-core': [helper.get_cosa_path()],
            'srcs': [quote_join(src_files)],
            'hdrs': [quote_join(hdr_files)],
            'board': project_data['board'],
            'custom-definitions': lib_defs,
        }

        for line in content:
            line = line[2:len(line) - 3]
            str_to_return += helper.fill_template(line, data) + '\n'

    if str_to_return == '':
        str_to_return = '# no libraries to include at the moment'

    return str_to_return.strip(' ').strip('\n') + '\n'


def cosa_search(content, project_data):
    """searches for cosa library search paths"""

    str_to_return = ''

    # go through all files and generate cmake tags
    data = {'wcosa-core': [helper.get_cosa_path() + '/cores/cosa'],
            'wcosa-board': [helper.get_cosa_path() + '/variants/arduino/' + project_data['board']]}

    for line in content:
        line = line[2:len(line) - 3]
        str_to_return += helper.fill_template(line, data) + '\n'

    return str_to_return.strip(' ').strip('\n') + '\n'


def firmware_gen(content, project_data):
    """searches for src files and then generates the firmware code for linking and building the project"""

    curr_lib_path = project_data['current-path'] + '/lib'
    str_to_return = ''

    lib_files = ' '.join(get_dirnames(curr_lib_path))

    data = {
        'name': project_data['project-name'],
        'libs': lib_files,
        'cosa-libraries': project_data['cosa-libraries'],
        'port': project_data['port'],
        'board': project_data['board'],
        'custom-definitions': project_data['custom-definitions'],
    }

    for line in content:
        line = line[2:len(line) - 3]
        str_to_return += helper.fill_template(line, data) + '\n'

    if project_data['cosa-libraries'] == '':
        str_to_return = str_to_return.replace('\tARDLIBS \n', '')

    return str_to_return.strip(' ').strip('\n') + '\n'


def def_search(content, project_data):
    """adds custom definitions"""

    definitions = project_data['custom-definitions'].strip(' ').strip('\n').split(' ')
    str_to_return = ''

    if definitions[0] == '':
        return '# no user definitions\n'

    for definition in definitions:
        for line in content:
            data = {'definition': [definition]}
            line = line[2:len(line) - 3]
            str_to_return += helper.fill_template(line, data) + '\n'

    return str_to_return.strip(' ').strip('\n') + '\n'


def get_elements(tpl_str, curr_index):
    """gather elements from the template block"""

    content = []

    # gather all the lines inside the loop block
    content_index = curr_index + 1
    while True:
        line = tpl_str[content_index]
        compare_tag = line.strip('\n').strip(' ')

        if compare_tag == END_TAG:
            break
        else:
            content.append(line)

        content_index += 1

    return content, content_index


def parse_update(tpl_path, project_data):
    """reads the cmake template file and completes it using project data"""

    tpl_path = helper.linux_path(tpl_path)
    tpl_file = open(tpl_path)
    tpl_str = tpl_file.readlines()
    tpl_file.close()

    new_str = ''
    index = 0
    while index < len(tpl_str):
        curr_line = tpl_str[index]
        compare_tag = curr_line.strip('\n').strip(' ')

        # handle loop statements
        if compare_tag == LIB_SEARCH_TAG:
            result = get_elements(tpl_str, index)

            new_str += lib_search(result[0], project_data)
            index = result[1]
        elif compare_tag == COSA_SEARCH_TAG:
            result = get_elements(tpl_str, index)

            new_str += cosa_search(result[0], project_data)
            index = result[1]
        elif compare_tag == FIRMWARE_GEN_TAG:
            result = get_elements(tpl_str, index)

            new_str += firmware_gen(result[0], project_data)
            index = result[1]
        elif compare_tag == DEF_SEARCH_TAG:
            result = get_elements(tpl_str, index)

            new_str += def_search(result[0], project_data)
            index = result[1]
        else:
            new_str += helper.fill_template(curr_line, project_data)
        index += 1

    tpl_file = open(tpl_path, 'w')
    tpl_file.write(new_str)
    tpl_file.close()
