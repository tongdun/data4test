#!/bin/bash

truncate goadmin_operation_log;
truncate scene_test_history;
truncate scene_data_test_history;
truncate api_test_detail;
truncate api_test_result;

#mysqldump -uUsername -h xx.xx.xx.xx -P3306 DBname t_name --set-gtid-purged=off -p >t_name.sql(表名)
### utf8_unicode_ci

# 备份全部数据表结果，排除goadmin_operation_log这张表
mysqldump -h 127.0.0.1 -u root -p data4test --ignore-table=data4test.goadmin_operation_log > data4test_all.sql
mysqldump -h 127.0.0.1 -u root -p data4test goadmin_operation_log >> data4test_all.sql

#备份表结构
mysqldump -h 127.0.0.1 -u root -p data4test -d > data4test_table.sql

# 备份表的数据
mysqldump -h 127.0.0.1 -u root -p data4test filemanager_setting goadmin_menu goadmin_permissions goadmin_site -t > data4test_init.sql

# 收集有内置的系统数据
mysqldump -h 127.0.0.1 -u root -p data4test goadmin_role_menu goadmin_role_permissions goadmin_role_users goadmin_roles goadmin_users goadmin_user_permissions -t >> data4test_init.sql

# 备份有内置的示例数据
mysqldump -h 127.0.0.1 -u root -p data4test assert_template env_config sys_parameter scene_data  product playbook -t >> data4test_init.sql

# 备份有内置的任务数据
mysqldump -h 127.0.0.1 -u root -p data4test schedule -t >> data4test_init.sql
