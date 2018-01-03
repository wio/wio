#!/usr/bin/env python

from setuptools import setup, find_packages

setup(name='WCosa',
      version='0.1.dev1',
      description='Create, Build, Upload and Monitor AVR Cosa Projects',
      author='Deep Dhillon, Jeff Niu, Ambareesh Balaji',
      author_email='deep.dhill6@gmail.com, jeffniu22@gmail.com, ambareeshbalaji@gmail.com',
      long_description=open('README.md').read(),
      license='MIT',
      packages=find_packages(),
      install_requires=[
          'colorama', 'serial'
      ],
      classifiers=[
          'Development Status :: 4 -Beta',
          'License :: MIT license',
      ]
      )
