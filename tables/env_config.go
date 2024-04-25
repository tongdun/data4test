package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"html/template"
	"strings"
)

func GetEnvConfigTable(ctx *context.Context) table.Table {

	envConfig := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	info := envConfig.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("产品名称", "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products)
	info.AddField("应用名称", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("环境IP/域名", "ip", db.Char).
		FieldSortable()
	info.AddField("请求协议", "protocol", db.Enum)
	info.AddField("路由前缀", "prepath", db.Varchar)
	info.AddField("是否并发", "threading", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return "是"
		}
		if model.Value == "no" {
			return "否"
		}
		return "否"
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "yes", Text: "是"},
		{Value: "no", Text: "否"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "yes", Text: "是"},
		{Value: "no", Text: "否"},
	})

	info.AddField("并发数", "thread_number", db.Int)

	info.AddField("鉴权信息", "auth", db.Longtext).FieldWidth(300).
		FieldHide()
	info.AddField("测试模式", "testmode", db.Enum).
		FieldFilterable()
	info.AddField("Swagger路径", "swagger_path", db.Text).
		FieldHide()
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddActionButton("查看报告", action.Jump("/admin/app_dashboard?id={{.Id}}"))

	info.AddButton("导入接口", icon.Android, action.Ajax("autogeneration_batch",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再导入"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if failCount, err := biz.GetSwagger(id); err == nil {
					if failCount > 0 {
						status = fmt.Sprintf("导入完成，%d个接口不符合规范，请前往[接口定义]列表查看", failCount)
					} else {
						status = "导入完成，请前往[接口定义]列表查看"
					}
					_ = biz.UpdateApiChangeByAppId(id)
				} else {
					status = fmt.Sprintf("导入失败：%s", err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("导入接口", action.Ajax("autogeneration",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if failCount, err := biz.GetSwagger(id); err == nil {
				if failCount > 0 {
					status = fmt.Sprintf("导入完成，%d个接口不符合规范，请前往[接口定义]列表查看", failCount)
				} else {
					status = "导入完成，请前往[接口定义]列表查看"
				}
				_ = biz.UpdateApiChangeByAppId(id)
			} else {
				status = fmt.Sprintf("导入失败：%s", err)
			}

			return true, status, ""
		}))

	info.AddButton("更新状态", icon.Android, action.Ajax("autoupdatestatus_batch",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再更新"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.UpdateApiAutoStatus(id); err == nil {
					status = "接口是否已自动化状态更新完成，请前往[接口定义]列表查看"
				} else {
					status = fmt.Sprintf("更新失败：%s", err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("更新状态", action.Ajax("autoupdatestatus",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.UpdateApiAutoStatus(id); err == nil {
				status = "接口是否已自动化状态更新完成，请前往[接口定义]列表查看"
			} else {
				status = fmt.Sprintf("更新失败：%s", err)
			}

			return true, status, ""
		}))

	info.AddButton("更新鉴权", icon.Android, action.Ajax("update_batch_auth",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			//idStr := ctx.FormValue("ids")
			var status string
			status = "功能暂未合入"
			//if idStr == "," {
			//	//status = "请先选择数据再更新"
			//	status = "功能暂未合入"
			//	return false, status, ""
			//}
			//ids := strings.Split(idStr, ",")
			//for _, id := range ids {
			//	if len(id) == 0 {
			//		//status = "更新完成，请开始接口测试之旅"
			//		status = "功能暂未合入"
			//		continue
			//	}
			//	if err := biz.GetAuth(id); err == nil {
			//		status = "更新完成，请开始接口测试之旅"
			//	} else {
			//		status = fmt.Sprintf("更新失败：%s", err)
			//		return false, status, ""
			//	}
			//}
			return false, status, ""
		}))

	info.AddActionButton("更新鉴权", action.Ajax("update_auth",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			//id := ctx.FormValue("id")
			var status string
			status = "功能暂未合入"
			//if err := biz.GetAuth(id); err == nil {
			//	status = "更新完成，请开始接口测试之旅"
			//} else {
			//	status = fmt.Sprintf("更新失败：%s", err)
			//}

			return false, status, ""
		}))

	info.AddSelectBox("关联产品", products, action.FieldFilter("product"))
	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))
	info.SetTable("env_config").SetTitle("应用配置").SetDescription("应用配置")

	formList := envConfig.GetForm()

	helpMsg := template.HTML("配置接口JSON格式数据路由，e.g.:<br> 直接触发：http://10.X.X.X:80XX/api/v1/api-docs?group=groupName<br>CICD触发自动替换占位符:http://{host:port}/api/v2/api-docs?group=groupName<br>为空时CICD触发自动填充: config.json中swagger_path配置项")
	authHelpMsg := template.HTML("配置请求数据Header头，e.g.: <br>{\"Content-Type\":\"application/x-www-form-urlencoded\",\"Cookie\":\"xxx\",\"Referer\":\"http://x.x.x.x:80xx\",\"X-Cf-Random\":\"xxx\"}")
	appHelpMsg := template.HTML("应用名称如果有gitlab仓库，请与仓库名保持一致，或直接通过CICD自动创建")
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
	formList.AddField("产品名称", "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("应用名称", "app", db.Varchar, form.Text).
		FieldDisableWhenUpdate().
		FieldHelpMsg(appHelpMsg)
	formList.AddField("环境IP/域名", "ip", db.Char, form.Ip)
	formList.AddField("请求协议", "protocol", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "http", Value: "http"},
			{Text: "https", Value: "https"},
		}).FieldDefault("http")

	formList.AddField("路由前缀", "prepath", db.Varchar, form.Text)
	formList.AddField("是否并发", "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "是", Value: "yes"},
			{Text: "否", Value: "no"},
		}).FieldDefault("no")
	formList.AddField("并发数", "thread_number", db.Int, form.Number).FieldDefault("1")
	formList.AddField("鉴权信息", "auth", db.Longtext, form.TextArea).
		FieldHelpMsg(authHelpMsg)
	formList.AddField("测试模式", "testmode", db.Enum, form.Radio).FieldOptions(types.FieldOptions{
		{Text: "自定义", Value: "custom"},
		{Text: "模糊", Value: "fuzzing"},
		{Text: "全部", Value: "all"},
	}).FieldDefault("custom")
	formList.AddField("Swagger路径", "swagger_path", db.Varchar, form.Text).
		FieldHelpMsg(helpMsg)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("env_config").SetTitle("应用配置").SetDescription("应用配置")

	return envConfig
}
