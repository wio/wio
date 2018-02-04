#!/usr/bin/env python

import os

from setuptools import find_packages, setup

install_requires = [
    'colorama',
    'pyserial>=3,<4',
    'six',
]


def package_files(directory):
    paths = []
    for (path, directories, filenames) in os.walk(directory):
        for filename in filenames:
            paths.append(os.path.join('..', path, filename))
    return paths


setup(
    name='WCosa',
    version='1.0.5',
    description='Create, Build, Upload and Monitor AVR Cosa Projects',
    author='Deep Dhillon, Jeff Niu',
    author_email='deep@deepdhillon.com, jeffniu22@gmail.com, ambareeshbalaji@gmail.com',
    setup_requires=['setuptools-markdown'],
    long_description_markdown_filename='README.md',
    license='MIT',
    packages=find_packages(),
    install_requires=install_requires,
    package_data={
        '': package_files('toolchain') + package_files('templates'),
        'wcosa': ['*.json'],
    },
    entry_points={
        'console_scripts': [
            'wcosa = wcosa.wcosa:main',
        ],
    },
    classifiers=[
        'Development Status :: 4 - Beta',
    ],
    keywords=[
        'iot', 'embedded', 'arduino', 'avr', 'fpga', 'firmware',
        'hardware', 'microcontroller', 'debug', 'cosa', 'tool',
    ],
)
