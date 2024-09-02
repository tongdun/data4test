#### Usage Example:

1. **Product List**: Create a new example environment named "Demo Product".

2. **File - Test Case Management Page**: Import a test case file (DemoProduct.xmind).

3. **Product List**: Select the corresponding environment name, click "Import Test Cases from XMind", and you will see the corresponding test case data under the Test Cases menu.

#### XMind File

- Demo: `doc/file/project_V1.0.0_testcase_demo.xmind`

- Field Descriptions:
  - Fields are separated by underscores (_).
  - `project` should match the project name entered in the test environment.
  - Version `V1.0.0` is used to set the case's version, reflected through `case_id`.
  - Module names in the demo use Chinese and English separated by a hyphen (-), reflected through `case_id`. If not set, defaults to `other`.

#### Test Case Composition

- **Case ID**: Classified by module for data statistics and analysis.
- **Case Name**: Focuses on functionality, with a testing emphasis for each case.
- **Case Type**: Functional/Process/UI/Exception/Usability/Regression/Performance/Security/Compatibility/Scenario/Stress/Long-Running/Environment/Data/Copywriting/Style/Interaction/Boundary, customizable.
- **Priority**: P1/P2/P3/P4 or High/Medium/Low, customizable.
- **Preconditions**: Setup required before the test.
- **Test Scope**: Defines the area of the system being tested.
- **Test Steps**: Detailed actions to be performed during the test.
- **Expected Results**: Outcomes expected after executing the test steps.
- **Test Process**: Records of the testing process, typically screenshots.
- **Automation**: Yes/No, indicating if the test case can be automated.
- **Feature Developer**: Person responsible for developing the tested feature.
- **Case Designer**: Person who designed the test case.
- **Case Executor**: Person who executed the test case.
- **Test Time**: Record of when the test was executed.
- **Test Result**: Pass/Partially Pass/Fail/Not Tested/Deprecated.
- **Case Module**: Sets the module the test case belongs to for data statistics and analysis.
- **Introduced Version**: Sets the version information referenced by the test case.
- **Associated Scenario**: Can be linked to automated data/scenarios for execution via automation.
- **Associated Product**: Links to the execution environment or product the test case belongs to.
- **Notes**: Clarifications, issue tickets, change reasons, etc.