#### 依赖框架
- 数据库ORM框架：GORM v1.9.16
- Web框架：Gin v1.7.2
- 文件管理框架：Filemanager v0.0.7
- 后台数据管理框架：GoAdmin  v1.2.24
- UI组件库：Iview v2.7.4
- 控制台开发框架：Vue 2.5.2
- 控制台前端：在goman的前端框架上进行二次开发
- Go 工具链：1.27.3（go.mod go 指令 1.23）

#### 运行形态
- Web 形态：完整 Web 服务（GoAdmin + Gin），管理数据存 MySQL，执行历史存服务器文件
- CLI 形态：CLI/SQLite 独立运行模式，使用 SQLite 存储，CLI执行数据/场景

#### 架构图
<img src="../image/arch.jpg">

#### 全局使用流程使用
<img src="../image/全局使用流程图.jpg">

#### 数据流程图
<img src="../image/数据流程图.png">

#### 系统数据关系图
<img src="../image/系统数据关系图.jpg">