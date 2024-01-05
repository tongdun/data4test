### 常用渲染变量：
说明：支持自动生成数据：常用渲染变量（可参考Jmeter支持的变量）
     变量组成支持支持英文，数字，下划线，中划线，e.g.: {name_test-1}

#### 上下文同名变量关联：
- {self}: 需要其他用例提供的参数

#### 产品环境专用参数：（在对应产品环境下跑的场景，均可使用此处变量）
- 在产品配置里，进行一些专用参数定义， JSON格式，e.g.: {"name1": "value1", "name2": "value2"}

#### 支持多语种的系统内置参数, 目前支持中文，英文
- {Name}
- {Email}
- {Phone}
- {Mobile}
- {Contact}
- {Gender}
- {Sex}
- {CardNo}
- {BankNo}
- {SSN}
- {IDNo}
- {Address}
- {Company}
- {Country}
- {Province}
- {State}
- {City}

#### 范围数据自动生成：
- {Int(a,b)} 数字, e.g.: {Int(1000,2000)}，生成1000到2000之间的随机数
- {Date(a,b)} 日期，e.g.: {Date(-30,1)}，生成30天前到1天后的随机日期，2006-01-02
- {Time(a,b)} 日期，e.g.: {TIme(-30,1)}，生成30天前到1天后的随机日期，2006-01-02 10:59:05

#### 指定长度数据自动生成：
- {Str(length)} 字符串, e.g.: {Str(64)}，生成长度为64的字符串
- {Rune(length)} 汉字, e.g.: {Rune(128)}，生成长度为128的中文字符串
- {Date(int)}    日期，e.g.: {Date(-3)}，生成3天前的日期, 2022-01-02
- {Month(int)}    日期，e.g.: {Month(-3)}，生成3个月前的月份, 2022-01
- {Year(int)}    日期，e.g.: {Year(-3)}，生成3年前的年份, 2022
- {IntStr(12)}  生成长度为12的数字字符串
- {UpperStr(12)}  生成长度为12的全大写字母字符串
- {LowerStr(12)}  生成长度为12的全小写字母字符串
- {UpperLowerStr(12)}  生成长度为12的含大小写字母字符串
- {Str(12)}  生成长度为12的含大小写字母，数字字符串
- {TimeGMT(-2)}  GMT日期，e.g.: {TimeGMT(-3)}，生成3天前的GMT日期, Fri Aug 26 2022 23:59:59 GMT+0800
- {DayBegin(1)} 日期，e.g.: {DayBegin(-3)}，生成3天前的一天开始日期, 2006-01-02 00:00:00
- {DayEnd(-1)} 日期，e.g.: {DayEnd(3)}，生成3天后的一天结束日期, 2006-01-02 23:59:59
- {MonthBegin(-1)} 日期，e.g.: {MonthBegin(3)}，生成3个月后的一月开始日期, 2006-01-01 00:00:00
- {MonthEnd(-1)}   日期，e.g.: {MonthEnd(-3)}，生成3个月前的一月结束日期, 2006-01-31 23:59:59
- {YearBegin(-1)}  日期，e.g.: {YearBegin(-1)}，生成1个年前开始日期, 2006-01-01 00:00:00
- {YearEnd(-1)}    日期，e.g.: {YearEnd(1)}，生成1年后的结束日期, 2006-12-31 23:59:59
- {DayStampBegin(-180)}  时间戳，原理同上
- {DayStampEnd(-2)}
- {MonthStampBegin(-1)}
- {MonthStampEnd(-1)}
- {YearStampBegin(-1)}
- {YearStampEnd(-1)}


#### 特征数据自动生成：
- {Date}   生成当前日期  2006-01-02
- {Time}   生成当前日期,带时分秒 2006-01-02 10:59:05
- {DeviceType}
- {QQ}       生成QQ号
- {Age}      生成年龄
- {Diploma}   生成学历
- {IntStr10}  生成长度为10的数字字符串
- {Int}     生成整数
- {Timestamp}   生成时间戳
- {TimeFormat(2006-01-02 03:04:05.000 PM Mon Jan)}  // 2021-06-25 10:59:05.410 AM Fri Jun
- {TimeFormat(15:04 2006/01/02)}  // 10:59 2021/06/25
- {TimeFormat(2006-1-2 15:04:05.000)}  // 2021-6-25 10:59:05.410
- {TimeFormat(XXX)} // 格式根据需要组装
- {TimeGMT} // Fri Aug 26 2022 23:59:59 GMT+0800
- {DayBegin}   生成当天开始日期时间2006-01-02 00:00:00
- {DayEnd}     生成当天结束日期时间2006-01-02 23:59:59
- {MonthBegin}  生成当月开始日期时间2006-01-01 00:00:00
- {MonthEnd}    生成当月结束日期时间2006-01-31 23:59:59
- {YearBegin}    生成当年开始日期时间2006-01-01 00:00:00
- {YearEnd}      生成当年结束日期时间2006-12-31 23:59:59
- {DayStampBegin} 生成时间戳，原理同上
- {DayStampEnd}
- {MonthStampBegin}
- {MonthStampEnd}
- {YearStampBegin}
- {YearStampEnd}
- {Int(100000,999999)}{Year(-58)}{Int(10,13)}{Int(10,28)}{Int(1000,9999)}  组合数据，生成58岁的身份证号

#### 特征数据自由定义：（系统参数中进行定义，支持多语种定义）
- JSON格式：e.g.: ```{"default": "v1,v2,v3,v4,……", "ch": "v1,v2,v3,v4,……", "en": "v1,v2,v3,v4,……"}```
- 普通格式定义：e.g.: ```v1,v2,v3,v4,……```

#### 级联数据自由定义
说明：（系统参数中先进行级联数据定义，然后通过{TreeData_Name}关联取用，格式参考China变量，目前只支持3级，如需更多级，联系管理员进行支持）  
- {TreeData_China[1]}      生成级联数据，e.g: 生成国家-中国  
- {TreeData_China[2]}      生成级联数据，e.g: 生成省份-浙江省  
- {TreeData_China[3]}      生成级联数据，e.g: 生成城市-杭州市  
