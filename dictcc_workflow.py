#!/usr/bin/python
# encoding: utf-8

import sys

from dictcc_argument_parser import ArgumentParser
from workflow import Workflow
from dictcc_translator import Dict


def main(wf):

    parser = ArgumentParser(wf.args)

    result = Dict(parser.search_string, parser.from_language, parser.to_language).translate()

    for translation in result.translation_tuples:
        wf.add_item(translation[1], translation[0], valid="yes", arg=translation[0])

    if result.n_results == 0:
        wf.add_item('"' + parser.search_string + '" not found', generate_subtitle(parser))

    wf.send_feedback()


def generate_subtitle(parser):
    if parser.from_language.isValid() and parser.to_language.isValid():
        return parser.from_language.name.capitalize() + " -> " + parser.to_language.name.capitalize()
    else:
        return ""

if __name__ == '__main__':
    wf = Workflow()
    sys.exit(wf.run(main))
