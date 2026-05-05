##### Action Types:
- sleep
- create_csv
- create_xls / create_excel / create_xlsx (Same action, generates batch data in Excel format)
- create_hive_table_sql (Generates Hive create table statements, field data types need refinement)
- record_csv (Records the current request data as a CSV file)
- record_xls / record_excel / record_xlsx (Same action, records the current request data as an Excel file)
- modify_file (Generates a new file with request data based on the template file)

##### Usage Examples:
```action:
- type: sleep
  value: 1    // Represents a 1-second wait, time can be automatically set according to needs, unit is in seconds
- type: create_csv
  value: all_type_data_10w:100000    // Before the colon is the file name, after the colon is the data quantity. If not specified, defaults to generating 100 records.
- type: create_xls
  value: all_type_data_10w:100000    // Before the colon is the file name, after the colon is the data quantity. If not specified, defaults to generating 100 records.
- type: create_hive_table_sql
  value: name.sql   // Generates a storage file name for the table form
- type: record_csv
  value: name.csv   // Records the current request body data to name.csv
- type: record_xls
  value: name.xls   // Records the current request body data to name.xls
- type: modify_file
  value: record_template.json:record_template_{phoneno}.json   // Before the colon is the template file, placeholders for the fields to be replaced. After the colon is the file to save after replacing data. {certid} takes the value of the certid variable in the request data.
- type: modify_file
  value: record_template.yml:record_template_{certid}.xml  // Template file name: Generated file name; It is best to make the placeholders used in the generated file name unique, otherwise data may be overwritten.
- type: change_output
  value: parameter_name:old:new:num  // Output parameter names:String to be replaced:Replacement string:Replacement scope, If the replacement scope is empty or -1, it indicates a global replacement; If the replacement scope is a positive integer, it means replacing the first N matched strings; TODO: The "old" parameter can be a regular expression, and the parts to be replaced should be matched using parentheses.
- type: loop     // todo
  value: 5:10    // Before the colon is the interval time, after the colon is the number of cycles. If the number of cycles is not specified, it defaults to 3 times. Request the interface every 5 seconds. Exit the loop if the assertion passes, for a total of 10 cycles. If the assertion fails on the last request, the entire process is considered a failure.
```

##### Example of modify_file Action Template Files
- [JSON type](../../../upload/record_template.json)
- [XML type](../../../upload/record_template.xml)
- [TXT type](../../../upload/record_template.txt)
- [YAML type](../../../upload/record_template.yml)

##### Next Steps
- Further actions can be developed based on needs.