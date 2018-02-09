"""
A lightweight package management system for WCosa.
*_many functions operate on lists of data (common scenario, better efficiency).
"""

import json
import os
import re
import sys

import git

from wcosa.utils.output import write, writeln


class Package:
    def __init__(self, name, url, branch, version, path):
        self.unqualified_name = name
        self.url = url
        self.branch = branch
        self.version = version
        self.path = path
        self.name = (self.unqualified_name +
                     ('-' + self.branch if self.branch else '') +
                     ('-' + self.version if self.version else ''))

    def __repr__(self):
        return ('name: %s, url: %s, branch: %s, version: %s, path: %s' %
                (self.name, self.url, self.branch, self.version, self.path))


class PackageFormatError(Exception):
    def __init__(self, package_string):
        self.package_string = package_string

    def __str__(self):
        return 'Bad package format: ' + self.package_string


class GitFetchException(Exception):
    def __init__(self, package):
        self.url = (package.url +
                    (':' + package.branch if package.branch else '') +
                    ('@' + package.version if package.version else ''))

    def __str__(self):
        return 'Could not fetch submodule from %s' % self.url


class AlreadyInstalledException(Exception):
    def __init__(self, link_updated):
        self.link_updated = link_updated


URL = r'(?P<url>https?://\S+/(?P<name>\S+))'
GITHUB = r'(?P<github>[\w\-]+/(?P<name>[\w\-]+))'
BRANCH = r'(:(?P<branch>[\w\-]+))?'
VERSION = r'(@(?P<version>\S+))?'
PATH = r'( at (?P<path>\S+))?'
VALID_SCHEMAS = [re.compile('^' + URL + BRANCH + VERSION + PATH + '$'),
                 re.compile('^' + GITHUB + BRANCH + VERSION + PATH + '$')]


def package_dir_path(path):
    """Return package path to package install directory"""
    return path + '/.pkg'


def package_string_parse_many(package_strings):
    """
    Convert package strings to package entities.
    Package strings must match (URL|GITHUB)[:BRANCH][@VERSION][ at PATH]
    where:
        URL is a valid URL pointing to a git repository
        GITHUB is of the form 'username/reponame'
        BRANCH [optional] is the branch to track
        VERSION [optional] is a tag on the given branch to check out
        PATH [default 'lib/NAME'] is the relative path to install location
    """
    packages = []
    for package_string in package_strings:
        for schema in VALID_SCHEMAS:
            match = re.match(schema, package_string)
            if match:
                groups = match.groupdict()
                break
        if not match:
            raise PackageFormatError(package_string)
        if 'github' in groups:  # only a group if matched with github format
            url = 'https://github.com/' + groups['github']
        else:
            url = groups['url']
        name = groups['name']
        branch = '' if not groups['branch'] else groups['branch']
        version = '' if not groups['version'] else groups['version']
        path = 'lib/' + name if not groups['path'] else groups['path']
        packages.append(Package(name, url, branch, version, path))
    return packages


def package_list_read(pkgpath):
    """Read package list"""
    try:
        with open(pkgpath + '/pkglist', 'r') as pkglistfile:
            return json.loads(pkglistfile.read())
    except Exception:
        return []


def package_list_add_many(pkgpath, packages):
    """Add given packages to package list"""
    if not packages:
        return  # Nothing to write
    repo = package_repo_open(pkgpath)
    newentries = []
    updentries = []
    with open(pkgpath + '/pkglist', 'r+') as pkglistfile:
        pkglist = json.loads(pkglistfile.read())
        pkgnames = list(map(lambda x: x['name'], pkglist))
        for package in packages:
            if package.name in pkgnames:
                index = pkgnames.index(package.name)
                if package.path in pkglist[index]['paths']:
                    continue
                pkglist[index]['paths'].append(package.path)
                updentries.append(package.name)
            else:
                pkglist.append(package.__dict__)
                pkglist[-1]['paths'] = [package.path]
                del pkglist[-1]['path']
                newentries.append(package.name)
        pkglistfile.seek(0)
        pkglistfile.write(json.dumps(pkglist))
    repo.index.add(['pkglist'])
    if repo.is_dirty():  # Something has changed
        repo.index.commit('Updated package list\n\n' +
                          ('New: %s\n' % ', '.join(newentries)
                           if newentries else '') +
                          ('Changed: %s\n' % ', '.join(updentries)
                           if updentries else ''))


