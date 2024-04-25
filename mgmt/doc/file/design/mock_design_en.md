##### Supported Mock Routes:
- http://{IP}:{PORT}/mock/file/name.xml?lang=en
    - Returns the content of the uploaded file name.xml. If there are placeholders in name.xml, data with corresponding features will be generated. Defaults to Chinese if lang is not specified.
    - Supported file types for return: .xml, .json, .yml, .txt.
    - File template examples:
        - [JSON Type](../../../upload/create_template.json)
        - [XML Type](../../../upload/create_template.xml)
        - [TXT Type](../../../upload/create_template.txt)
        - [YAML Type](../../../upload/create_template.yml)
- http://{IP}:{PORT}/mock/data/quick
    - Returns data in JSON format as per the example information below.
    - Might be deprecated in future versions and merged with other functionalities.
- http://{IP}:{PORT}/mock/data/slow?sleep=N
    - N is an integer representing the wait time. Simulates a delay in response similar to a real third-party request. Returns data in JSON format.
    - Might be deprecated in future versions and merged with other functionalities.

##### Example Information (XX represents random data):
```
{
"addr": "XX City, XX Road, No. XX, XX Community, Unit XX, Room XX",
"age": "XX",
"bankCard": "XXXXX",
"bool": "XX",
"grade": "X",
"id_no": "XXXXX",
"mobile": "XX",
"name": "XXXXX",
"sex": "XXX",
"weight": "XX",
"yesOrNo": "XX"
}
```
- http://{IP}:{PORT}/mock/data/certid/:idno
    - Provide a Mainland China ID card number, and basic information will be returned.
    - Generated information is all test data and has no real significance for reference:
```
{
"city": "XX",
"sex": "XX",
"code": "XXXXX",
"birthday": "XX",
"district": "X",
"province": "XXXXX",
"address": "XX City, XX Road, No. XX, XX Community, Unit XX, Room XX",
"country": "China"
}
```
- http://{IP}:{PORT}/mock/systemParameter/:name?lang=en
    - Provide the name of the system parameter, and a randomly selected value from the defined value list will be returned. Supports multiple languages.