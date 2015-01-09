import sys
import re
import subprocess

import libmopreter as lm

rpminfomgr = None

class RPMInfoManager(object):
    def rpm_qa(self):
        lm.parser_output('Gathering installed packages...\n')
        cmd = ['rpm', '-qa']
        o = subprocess.Popen(cmd, stdout=subprocess.PIPE,
            stderr=subprocess.PIPE)
        out, err = o.communicate()
        lns = out.split('\n')
        for i in lns:
            self.rpm_pkgs.append(i)

    def rpm_version_string(self, p):
        matchexp = '%s-(.*)' % re.escape(p)
        for i in self.rpm_pkgs:
            ret = re.search(matchexp, i)
            if ret == None:
                continue
            return ret.group(1)
        return None

    def __init__(self):
        self.rpm_pkgs = []

        lm.parser_output('RPMInfoManager is initializing\n')

        self.rpm_qa()

class RPMInfoObject(lm.OvalObject):
    def __init__(self, et, checks):
        super(RPMInfoObject, self).__init__(et, checks)

        rpminfomgr_prod()

        self.installed = False
        self.rpminfo_name = None
        self.rpm_version_string = None

        name = et.find(lm.OvalParserHints.def_rpminfo_name)
        self.rpminfo_name = name.text

    def prepare(self):
        if self.prepared:
            return
        lm.parser_output('        Preparing RPM object %s for evaluation\n' % \
            self.rpminfo_name)
        self.rpm_version_string = \
            rpminfomgr.rpm_version_string(self.rpminfo_name)
        if self.rpm_version_string == None:
            self.installed = False
        else:
            lm.parser_output('        Package is installed\n')
            self.installed = True
        self.prepared = True

class RPMInfoState(lm.OvalState):
    def __init__(self, et, checks):
        super(RPMInfoState, self).__init__(et, checks)

        rpminfomgr_prod()

        # There are a few types of evaluations we support right now, some are
        # just ignored and some are handled, others will result in an exception

    def evaluate(self, obj):
        pass

class RPMInfoTest(lm.OvalTest):
    def __init__(self, et, checks):
        super(RPMInfoTest, self).__init__(et, checks)

        rpminfomgr_prod()

        self.parse_linux_object_state(et)

    def execute(self):
        if self.state_ref not in self.checks.states:
            raise OvalParserException('referenced state %s not found' % \
                self.state_ref)
        if self.object_ref not in self.checks.objects:
            raise OvalParserException('referenced object %s not found' % \
                self.object_ref)
        obj = self.checks.objects[self.object_ref]
        self.checks.states[self.state_ref].evaluate(obj)

def rpminfomgr_prod():
    global rpminfomgr

    if rpminfomgr == None:
        rpminfomgr = RPMInfoManager()