def package_list_remove_many(pkgpath, packages):
    """Remove given packages from package list"""
    if not packages:
        return  # Nothing to remove
    repo = package_repo_open(pkgpath)
    uninstalled = []
    unlinked = []
    with open(pkgpath + '/pkglist', 'r+') as pkglistfile:
        pkglist = json.loads(pkglistfile.read())
        pkgnames = list(map(lambda x: x['name'], pkglist))
        for package in packages:
            assert package.name in pkgnames
            index = pkgnames.index(package.name)
            pkglist[index]['paths'].remove(package.path)
            if pkglist[index]['paths']:
                unlinked.append(package.name)
            else:
                del pkglist[index]
                uninstalled.append(package.name)
    with open(pkgpath + '/pkglist', 'w') as pkglistfile:
        pkglistfile.write(json.dumps(pkglist))
    repo.index.add(['pkglist'])
    if repo.is_dirty():  # Something has changed
        repo.index.commit('Updated package list\n\n' +
                          ('Uninstalled: %s\n' % ', '.join(uninstalled)
                           if uninstalled else '') +
                          ('Changed: %s\n' % ', '.join(unlinked)
                           if unlinked else ''))


def package_repo_open(pkgpath):
    """Try to open package repo; initalize upon failure"""
    try:
        return git.Repo(pkgpath)
    except Exception:
        return package_repo_init(pkgpath)


def package_repo_init(pkgpath):
    """Initialize package repo"""
    write('Initializing package repository... ')
    sys.stdout.flush()
    pkgrepo = git.Repo.init(pkgpath)

    with open(pkgpath + '/pkglist', 'w+') as pkglist:
        pkglist.write('[]')  # Start with empty package list

    pkgrepo.index.add(['pkglist'])
    pkgrepo.index.commit('Initialized repository')
    writeln('Done')

    return pkgrepo


def package_link(path, package):
    """Link package directory from pkgpath to package.path"""
    install_path = os.path.abspath(package_dir_path(path) + '/' + package.name)
    link_path = os.path.abspath(path + '/' + package.path)
    link_basedir = '/'.join(link_path.split('/')[:-1])
    try:
        os.mkdir(link_basedir)
    except Exception:
        pass  # Already exists or failed (then next try will fail)
    try:
        os.symlink(install_path, link_path)
    except Exception as e:
        try:  # Maybe the path is already linked
            current_path = os.readlink(link_path)
            if current_path == install_path:
                return  # Then we're done
        except Exception:
            pass
        raise (type(e))('Could not link package: ' + str(e))


def _package_install_unsafe(path, package, pkgrepo, pkglist, pkgnames):
    """
    NOT A PUBLIC INTERFACE: use package_install[_many] instead.

    Try to install a package and forward exceptions to the caller.
    Will leave package repository in dirty state.
    Returns
    """
    write('Installing %s... ' % package.name)
    sys.stdout.flush()
    if package.name in pkgnames:
        index = pkgnames.index(package.name)
        if package.path in pkglist[index]['paths']:
            writeln('Already installed.')
            raise AlreadyInstalledException(link_updated=False)
        else:
            write('Already installed, linking to %s... ' % package.path)
            sys.stdout.flush()
            package_link(path, package)
            writeln('Linked.')
            raise AlreadyInstalledException(link_updated=True)
    # If the above did not return, we need to actually install the package
    try:
        if package.branch:
            pkgrepo.create_submodule(package.name, package.name,
                                     url=package.url, branch=package.branch)
        else:
            pkgrepo.create_submodule(package.name, package.name,
                                     url=package.url)
    except Exception:  # Default message is cryptic
        raise GitFetchException(package)
    package_link(path, package)
    writeln('Installed.')


