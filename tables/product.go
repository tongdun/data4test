package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
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

func GetProductTable(ctx *context.Context) table.Table {

	product := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	apps := biz.GetApps()
	info := product.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField("唯一标识", "id", db.Int).
		FieldFilterable()
	info.AddField("产品名称", "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("环境IP/域名", "ip", db.Char).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("请求协议", "protocol", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("是否并发", "threading", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return "yes"
		}
		if model.Value == "no" {
			return "no"
		}
		return "unknown"
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
	info.AddField("测试模式", "testmode", db.Enum)
	info.AddField("关联应用", "apps", db.Longtext).FieldWidth(250)
	info.AddField("环境类型", "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "开发"
			} else if model.Value == "2" {
				return "测试"
			} else if model.Value == "3" {
				return "预发"
			} else if model.Value == "4" {
				return "演示"
			} else if model.Value == "5" {
				return "生产"
			}
			return ""
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "开发"},
		{Value: "2", Text: "测试"},
		{Value: "3", Text: "预发"},
		{Value: "4", Text: "演示"},
		{Value: "5", Text: "生产"},
	})
	info.AddField("专用参数", "private_parameter", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldHide()
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()
	info.AddButton("复制", icon.Android, action.Ajax("product_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = "请先选择数据再复制"
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.CopyProduct(id, userNameSub); err == nil {
					status = "复制成功，请刷新列表查看"
				} else {
					status = fmt.Sprintf("复制失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("复制", action.Ajax("product_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopyProduct(id, userNameSub); err == nil {
				status = "复制成功，请刷新列表查看"
			} else {
				status = fmt.Sprintf("复制失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("查看报告", action.Jump("/admin/product_dashboard?id={{.Id}}"))

	info.AddButton("导入场景", icon.Android, action.Ajax("import_batch_scene",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择产品再导入"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}

				if newCount, oldCount, err := biz.ImportPlaybookFromExcel(id); err == nil {
					status = fmt.Sprintf("导入完成, 新增: %d条，已存在: %d条", newCount, oldCount)
				} else {
					status = fmt.Sprintf("同步失败: %s, 新增: %d条，已存在: %d条", err, newCount, oldCount)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("导入场景", action.Ajax("import_scene",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string

			if newCount, oldCount, err := biz.ImportPlaybookFromExcel(id); err == nil {
				status = fmt.Sprintf("导入完成, 新增: %d条，已存在: %d条", newCount, oldCount)
			} else {
				status = fmt.Sprintf("同步失败: %s, 新增: %d条，已存在: %d条", err, newCount, oldCount)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddButton("XMind导入测试用例", icon.Android, action.Ajax("testcase_xmind_batch_import",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.GetJSON(id); err == nil {
					status = "导入完成"
				} else {
					status = fmt.Sprintf("导入失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("XMind导入测试用例", action.Ajax("testcase_xmind_import",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.GetJSON(id); err == nil {
				status = "导入完成"
			} else {
				status = fmt.Sprintf("导入失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.SetTable("product").SetTitle("产品配置").SetDescription("产品配置")
	helpMsg := template.HTML("JSON格式，e.g.: {\"name1\": \"value1\", \"name2\": \"value2\"}")

	formList := product.GetForm()
	formList.AddField("唯一标识", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("产品名称", "product", db.Varchar, form.Text)
	formList.AddField("环境IP/域名", "ip", db.Char, form.Ip)
	formList.AddField("请求协议", "protocol", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "http", Value: "http"},
			{Text: "https", Value: "https"},
		}).FieldDefault("http")
	formList.AddField("是否并发", "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "是", Value: "yes"},
			{Text: "否", Value: "no"},
		}).FieldDefault("no")
	formList.AddField("并发数", "thread_number", db.Int, form.Number).FieldDefault("1")
	formList.AddField("鉴权信息", "auth", db.Longtext, form.TextArea)
	formList.AddField("测试模式", "testmode", db.Enum, form.Radio).FieldOptions(types.FieldOptions{
		{Text: "自定义", Value: "custom"},
		{Text: "模糊", Value: "fuzzing"},
		{Text: "全部", Value: "all"},
	}).FieldDefault("custom")
	formList.AddField("关联应用", "apps", db.Longtext, form.Select).
		FieldOptions(apps)
	formList.AddField("环境类型", "env_type", db.Int, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "开发", Value: "1"},
			{Text: "测试", Value: "2"},
			{Text: "预发", Value: "3"},
			{Text: "演示", Value: "4"},
			{Text: "生产", Value: "5"},
		}).FieldDefault("2")
	formList.AddField("专用参数", "private_parameter", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("product").SetTitle("产品配置").SetDescription("产品配置")

	//formList.SetPostHook(func(values form2.Values) (err error) {  // 存量数据升级函数
	//	err = biz.ModifyPlaybookContent()
	//	return
	//})

	return product
}
