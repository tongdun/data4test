#### 任务类型
- 自定义任务: Crontab表达式，按需设置执行频率
- 一次任务: 执行一次即停止，可反复执行
- 每天任务: 按天维度设置执行任务
- 每周任务: 按周维护设置执行任务

#### 任务执行范围
- 数据类型任务
- 场景类型任务

#### 任务执行环境
- 可关联单个环境：功能测试，性能测试等
- 可关联多个环境：信创多环境，可一次执行
- 可不关联环境：直接使用场景关联的环境

#### Crontab表达式
```
* * * * *
- - - - -
| |  | | |
| |  | | +-------------- 星期中星期几 (0 - 6) (星期天 为0)
| |  | +---------------- 月份 (1 - 12)
| |  +------------------ 一个月中的第几天 (1 - 31)
| +---------------------- 小时 (0 - 23)
+------------------------- 分钟 (0 - 59)
```