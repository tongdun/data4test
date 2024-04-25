Explanation:
When the source data definition is set to "raw", the returned data is treated as a whole for data verification.

##### All Return Assertion Types:
- equal: Strings are equal
- not_equal: Strings are not equal
- contain: Contains a specific string
- not_contain: Does not contain a specific string
- null: Is null
- not_null: Is not null
- in: Contains
- not_in: Does not contain
- empty: Is empty
- not_empty: Is not empty
- '=': Equals
- '!=': Does not equal
- re: Regular expression match
- regex: Regular expression match
- regexp: Regular expression match
- output_re: Regular expression match to extract values enclosed in (). If multiple values match, all will be extracted.
    e.g.1: {type: output_re, source: '\\"taskId\\":\\"([0-9]+)\\"', value: taskId}
    Define an output variable, assign the value matched within (.+) to taskId for use by other interface dependencies.
    e.g.2: {type: output_re, source: '\\"taskId\":\\"([a-zA-Z0-9]+)\\"', value: taskId}
    Define an output variable, assign the value matched within ([a-zA-Z0-9]+) to taskId for use by other interface dependencies.

##### Field Value Assertion Types:
- output: {type: output, source: data-contents*uuid, value: uuid}, Define an output variable for use by other interface dependencies.
- equal: {type: equal, source: data-total, value: 0}
- not_equal: {type: not_equal, source: code, value: 200}
- contain: {type: contain, source: message, value: Already exists}
- not_contain:
- null: Is null
- not_null: Is not null
- in: Contains
- not_in: Does not contain
- empty: Is empty
- not_empty: Is not empty
- '=': Equals
- '!=': Not equal to
- in: Contains
- '!in': Does not contain
- not_in: Does not contain
- re: Regular expression match
- regex: Regular expression match
- regexp: Regular expression match
- '>': Greater than
- '<': Less than
- '>=': Greater than or equal to
- '<=': Less than or equal to

##### Special Notes:
For the output type, values can be extracted based on array indices, e.g.,:
```
{type: output, source: data-contents*uuid[-1], value: codeUuid}, Extract the last uuid from the uuid array and assign it to codeUuid.
For the output type, all values can be extracted at once from an array in the returned information, e.g.:
{type: output, source: data-contents**uuid, value: codeUuid}, ** means treat uuid as a whole data and assign it to codeUuid.
```

##### Assertion Value Templates, Support for Multilingual Definitions:
- JSON format: e.g.:
```{"default": "v1,v2,v3,v4,...", "ch": "v1,v2,v3,v4,...", "en": "v1,v2,v3,v4,..."}```
- Standard format definition: e.g.:
```v1,v2,v3,v4,...```
- When the assertion value is a template, use placeholders. During execution, it will automatically retrieve the value. If multiple languages are set, it will retrieve the corresponding value based on the language, e.g.:
```{type: re, source: data-message, value: {successTemplate}```
- Add a template named "successTemplate" to the "Environment-Assertion Value Templates" list:
  ```{"ch": "成功|重复|已存在|已经存在", "en": "success|Success|exist|duplicate"}```