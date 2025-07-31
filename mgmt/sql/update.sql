# 更新SQL维护
# 2024年1月19日
alter table playbook modify scene_type enum('1', '2', '3', '4', '5') default '1' null comment '场景类型，1:串行中断，2:串行比较，3:串行继续，4:普通并发, 5:并发比较';
alter table scene_test_history modify scene_type enum('1', '2', '3', '4', '5') default '1' null comment '场景类型，1:串行中断，2:串行比较，3:串行继续，4:普通并发, 5:并发比较';

# 2024年2月25日
alter table scene_data
    add file_type ENUM ('1', '2', '3', '4', '5', '99') default '1' comment '文件类型，1:标准数据，2:Python脚本，3:Shell脚本，4:Bat脚本，5:Jmeter脚本，99:其他脚本' after file_name;

# 2024年4月15日
update api_definition set auto='1' where auto='是';
update api_definition set auto='0' where auto='否';
alter table api_definition change auto is_need_auto enum('1', '0') default '1' null comment '是否需自动化';
alter table api_definition
    add is_auto enum('1', '0') default '0' null comment '是否已自动化';
alter table api_definition modify api_status enum('1', '2', '3', '4') default '1' null comment '接口状态, 1:新增,2:被删除,3:被修改,4:保持原样';
alter table api_definition modify is_auto enum('1', '0') default '0' null comment '是否已自动化';
alter table api_definition modify is_auto enum('1', '-1') default '-1' null comment '是否已自动化,1:已自动化，-1:未自动化';
alter table api_definition modify is_need_auto enum('1', '-1') default '-1' null comment '是否需自动化,1:需自动化，-1:无需自动化';

# 2025年2月6日
alter table api_definition modify api_status enum('1', '2', '3', '4', '30', '31', '32', '33', '34') default '1' null comment '接口状态, 1:新增,2:被删除,3:被修改,4:保持原样,30:Header被修改，31:Path被修改，32:Query被修改，33:Body被修改，34:Resp被修改';

