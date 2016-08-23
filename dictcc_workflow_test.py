import unittest
import dictcc_workflow as wf

class ModuleImportTests(unittest.TestCase):

    def testAlfredWorkflowModule(self):
        alfredWorkflow = wf.Workflow()
        self.failIf(alfredWorkflow is None)

    def testDictccModule(self):
        self.failIf(wf.Dict() is None)

def main():
    unittest.main()

if __name__ == '__main__':
    main()