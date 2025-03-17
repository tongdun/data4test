### 支持的断言大类
#### 字段值断言
- 从JSON中取特定的字段进行断言
- 数据源定义(source)为具体字段的关联关系

#### 整体值断言
- 从返回信息Response当作整体断言
- 数据源定义(source)为raw或ResponseBody

#### 性能值断言(待实现)
- 从请求到得到Response的返回时间
- 数据源定义(source)为RT/ResponseTime

#### 文件值断言
- 从接口获取下载的文件，对文件的内容进行断言
- 数据源定义(source)为FileType:line:column:split
- 具体FileType支持如下
  - File:TXT:line:column:split  //split若未定义，默认为英文逗号,
  - File:CSV:line:column:split  //split若未定义，默认为英文逗号,
  - File:EXCEL:line:column     // line为行号，值需为整数，column为列信息，可以为列号，也可以为列名
  - File:JSON:data.total[0]    // 与字段值取值使用规则一致
  - File:YML:data.total        // 与字段值取值使用规则一致
  - File:XML:未细化             // 待有需要再实现
  -  ```File:Other:'\\"taskId\\":\\"(.+)\\"' ```   //待实现，提取(.+)中匹配到的值断言
  -  ```File:JSON:'\\"taskId\":\\"([a-zA-Z0-9]+)\\"' ```   //待实现，提取([a-zA-Z0-9]+)中匹配到的值断言

### 断言支持的类型
##### 整体返回支持的断言类型
- '=': 等于
- '!=': 不等于
- in: 包含
- '!in':不包含
- not_in: 不包含
- re:  正则匹配
- regex:  正则匹配
- regexp:  正则匹配
- equal: 字符串相等
- not_equal: 字符串不相等
- contain: 包含某个字符串  
- not_contain: 不包含某个字符串
- null: 为空
- not_null: 不为空
- in: 包含
- not_in: 不包含
- empty: 为空
- not_empty: 不为空
- output_re:   正则匹配 ()匹配到的值取出来，如果匹配到多个值，均会进行提取
   - e.g.1:  ```{type: output_re, source: '\\"taskId\\":\\"(.+)\\"', value: taskId}```
     - 定义输出变量，(.+)中匹配到的值赋值给taskId, 提供给其他接口依赖使用
   - e.g.2:  ```{type: output_re, source: '\\"taskId\":\\"([a-zA-Z0-9]+)\\"', value: taskId}```
     - 定义输出变量, ([a-zA-Z0-9]+)中匹配到的值赋值给taskId, 提供给其他接口依赖使用

##### 字段值支持的断言类型
- '=': 等于
- '!=': 不等于
- in: 包含
- '!in':不包含
- not_in: 不包含
- re:  正则匹配
- regex:  正则匹配
- regexp:  正则匹配
- '>': 大于
- '<': 小于
- '>=': 大于等于
- '<=': 小于等于
- equal: {type: equal, source: data.total, value: 0},  .表示取map中的值，[int]表示取数组中的值，与写代码的语法一致，int可以为负数，表示从后往前取，-1表示取最后一个值
- not_equal: {type: not_equal, source: code, value: 200}
- contain: {type: contain, source: message, value: 已存在}
- not_contain:
- null: 为空
- not_null: 不为空
- in: 包含
- not_in: 不包含
- empty: 为空
- not_empty: 不为空
- output: {type: output, source: data.contents[0].uuid, value: uuid}, 定义输出变量，提供给其他接口依赖使用

### 其他特性说明
##### source提取目标字段说明
- 可针对数组下索引取值，从前往后或从后往前，e.g,:
  - {type: output, source: data.contents[-1].uuid, value: codeUuid}, 取contents数组最后一个值的uuid赋值给codeUuid  
  - {type: output, source: data.contents[0].uuid, value: codeUuid}, 取contents数组第一个值的uuid赋值给codeUuid  
  - {type: output, source: data.contents[1].uuid, value: codeUuid}, 取contents数组第二个值的uuid赋值给codeUuid
  - {type: output, source: contents[:].uuid, value: codeUuid, 取contents全部数组的uuid赋值给codeUuid。不支持多层嵌套的取，如contents[:].uuidList[:].uuid
  - {type: output, source: data.contents.uuid, value: codeUuid}, 取contents数组第一个值的uuid赋值给codeUuid，若contents为数组，但未定义索引，会默认取第一个值
  - {type: output, source: data.contents.uuid, value: codeUuid}, 若uuid是一个数组，会当做一个整体数据赋值给codeUuid
  - {type: =, source: contents[0].uuid, value: "{codeUuid}", 取contents数组第一个uuid的值与变量codeUuid比较
  - {type: re, source: data.contents[-1].name, value: XXX}, 取contents数组第一个值的name,与XXX正常行正则匹配

- 可针对数组下字典属性取值，e.g,:
  - {type: output, source: data.contents[@name=XXX&&@status=1].uuid, value: codeUuid}, 取contents数组下，name属性值为XXX且status属性值为1，取对应的uuid，赋值给codeUuid, 取到第一个值后即不再判断后续的对象
  - {type: output, source: data.contents[@name=XXX||@status=1].uuid, value: codeUuid}, 取contents数组下，name属性值为XXX或status属性值为1，取对应的uuid，赋值给codeUuid,第一个条件满足后即不再判断后续的条件
  - {type: =, source: data.contents[@name=XXX&&@status=1].uuid, value: "{codeUuid}"}, 取contents数组下，name属性值为XXX且status属性值为1，取对应的uuid，跟codeUuid变量值进行比较, 取到第一个值后即不再判断后续的对象
  - {type: re, source: data.contents[@name=XXX||@status=1].remark, value: "XXX}, 取contents数组下，name属性值为XXX或status属性值为1，取对应的remark，与XXX进行正则匹配,第一个条件满足后即不再判断后续的条件

##### value比较字段说明
- 支持断言模板，支持多语种比较
  - {type: re, source: data.message, value: {successTemplate},  successTemplate定义为多语种，运行时会根据语种进行比较
- 支持四则运算+、-、*、/
  - {type: =, source: data.total, value: {lastTotal} + 1, 获取当前接口的total值, 与lastTotal变量值+1进行相等比较
- 支持常用数据统计函数: Max、Min, Sum, Avg, Floor, Pow, Abs, Ceil, Round, Remainder, Exp, Log, Mean, Variance, Median
  - {type: =, source: data.value, value: Max({Value[0]}, {Value[1]}, {Value[3]}), 获取当前接口的value值，与计算出来的最大值进行相等比较
  - {type: >, source: data.value, value: Min({Value[0]}, {Value[1]}, {Value[3]}), 获取当前接口的value值，比计算出来的最小值进行小于比较

##### 断言值模板，支持多语种定义  
- JSON格式：e.g.:
  - {"default": "v1,v2,v3,v4,……", "ch": "v1,v2,v3,v4,……", "en": "v1,v2,v3,v4,……"}
  
- 普通格式定义：e.g.:ß
  - v1,v2,v3,v4,……
  
- 当断言值为模板时，用占位符，执行时会自动获取，如果有设置多语种，会根据语言获取对应的值，e.g.:
   - {type: re, source: data.message, value: {successTemplate}
  
- 在"环境-断言值模板"列表增加名为successTemplate的模板信息：
  - {"ch": "成功|重复|已存在|已经存在", "en": "success|Success|exist|duplicate"}