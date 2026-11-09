### Todo:
#### System Functionality
1. Provide standardized data use case templates: CRUD, exceptional scenarios, concurrent scenarios, performance testing, etc.
2. Complex scenario interfaces tested in proportion, allowing for controlled data composition.
3. Support statistical data output: sum, count, avg, etc.
4. Export to Swagger documentation (low priority)
5. Save and select assertion templates (assertion value templates already implemented)
6. Mock functionality (already partially implemented):
- Mock call recording
- Automatic generation of mock data
- Automatic generation of mock data with specified characteristics (V)
- Automatic generation of boundary mock data
7. Interface status management (research and development process control):
- Interface change prompts for change status after interface data modification
- Interface status history recording
- Automatic interface status maintenance (definition, docking, debugging, stability, testing, release):
    - Definition: Edited interfaces
    - Docking: Interfaces that have been mocked
    - Debugging: Interfaces that have been run in the development environment
    - Stability: Interfaces that have passed testing with N test data in the development environment
    - Testing: Interfaces that have begun running in the testing environment
    - Release: Interfaces that have passed testing with N test data in the testing environment
8. Data statistics:
- Increase data statistics for interface process management
- Single interface test frequency statistics
- Project interface test statistics
- Concurrent testing environment data collection and display
- Concurrent testing performance data collection and display

#### Usability Features
1. Redirect to a guide document on first entry, then redirect to the official login page (low priority)
2. Control console, scenario domain: Add an entry for full-text search of scenarios.
3. Data history domain: Format request and response data in the console for easier result viewing.
4. Return format errors back to the front end when saving data with improper formatting.
5. Control console supports file type uploads (currently done through direct path association).
6. Add upstream output parameter records to historical data files for easier issue localization and troubleshooting.
7. Robustness when body request data is raw in the control console.
8. When the application in a data file is inconsistent with the application in the data list, how should priority be handled? When fetching public environment info, which takes precedence? In a scenario, only data files exist, so it can only be associated with the data file's application.

#### Performance Optimization
1. Optimize SQL statements for faster report loading when dealing with large amounts of data (done via async).
2. Design parallel data generation to improve performance and control concurrency; currently serial, with controllable concurrency and batch writes.
3. Optimize variable substitution to a whole-string replacement approach (instead of post-structure replacement) to support variables globally.

#### User Feedback Features
1. Evaluate and potentially remove or migrate the parameter definition module functionality.
2. Compatibility with OPTIONS, HEAD, PATCH, and TRACE request types in some Swagger documents.
3. Check the environment's authentication and other non-functional status before testing.
4. Support pausing for one-time execution tasks: when data or scenario execution reaches a certain count, automatically create a task that can be paused and managed.
5. Consider recording process variables for better issue localization; log by level to control output levels (low priority).
6. Add Swagger visualization to view interface information.
7. [Currently not handled; use the top-level title filter if a dropdown is needed] In the interface definition list, when adding a new one, the application is not a dropdown association (medium priority).
8. Add scenario execution types and remove the concurrency switch from products and applications, letting scenarios and data control it themselves (high priority).
9. Design a scheme for iterating over enumerated fields in system parameters (medium priority).
10. [Overall supported, indexed one pending] Support direct associated references for List data, and also range references [0,10], (0,9) (high priority).
11. postman file transfer, pytest data file transfer (low priority).
12. Continue optimizing log information for easier problem localization (doing).
13. Data history modification has issues (see if it needs to be exposed for quantitative data repair and replay).
14. Support RT assertions, covering normal function return time; timeout assertions need handling.
15. Add correctness validation when creating new data to reduce exceptions caused by incorrect names (global, when creating data).
16. When a newly created data has an empty associated scenario, re-entering shows some data not displayed.
17. Remove product information from the application list to decouple permissions (Done).
18. Show whether parameters are required in the console interface definition and running domains.
19. After authentication expires, console requests do not automatically redirect to the login page.
20. When a task associates a scenario and the first data file has a variable that does not exist, after execution the scenario history and data history are not recorded.
21. Console data domain does not support saving List JSON format; enhance JSON format.
22. When interface request data is an array, "run again" needs compatibility.
23. Console data domain: when a path variable has multiple values, only one request is sent.
24. The array management component uses features the console found unsupported; schedule time to merge (e.g., complex JSON format requests).
25. Optimize regex matching to end at the next top-level element instead of a specified keyword; currently misjudges; finally, except for the assert section, re-match undefined parameters.
26. Design a scheme for the correctness judgment of creation time, modification time, etc.
27. Add a test_time field to data file results to distinguish data-list debug records from scenario-list result records that were not executed.
28. Offline tasks need callback actions for polling: the data-layer action polling engine and scenario-layer branch execution engine.
29. Change environment authentication maintenance from JSON to YML format to reduce JSON configuration errors.
30. Export a product's full set of use cases to Xmind format.
31. If a data file's request parameters include a file, support direct upload.
32. When a task associates data, one-click export has issues; enhance export compatibility.
33. When selecting a product in the console, if the product has a dedicated application prefix, it is not auto-populated; needs synchronized modification.
34. Optimize the acquisition of application prefixes and data interaction during task/scenario/data execution, avoiding a database query when it already exists.
35. If a data file cannot find an action, do not error; automatically add an empty placeholder.
36. For one-time tasks that are running, restarting the service does not support re-running (status cannot be repaired; the feature cannot be supported).
37. When re-importing one-click exported SQL with an inconsistent product, a new scenario is generated; review the mechanism.
38. Test case detail view: result history not displayed correctly (normal after verification).
39. Consider whether to lift the limit on loop execution counts for tasks associating scenarios/data.
40. Actions support Excel template data replacement to generate new import file data (low priority).
41. In the test case list, modifying non-test-result fields automatically updates the case executor (pending reproduction).
42. In scheduled task reports, the execution trend pass/fail data statistics have issues (pending reproduction).
43. The test case import/export template probably needs to be separated from system parameters.
44. Test case statistics definition page Excel export can support setting first-level and second-level menu directories, exporting by hierarchy.
45. Evaluate and plan integrating AI for one-click automatic translation (Done).
46. CLI command-line help does not support English.

#### Code Optimization
1. Switch to a pipelined communication mechanism for concurrency improvement.

#### Operation and Maintenance Features:
1. Package all files for deployment as a whole instead of individual files; perform file-level changes separately (low priority).
2. View historical versions to identify who modified the data (low priority).
3. Regularly clean up historical records to conserve environment resources (low priority).
4. Regularly backup management domain data in MySQL.
5. Regularly backup files under the upload directory.
6. Monitor the system under test environment (low priority).

---------------------------------------------------------------------------------------
#### BUG List
1. In the console data domain, when the request format is raw, Content-Type is set to raw, and the frontend parameter passing has issues [console].
2. [Verified OK, keep tracking for a while] After changing the environment in the console scenario domain, selecting a scenario does not sync associated data changes [console].
3. [Only allow saving as .yml] yaml files auto-saved as .yaml cannot be found; save as .yml first, then improve compatibility and fix the file validation logic.
4. The request order becomes scrambled when running requests in a scenario: Decision-Strategy Management-Scenario Strategy-Modify Strategy-illegal modification.
5. [Verified OK, keep tracking] Execution has issues when use cases are stored as .json files (high priority).
6. When scenario management associates too many data files, saving logs out (medium priority).
7. When a scenario associates the same data file multiple times and the pre-data file fails, the historical record interception position is wrong, resulting in missing data.
8. When a data file lacks an action:[] placeholder, a "slice bounds out of range [:-1]" error occurs.
9. Scenario history list "run again" execution is abnormal.
