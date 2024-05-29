### 一、公共原则
#### 1.1 定义执行引擎
##### 第一种：在系统参数中定义公共的执行引擎
1、参数名-scriptRunEngine:   
```{".py": "/usr/local/bin/python3", ".sh": "/bin/sh", ".jmx": "xxxx", ".bat": "xxxx"}```  
说明：jmeter脚本执行时需要配置不同的参数，具体功能待合入

##### 第二种：在脚本文件中定义执行引擎
例如：  
1、shell脚本的第一行示例：#!/bin/bash  #!/bin/sh  根据编写的执行引擎，如实定义即可  
2、python脚本的第一条示例：#!/usr/bin/env python  

#### 1.2 脚本执行引擎优先级
1、系统参数 > 脚本定义  
2、执行引擎以文件后缀为key可任意扩展  
3、执行引擎需自行配置环境，若无，执行会报错  

#### 1.3 需添加脚本维护信息
1、脚本从第二行开始，需要添加变更的相关信息，方便脚本的长期维护  
   示例：文件名称，编写者，编写者联系方式，创建时间，脚本描述用例，使用方法，变更历史等信息

#### 1.4 支持占位符变更替换
1、所有脚本支持占位符变量替换，同标准数据使用使用方法一致    
   示例："{host}"，环境信息，专用参数，公共参数，均可直接使用

#### 1.5 编写的脚本尽量有help信息
1、提供帮助信息给使用者查看，指导使用

### 二、示例脚本
#### 2.1 shell脚本
```
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
echo "test"
echo $HOSTIP
echo $TableName

function testParameter()
{
    echo $0;
    echo $1
    echo $2
}

function test()
{
    cmd1；
    cmd2;
    return 整数;
}

# ============================= MAIN ============================================
test
testParameter "paramter1" "parameter2"
```

#### 2.2 python脚本
```
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
    do somthing
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