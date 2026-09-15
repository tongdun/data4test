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

#### Import Test Cases

##### Import from XMind
- On the Test Case Management page, click **Import** → **Import XMind** and upload the `.xmind` file. See "Usage Example" and "XMind File" above for the required format.

##### Import from Excel
- Click **Import** → **Import Excel** to open the import page.
- The import uses the same template mechanism as Excel export: the template is defined by the `testCaseExportTemplates` parameter, and the Excel header row (column titles) maps to case fields.
- Fill in the cases, then upload the Excel file; you may also upload an image package (`.zip`) to import test-process screenshots.
- Image package naming: each image is named `CaseNumber_N.format` (e.g. `CASE_001_1.png`) to associate it with the corresponding case.
- On submit, the system checks conflicts by `case number + module`:
  - New cases are inserted directly.
  - Existing cases can be skipped or overwritten (overwrite only updates non-empty fields).
- The sample row in the generated template (case number starting with `示例`) is skipped automatically.

#### Export Test Cases

##### Export Markdown / XMind
- Select cases, then click **Export** → **Export Markdown** / **Export XMind** to export in the corresponding format.

##### Export Excel
- Click **Export** → **Export Excel** and configure:
  - **Template**: select a predefined export template (defines columns, their order, and headers).
  - **Language**: choose the export language; localized fields are exported in that language, falling back to the default value if a translation is missing.
  - **Image Mode**: `Export paths (package images)` writes screenshot paths into cells and packages the images, or `Embed images` embeds screenshots directly into cells.
  - **Package Format**: `tgz` or `zip` (applies only to the "path" image mode).
  - **Filters**: product / module / introduced version / case designer / creation time — or export the checked cases directly.
- The result is a downloadable `.xlsx`, or a `.tgz`/`.zip` package containing the `.xlsx` and the referenced images.