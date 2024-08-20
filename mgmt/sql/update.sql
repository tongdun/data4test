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