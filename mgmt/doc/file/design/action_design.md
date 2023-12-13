##### 动作类型：
- sleep
- create_csv
- create_xls

##### 使用示例：
action:
- type: create_csv
  value: all_type_data_10w:100000    // 冒号前是文件名称，冒号后是数据数量，数量不写默认生成100条
- type: create_xls
  value: all_type_data_10w:100000    // 冒号前是文件名称，冒号后是数据数量，数量不写默认生成100条
- type: sleep
  value: 10   // 单位秒

##### 后续
- 可以根据需要进行更多动作的开发