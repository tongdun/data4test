### Test Case

- Scene test case name should not contain / or other special characters. The route request will directly use the name of the scene test case. / is a special character in the route, and it will cause parsing exceptions.
- The data test case file name should not contain &, /, () or other special characters. These characters are reserved for specific purposes in the Linux file system, and they will cause parsing exceptions.

### Scenario Test Case

### Common Issues in Writing Data Files

#### General Rules
- Due to the syntax of YAML files, when writing content with special characters, you need to use quotes to indicate that it is a string.

#### Specific Examples:
- If the assertion value is of boolean type, use quotes, e.g.:
  ```
  - source: success
    type: re
    value: "true"
  ```
- If the assertion value is a template, use quotes, e.g.:
  ```
  - source: success
    type: re
    value: "{sucessTemplate}"
  ```