#!/usr/bin/python
# encoding: utf-8

import os
from dictcc_language import Language, UnavailableLanguageError


class ArgumentParser(object):

    def __init__(self, args):
        self.search_string = " ".join(args)
        self.from_language = Language(os.getenv('from_language'))
        self.to_language = Language(os.getenv('to_language'))

        if len(args) < 3:
            return

        try:
            self.to_language = Language(args[1])
            self.from_language = Language(args[0])
        except UnavailableLanguageError:
            return

        self.search_string = " ".join(args[2:])
