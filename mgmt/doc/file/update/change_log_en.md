# Release Notes
## 2024-01-05
1. Added the ability to generate test scenarios using the list of API interfaces. The feature allows users to select all associated data files with an interface for scenario creation.
2. Fixed an issue where assertion values with expressions caused the template to be retrieved prematurely.
3. Fixed an issue where assertion values of the form '=', when set to negative numbers, displayed an error message indicating an incorrect expression.

## 2024-01-04
1. Added support for indexing when multiple placeholders exist within a single field value. This allows users to retrieve input parameter values from upstream test cases based on the index of the placeholder.
2. Fixed an issue where multiple assertions on a single data file caused the number of results to exceed the number of data records. Now, each assertion is associated with a specific data record.
3. Optimized data request handling for GET requests. If the input parameter value is empty, it will not be wrapped in a request payload, allowing for more streamlined processing.
4. Optimized console request formatting, making it easier to view and compose requests.
5. Optimized the console domain selection feature, ensuring that the selected domain updates the application prefix information automatically.

## 2024-01-03
1. Fixed an issue where console data domain executions failed to record historical information when assertion failures occurred.

## 2023-12-28
1. Increased the length limit for data and scenario descriptions in the console input boxes.

## 2023-12-25
1. Added support for multiple languages. Fixed an issue where query variable language and value were swapped, resulting in query placeholder variables not being replaced properly.
2. Optimized regex pattern matching for placeholder variables, supporting English, numbers, underscores, and hyphens.

## 2023-12-18
1. Optimized data domain filtering in the console, allowing users to search for specific data records or perform full-data filtering for both data names and file names.
2. Optimized scenario domain filtering in the console, allowing users to search for specific scenarios or perform full-scenario filtering based on scenario names.

## 2023-12-16
1. Changed the local file manager dependency in go.mod to a remote library for improved maintainability and accessibility.

## 2023-12-14
1. Added a new feature to record test data automatically for POST requests. The action "record_csv" and "record_xls" can be used to record response data as CSV or XLS files, ensuring offline and real-time data consistency.
2. Added a new feature to generate test data based on templates using the action "modify_file". It automatically records response data into new files, suitable for POST requests and external data generation.
3. Added support for placeholder replacement in Mock API responses (/mock/file/name.xml?lang=en). It supports multi-language data generation and only accepts system-generated or predefined variables as placeholder variables.
4. Fixed an issue where task status, type, mode, and other filters displayed incorrectly in the task list.
5. Fixed an issue where saving a data file's annotation would erase any linked files associated with it.
6. Optimized console execution by returning a status code of "pass" instead of "%!s()" upon successful execution.
7. Fixed an issue where setting the request format to "raw" would cause the Content-Type header to be set incorrectly in the console's data domain editor.

##### February 1st, 2023
1. [Optimize] Nested array data extraction supported for returned data, e.g., data:[[{"msg": "Passed via XXX", "timeUnit": "MONTH"}]], extract time unit data via data**timeUnit
2. [Bug] Fixed issue with Single Body request data losing nested data

##### January 29th, 2023

1. [Optimize] Interface ID field in data list page supports fuzzy search

##### January 28th, 2023

1. [Optimize] Header automatically written after selecting data type in console (100% completion), synchronized changes on front and back ends, including defined domain, running domain, data domain, and data history domain

##### January 12th, 2023

1. [Optimize] Header automatically written after selecting data type in console (20% completion)

##### January 10th, 2023

1. [Feature] Added new parameters {Month(int)}, {Year(int)}, {Month}, {Year} for automatic generation

##### January 4th, 2023

1. [Optimize] All file formats are written back, complex JSON files are replaced with placeholders after writing, and placeholders are retrieved from the library for the next execution
2. [Optimize] Data is saved to the data file from the library and then executed during scenario execution
3. [Optimize] Compatible with extracting values of []interface{} type, default to the first one and convert to string

##### January 3rd, 2023

1. [Optimize] If there is an output for application/json format data, it will be written back

##### December 5th, 2022

1. [Feature] Nested List dictionary parameters supported for loop replacement

##### November 30th, 2022

1. [Feature] Added automation field to interface definition list, synchronized changes to statistical report logic

##### November 29th, 2022

1. [Bug] Fixed issue where string parameters with list JSON format in Single are cleared when there is no placeholder

