### 数据生成和性能测试数据库设计

####  数据库

- 数据库名：data4test

####  表单设计

##### 环境配置

- 表名：env_config
- 主键：id
- 唯一键： product
- 表创建SQL:
```
create table if not exists env_config
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    product    varchar(64)                         not null comment '产品名称',
    ip         char(39)                            null comment '环境域名 / IP / IP:Port',
    protocol   enum('http','https') default 'http' comment '请求协议',
    prepath    varchar(255)                        null comment '路由前缀',
    threading  enum('yes','no') default 'no'       not null comment '是否并发',
    auth      longtext                            null comment '鉴权信息',
    testmode   enum('custom','fuzzing', 'all') default 'fuzzing'                           not null comment '测试模式',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间',
    constraint product
        unique (product)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='环境配置';
```



##### 接口定义

- 表名：api_definition
- 主键：id
- 唯一键： api_id, product
- 表创建SQL:
```
create table if not exists api_definition
(
    id             int auto_increment
        primary key comment '自增主键',
    api_id        varchar(255)                        not null comment '接口ID',
    api_module         varchar(255)                        null comment '所属模块',
    api_desc    varchar(255)                        null comment '接口描述',
    http_method     enum('get','post', 'put', 'delete') default 'get'                         null comment '请求方法',
    path           varchar(255)                        null comment '请求路径',
    header         longtext                            null comment 'Header参数',
    path_variable   longtext                            null comment 'Path参数',
    query_parameter longtext                            null comment 'Query参数',
    body           longtext                            null comment 'Body参数',
    response       longtext                            null comment 'Resp参数',
    product        varchar(64)                         null comment '所属产品',
    created_at     timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at     timestamp                           null comment '更新时间', 
    deleted_at     timestamp                           null comment '删除时间',
    constraint api_id_product
        unique (api_id, product)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='接口定义';
```



##### 接口关系

- 表名：api_relation
- 主键：id
- 唯一键：api_id ,  product
- 表创建SQL:
```
create table api_relation
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    api_id    varchar(255)                         not null comment '接口ID',
    api_desc  varchar(255)                         null comment '接口描述',
    api_module     varchar(255)                    null comment '所属模块',
    product     varchar(64)                         null comment '所属产品',
    auto  enum('yes','no') default 'yes'       not null comment '是否自动化',
    pre_apis    varchar(255)                       null comment '前置接口',
    out_vars    varchar(255)                        null comment '提供变量关系',
    check_vars  varchar(255)                      null comment '校验变量转换关系',
    param_apis  varchar(255)                  null comment '依赖参数关联接口',
    created_at timestamp default CURRENT_TIMESTAMP not null  comment '创建时间',
    updated_at timestamp                           null  comment '更新时间',
    deleted_at timestamp                           null  comment '删除时间',
    constraint api_id_product
        unique (api_id, product)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='接口关系';
```

##### 测试数据

- 表名：api_test_data
- 主键：id
- 唯一键：无
- 表创建SQL:
```
create table api_test_data
(
    id              int auto_increment
        primary key comment '自增主键',
    data_desc       varchar(255)                        not null comment '数据描述',
    api_desc     varchar(255)                           null comment '接口描述',
    api_id         varchar(255)                         not null comment '关联接口',
    api_module          varchar(255)                    null comment '所属模块',
    header  longtext                                    null comment 'Header',
    url_query        longtext                           null comment 'UrlQuery',
    body            longtext                            null comment 'Body',
    run_num         int       default 1                 null comment '执行次数',
    expected_result varchar(255)                        null comment '预期结果',
    actual_result   varchar(255)                        null comment '实际结果',
    result          varchar(255)                        null comment '测试结果',
    fail_reason     longtext                            null comment '失败原因',
    response        longtext                            null comment 'Response',
    product         varchar(64)                         null comment '关联产品',
    created_at      timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at      timestamp                           null comment '更新时间',
    deleted_at      timestamp                           null  comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试数据';
```

##### 模糊数据

- 表名：api_fuzzing_data
- 主键：id
- 唯一键：无
- 表创建SQL:

```
create table api_fuzzing_data
(
    id              int auto_increment
        primary key comment '自增主键',
    data_desc       varchar(255)                        not null comment '数据描述',
    api_desc     varchar(255)                           null comment '接口描述',
    api_id         varchar(255)                         not null comment '关联接口',
    api_module          varchar(255)                    null comment '所属模块',
    header  longtext                                    null comment 'Header',
    url_query        longtext                           null comment 'UrlQuery',
    body            longtext                            null comment 'Body',
    run_num         int       default 1                 null comment '执行次数',
    expected_result varchar(255)                        null comment '预期结果',
    actual_result   varchar(255)                        null comment '实际结果',
    result          varchar(255)                        null comment '测试结果',
    fail_reason     longtext                            null comment '失败原因',
    response        longtext                            null comment 'Response',
    product         varchar(64)                         null comment '关联产品',
    created_at      timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at      timestamp                           null comment '更新时间',
    deleted_at      timestamp                           null  comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模糊数据';
```



