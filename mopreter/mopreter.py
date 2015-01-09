#!/usr/bin/python2

import sys
import getopt

import libmopreter

def check_args(args):
    argspec = {
        'list': (
            1,
            'usage: mopreter.py list oval_spec_path\n'
        )
        }

    if args[0] not in argspec:
        usage()
    if (len(args) - 1) < argspec[args[0]][0]:
        sys.stdout.write(argspec[args[0]][1])
        sys.exit(0)

def do_list(path):
    libmopreter.parse_oval_checks(path)

def usage():
    sys.stdout.write('usage: mopreter.py [-h] command arguments...\n')
    sys.exit(0)

def mopreter():
    try:
        opts, args = getopt.getopt(sys.argv[1:], 'h')
    except getopt.GetoptError:
        usage()
    for o, a in opts:
        if o == '-h':
            usage()
    if len(args) < 1:
        usage()
    check_args(args)

    if args[0] == 'list':
        do_list(args[1])

if __name__ == '__main__':
    mopreter()

sys.exit(0)
