说明：源数据定义(source)为raw时，返回数据当做一个整体做数据校验

##### 全部返回断言类型
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
- '=': 等于
- '!=': 不等于
- in: 包含
- '!in':不包含
- not_in: 不包含
- re:  正则匹配
- regex:  正则匹配
- regexp:  正则匹配

##### 字段值断言类型
- output: {type: output, source: data-contents*uuid, value: uuid}, 定义输出变量，提供给其他 API依赖使用
- equal: {type: equal, source: data-total, value: 0}
- not_equal: {type: not_equal, source: code, value: 200}
- contain: {type: contain, source: message, value: 已存在}
- not_contain:
- null: 为空
- not_null: 不为空
- in: 包含
- not_in: 不包含
- empty: 为空
- not_empty: 不为空
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

##### 特别说明
- output类型时，可以针对数组下标取值，e.g,:
  ```
  {type: output, source: data-contents*uuid[-1], value: codeUuid}, 取uuid数组最后一个赋值给codeUuid
  {type: output, source: data-contents*uuid[0], value: codeUuid}, 取uuid数组第一个值赋值给codeUuid
  {type: output, source: data-contents*uuid[1], value: codeUuid}, 取uuid数组第二个值赋值给codeUuid
  ```
- output类型时，可以针对数组一次全部提（返回的信息中uuid是一个数组）,e.g.:  
  ```{type: output, source: data-contents**uuid, value: codeUuid}, **表示把数据uuid当做一个整体数据赋值给codeUuid```

##### 断言值模板，支持多语种定义
- JSON格式：e.g.:  
```{"default": "v1,v2,v3,v4,……", "ch": "v1,v2,v3,v4,……", "en": "v1,v2,v3,v4,……"}```
- 普通格式定义：e.g.:  
```v1,v2,v3,v4,……```
- 当断言值为模板时，用占位符，执行时会自动获取，如果有设置多请种，会根据语言获取对应的值，e.g.:  
```{type: re, source: data-message, value: {successTemplate}}```
- 在"环境-断言值模板"列表增加名为successTemplate的模板信息：  
```{"ch": "成功|重复|已存在|已经存在", "en": "success|Success|exist|duplicate"}```
