#!/usr/bin/python
# encoding: utf-8


import urllib2
import requests
import re
from bs4 import BeautifulSoup


class Result(object):
    def __init__(self, from_lang=None, to_lang=None, translation_tuples=None):
        self.from_lang = from_lang
        self.to_lang = to_lang
        self.translation_tuples = list(translation_tuples) \
            if translation_tuples else []

    @property
    def n_results(self):
        return len(self.translation_tuples)

    @property
    def from_words(self):
        return map(lambda tuple: tuple[0], self.translation_tuples)

    @property
    def to_words(self):
        return map(lambda tuple: tuple[1], self.translation_tuples)

    @property
    def from_words_lowercase(self):
        return map(lambda tuple: tuple[0].lower(), self.translation_tuples)

    @property
    def to_words_lowercase(self):
        return map(lambda tuple: tuple[1].lower(), self.translation_tuples)


class Dict(object):

    def __init__(self, search_string, from_language, to_language):
        self.search_string = search_string
        self.from_language = from_language
        self.to_language = to_language

    @property
    def request_subdomain(self):
        subdomain = self.from_language.subdomain.lower() + self.to_language.subdomain.lower()

        if len(subdomain) > 4:
            return "www"
        else:
            return subdomain

    def translate(self):
        response = self.get_response()
        result = self.parse_response(response.content)
        return self.correct_translation_order(result)

    def get_response(self):
        subdomain = self.request_subdomain

        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 6.3; WOW64; rv:30.0) Gecko/20100101 Firefox/30.0'
        }

        params = {
            "s": self.search_string
        }

        return requests.get("https://" + subdomain + ".dict.cc", params=params, headers=headers)

    def parse_response(self, response_body):

        in_list = []
        out_list = []

        def sanitize(word):
            return re.sub("[\\\\\"]", "", word)

        javascript_list_pattern = "\"[^,]+\""

        for line in response_body.split("\n"):
            if "var c1Arr" in line:
                in_list = map(sanitize, re.findall(javascript_list_pattern, line))
            elif "var c2Arr" in line:
                out_list = map(sanitize, re.findall(javascript_list_pattern, line))

        if not any([in_list, out_list]):
            return Result()

        soup = BeautifulSoup(response_body, "html.parser")

        left_lang = soup.find_all("td", width="307")[0].b.text
        right_lang = soup.find_all("td", width="306")[0].b.text

        in_list = map(lambda word: unicode(word, 'utf-8'), in_list)
        out_list = map(lambda word: unicode(word, 'utf-8'), out_list)

        return Result(
            from_lang=left_lang,
            to_lang=right_lang,
            translation_tuples=zip(in_list, out_list),
        )

    def correct_translation_order(self, result):

        if not result.translation_tuples:
            return result

        left_occurrences = len(filter(lambda word: word.count(self.search_string.lower()), result.from_words_lowercase))
        right_occurrences = len(filter(lambda word: word.count(self.search_string.lower()), result.to_words_lowercase))

        if left_occurrences >= right_occurrences:
            return result
        else:
            return Result(from_lang=result.to_lang,
                          to_lang=result.from_lang,
                          translation_tuples=zip(result.to_words, result.from_words)
                          )
