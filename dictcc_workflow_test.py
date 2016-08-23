import unittest
import dictcc_workflow as wf

# Here's our "unit".
def IsOdd(n):
    return n % 2 == 1

# Here's our "unit tests".
class IsOddTests(unittest.TestCase):

    def testOne(self):
        self.failUnless(IsOdd(1))

    def testTwo(self):
        self.failIf(IsOdd(2))


class ModuleImportTests(unittest.TestCase):

    def testAlfredWorkflowModule(self):
        self.failIf(wf.Workflow() is None)

def main():
    unittest.main()
    wf.printer()

if __name__ == '__main__':
    main()