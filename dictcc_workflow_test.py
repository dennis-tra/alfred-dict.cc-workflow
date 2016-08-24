#!/usr/bin/python
# encoding: utf-8

import unittest
import dictcc_workflow as wf
from dictcc_language import Language, UnavailableLanguageError
from dictcc_argument_parser import ArgumentParser
from dictcc_translator import Dict, Result

class ModuleImportTests(unittest.TestCase):

    def testAlfredWorkflowModule(self):
        alfredWorkflow = wf.Workflow()
        self.failIf(alfredWorkflow is None)


class ArgumentParserTests(unittest.TestCase):

    def testNoLanguagesGiven(self):
        parser = ArgumentParser(["string"])

        self.assertEqual(parser.from_language, Language(None))
        self.assertEqual(parser.to_language, Language(None))
        self.assertEqual(parser.search_string, "string")

    def testNoLanguagesTwoWordsGiven(self):
        parser = ArgumentParser(["guitar", "string"])

        self.assertEqual(parser.from_language, Language(None))
        self.assertEqual(parser.to_language, Language(None))
        self.assertEqual(parser.search_string, "guitar string")

    def testOneLanguageTwoWordsGiven(self):
        parser = ArgumentParser(["en", "guitar", "string"])

        self.assertEqual(parser.from_language, Language(None))
        self.assertEqual(parser.to_language, Language(None))
        self.assertEqual(parser.search_string, "en guitar string")

    def testTwoLanguagesOneWordGiven(self):
        parser = ArgumentParser(["en", "ger", "string"])

        self.assertEqual(parser.from_language, Language("eng"))
        self.assertEqual(parser.to_language, Language("de"))
        self.assertEqual(parser.search_string, "string")

    def testTwoLanguagesTwoWordsGiven(self):
        parser = ArgumentParser(["swe", "ru", "öl", "bier"])

        self.assertEqual(parser.from_language, Language("sv"))
        self.assertEqual(parser.to_language, Language("rus"))
        self.assertEqual(parser.search_string, "öl bier")

class DictccTranslator(unittest.TestCase):

    def testOrderCorrection(self):
        dict = Dict("Hallo", Language("ger"), Language("eng"))

        in_list = ['Hello!', 'Hi!', 'Howdy!', "G'day!", 'Hiya!', 'Hallo!', 'Heyday!', 'Gidday!', 'Hullo!']
        out_list = ['Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!']

        tuple = zip(in_list, out_list)

        res = Result("Deutsch", "Englisch", tuple)
        result = dict.correct_translation_order(res)

        self.assertEqual(result.from_words[0], "Hallo!")

    def testOrderCorrection2(self):
        dict = Dict("Hallo", Language("ger"), Language("eng"))

        out_list = ['Hello!', 'Hi!', 'Howdy!', "G'day!", 'Hiya!', 'Hallo!', 'Heyday!', 'Gidday!', 'Hullo!']
        in_list = ['Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!', 'Hallo!']

        tuple = zip(in_list, out_list)

        res = Result("Deutsch", "Englisch", tuple)
        result = dict.correct_translation_order(res)

        self.assertEqual(result.from_words[0], "Hallo!")

class LanguagesTests(unittest.TestCase):

    def testClassVariablesHaveSameLength(self):
        assert len(Language.abbreviations_list) == len(Language.names) == len(Language.subdomains)

    def testValidLanguageSetAppropriately(self):
        lang = Language("en")
        self.assertEqual(lang.name, "english")
        self.assertEqual(lang.subdomain, "en")

    def testValidLanguageSetAppropriately2(self):
        lang = Language("rus")
        self.assertEqual(lang.name, "russian")
        self.assertEqual(lang.subdomain, "ru")

    def testInvalidLanguageRaisesError(self):
        self.assertRaises(UnavailableLanguageError, Language, "asdf")

    def testNoneInit(self):
        lang = Language(None)
        self.assertEqual(lang.name, None)
        self.assertEqual(lang.subdomain, "www")


def main():
    unittest.main()

if __name__ == '__main__':
    main()