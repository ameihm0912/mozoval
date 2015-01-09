import sys

import libmopreter as lm

class RPMInfoObject(lm.OvalObject):
    def __init__(self, et, checks):
        super(RPMInfoObject, self).__init__(et, checks)

        self.rpminfo_name = None

        name = et.find(lm.OvalParserHints.def_rpminfo_name)
        self.rpminfo_name = name.text

class RPMInfoState(lm.OvalState):
    def __init__(self, et, checks):
        super(RPMInfoState, self).__init__(et, checks)

class RPMInfoTest(lm.OvalTest):
    def __init__(self, et, checks):
        super(RPMInfoTest, self).__init__(et, checks)

        self.parse_linux_object_state(et)
