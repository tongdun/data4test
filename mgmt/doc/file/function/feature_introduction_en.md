#### Basic Features
- Free combination of HTTP interfaces: Interfaces can be freely combined across environments and applications to form scenarios for testing long-chain links.
- Definition of data repetition calls: The same data file can be repeatedly called by different scenarios and by the same scenario.
- Data file definition in declarative style: After editing test data through the interface, the data is saved as a declarative YAML file, supporting rapid batch writing of test data.
- Automatic generation of request data: Through simple definitions, the system can automatically generate rich and diverse feature test data, reducing manual intervention and construction of data.
- Precise parsing of return data: Through the definition of return data structures, rich assertions and data utilization can be performed, and provided to upstream and downstream users.
- Rich action types: Supports generating or recording as CSV, EXCEL, and supports generating test data from template files, aligning multiple data features.
- Rich assertion types: Through the definition of diverse assertion definitions, multiple aspects of returned data can be validated, with anti-stupid mechanisms.
- Idempotent execution of data files: Through the assertion of data definitions, the same data file can be repeatedly executed in the same environment to achieve idempotence.
- Variable context unawareness within scenarios: Different data files are combined into scenarios, and variables in requests and variables parsed from responses are used without awareness.

#### Extended Features
- Language judgment based on request headers, system automatically generates corresponding language test data, system parameters and assertion value templates can be customized for multi-language test data.
- Supports mock interfaces, returns file data, and supports the generation of placeholder data.
- Supports interface definition standard checks to quickly sense changes in interface information.