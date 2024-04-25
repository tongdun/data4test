### Todo:
#### System Functionality
1. Provide standardized data use case templates: CRUD, exceptional scenarios, concurrent scenarios, performance testing, etc.
2. Complex scenario interfaces tested in proportion, allowing for controlled data composition.
3. Support statistical data output: sum, count, avg, etc.
4. Export to Swagger documentation (low priority)
5. Support script execution types: shell, Python, etc.
6. Save and select assertion templates (already implemented assertion value templates)
7. Mock functionality (already partially implemented):
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
4. One-click access to unautomated interface data.
5. Add failure reason records to historical data files for easier historical information tracing.
6. Format errors and redirect them back to the front end when saving data with improper formatting.
7. Control console supports file type uploads (currently done through direct path association).
8. Add upstream output parameter record to historical data files for easier issue localization and troubleshooting.
9. Robustness when body request data is raw in the control console.
10. Improve the usability of writing content in the management domain data use case files.

#### Performance Optimization
1. Optimize SQL statements for faster report loading when dealing with large amounts of data.
2. Design parallel data generation to improve performance and control concurrency.

#### Other Features
1. Rich file editing mode for scenario list and scenario editing (low priority)
2. Evaluate and potentially remove or migrate the parameter definition module functionality.
3. Compatibility with OPTIONS, HEAD, PATCH, and TRACE request types in the Swagger documentation.
4. Check the environment's authentication status during testing.
5. Support pausing for one-time execution tasks, allowing creation of tasks through execution and management through tasks.
6. Implement error fine management for task execution scenarios and data: if it's a connection-level error, subsequent executions are skipped; if it's a scenario execution error, execution continues.
7. Consider recording process variables for better issue localization; log by level to control output levels (low priority).
8. Generate comprehensive product testing reports in formats such as Excel, PDF, HTML, etc.
9. Add Swagger visualization to view interface information.
10. The application should not be a dropdown for interface definition lists when adding new ones.

#### Code Optimization
1. Switch to a pipelined communication mechanism for concurrency improvement.
2. Implement variable substitution in YAML documents using a template engine similar to Jinja's template substitution approach.

#### Operation and Maintenance Features:
1. Package all files for deployment as a whole instead of individual files; perform file-level changes separately (low priority).
2. View historical versions to identify who modified the data (low priority).
3. Regularly clean up historical records to conserve environment resources (low priority).
4. Regularly backup management domain data in MySQL.