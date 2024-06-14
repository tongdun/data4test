#### 参数定义
- 内置变量：
    - uniVar: 需要唯一的变量，例如：名称，标识, ……
    - noChkVar: GET方法请求后，某些字段不需要检查，例如:数据库递增的id号
- 自定义变量支持多种形式：
    - name = "string"
    - name = ["string1", "string2", ……] 不要定义太多值，组装数据做全排列，赋值太多，请求数据会呈指数增长， 是否要做全排列？？
    - /path = {"name":["string1", "string2", ……]}  同一变量定义多种形式，该方式优先，httpMethod_/path or /path ？？
    - 可以由其他API可以动态获取到的值，使用依赖关系

#### 随机数据定义：
  - string类型的变量，模糊测试模式自动遍历 ["", " ", "8长度", "256长度", "……"] 
  - integer类型的变量，模糊测试模式遍历 [-1, 65536, ……] 
  - array类型的变量，模糊测试模式遍历[[], {}, '']
  - bool类型的变量，模糊测试模式遍历["true", "false"], 赋值时，接口是字符串形式还是true or false类型得看接口……
  - ……

#### 接口关系：
- 前置用例：
  - pre_apis = ["get_/path1", "post_/path2", ……]: 执行该API时，先执行此处关联的API
- 提供变量：
  - out_vars = {"name": "content-record-id", "id": "content-record-id"} :其他API依赖的数据由此处定义，返回数据全是字典的用"-"分割, 后期可以页面进行选择，直接生成这个格式的数据
  - out_vars = {"name": "content-record*id", "ids": "content-record*id"} :其他API依赖的数据由此处定义，返回数据中有列表的用"*"分割， 后期可以页面进行选择，直接生成这个格式的数据
- 名称转义：
  - check_vars = {"name_id": "content-id"} : 如果post请求的数据,get查询时名称一样的，无需定义，名字不一样的需要转义

