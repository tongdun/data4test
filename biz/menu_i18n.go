package biz

import (
	"data4test/models"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"strings"
)

func InitMenuI18n() {
	models.Orm.Table("goadmin_menu").Where("type = 0").UpdateColumn("type", 1)

	var titles []string
	models.Orm.Table("goadmin_menu").Where("id > 1").Pluck("title", &titles)

	if len(titles) == 0 {
		return
	}

	enMap := make(map[string]string)
	cnMap := make(map[string]string)

	for _, title := range titles {
		cnMap[title] = title
		if en := menuTitleEn(title); en != "" {
			enMap[title] = en
		}
	}

	if language.Lang != nil {
		if language.Lang[language.EN] == nil {
			language.Lang[language.EN] = make(language.LangSet)
		}
		for k, v := range enMap {
			language.Lang[language.EN][strings.ToLower(k)] = v
		}

		if language.Lang[language.CN] == nil {
			language.Lang[language.CN] = make(language.LangSet)
		}
		for k, v := range cnMap {
			language.Lang[language.CN][k] = v
		}
	}

	//Logger.Info("菜单国际化初始化完成: %d 个菜单标题已注册翻译", len(titles))
}

// RefreshMenuI18n 刷新菜单翻译注册，获取数据库中最新菜单列表重新注册翻译
// 在新增/修改菜单后调用，确保新菜单也有对应的翻译
func RefreshMenuI18n() {
	if models.Orm == nil {
		Logger.Warning("RefreshMenuI18n: models.Orm 未初始化，跳过菜单翻译刷新")
		return
	}
	InitMenuI18n()
}

func menuTitleEn(cnTitle string) string {
	m := map[string]string{
		"盾测-自动化":    "Data4Test",
		"控制台":       "Console",
		"统计报告":      "Dashboard",
		"全局":        "Overview",
		"概览":        "Overview",
		"产品统计报告":    "Product Dashboard",
		"产品维度":      "By Product",
		"产品配置":      "Product Config",
		"应用统计报告":    "App Dashboard",
		"应用维度":      "By App",
		"应用配置":      "App Config",
		"环境":        "Env",
		"参数定义":      "Parameter Defs",
		"系统参数":      "System Params",
		"断言值模板":     "Assert Template",
		"模糊因子":      "Fuzz Factor",
		"接口":        "API",
		"接口定义":      "API Define",
		"接口关系":      "API Relation",
		"接口关联关系":    "API Relation",
		"接口测试数据":    "API Test Data",
		"接口记录表":     "API Record Table",
		"接口ID统计":    "API ID Stats",
		"接口测试结果":    "API Test Result",
		"结果详情":      "Detail Result",
		"结果":        "Result",
		"测试数据":      "Test Data",
		"模糊数据":      "Fuzz Data",
		"模糊测试":      "Fuzz Testing",
		"模糊定义":      "Fuzz Definition",
		"Swagger文件": "Swagger Files",
		"场景":        "Playbook",
		"场景列表":      "Playbook List",
		"场景配置":      "Playbook Config",
		"场景测试历史":    "Playbook Test History",
		"数据文件":      "Data File",
		"数据列表":      "Data List",
		"数据测试历史":    "Data Test History",
		"结果报告列表":    "Result Report List",
		"产品统计":      "Product Stats",
		"任务":        "Task",
		"定时任务":      "Schedule Task",
		"任务执行报告":    "Execution Report",
		"文件":        "File",
		"公共文件":      "Public File",
		"用例文件":      "Test Case File",
		"API文件":     "API File",
		"上传文件":      "Upload File",
		"下载文件":      "Download File",
		"日志文件":      "Log File",
		"历史记录":      "History Record",
		"历史版本":      "History Version",
		"日志":        "Log",
		"用例":        "TestCase",
		"用例管理":      "Case Mgmt",
		"测试用例":      "Test Case",
		"自定义接口":     "Custom API",
		"接口变更日志":    "API Changelog",
		"变更日志":      "Changelog",
		"系统变更":      "System Change",
		"最新发布":      "Latest Release",
		"智能用例":      "AI Test Case",
		"智能数据":      "AI Test Data",
		"智能场景":      "AI Playbook",
		"智能分析":      "AI Analysis",
		"智能报告":      "AI Report",
		"智能模板":      "AI Template",
		"AI任务":      "AI Task",
		"AI模板":      "AI Template",
		"AI生成配置":    "AI Creation",
		"AI优化":      "AI Optimization",
		"智能测试分析":    "AI Test Analysis",
		"智能数据测试":    "AI Data Test",
		"助手":        "Assistant",
		"智能助手":      "AI Assistant",
		"研发事项":      "Development",
		"开发须知":      "Dev Note",
		"开发待办":      "Dev Todo",
		"数据库设计":     "Database Design",
		"系统设计":      "System Design",
		"架构图":       "Architecture",
		"发展规划":      "Roadmap",
		"功能特性":      "Feature",
		"模块特性":      "Module Feature",
		"性能测试设计":    "Perf Test Design",
		"模糊测试设计":    "Fuzz Test Design",
		"说明":        "Doc",
		"脚本说明":      "Script Guide",
		"断言说明":      "Assert Guide",
		"动作说明":      "Action Guide",
		"参数说明":      "Parameter Guide",
		"系统使用":      "User Guide",
		"控制台使用":     "Console Usage",
		"接口管理":      "API Mgmt",
		"场景管理":      "Playbook Mgmt",
		"数据管理":      "Data Mgmt",
		"任务管理":      "Task Mgmt",
		"Mock使用":    "Mock Usage",
		"常见问题":      "FAQ",
		"知识库":       "Knowledge Base",
		"配置":        "Config",
		"统计":        "Stats",
		"工具":        "Tool",
	}

	if en, ok := m[cnTitle]; ok {
		return en
	}
	return cnTitle
}
