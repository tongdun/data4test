### 数据文件类型
- 标准数据文件
- Python脚本文件
- Shell脚本文件
- Python脚本文件
- Dos脚本文件
- Jmeter脚本文件
- 其他自定义脚本文件

### 数据文件适用场景
#### 标准数据文件适用场景
- 接口功能测试：多个功能点可一次维护
- 接口异常值测试：多条异常值可一次维护
- 接口并发测试：多条数据并发执行
- 接口边界值测试：多条边界值可一次维护
- 接口大数据量测试：数据文件中的单条或多条数据，或反复执行
- 业务数据沉淀和管理：多条数据可一次维护


#### 脚本文件适用场景
- 系统业务功能在接口层面上串不起来，需要特殊处理
- 加密的东西，涉及密钥等安全信息的，需要单独处理
- 日志类的敏感信息测试
- 部分业务逻辑在前端的，进行UI自动化
- 环境升级管理类的
- 统计报告文档生成类的
- 宗旨：标准数据文件实现不了的，均可通过脚本的形式实现自动化

### 标准数据文件使用
#### 标准数据文件组成
- 用例基本信息：名称，版本，关联api_id, 是否跑前置API，是否跑后置API, 是否并行
- 环境信息: 协议类型，环境IP, 统一路由前缀
- API信息: 描述，所属模块，所属应用，请求方法，请求路由，前置 API,后置 API，提供的参数关系
- 单值参数定义:由使用者定义，或自动生成，选择单条，生成单值模式
- 多值参数定义:由使用者定义，或自动生成，选择多条，生成多值模式，一个字段多值，一一对应关系
- 输出参数：由提供的参数关系定义解析出来
- 动作定义: 支持暂停，生成CSV, EXCEL海量数据，也可记录实时请求的数据为CSV, EXCEL，也可根据文件模板记录实时请求的数据等
- 断言关系：由使用者定义，支持为空，不为空，等于，不等于，包含， 不包含，关键字等等
- 测试结果：通过断言关系判断，符合要求则 PASS, 无断言定义则接口请求成功即为 PASS， 多值多结果
- 返回信息: 记录接口返回的信息，多值多返回

#### header设置
- header：鉴权数据， 请求类型等, application/json， form/data, 文件等
- respHeader：下载文件, 如果文件名为前端生成，对Content-Disposition进行设置，filename=XXX.filetype, 可根据实际情况进行设置,    
  Content-Type根据接口请求查看，按实际填写，请求后，response会自动置为XXX.filetype,然后使用文件类型进行断言判断
- respHeader:  
  Content-Disposition: attachment; filename=XX模块-XX管理-XX配置-XX导出.csv  
  Content-Type: application/csv
- 如果下载文件名为后端生成，直接断言进行设置即可，可以通过返回的文件名进行断言判断
    - {type: =, source: raw, value: XX模块-XX管理-XX配置-XX导出.csv}
    - {type: =, source: ResponseBody, value: XX模块-XX管理-XX配置-XX导出.csv}

#### 标准数据文件版本管理
- 根据历史记录自动累加版本，每次更新数据文件，备份一份至历史版本目录下

#### 标准数据文件子项详情
- [参数应用](./parameter_design.md)
- [动作应用](./action_design.md)
- [断言应用](./assert_design.md)

#### 标准数据文件鉴权优先级
- 单个文件 > 环境配置， 默认环境配置中的数据，后续扩展：开发成单独执行工具，命令行参数 > 全局配置文件 > 单个配置文件
- 单个文件中有定义的 header参数，环境配置中无，会汇合成新的header头，兼容文件类型的上传与下载


### 脚本数据文件使用
#### 脚本数据文件组成
- 执行引擎定义
- 执行内容

#### 脚本数据文件编写详情
- [脚本应用](./script_design.md)

### 数据文件参考示例
<details>
<summary>标准数据文件示例</summary>

