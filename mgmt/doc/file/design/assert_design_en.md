### Supported Assertion Categories

#### Field Value Assertions
- **Purpose**: Validate specific fields within a JSON response.
- **Data Source Definition (source)**: Defines the relationship to the specific field being asserted.

#### Overall Value Assertions
- **Purpose**: Validate the entire response body or raw response.
- **Data Source Definition (source)**: Can be `raw` or `ResponseBody`.

#### Performance Value Assertions (To Be Implemented)
- **Purpose**: Measure the time taken from request to response.
- **Data Source Definition (source)**: `RT` or `ResponseTime`.

#### File Value Assertions
- **Purpose**: Validate the content of files downloaded from APIs.
- **Data Source Definition (source)**: `FileType:line:column:split`
    - Supported File Types:
        - `File:TXT:line:column:split` - Split defaults to comma `,` if not defined.
        - `File:CSV:line:column:split` - Split defaults to comma `,` if not defined.
        - `File:EXCEL:line:column` - Line is an integer, column can be a number or name.
        - `File:JSON:data-total[0]` - Uses the same rules as field value extraction.
        - `File:YML:data-total` - Uses the same rules as field value extraction.
        - `File:XML` - Not detailed, to be implemented as needed.
        - `File:Other` - Supports regex extraction patterns, e.g., `'\"taskId\":\"(.+)\"'` (to be implemented).
        - Advanced JSON and YML patterns for extracting specific values, e.g., `'\"taskId\":\"([a-zA-Z0-9]+)\"'` (to be implemented).

#### Supported Assertion Types

##### Overall Response Assertions
- `equal`: Checks if the strings are equal.
- `not_equal`: Checks if the strings are not equal.
- `contain`: Checks if the string contains a specified substring.
- `not_contain`: Checks if the string does not contain a specified substring.
- `null`: Checks if the value is null.
- `not_null`: Checks if the value is not null.
- `in`: Checks if the value is in a specified list.
- `not_in`/`!in`: Checks if the value is not in a specified list.
- `empty`: Checks if the value is empty.
- `not_empty`: Checks if the value is not empty.
- `=`, `!=`, `re`, `regex`, `regexp`: Aliases for equal, not_equal, and regex matching.
- `output_re`: Extracts and validates values matched by a regex, supporting multiple matches.

##### Field Value Assertions
- Same as overall response assertions, plus:
    - `output`: Defines an output variable for use in subsequent API calls.
    - `>`, `<`, `>=`, `<=`: Compares numeric values.

#### Other Features

##### Output Type Explanation
- The `output` type allows extracting values from arrays by index or as a whole.
    - Examples:
        - `{type: output, source: data-contents*uuid[-1], value: codeUuid}` - Extracts the last item in the `uuid` array.
        - `{type: output, source: data-contents*uuid[0], value: codeUuid}` - Extracts the first item.
        - `{type: output, source: data-contents**uuid, value: codeUuid}` - Extracts the entire `uuid` array as a single value.

##### Assertion Value Templates
- Supports multi-language definitions for assertion values.
- JSON Format Example:
  ```json
  {"default": "v1,v2,v3,v4,…", "ch": "v1,v2,v3,v4,…", "en": "v1,v2,v3,v4,…"}
  ```
- Plain Text Format Example:
  ```
  v1,v2,v3,v4,…
  ```
- Placeholder support for dynamic values based on templates and language settings.
- Example usage with template:
  ```json
  {type: re, source: data-message, value: {successTemplate}}
  ```
  Where `successTemplate` is defined in the "Environment-Assertion Value Templates" list as:
  ```json
  {"ch": "成功|重复|已存在|已经存在", "en": "success|Success|exist|duplicate"}
  ```