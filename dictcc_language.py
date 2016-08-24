#!/usr/bin/python
# encoding: utf-8


class UnavailableLanguageError(Exception):
    def __str__(self):
        return "Languages have to be in the following list: " + str(Language.abbreviations_list)


class Language:

    abbreviations_list = [
        [None],
        ["en", "eng"],
        ["de", "ger"],
        ["fr", "fra"],
        ["sv", "swe"],
        ["es", "esp"],
        ["bg", "bul"],
        ["ro", "rom"],
        ["it", "ita"],
        ["pt", "por"],
        ["ru", "rus"]
    ]

    names = [
        None,
        "english",
        "german",
        "french",
        "swedish",
        "spanish",
        "bulgarian",
        "romanian",
        "italian",
        "portuguese",
        "russian"
    ]

    subdomains = [
        "www",
        "en",
        "de",
        "fr",
        "sv",
        "es",
        "bg",
        "ro",
        "it",
        "pt",
        "ru"
    ]

    def __init__(self, language):

        index = self.find_index(language)

        self.name = Language.names[index]
        self.subdomain = Language.subdomains[index]

    def find_index(self, language):
        for index, abbreviations in enumerate(Language.abbreviations_list):
            if abbreviations.count(language) > 0:
                return index

        raise UnavailableLanguageError

    def isValid(self):
        return self.name is not None

    def __eq__(self, other):
        return isinstance(other, self.__class__) and self.__dict__ == other.__dict__

    def __ne__(self, other):
        return not self.__eq__(other)