##### 结果详情

- 表名：api_test_detail
- 主键：id
- 唯一键：无
- 表创建SQL:
```
create table api_test_detail
(
    id          int auto_increment
        primary key comment '自增主键',
    api_id     varchar(255)                        not null comment '接口ID',
    api_desc varchar(255)                        not null comment '接口描述',
    data_desc varchar(255)                        not null comment '数据描述',
    header  longtext                                    null comment 'Header',
    url         varchar(255)                        not null comment 'URL',
    body        longtext                            null comment 'Body',
    response    longtext                            null comment 'Response',
    fail_reason longtext                            null comment '失败原因',
    test_result varchar(255)                        not null comment '测试结果',
    product     varchar(64)                         null comment '关联产品',
    created_at  timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at  timestamp                           null comment '更新时间',
    deleted_at  timestamp                           null comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试详情';
```


##### 变量提供

- 表名：api_test_result
- 主键：id
- 唯一键：无，取最新的值
- 表创建SQL:
```
create table api_test_result
(
    id          int auto_increment
        primary key comment '自增主键',
    api_id     varchar(255)                        not null comment '接口ID',
    out_vars     longtext                            null comment '提供变量',
    product     varchar(255)                        null comment '关联产品',
    created_at  timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at  timestamp                           null comment '更新时间',
    deleted_at  timestamp                           null comment '创建时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='变量提供';
```


##### 接口统计

- 表名：api_id_count
- 主键：id
- 唯一键：无
- 表创建SQL:
```
create table api_id_count
(
    id           int unsigned auto_increment
        primary key comment '自增主键',
    api_id      varchar(255)                        not null comment '接口ID',
    api_desc varchar(255)                        null comment '接口描述',
    run_times    int                                 not null comment '执行次数',
    test_times   int                                 not null comment '测试次数',
    pass_times   int                                 not null comment '通过次数',
    fail_times   int                                 not null comment '失败次数',
    untest_times int                                 not null comment '未测试次数',
    test_result  char(8)                             not null comment '测试结果',
    fail_reason  longtext                            null comment '失败原因',
    product      varchar(64)                         not null comment '关联产品',
    created_at   timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at   timestamp                           null comment '更新时间',
    deleted_at   timestamp                           null comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='接口统计';
```


##### 产品统计

- 表名：product_count
- 主键：id
- 唯一键：无
- 表创建SQL:

```
create table product_count
(
    id                  int unsigned auto_increment
        primary key comment '自增主键',
    all_count           int                                 null comment '接口总数', 
    automatable_count   int                                 null comment '可自动化数',
    unautomatable_count int                                 null comment '不可自动化数',
    auto_test_count     int                                 null comment '自动化测试总数',
    untest_count        int                                 null comment '未测试总数',
    pass_count          int                                 null comment '通过总数',
    fail_count          int                                 null comment '失败总数',
    auto_per            double                              null comment '自动化率',
    pass_per            double                              null comment '通过率',
    fail_per            double                              null comment '失败率',
    product             varchar(64)                         null comment '关联产品',
    created_at          timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at          timestamp                           null comment '更新时间',
    deleted_at          timestamp                           null  comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品统计';
```



##### 参数定义

- 表名：parameter_definition
- 主键：id
- 唯一键：api_id ,  product
- 表创建SQL:
```
create table parameter_definition
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    name       varchar(64)                         not null comment '参数名称/接口ID',
    value      longtext                            not null comment '值, e.g.: string, list, dict',
    product    varchar(64)                         null comment '关联产品',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间',
    constraint name_product
        unique (name, product)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参数定义';
```


##### 随机数据定义

- 表名：fuzzing_definition
- 主键：id
- 唯一键：无
- 表创建SQL:
```
create table fuzzing_definition
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    name       varchar(64)                         not null comment '数据描述',
    value      longtext                            not null comment '值, e.g.: string, integer, bool',
    type       enum('string','integer', 'bool', 'list', 'dict') default 'string' not null comment '值类型',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='随机数据定义';
```

#### 场景列表

- 表名： playbook
- 主键：id
- 唯一键：name, product
- 表创建SQL:
```
create table playbook
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    name       varchar(64)                         not null comment '场景描述',
    api_list      longtext                         null comment 'API列表',
    product    varchar(64)                         null comment '所属产品',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景列表';
```

#### 场景数据

- 表名： scene_data
- 主键：id
- 唯一键：app_id, app
- 表创建SQL:
```
create table scene_data
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    name       varchar(64)                         null comment '数据描述',
    api_id      varchar(255)                          not null comment '接口ID',
    app    varchar(64)                         null comment '所属应用',
    content    longtext                         null comment '.yml文件内容',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间',
    UNIQUE KEY `api_id_app_name` (`api_id`,`app`, "name")
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景数据';
```

#### 说明
 - 此处的信息，为最初的设计
 - 随着版本的迭代，部分表单的SQL设计未全部维护至此