# 2025年4月1日
# 智能用例
CREATE TABLE IF NOT EXISTS ai_case
(
    id            int auto_increment comment '自增主键'
        primary key,
    case_number   varchar(255)                                   not null comment '用例编号',
    case_name     varchar(255)                                   not null comment '用例名称',
    case_type     varchar(16)                                    null comment '用例类型',
    priority      varchar(6)                                     null comment '优先级',
    pre_condition text                                           null comment '前置条件',
    test_range    text                                           null comment '测试范围',
    test_steps    text                                           null comment '测试步骤',
    expect_result text                                           null comment '预期结果',
    auto          enum ('0', '1', '2') default '1'               null comment '是否自动化，0:否, 1:是, 2:部分是',
    module        varchar(255)                                   null comment '所属模块',
    intro_version varchar(64)                                    null comment '引入版本',
    case_version  int                  default 1                 null comment '用例版本',
    product       varchar(64)                                    null comment '关联产品',
    source        varchar(16)                                    null comment '生成来源: DeepSeek/OpenAi/Kimi等',
    use_status    enum ('1', '2', '3') default '1'               null comment '取用状态, 1:初始, 2:取用, 3:废弃',
    modify_status enum ('1', '2', '3') default '1'               null comment '改造状态, 1:初始, 2:人工改造, 3:自动改造',
    create_user   varchar(100)                                   null comment '创建人',
    created_at    timestamp            default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at    timestamp            ON UPDATE CURRENT_TIMESTAMP  null comment '更新时间',
    deleted_at    timestamp                                      null comment '删除时间',
    constraint case_number_case_name_source_product
        unique (case_number, case_name, source, product)
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能用例';


# 智能数据
CREATE TABLE IF NOT EXISTS ai_data (
                                       id INT AUTO_INCREMENT PRIMARY KEY comment '自增主键',
                                       data_desc varchar(255) not null comment '数据描述',
                                       api_id varchar(255) not null comment '接口ID',
                                       app varchar(64) null comment '所属应用',
                                       source varchar(16) DEFAULT NULL COMMENT '生成来源: DeepSeek/OpenAi/Kimi等',
                                       file_name varchar(255) COMMENT '文件名称',
                                       file_type ENUM ('1', '2', '3', '4', '5', '99') default '1' comment '文件类型，1:标准数据，2:Python脚本，3:Shell脚本，4:Bat脚本，5:Jmeter脚本，99:其他脚本',
                                       file_content TEXT COMMENT '文件内容',
                                       result varchar(5) DEFAULT NULL COMMENT '测试结果',
                                       fail_reason TEXT COMMENT '失败原因',
                                       use_status ENUM('1', '2', '3') DEFAULT '1' COMMENT '取用状态, 1:初始, 2:取用, 3:废弃',
                                       modify_status ENUM('1', '2', '3') DEFAULT '1' COMMENT '改造状态, 1:初始, 2:人工改造, 3:自动改造',
                                       create_user varchar(100) DEFAULT NULL COMMENT '创建人',
                                       created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
                                       updated_at timestamp null comment '更新时间',
                                       deleted_at timestamp null comment '删除时间',
                                       INDEX data_desc_app (data_desc, app),
                                       UNIQUE KEY `api_id_app_data_desc` (`api_id`, `app`, `data_desc`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '智能数据';

# 智能场景
CREATE TABLE IF NOT EXISTS ai_playbook
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    playbook_desc       varchar(255)  not null comment '场景描述',
    data_file_list      text                         null comment '数据文件列表',
    playbook_type enum('1', '2', '3', '4', '5') default '1' null comment '场景类型，1:串行中断，2:串行比较，3:串行继续，4:普通并发, 5:并发比较',
    priority int(11) DEFAULT NULL COMMENT '优先级',
    source    varchar(16) DEFAULT NULL COMMENT '生成来源: DeepSeek/OpenAi/Kimi等',
    use_status ENUM('1','2','3')   DEFAULT '1' COMMENT '取用状态, 1:初始, 2:取用, 3:废弃',
    modify_status   ENUM('1','2','3')   DEFAULT '1' COMMENT '改造状态, 1:初始, 2:人工改造, 3:自动改造',
    result varchar(5)  DEFAULT NULL COMMENT '测试结果',
    fail_reason TEXT COMMENT '失败原因',
    product    varchar(64)                         null comment '所属产品',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能场景';

# 智能任务
CREATE TABLE IF NOT EXISTS ai_task
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    task_name varchar(255) NOT NULL COMMENT '任务名称',
    task_mode enum('once')  NOT NULL DEFAULT 'once' COMMENT '任务模式',
    task_type enum('data','playbook')  DEFAULT 'playbook' COMMENT '任务类型',
    source    varchar(16) DEFAULT NULL COMMENT '生成来源: DeepSeek/OpenAi/Kimi等',
    task_status enum('running','stopped','finished','not_started')  DEFAULT 'not_started' COMMENT '任务状态',
    data_list text  COMMENT '关联数据',
    playbook_list text COMMENT '关联场景',
    use_status ENUM('1','2','3')   DEFAULT '1' COMMENT '取用状态, 1:初始, 2:取用, 3:废弃',
    modify_status   ENUM('1','2','3')   DEFAULT '1' COMMENT '改造状态, 1:初始, 2:人工改造, 3:自动改造',
    product    varchar(64) null comment '所属产品',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能任务';

# 智能分析
CREATE TABLE IF NOT EXISTS ai_issue (
                                        id                  INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主键',
                                        issue_name          VARCHAR(255)                        NOT NULL COMMENT '问题名称',
                                        issue_level         VARCHAR(16)   NULL  COMMENT '问题级别: P0~P4或严重/轻微等',
                                        issue_source        ENUM('1','2','3')   NULL COMMENT '问题来源: 1:数据, 2:场景, 3:手动录入',
                                        source_name         VARCHAR(255)                      NULL COMMENT '来源名称: 数据名称/场景名称',
                                        request_data        TEXT                          NOT NULL COMMENT '请求数据: url, header, request等',
                                        response_data       TEXT                                NULL COMMENT '返回数据: response',
                                        issue_detail        TEXT                                NULL COMMENT '问题详情:可含预期结果，复现步骤等',
                                        confirm_status      ENUM('1','2','3')           DEFAULT '1' COMMENT '确认状态: 1:BUG, 2:优化, 3:误判',
                                        root_cause          TEXT                     NULL COMMENT '问题原因推测',
                                        impact_scope_analysis    TEXT                        NULL COMMENT '影响范围分析',
                                        impact_playbook  TEXT  NULL COMMENT '受影响的场景推测',
                                        impact_data  TEXT  NULL COMMENT '受影响的数据推测',
                                        resolution_status   ENUM('1','2','3','4') DEFAULT '1' COMMENT '解决状态, 1:创建, 2:解决中, 3:修复完成, 4:验证完成',
                                        again_test_result   ENUM('0','1', '2')   DEFAULT '2' COMMENT '回归测试结果, 0:失败，1:成功，2:未知',
                                        impact_test_result   ENUM('0','1', '2')   DEFAULT '2' COMMENT '受影响模块回归测试结果, 0:失败，1:成功，2:未知',
                                        product_list    TEXT null comment '关联产品， 可多产品环境回归执行',
                                        create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
                                        modify_user varchar(100)  DEFAULT NULL COMMENT '修改人',
                                        created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
                                        updated_at          TIMESTAMP                           NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                        deleted_at          TIMESTAMP                           NULL COMMENT '删除时间',
                                        INDEX idx_tracking (issue_source, source_name),
                                        UNIQUE KEY uk_issue_identity (issue_name, source_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能问题跟踪表';


# 智能模板
CREATE TABLE IF NOT EXISTS ai_template
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    template_name       varchar(255)  not null comment '模板名称',
    template_type enum('1', '2', '3', '4', '5', '6') default '1' null comment '模板类型，1:用例，2:数据，3:场景，4:任务, 5:Issue, 6: 报告',
    template_content text null comment '模板内容',
    use_status enum('apply', 'edit') default 'edit' null  COMMENT '生效状态: 启用/编辑',
    applicable_platform varchar(64)  NULL DEFAULT '通用' COMMENT '适用平台',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    modify_user varchar(100)  DEFAULT NULL COMMENT '修改人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间',
    UNIQUE KEY `template_name_template_type_applicable_platform` (`template_name`,`template_type`, `applicable_platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能模板';

# 智能报告
CREATE TABLE IF NOT EXISTS ai_report
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    report_name       varchar(255)  not null comment '报告名称',
    demand       text  not null comment '报告需求',
    `source`   varchar(16) DEFAULT NULL COMMENT '生成来源: DeepSeek/OpenAi/Kimi等',
    use_status ENUM('1','2','3')   DEFAULT '1' COMMENT '取用状态, 1:初始, 2:取用, 3:废弃',
    modify_status   ENUM('1','2','3')   DEFAULT '1' COMMENT '改造状态, 1:初始, 2:人工改造, 3:自动改造',
    report_link     varchar(255)   null comment '报告详情：超链形式',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    modify_user varchar(100)  DEFAULT NULL COMMENT '修改人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='智能报告';


# 智能生成
CREATE TABLE IF NOT EXISTS ai_create
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    create_desc        text  not null comment '生成指令',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='生成指令';


# 智能优化
CREATE TABLE IF NOT EXISTS ai_optimize
(
    id         int unsigned auto_increment
        primary key comment '自增主键',
    optimize_desc       text not null comment '优化指令',
    create_user varchar(100)  DEFAULT NULL COMMENT '创建人',
    created_at timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    updated_at timestamp                           null comment '更新时间',
    deleted_at timestamp                           null comment '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优化指令';

alter table ai_template
    add template_content text null comment '模板内容';
alter table ai_template modify template_content text null comment '模板内容' after template_type;
# alter table ai_case
#     add pre_condition text null comment '前置条件' after priority;
# alter table ai_case
#     add test_range text null comment '测试范围' after pre_condition;
# alter table ai_case
#     add test_steps text null comment '测试步骤' after test_range;
# alter table ai_case
#     add expect_result text null comment '预期结果' after test_steps;

alter table test_case modify auto enum('0', '1', '2') default '1' null comment '是否自动化: 0:否,1:是,2:部分是';
alter table test_case
    add intro_version varchar(64) null comment '引入版本' after module;

alter table test_case modify updated_at timestamp null ON UPDATE CURRENT_TIMESTAMP comment '更新时间';
alter table test_case modify test_time timestamp null comment '测试时间';

alter table ai_template
    add append_conversion TEXT null comment '追加对话' after template_content;

alter table ai_data change data_desc name varchar(255) not null comment '数据描述';
alter table ai_playbook
    add last_file varchar(255) null comment '最近数据文件' after data_file_list;
alter table ai_data drop column file_content;

alter table ai_template modify template_type enum('1', '2', '3', '4', '5', '6', '7') default '1' null comment '模板类型，1:用例，2:数据，3:场景，4:任务, 5:Issue, 6: 报告， 7: 分析';

alter table ai_issue
    add `source` varchar(16) DEFAULT NULL COMMENT '生成来源: DeepSeek/OpenAi/Kimi等' after source_name;
alter table ai_data modify content longtext null comment '文件内容';
alter table ai_data modify file_name text null comment '文件名称';
alter table ai_issue modify request_data longtext not null comment '请求数据: url, header, request等';
alter table ai_issue modify response_data longtext null comment '返回数据: response';
alter table ai_issue modify resolution_status enum('1', '2', '3', '4', '5') default '1' null comment '解决状态, 1:创建, 2:解决中, 3:修复完成, 4:验证完成, 5:不处理';

# 2025年6月17日
alter table scene_data_test_history
    add file_type enum('1', '2', '3', '4', '5', '99') null comment '环境类型' after env_type;

# 2025年7月25日
alter table playbook change api_list data_file_list text null comment '数据文件列表';
alter table scene_test_history change api_list data_file_list text null comment '数据文件列表';
alter table ai_playbook change playbook_desc name varchar(255) not null comment '场景描述';
alter table ai_playbook change playbook_type scene_type enum('1', '2', '3', '4', '5') default '1' null comment '场景类型，1:串行中断，2:串行比较，3:串行继续，4:普通并发, 5:并发比较';
