### 接口管理
#### 接口变更及追踪
1. 新增接口定义进行规范检查，进行变更检查
   - 规范检查（会根据接口规范持续增加新的检查点）
     - 新增接口是否返回唯一标识
     - 查询数据信息是否返回data内容
     - 检查参数位置是否错误， e.g.: GET 设计了Body参数，POST设计了Query参数等均会失败
   - 变更检查
     - 应用层：新增接口数，被删除接口数， 被修改接口数， 保持原样接口数
     - 接口层：新增字段参数， 被删除字段参数，字段参数类型变更
2. 通过接口控制台或单独的管理页面，进行接口变更均会进行规范检查和变更检查
3. 规范检查结果历史记录支持查询
4. 变更检查历史记录支持查询


#### 接口数据复用
- 选择接口，自动生成测试数据，YAML格式/JSON格式的标准数据文件