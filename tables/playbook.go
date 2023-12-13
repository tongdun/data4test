package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/PuerkitoBio/goquery"
	template2 "html/template"
	"strings"
)

func GetPlaybookTable(ctx *context.Context) table.Table {

	playbook := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	info := playbook.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("唯一标识", "id", db.Int).
		FieldFilterable().FieldWidth(150)
	info.AddField("场景描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldWidth(220)
	info.AddField("数据文件列表", "api_list", db.Longtext).FieldWidth(600).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
		})

	info.AddField("最近数据文件", "last_file", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return template.Default().
			Link().
			SetURL("/admin/fm/data/preview?path=/" + value.Value).
			SetContent(template2.HTML(value.Value)).
			OpenInNewTab().
			SetTabTitle(template.HTML("数据文件")).
			GetContent()
	}).FieldWidth(160)

	info.AddField("类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "默认"
			}
			if model.Value == "2" {
				return "比较"
			}
			return "默认"
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "默认"},
		{Value: "2", Text: "比较"},
	})
	info.AddField("优先级", "priority", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).
		FieldSortable().FieldWidth(80).
		FieldEditAble(editType.Text)
	info.AddField("执行次数", "run_time", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).FieldSortable().
		FieldEditAble(editType.Text)

	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}).FieldWidth(80)
	info.AddField("失败原因", "fail_reason", db.Longtext).FieldWidth(160)
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("所属产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products).FieldWidth(120)
	info.AddField("创建人", "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("复制", icon.Android, action.Ajax("playbook_batch_copy",
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
				if err := biz.CopyPlaybook(id, userNameSub); err == nil {
					status = "复制成功，请刷新列表查看"
				} else {
					status = fmt.Sprintf("复制失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("复制", action.Ajax("playbook_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopyPlaybook(id, userNameSub); err == nil {
				status = "复制成功，请刷新列表查看"
			} else {
				status = fmt.Sprintf("复制失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("测试任务", icon.Android, action.Ajax("playbook_batch_task_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}

			err := biz.AutoCreateSchedule(idStr, "scene")
			if err != nil {
				status = "发起测试任务失败"
				return false, status, fmt.Sprintf("%s", err)
			}
			//status = "已开启任务执行，请稍后查看测试结果，可在任务列表管理任务"
			status = "已创建任务，请前往任务列表查看，确认后再发起任务执行"
			return true, status, ""

		}))

	info.AddButton("测试", icon.Android, action.Ajax("playbook_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			ids, err := biz.GetPriority(idStr)
			if err != nil {
				status = fmt.Sprintf("%s", err)
				return false, status, ""
			}

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}

				if err := biz.RepeatRunPlaybook(id, "", "", "playbook"); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("测试", action.Ajax("playbook_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RepeatRunPlaybook(id, "", "", "playbook"); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("继续", icon.Android, action.Ajax("playbook_batch_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再继续"
				return false, status, ""
			}
			ids, err := biz.GetPriority(idStr)
			if err != nil {
				status = fmt.Sprintf("%s", err)
				return false, status, ""
			}

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.RepeatRunPlaybook(id, "continue", "", "playbook"); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("继续", action.Ajax("playbook_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string

			if err := biz.RepeatRunPlaybook(id, "continue", "", "playbook"); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddSelectBox("关联产品", products, action.FieldFilter("product"))

	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))
	info.SetTable("playbook").SetTitle("场景列表").SetDescription("场景列表")

	files := biz.GetFilesFromMySQL()

	formList := playbook.GetForm()
	formList.AddField("唯一标识", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("场景描述", "name", db.Varchar, form.Text)

	dataHelp := template.HTML("关联数据必填")
	formList.AddTable("关联数据", "data_table", func(panel *types.FormPanel) {
		panel.AddField("序号/标签", "data_number", db.Varchar, form.Text).
			FieldHideLabel().
			FieldDisplay(func(model types.FieldModel) interface{} {
				return strings.Split(model.Value, ",")
			})
		panel.AddField("关联数据", "api_list", db.Varchar, form.SelectSingle).
			FieldHideLabel().
			FieldOptions(files).
			FieldDisplay(func(model types.FieldModel) interface{} {
				var afterTxt string
				if !strings.Contains(model.Value, "</a>") && !strings.Contains(model.Value, "</p>") {
					afterTxt = model.Value
				} else {
					doc, _ := goquery.NewDocumentFromReader(strings.NewReader(model.Value))
					handle := doc.Text()
					afterTxt1 := strings.Replace(handle, ".yml", ".yml,", -1)
					afterTxt2 := strings.Replace(afterTxt1, ".json", ".json,", -1)
					afterTxt = strings.Replace(afterTxt2, ".yaml", ".yaml,", -1)
				}
				dataList := strings.Split(afterTxt, ",")

				return dataList
			})
		panel.SetInputWidth(10)
	}).FieldHelpMsg(dataHelp)

	formList.AddField("最近数据文件", "last_file", db.Varchar, form.SelectSingle).
		FieldOptions(files)
	formList.AddField("场景类型", "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "默认", Value: "1"},
			{Text: "比较", Value: "2"},
		}).FieldDefault("1")
	formList.AddField("优先级", "priority", db.Int, form.Number).FieldDefault("1")
	formList.AddField("执行次数", "run_time", db.Int, form.Number).FieldDefault("1")
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.TextArea)
	formList.AddField("备注", "remark1", db.Longtext, form.TextArea)
	formList.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("创建人", "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("playbook").SetTitle("场景列表").SetDescription("场景列表")

	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		apiList := values["api_list"]
		numList := values["data_number"]
		err = biz.UpdatePlaybookApiList(id, apiList, numList)
		return
	})

	detail := playbook.GetDetail()
	detail.AddField("唯一标识", "id", db.Int)
	detail.AddField("场景描述", "name", db.Varchar)
	detail.AddField("数据文件列表", "api_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return model.Value
		})
	detail.AddField("最近数据文件", "last_file", db.Varchar)
	detail.AddField("场景类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "默认"
			}
			if model.Value == "2" {
				return "比较"
			}
			return "默认"
		})
	detail.AddField("优先级", "priority", db.Int)
	detail.AddField("执行次数", "run_time", db.Int)
	detail.AddField("测试结果", "result", db.Varchar)
	detail.AddField("失败原因", "fail_reason", db.Longtext)
	detail.AddField("备注", "remark", db.Longtext)
	detail.AddField("所属产品", "product", db.Varchar)
	detail.AddField("创建人", "user_name", db.Varchar)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp)

	detail.SetTable("playbook").SetTitle("场景列表").SetDescription("场景列表")

	return playbook
}
