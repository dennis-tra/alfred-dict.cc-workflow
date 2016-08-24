#!/usr/bin/python
# encoding: utf-8

import sys

from dictcc_argument_parser import ArgumentParser
from workflow import Workflow
from dictcc import Dict, AVAILABLE_LANGUAGES


def main(wf):

    parser = ArgumentParser(wf.args)

    wf.send_feedback()

if __name__ == '__main__':
    wf = Workflow()
    sys.exit(wf.run(main))
