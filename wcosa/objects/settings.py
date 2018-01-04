"""@package parsers
Gives access to the settings file data. This contains settings for project that we can change
"""

import json

from wcosa.utils import helper


class Settings:
    """Settings class to parse settings.json file"""

    def __init__(self):
        with open(helper.get_settings_path()) as f:
            self.settings_data = json.load(f)


__settings = Settings()


def get_settings_value(key):
    """Returns the value of settings key"""

    return __settings.settings_data[key]