##### November 25th, 2022

1. [Optimize] Execution request parameters and header information are not returned when encountering errors in console data domain, logic optimized and adjusted accordingly

##### November 22nd, 2022

1. [Bug] Fixed issue where console scenario domain cannot be properly maintained due to mismatched field docking types

##### November 17th, 2022

1. [Feature] Added scenario history domain to console, including historical data viewing, running, and saving functionalities
2. [Optimize] Fixed issue where Xmind import test cases display failure but actually succeed due to empty product name error, filtered directly
3. [Bug] Fixed issue with test case management filter function abnormality
4. [Bug] Scenario test history scene types do not display Chinese names correctly
5. [Bug] Result is not synchronized after executing comparison type scenarios

##### November 8th, 2022

1. [Optimize] Improved logic and sequence of scenario domain execution errors and exceptions
2. [Bug] Fixed issue where saving modified scenarios in the scenario list and loading them in the console does not display correctly due to file not found issues during execution in the scenario list and domain fields' logic re-ordered for loading file by name first before saving as default for execution fields that require input of the current user and for multi-selection scenarios. The logic of automatically saving scenarios has been changed to only occur when there are pre- or post-related data configurations in the data domain fields. This change ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. The change also ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. This change also ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. This change also ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. This change also ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. This change also ensures that the correct file is loaded when executing a multi-selection scenario and that the correct file is saved when saving a multi-selection scenario. This change also ensures

##### September 13, 2022
1. [Feature] Added support for cascading data generation. For example: China-Zhejiang-Jinhua, define the cascading data definition first, and then access the data through {TreeData_China[3]},{TreeData_China[1]},{TreeData_China[2]}. {TreeData_China[3]} represents the third-level data: Jinhua, {TreeData_China[1]} accesses the first-level data: China, and {TreeData_China[2]} represents the second-level data: Zhejiang. Cascading data can be defined by users, currently supporting up to three levels of access, and arbitrary three-level cascading data can be defined for use.

##### August 16, 2022
1. Optimized feature data generation code, fixed bugs, and added several more feature data generations.

##### August 15, 2022
1. Added support for body data that is directly a List request (only file writing, console not yet adapted).
2. Added support for complex multi-layer body request data. The first and second layers support placeholder replacement, while deeper layers are requested as-is.
3. Added support for delete operations with query parameters.
4. When a data file does not exist, it is automatically created at runtime (applicable for platforms with SQL data migration, allowing the elimination of data packages and on-demand creation during execution).

##### August 12, 2022

1. Added support for generating GMT date formats several days ago or days later using TimeGMT(-3).

##### August 11, 2022

1. Added support for generating various formatted time strings using the keyword TimeFormat.
2. Added support for generating "Fri Aug 26 2022 23:59:59 GMT+0800" formatted time strings using TimeGMT.

##### August 2, 2022

1. Added support for variable placeholder replacement in Body parameters with Single when the variable contains a dictionary.

##### August 1, 2022

1. Fixed the issue where variable replacement was not successful in Body parameters with Single.

##### July 28, 2022

1. Fixed the issue where the creator would be incorrectly copied when copying data from the scenario list or data list due to variable scope issues.

##### July 21, 2022

1. Fixed the issue where multiple variable replacements were not fully replaced in some cases.

##### July 19, 2022

1. [Front-end] The prepath path in the data domain now updates dynamically based on the content of the data file.

##### July 14, 2022

1. Changed the sleep duration for CICD node interface check from sleep2s to sleep10s when fetching interface data.

##### July 13, 2022

1. Interrupt subsequent executions if test failures occur during multiple executions of scenarios or data. (Note: Details on this change were not provided.)

##### July 5, 2022

1. Fixed the issue where the dashboard page could not be loaded properly when no data was running in the system.
2. Fixed the issue where cross-module interface description information was not cleared correctly in the control console's Run domain.
3. Modified the control console's Run domain to allow modifying existing data description information without clearing it.
4. Changed the logo from "Data4Perf" to "盾测-自动化" (DunCe - Automation).

##### June 13, 2022

1. Added the ability to write environment type information to the scenario test history form and data test history form in the system.