```---
# 用例信息
name: 示例-用户管理-新建用户 # 数据用例名称，e.g.: 类型-模块-用例， 类型：功能/性能/异常/内置/……， 模块：用户管理/规则管理/……
api_id: post_/path        # 用例ID, method_path组合，后续做数据联动使用，数据统计使用
version: 1.0              # 数据用例版本，后续可以进行数据升级
is_run_pre_apis: "no"     # 是否跑前置用例，选项：yes / no,  默认 no， 功能未开发
is_run_post_apis: "no"    # 是否跑后置用例，选项：yes / no,  默认 no， 功能未开发
is_parallel: "no"         # 是否并行跑数据，选项：yes / no,  默认 no，
is_use_env_config: "yes"  # 是否使用公共环境，选项：yes / no,  默认 yes

# 环境信息
env:
  protocol: http        # http 或 https，请求协议
  host: X.X.X.X:8088    # 环境IP 或 环境域名 或 环境IP:端口
  prepath: /prefix      # 路由前缀，公共部分可以抽出来

# API 基本信息
api:
  description: 新建用户   # API用途
  module: 用户管理        # API所属模块
  app: appName           # API所属应用
  method: post           # （注意：保证正确） API请求方法
  path: /path            # （注意：保证正确）API请求路由，路由前缀抽离到prepath下时或公共环境中已定义prepath时，这里无需再写路由前缀
  pre_apis: []           # 调试时，依赖前置用例时，可以把关联前置文件写上，功能未充分验证
  param_apis: []         # 调试时，依赖其他用例的参数时，可以把关联文件写上，功能未充分验证
  post_apis: []          # 调试时，测试跑完后需要跑的用例，可以把关联文件写上，功能未充分验证

# 定义单值参数，如果is_use_env_config值为no, 需要定义此处的 header
single:
  header:
    Content-Type: multipart/form-data   # 如果api为导入文件功能，需要把Content-Type定义为multipart/form-data进行公用环境值的覆盖，优化级：数据文件>应用配置>产品配置
  respHeader:
    Content-Disposition: attachment; filename=XX模块-XX管理-XX配置-XX导出.csv  # 如果文件名为前端生成，对filename进行设置
    Content-Type: application/csv  # 根据接口请求查看，按实际填写，请求后，response会自动置为XXX.filetype,然后使用文件类型进行断言判断
  query: {}                             # GET请求时，请求参数定义，只定义一个值，共用的参数放在这里，无需反复定义
  path: {}                              # PATH 变量参数定义，只定义一个值
  body:
    condition: '{"children":[{"name":"definitionList","type":"string","value":"{nameList}"}]'  # {nameList} 代表字符串里有需替换的变更，nameList为 ouput 中输出的参数名字，在前置的用例中有定义同名变量，即会替换
    vaLue: '{FlowType}'   # 在'系统参数'菜单下，进行参数定义定义，支持多语种定
    name: '*{Name}*'      # 引用上文Name变量，当做一个整体，JSON格式
    name2: '**{Name}**'   # 引用上文Name变量，当做一个整体，字符串格式
    XXName: '{self}'      # 引用上文XXName变量的值，[{self}值变量将逐步废弃，不要再使用，已有的，尽快替换为具体的变量名
  bodyList:              # 当请求body直接是List时，相关请求参数放到bodyList下
    - name: '{Name}'
      sex: '{Sex}'
    - name: '{Rune(4)}'
      sex: '{Sex}'
# 定义多值参数
multi:
  query: {}                   # GET请求时，请求参数定义，定义的值为列表
  path: {}                    # PATH 变更参数定义，定义的值为列表
  body:
    description:              # 定义多值时，取各项定义的个数最少的数据，一一对应
    - '{Rune(128)}'    # 获取设置长度的汉字
    - '{Str(64)}'      # 获取设置长度的字符串
    - '{Int(10,100)}'  # 获取设置范围内的整数
    displayName:
    - '{Date(-2)}'      # 获取两天前的日期
    - '{Date(2)}'       # 获取两天后的日期
    - '{Timestamp(-2)}' # 获取两天前的时间戳
    name:
    - '{IDNo}'          # 获取身份证字符串
    - '{Name}'          # 获取名字字符串
    - '{Address}'       # 获取地址字符串
    - '{BankNo}'        # 获取银行卡号字符串

# 断言，数据校验，根据需要写不同类型的断言，不写断言，只要返回为200，即算 PASS
assert:
- type: equal   # 验证code的值等于200
  source: code    # 返回的json信息，取key为code的值
  value: 200 
- type: "!=" # 验证code的值不等于200
  source: code    # 返回的json信息，取key为code的值
  value: 200 
- type: ">="    # 验证source字段大于等于1
  source: data-total     # 返回的json信息，data字典-取出productDesc的值
  value: 1
- type:  contain
  source: data-contents*productDesc  # 返回的json信息，data字典-content字典-字典列表，取出productDesc的值, 并校验是否包含 value中的值
  value: 待删除的产品描述
- type: "!in"   # 验证取到的productName的值包含删除
  source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值, 不包含 value 中的值
  value: 删除
- type: not_contain   # 验证取到的productName的值不包含删除
  source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值
  value: 产品
- type: re
  source: message
  value: 成功|重复|已存在
- type: re
  source: message
  value: '{successTemplate}'  # 在'断言值模板'菜单下，进行断言值模板定义，支持多语种
- type: output  # 从返回的json 信息取取出 uuid 的值，并命名为uuid
  source: data-contents*uuid
  value: uuid
- type: output  # 从返回的json信息取出uuid的值，并重命名为ProductUuid
  source: data-contents*uuid
  value: ProductUuid
- type: output_re  # 从整体返回中进行正则匹配提取，并重命名为taskId，()中匹配到的值取出来，如果匹配到多个值，均会进行提取
  source: '\\"taskId\\":\\"(.+)\\"'
  value: taskId
- type: output_re  # 定义输出变量, ([a-zA-Z0-9]+)中匹配到的值赋值给taskId, 提供给其他接口依赖使用
  source: '\\"taskId\":\\"([a-zA-Z0-9]+)\\"'
  value: taskId
- type: output  # 返回值为文件时，从输出的文件中取第一行第一列的值，赋值给taskId
  source: File:TXT:1:1:,
  value: taskId
- type: output  # 返回值为文件时，从输出的文件中取第一行第一列的值，赋值给taskId
  source:  File:CSV:1:1:|
  value: taskId
- type: output  # 返回值为文件时，从输出的文件中取第一行第一列的值，赋值给taskId
  source:  File:CSV:1:1:|
  value: taskId
- type: output  # 返回值为文件时，从输出的文件中取第一行第一列的值，赋值给taskId
  source:  File:EXCEL:1:1
  value: taskId
- type: output  # 返回值为文件时，从输出的文件中取data字典下total的值赋值给XXXCount, 取值与标准文件的取值规则一致
  source: File:JSON:data-total    #
  value: XXXCount
- type: output  # 返回值为文件时，从输出的文件中取data字典下total的值赋值给XXXCount, 取值与标准文件的取值规则一致
  source: File:YML:data-total
  value: XXXCount

# 数据执行后的动作
action:
  - type: sleep
    value: 1    // 表示等待1秒种，时间可根据需要自动设置，单位为秒
  - type: create_csv
    value: name:number    // 生成文件名:生成的数据条数，默认生成10条
  - type: create_xls
    value: name:number    // 生成文件名:生成的数据条数, 默认生成10条
  - type: record_csv
    value: name.csv    // 记录实时请求的body数据,title为body请求的字段名，如果多个数据输出到一个记录文件中,自动追加
  - type: record_xls
    value: name.xls    // 记录实时请求的body数据,title为body请求的字段名，如果多个数据输出到一个记录文件中,自动追加
  - type: modify_file
    value: name.xml:name_{certid}.xml  // 冒号前为模板文件，文件的内容中需要替换的字段用占位符，冒号后为替换数据后保存的文件，{certid}为取请求数据中certid变量的值，区分生成的数据和记录
  - type: modify_file
    value: name.txt:{phoneno}.txt  // 模板文件名称:生成文件名称；生成文件名用的占位符取值最好是唯一的，否则数据会发生覆盖
  - type: modify_file
    value: name.json:{phoneno}.json  // 模板文件名称:生成文件名称；生成文件名用的占位符取值最好是唯一的，否则数据会发生覆盖
  - type: modify_file
    value: name.yaml:{phoneno}.yaml  // 模板文件名称:生成文件名称；生成文件名用的占位符取值最好是唯一的，否则数据会发生覆盖
  - type: modify_file
    value: name.yml:{phoneno}.yml  // 模板文件名称:生成文件名称；生成文件名用的占位符取值最好是唯一的，否则数据会发生覆盖

# 输出其他接口需要的依赖数据, 由断言中类型为 ouput 定义，自动生成, 定义为'{self}'或 '{uuid}' 从此处取值
output:
  uuid:
  - XXX
  - XXX

# 测试结果：pass, fail, untest, 自动生成，断言全部符合要求设为pass, 请求若返回非200，直接置为 fail, 如果执行次数测试为0，测置为 untest
# 保留最新测试结果
test_result:
- pass
- fail
- untest

# 请求 URL，自动生成， 保留最新测试结果
urls:
- http://X.X.X.X:8089/prefix/path

# 请求数据，body, query, 自动生成, 保留最新测试结果
requests:
  - '{"curPage":"1","endEntryTime":"1627095420000","pageSize":"10","searchOption":"{}"startEntryTime":"1626749820000","timeType":"1"}'

# 返回信息, 自动生成， 保留最新测试结果
response:
- "response1"
- "response2"
```
</details>

