#!/bin/bash

# 生成数据初始化 SQL：
#   1. 导出全库「表结构」(-d)
#   2. 整表导出框架数据（GoAdmin 菜单/权限/角色/用户等，全表均为内置数据）
#   3. 按「特定名称」过滤导出业务内置数据，避免把运行中用户新建的数据带入初始化 SQL
#
# 维护内置数据：只需改下面两个变量，无需改下方逻辑。
#
# 参考：
#   # 单表备份：mysqldump -uUsername -h xx.xx.xx.xx -P3306 DBname 表名 --set-gtid-purged=off -p > 表名.sql
#   # 字符集：utf8_unicode_ci
#   # 运行时日志表（无需内置数据）：goadmin_operation_log / scene_test_history / scene_data_test_history / api_test_detail / api_test_result

export HOST_IP="127.0.0.1"
export USER_NAME="root"
export DB_PORT="3306"
export DATABASE="data4test"
export BAK_FILE_NAME="data4test_init_$(date +%Y%m%d).sql"

# ===== 框架数据表：整表导出（全表均为内置数据）=====
FRAMEWORK_TABLES="filemanager_setting goadmin_menu goadmin_permissions goadmin_site goadmin_role_menu goadmin_role_permissions goadmin_role_users goadmin_roles goadmin_users goadmin_user_permissions"

# ===== 业务内置数据表：仅导出「特定名称」的数据 =====
# 格式："表名|过滤列|内置名称(多个用英文逗号分隔)"
BUILTIN_TABLES=(
  "assert_template|name|successTemplate"
  "env_config|app|exampleApp"
  "sys_parameter|name|fileName,scriptRunEngine,aiRunEngine,aiPlatform,Router4Add,TestCaseType,RUID,testCaseExportTemplates,caseProduct,supportLanguages"
  "scene_data|name|示例-用户管理-新建用户"
  "product|product|示例产品"
  "playbook|name|示例场景"
  "schedule|task_name|示例任务-关联场景-一次性,示例任务-关联数据-自定义执行,示例任务-关联场景-每天0点和12点执行,示例任务-关联场景-每周六和每周日0点和20点执行"
  "ai_template|template_name|接口定义生成测试用例,需求文档生成测试用例,接口定义生成测试数据,根据测试数据分析结果数据"
)

# 1. 备份全库表结构
mysqldump -h "$HOST_IP" -u "$USER_NAME" -P"$DB_PORT" -p "$DATABASE" -d > "$BAK_FILE_NAME"

# 2. 整表导出框架数据
echo "-- ===== 框架数据（整表）=====" >> "$BAK_FILE_NAME"
# shellcheck disable=SC2086  # 表名需按空格拆分为多个参数
mysqldump -h "$HOST_IP" -u "$USER_NAME" -P"$DB_PORT" -p "$DATABASE" $FRAMEWORK_TABLES -t >> "$BAK_FILE_NAME"

# 3. 按名称过滤导出业务内置数据
for entry in "${BUILTIN_TABLES[@]}"; do
  IFS='|' read -r table col names <<< "$entry"

  # 把逗号分隔的名称拼成 SQL IN 列表：'a','b',...
  sql_names=""
  IFS=',' read -r -a name_arr <<< "$names"
  for n in "${name_arr[@]}"; do
    [ -n "$sql_names" ] && sql_names="${sql_names},"
    sql_names="${sql_names}'${n}'"
  done

  where="${col} IN (${sql_names})"
  echo "-- ===== 内置数据：${table} WHERE ${where} =====" >> "$BAK_FILE_NAME"
  mysqldump -h "$HOST_IP" -u "$USER_NAME" -P"$DB_PORT" -p --where="$where" "$DATABASE" "$table" -t >> "$BAK_FILE_NAME"
done