##### May 25, 2022
1. When the service restarts, the running tasks will be automatically restarted.
2. When a panic occurs in the service, recovery will be performed.
3. The issue of occasional panic during the execution of scheduled tasks has been fixed.
4. The related scenarios of task management and related data files are sorted in descending order of time.
5. The CICD interface specification check node has added a retry mechanism to solve the issue where the upstream node service has not fully started but has already turned green.

##### April 10, 2022
1. Supports start time of day, month, year, end time, timestamp, etc.

##### February 11, 2022
1. Added an interface call without authentication to interface specification checks.

##### December 7, 2021

1. Fixed the issue where the domain method on the interface management console could not be displayed normally on the front end.

##### December 3, 2021

1. Added an interface management console.
2. Reformed the swagger import to adapt to the new data structure.

##### November 10, 2021

1. Modified the ownership of the login page to display "TongDun".

##### November 9, 2021

1. Fixed the issue with maintaining placeholders for complex structures under single (undo).

##### November 8, 2021

1. Supports the use of placeholders for complex structure data under single parameters.
2. Supports viewing creator information in the scenario list.
3. Supports viewing creator information in the data file list.
4. Supports viewing creator information in the task list.
5. Adjusted the width of some lists.

##### November 4, 2021

1. Captures exceptions for assertion validation structure mismatches.
2. Supports comparing defined variables and previous file data in assertions.
3. Supports selecting data files and scenarios to execute as tasks.
4. Prints header information instead of global environment information for failed test results.

##### November 3, 2021

1. Improved file import error handling by providing error messages when files are not found.

##### November 1, 2021

1. Added a related application field to the product list.
2. Added global dimension data statistics reports.
3. Added product dimension data statistics reports.
4. Added application dimension data statistics reports.

## 2021-10-21
1. Fixed issue where scene list couldn't find JSON files.

## 2021-10-19
1. Single-parameter query supports more data types.

## 2021-10-15
1. Fixed issue where test results of JSON format data files were written back as YAML format. Now applies to JSON format.
2. Fixed issue where there were problems with special characters in the body request data.

## 2021-10-14
1. Fixed issue with generating specified length Chinese characters.
2. Supports editing data files in JSON format and testing, as well as data writing back.

## 2021-10-09
1. Temporary fix for issue with copying scenes and importing SQL updates resulting in errors.

## 2021-09-22
1. Added more check items to the API specification check.

## 2021-09-14
1. Supports obtaining API documentation information from Swagger paths.
2. Added API specification check functionality.
3. Added page to display results of API specification checks.

## 2021-09-13
1. Fixed issue where generating scene data files from interface definitions, test data, and fuzzy data would overwrite existing files due to missing index names as combined unique conditions.

## 2021-09-13
1. Fixed issue with generating data file file names not correctly generated in interface definition section due to missing placeholders.
2. Fixed issue with generating scene data files from test data encountering integer types and failing due to mismatched data types.
3. Added recognition of `formData` types in Swagger API document parsing.
4. Fixed issue with generated history files not being viewable due to missing file suffix ".yml".
5. Adjusted list widths for better data display.

## 2021-09-10
1. Captured and threw error messages when data files are not found.
2. Removed legacy code for obtaining JSON format data files.
3. Supports executing tasks based on priority of associated scenes.
4. Supports appending error messages for return.
5. Scene list supports viewing modification times and synchronized update of write-back results.

## 2021-09-03
1. When generating scene data files from test data, module data, and interface definitions, appended random strings to names and file names to reduce chances of overwrite when using generated data.
2. Supports managing test cases.
3. Supports importing test cases from Xmind (requires xmind2case third-party tool).

## 2021-09-02
1. The database fields have been modified to ensure consistency between the code and the actual meaning.
2. The newline compatibility of the input data file for scenario lists has been optimized.
3. The issue of content not being copied properly when copying data files has been fixed.
4. The database initialization file has removed forced encoding, making it more compatible with various MySQL versions.
5. The built-in data in the database initialization file has been optimized.
6. The functionality and logic of interface definitions, test data, and scenario data file generation for fuzzy testing have been improved and optimized.
7. Support for JSON format data files has been removed, and related functional code has been deleted.

# 2021-08-31
1. The project was officially released to the public for the first time.

# 2021-08-04
1. The functionality of data files has been enriched and the logic has been improved.

# 2021-07-28
1. The project was deployed on the intranet.