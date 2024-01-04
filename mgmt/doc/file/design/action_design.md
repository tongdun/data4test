##### 动作类型：
- sleep
- create_csv
- create_xls / create_excel / create_xlsx  动作相同，生成excel形式的批量数据
- create_hive_table_sql  生成Hive创建表单，字段数据类型待完善
- record_csv    记录当前的请求数据为CSV文件
- record_xls / record_excel / record_xlsx   动作相同，记录当前的请求数据为EXCEL文件
- modify_file   根据模板文件内容生成新的带请求数据的文件

##### 使用示例：
```action:
- type: sleep
  value: 1    // 表示等待1秒种，时间可根据需要自动设置，单位为秒
- type: create_csv
  value: all_type_data_10w:100000    // 冒号前是文件名称，冒号后是数据数量，数量不写默认生成100条
- type: create_xls
  value: all_type_data_10w:100000    // 冒号前是文件名称，冒号后是数据数量，数量不写默认生成100条
- type: create_hive_table_sql
  value: name.sql   // 生成表单的存储文件名
- type: record_csv
  value: name.csv   // 把当前请求的body数据记录到name.csv中
- type: record_xls
  value: name.xls   // 把当前请求的body数据记录到name.xls中
- type: modify_file
  value: name.xml:name_{certid}.xml  // 冒号前为模板文件，需要替换的字段用占位符，冒号后为替换数据后保存的文件，{certid}为取请求数据中certid变量的值，区分生成的数据和记录
- type: modify_file
  value: name.xml:{phoneno}.xml  // 模板文件名称:生成文件名称；生成文件名用的占位符取值最好是唯一的，否则数据会发生覆盖
```

#### modify_file动作模板文件示例
- [JSON类型](../../../upload/record_template.json)
- [XML类型](../../../upload/record_template.xml)
- [TXT类型](../../../upload/record_template.txt)
- [YAML类型](../../../upload/record_template.yml)

##### 后续
- 可以根据需要进行更多动作的开发