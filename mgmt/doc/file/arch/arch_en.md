#### Dependency Frameworks
- Database ORM Framework: GORM v1.9.16
- Web Framework: Gin v1.7.2
- File Management Framework: Filemanager v0.0.7
- Backend Data Management Framework: GoAdmin v1.2.24
- UI Component Library: Iview v2.7.4
- Console Development Framework: Vue 2.5.2
- Console Frontend: Secondary development based on the frontend framework of goman
- Go Toolchain: 1.27.3 (go.mod `go` directive 1.23)

#### Run Modes
- Web mode: full web service (GoAdmin + Gin), management data stored in MySQL, execution history stored in server files.
- CLI mode: CLI/SQLite standalone run mode, uses SQLite storage, CLI execution of data/scenarios.

#### Architecture Diagram
<img src="../image/arch.jpg">

#### Global Usage Flow
<img src="../image/全局使用流程图.jpg">

#### Data Flow Diagram
<img src="../image/数据流程图.png">

#### System Data Relationship Diagram
<img src="../image/系统数据关系图.jpg">
