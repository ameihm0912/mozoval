#!/usr/bin/python2

import xml.etree.ElementTree as ET

import libmopreter as lm

class OvalParserException(Exception):
    pass

class OvalParserHints(object):
    schema_loc = '{http://oval.mitre.org/XMLSchema/oval-definitions-5}'
    linux_schema_loc = '{http://oval.mitre.org/XMLSchema/' \
        'oval-definitions-5#linux}'

    tag_obj = '%sobjects' % schema_loc
    tag_states = '%sstates' % schema_loc
    tag_tests = '%stests' % schema_loc
    tag_def = '*/%sdefinition' % schema_loc

    def_meta = '%smetadata' % schema_loc
    def_criteria = '%scriteria' % schema_loc
    def_criterion = '%scriterion' % schema_loc
    def_title = '%stitle' % schema_loc

    def_test_linux_object = '%sobject' % linux_schema_loc
    def_test_linux_state = '%sstate' % linux_schema_loc

    def_rpminfo_name = '%sname' % linux_schema_loc

class OvalState(object):
    @staticmethod
    def allocator(et):
        if 'rpminfo_state' in et.tag:
            return lm.RPMInfoState(et)
        return OvalState(et)

    def __init__(self, et):
        self.state_id = None

        if 'id' not in et.attrib:
            raise OvalParserException('state has no id')
        self.state_id = et.attrib['id']

class OvalObject(object):
    @staticmethod
    def allocator(et):
        if 'rpminfo_object' in et.tag:
            return lm.RPMInfoObject(et)
        return OvalObject(et)

    def __init__(self, et):
        self.object_id = None

        if 'id' not in et.attrib:
            raise OvalParserException('object has no id')
        self.object_id = et.attrib['id']

class OvalTest(object):
    @staticmethod
    def allocator(et):
        if 'rpminfo_test' in et.tag:
            return lm.RPMInfoTest(et)
        return OvalTest(et)

    def parse_linux_object_state(self, et):
        obj = et.find(lm.OvalParserHints.def_test_linux_object)
        if obj == None:
            raise OvalParserException('test has no object reference')
        self.object_ref = obj.attrib['object_ref']

        state = et.find(lm.OvalParserHints.def_test_linux_state)
        if state == None:
            raise OvalParserException('test has no state reference')
        self.state_ref = state.attrib['state_ref']

    def __init__(self, et):
        self.test_id = None
        self.object_ref = None
        self.state_ref = None
        self.comment = None

        if 'id' not in et.attrib:
            raise OvalParserException('test has no id')
        self.test_id = et.attrib['id']

        if 'comment' not in et.attrib:
            raise OvalParserException('comment has no id')
        self.comment = et.attrib['comment']

class OvalCriterion(object):
    def __init__(self, et):
        self.comment = 'Unknown'
        self.test_ref = None

        if 'test_ref' not in et.attrib:
            raise OvalParserException('criterion has no test reference')
        if 'comment' not in et.attrib:
            raise OvalParserException('criterion has no comment')

        self.test_ref = et.attrib['test_ref']
        self.comment = et.attrib['comment']

class OvalCriteria(object):
    TYPE_AND = 0
    TYPE_OR = 1

    def __init__(self, et):
        self.criteria_type = 0
        self.criterion_list = []
        self.criteria_list = []

        if 'operator' not in et.attrib:
            raise OvalParserException('criteria has no operator')
        opmode = et.attrib['operator']
        if opmode == 'AND':
            self.criteria_type = self.TYPE_AND
        elif opmode == 'OR':
            self.criteria_type = self.TYPE_OR
        else:
            raise OvalParserException('invalid criteria type %s' % opmode)

        c = et.findall(OvalParserHints.def_criteria)
        for i in c:
            self.criteria_list.append(OvalCriteria(i))

        c = et.findall(OvalParserHints.def_criterion)
        for i in c:
            self.criterion_list.append(OvalCriterion(i))

class OvalDefinition(object):
    def __init__(self, et):
        self.oval_id = None

        self.criteria = None

        self.meta_title = 'Unknown'

        self.parse_et(et)

    def parse_et(self, et):
        if 'id' not in et.attrib:
            raise OvalParserException('definition has no id')
        self.oval_id = et.attrib['id']

        meta = et.find(OvalParserHints.def_meta)
        if meta == None:
            raise OvalParserException('metadata element not found ' \
                'in definition')
        for i in meta:
            if i.tag == OvalParserHints.def_title:
                self.meta_title = i.text

        criteria = et.find(OvalParserHints.def_criteria)
        if criteria == None:
            raise OvalParserException('criteria element not found ' \
                'in definition')
        self.criteria = OvalCriteria(criteria)

class OvalChecks(object):
    def __init__(self):
        self.definitions = {}
        self.states = {}
        self.objects = {}
        self.tests = {}

    def add_definition(self, d):
        if d.oval_id in self.definitions:
            raise OvalParserException('duplicate id %s' % d.oval_id)
        self.definitions[d.oval_id] = d

    def add_state(self, s):
        if s.state_id in self.states:
            raise OvalParserException('duplicate id %s' % d.oval_id)
        self.states[s.state_id] = s

    def add_object(self, o):
        if o.object_id in self.objects:
            raise OvalParserException('duplicate id %s' % d.oval_id)
        self.objects[o.object_id] = o

    def add_test(self, t):
        if t.test_id in self.tests:
            raise OvalParserException('duplicate id %s' % d.oval_id)
        self.tests[t.test_id] = t

def parse_oval_checks(path):
    ret = OvalChecks()

    tree = ET.parse(path)

    root = tree.getroot()

    definitions = root.findall(OvalParserHints.tag_def)
    for d in definitions:
        new = OvalDefinition(d)
        ret.add_definition(new)

    objects = root.find(OvalParserHints.tag_obj)
    for o in objects:
        new = OvalObject.allocator(o)
        ret.add_object(new)

    states = root.find(OvalParserHints.tag_states)
    for s in states:
        new = OvalState.allocator(s)
        ret.add_state(new)

    tests = root.find(OvalParserHints.tag_tests)
    for t in tests:
        new = OvalTest.allocator(t)
        ret.add_test(new)