def package_install(path, package, batch=False, pkgrepo=None, pkglist=None):
    """
    Install a package or roll back to last coherent state upon failure.
    If batch is True, do not update package list (caller will update).
    Returns True on success, else (error or already installed) False.
    """
    pkgpath = package_dir_path(path)
    if pkgrepo is None:
        pkgrepo = package_repo_open(pkgpath)
    if pkglist is None:
        pkglist = package_list_read(pkgpath)
    pkgnames = list(map(lambda x: x['name'], pkglist))
    try:
        _package_install_unsafe(path, package, pkgrepo, pkglist, pkgnames)
        pkgrepo.index.add(['.gitmodules', package.name])
        pkgrepo.index.commit('Installed ' + package.name)
        if not batch:
            package_list_add_many(pkgpath, [package])
    except AlreadyInstalledException as e:
        return e.link_updated
    except Exception as e:  # Installation failed, roll back
        try:
            sm = pkgrepo.submodule(package.name)
            sm.remove()
        except Exception:
            pass
        pkgrepo.git.clean('-fdX')  # Remove all untracked files
        writeln('Install aborted.')
        writeln(str(e))
        return False
    return True


def package_install_many(path, packages):
    """Install a list of packages"""
    packages = package_string_parse_many(packages)
    installed_packages = []
    pkgpath = package_dir_path(path)
    pkglist = package_list_read(pkgpath)
    pkgrepo = package_repo_open(pkgpath)

    for package in packages:
        if package_install(path, package, True, pkgrepo, pkglist):
            installed_packages.append(package)  # To be written to database
    if installed_packages:
        package_list_add_many(pkgpath, installed_packages)


def package_update_all(path):
    """Update all installed packages"""
    repo = package_repo_open(package_dir_path(path))
    for sm in repo.submodules:
        write('Updating %s... ' % sm.name)
        sm.update()
        writeln('Done.')
        if repo.is_dirty():  # Something has changed
            repo.index.commit('Updated ' + sm.name)


def package_uninstall(path, package, batch=False, pkgrepo=None, pkglist=None):
    """
    Uninstall a package or unlink from given location if linked to multiple.
    If batch is True, do not update package list (caller will update).
    Returns True on success, else (on error) False.
    """
    pkgpath = package_dir_path(path)
    if pkgrepo is None:
        pkgrepo = package_repo_open(pkgpath)
    if pkglist is None:
        pkglist = package_list_read(pkgpath)
    pkgnames = list(map(lambda x: x['name'], pkglist))

    write('Uninstalling %s... ' % package.name)
    sys.stdout.flush()
    if package.name not in pkgnames:
        writeln('Not installed.')
        return False
    try:
        paths = pkglist[pkgnames.index(package.name)]['paths']  # Get paths
        try:
            os.unlink(path + '/' + package.path)
        except Exception:
            writeln('%s not linked at %s.' % (package.name, package.path))
            return False
        if len(paths) > 1:  # Installed in multiple locations
            writeln('Unlinked from %s.' % package.path)
            if not batch:
                package_list_remove_many(pkgpath, [package])
            return True
        else:
            sm = pkgrepo.submodule(package.name)
            sm.remove(force=True)  # Remove even if there are local changes
            pkgrepo.index.add(['.gitmodules'])
            pkgrepo.index.remove([package.name], r=True)  # Remove recursively
            pkgrepo.index.commit('Uninstalled ' + package.name)
        if not batch:
            package_list_remove_many(pkgpath, [package])
    except Exception as e:
        writeln('Failed to uninstall %s: %s' % (package.name, e))
        return False
    else:
        writeln('Uninstalled.')
        return True


def package_uninstall_many(path, packages):
    """Uninstall a list of packages"""
    packages = package_string_parse_many(packages)
    uninstalled_packages = []
    pkgpath = package_dir_path(path)
    pkgrepo = package_repo_open(pkgpath)

    for package in packages:
        if package_uninstall(path, package, True, pkgrepo):
            uninstalled_packages.append(package)
    if uninstalled_packages:
        package_list_remove_many(pkgpath, uninstalled_packages)