<details>
<summary>Python脚本文件示例</summary>

```---
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# =========================================================================
# FileName: XXX.py
# Creator: XXX
# Mail: XXX@qq.com
# Created Time: 20XX-0X-XX
# Description: Usage desc
# Usage:
# 1. XXX
# 2. XXX
# History:
# 202X-0X-0X/change log
# 202X-0X-0X//change log
#
# Copyright (c) 2024-20XX XXX Tech. All Right Reserved.
# =========================================================================
import argparse
import sys

def functon_something():
    return

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Check and Config OS Env")
    parser.add_argument('-d','--debug', dest="debug", action="store", default='N', help="default value is N")
    args = parser.parse_args()
 
    if args.debug.upper() == "Y":
        DEBUG = True
    else:
        DEBUG = False
    if args.target_host_ip.upper() == "N":
        parser.print_help()
        sys.exit(1)
    handler = functon_something()
 ```
</details>

<details>
<summary>Shell脚本文件示例</summary>

```---
#!/bin/bash
# =========================================================================
# FileName: XXX.sh
# Creator: XXX
# Mail: XXX@qq.com
# Created Time: 20XX-0X-XX
# Description: Usage desc
# Usage:
# 1. XXX
# 2. XXX
# History:
# 202X-0X-0X/change log 2
# 202X-0X-0X//change log 1
#
# Copyright (c) 2024-20XX XXX Tech. All Right Reserved.
# =========================================================================
HOSTIP="{host}"
TableName="{HiveName}"

function testParameter()
{
    echo $0
    echo $1
    echo $2
    echo "test"
    echo $HOSTIP
    echo $TableName
}

test()
{
   cmd1
   cmd2
   return 整数
}

# ============================= MAIN ============================================
test
testParameter "paramter1" "parameter2"
```
</details>