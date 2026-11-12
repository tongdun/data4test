#### Basic Features
- Free combination of HTTP interfaces: Interfaces can be freely combined across environments and applications to form scenarios for testing long-chain links.
- Definition of data repetition calls: The same data file can be repeatedly called by different scenarios and by the same scenario.
- Data file definition in declarative style: After editing test data through the interface, the data is saved as a declarative YAML file, supporting rapid batch writing of test data.
- Automatic generation of request data: Through simple definitions, the system can automatically generate rich and diverse feature test data, reducing manual intervention and construction of data.
- Precise parsing of return data: Through the definition of return data structures, rich assertions and data utilization can be performed, and provided to upstream and downstream users.
- Rich action types: Supports generating or recording as CSV, EXCEL, generating test data from template files, and rsa_encrypt encryption actions, aligning multiple data features.
- Rich assertion types: Through the definition of diverse assertion definitions (including output2env, etc.), multiple aspects of returned data can be validated, with anti-stupid mechanisms.
- Idempotent execution of data files: Through the assertion of data definitions, the same data file can be repeatedly executed in the same environment to achieve idempotence.
- Variable context unawareness within scenarios: Different data files are combined into scenarios, and variables in requests and variables parsed from responses are used without awareness.
- Support for calling LLMs: Generate test cases, test data, test scenarios, etc. by calling large language models.
- Support for CLI run mode: Runs offline with SQLite, invokes interfaces via command line to execute data/scenarios, and supports `--param` parameter injection.

#### High-Frequency Features
- 1. Task management: custom / one-time / daily / weekly.
- 2. Interface management: interface change detection and tracking, interface specification checks.
- 3. Standard data files for functional / exception / boundary value / internationalization / concurrency / high-volume testing of interfaces.
- 4. Third-party Mock data generation for invocation.
- 5. Scenario orchestration and testing.
- 6. Visualized console operation with a Postman-like interface.
- 7. Historical data and scenarios: re-run once more and continue scenario testing.
- 8. Rich assertion judgment types for result validation.
- 9. Rich built-in feature data auto-generation.
- 10. Script management and execution, e.g.: python, shell, jmeter, dos, etc., with extensible script execution engines.
- 11. Scenario types: serial-interrupt / serial-compare / normal-concurrency / normal-compare.
- 12. Standard data files support automated export file validation, implemented: CSV/EXCEL/YAML/JSON.
- 13. Scenarios and tasks support associating multiple environments for execution.
- 14. Generate polymorphic test data for upstream/downstream interfaces or functions, implemented: CSV/EXCEL/TXT/YAML/XML/JSON.
- 15. One-click view of non-automated interfaces.
- 16. Internationalization: multi-language test data auto-generation, customizable multi-language system parameters and assertion value templates; i18n sync supports AI translation.
- 17. Intelligent testing: test case generation, test data/scenario generation, execution analysis, and testing.
- 18. Test case statistics report: test-time trend chart (by day/week), plan-type data, extension-info dimension statistics, and click-through navigation.
- 19. Task export/import supports on-demand handling of environment information, assertions, system parameters, i18n corpus, etc.

#### Extended Features
- Language judgment based on request headers, system automatically generates corresponding language test data, system parameters and assertion value templates can be customized for multi-language test data.
- Supports mock interfaces, returns file data, and supports the generation of placeholder data.
- Supports interface definition standard checks to quickly sense changes in interface information.
