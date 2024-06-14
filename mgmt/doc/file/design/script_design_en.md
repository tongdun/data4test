### 一、Public Principles
#### 1.1 Defining Execution Engines
#####  First Method: Define common execution engines in system parameters
1. Parameter name
- `scriptRunEngine`:
```json
{
    ".py": "/usr/local/bin/python3",
    ".sh": "/bin/sh",
    ".jmx": "xxxx",
    ".bat": "xxxx"
}
```
Note: Specific functionality for configuring different parameters during JMeter script execution is pending integration.

##### Second Method: Define the execution engine in the script file itself
For example:
1. For a shell script, the first line might be: `#!/bin/bash` or `#!/bin/sh` depending on the chosen execution engine.
2. For a Python script, the first line might be: `#!/usr/bin/env python`

####  1.2 Priority of Script Execution Engines
1. System parameters have higher priority than script definitions.
2. Execution engines can be extended with any file extension as the key.
3. The environment for the execution engine needs to be configured by the user; otherwise, execution will result in an error.

#### 1.3 Adding Script Maintenance Information

Starting from the second line of the script, relevant information for changes should be added to facilitate long-term script maintenance.   
Example: file name, author, author's contact information, creation time, script description, usage instructions, change history, etc.

#### 1.4 Support for Placeholder Replacement
All scripts support placeholder variable replacement, similar to the usage of standard data.  
Example: `"{host}"` can be used directly for environment information, dedicated parameters, and common parameters.

#### 1.5 Include Help Information in Scripts
Provide help information for users to view and guide usage.

### 二、Example Scripts
#### 2.1 Shell Scripts
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

function testParameter()
{
    echo $0;
    echo $1
    echo $2
    echo "test"
    echo $HOSTIP
    echo $TableName
}

function test()
{
     cmd1
     cmd2
     return 整数
}

# ============================= MAIN ============================================
test
testParameter "paramter1" "parameter2"
```

#### 2.2 Python Script
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