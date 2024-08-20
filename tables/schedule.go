package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"html/template"
	"strings"
)

func GetScheduleTable(ctx *context.Context) table.Table {

	schedule := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	userName := auth.Auth(ctx).Name

	info := schedule.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("任务名称", "task_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(160)
	info.AddField("任务模式", "task_mode", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "cron" {
				return "自定义"
			}
			if model.Value == "once" {
				return "一次"
			}
			if model.Value == "day" {
				return "每天"
			}
			if model.Value == "week" {
				return "每周"
			}
			return "一次"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "cron", Text: "自定义"},
		{Value: "once", Text: "一次"},
		{Value: "day", Text: "每天"},
		{Value: "week", Text: "每周"},
	}).FieldWidth(100)
	info.AddField("选择星期", "week", db.Varchar).
		FieldHide()
	info.AddField("选择时刻", "time4week", db.Varchar).
		FieldHide()
	info.AddField("选择时刻", "time4day", db.Varchar).
		FieldHide()
	info.AddField("Cron表达式", "crontab", db.Varchar).
		FieldHide()

	info.AddField("是否并发", "threading", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
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

	info.AddField("任务类型", "task_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "data" {
				return "数据"
			}
			if model.Value == "scene" {
				return "场景"
			}
			return "场景"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "data", Text: "数据"},
		{Value: "scene", Text: "场景"},
	}).FieldWidth(60)

	info.AddField("关联数据", "data_list", db.Longtext).
		FieldHide()
	info.AddField("关联场景", "scene_list", db.Longtext).
		FieldHide()
	info.AddField("关联产品", "product_list", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products).FieldWidth(120)

	info.AddField("任务状态", "task_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "running" {
				return "运行中"
			}
			if model.Value == "stopped" {
				return "已暂停"
			}
			if model.Value == "finished" {
				return "已结束"
			}
			if model.Value == "not_started" {
				return "未开始"
			}
			return "未开始"
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "running", Text: "运行中"},
		{Value: "stopped", Text: "已暂停"},
		{Value: "finished", Text: "已结束"},
		{Value: "not_started", Text: "未开始"},
	}).FieldWidth(100)

	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("上次执行时间", "last_at", db.Timestamp).FieldWidth(100).FieldSortable()
	info.AddField("下次执行时间", "next_at", db.Timestamp).FieldWidth(100).FieldSortable()
	info.AddField("创建人", "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()
	dataNames := biz.GetDatas()
	sceneNames := biz.GetScenes()
	timeNos := biz.Get24No()
	weekNos := biz.Get7No()

	info.AddButton("一键导出", icon.Android, action.Ajax("schedule_batch_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再一键导出"
				return false, status, ""
			}
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if fileName, err := biz.ExportSchedule(idStr, userNameSub); err == nil {
				status = fmt.Sprintf("一键导出成功，请至[文件-下载文件]下载, 文件名为: %s", fileName)
			} else {
				status = fmt.Sprintf("一键导出失败：%s: %s", idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton("一键导出", action.Ajax("schedule_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if fileName, err := biz.ExportSchedule(id, userNameSub); err == nil {
				status = fmt.Sprintf("一键导出成功，请至[文件-下载文件]下载, 文件名为:[%s]", fileName)
			} else {
				status = fmt.Sprintf("一键导出失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("复制", icon.Android, action.Ajax("schedule_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再复制"
				return false, status, ""
			}
			user := auth.Auth(ctx)
			userNameSub := user.Name
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.CopySchedule(id, userNameSub); err == nil {
					status = "复制成功，请刷新列表查看"
				} else {
					status = fmt.Sprintf("复制失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("复制", action.Ajax("schedule_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopySchedule(id, userNameSub); err == nil {
				status = "复制成功，请刷新列表查看"
			} else {
				status = fmt.Sprintf("复制失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("执行", icon.Android, action.Ajax("run_batch_task",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再执行"
				return false, status, ""
			}
			if err := biz.RunTask(idStr, ""); err == nil {
				status = "任务已在后台执行"
			} else {
				status = fmt.Sprintf("任务执行失败：%s", err)
			}
			return true, status, ""
		}))

	info.AddActionButton("执行", action.Ajax("run_task",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTask(id, ""); err == nil {
				status = "任务已在后台执行"
			} else {
				status = fmt.Sprintf("任务执行失败：%s", err)
			}

			return true, status, ""
		}))

	info.AddButton("暂停", icon.Android, action.Ajax("task_batch_stop",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再暂停"
				return false, status, ""
			}
			if err := biz.StopTask(idStr); err == nil {
				status = "任务已暂停"
			} else {
				status = fmt.Sprintf("任务暂停失败：%s", err)
				return false, status, ""
			}
			//ids := strings.Split(idStr, ",")
			//for _, id := range ids {
			//	if len(id) == 0 {
			//		status = "请先选择任务再暂停"
			//		continue
			//	}
			//	if err := biz.StopTask(id); err == nil {
			//		status = "任务已暂停"
			//	} else {
			//		status = fmt.Sprintf("任务暂停失败：%s", err)
			//		return false, status, ""
			//	}
			//}
			return true, status, ""
		}))

	info.AddActionButton("暂停", action.Ajax("task_stop",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.StopTask(id); err == nil {
				status = "任务已暂停"
			} else {
				status = fmt.Sprintf("任务暂停失败：%s", err)
			}
			return true, status, ""
		}))

	info.AddSelectBox("关联产品", products, action.FieldFilter("product_list"))

	info.SetTable("schedule").SetTitle("定时任务").SetDescription("定时任务")
	cronHelp := template.HTML("<br>" +
		"*&emsp;*&emsp;*&emsp;*&emsp;*" +
		"<br>" +
		"-&emsp;-&emsp;-&emsp;-&emsp;-" +
		"<br>" +
		"|&emsp;|&emsp;&ensp;|&emsp;|&emsp;|" +
		"<br>" +
		"|&emsp;|&emsp;&ensp;|&emsp;|&emsp;+-------------- 星期中星期几 (0 - 6) (星期天 为0)" +
		"<br>" +
		"|&emsp;|&emsp;&ensp;|&emsp;+---------------- 月份 (1 - 12) " +
		"<br>" +
		"|&emsp;|&emsp;&ensp;+------------------ 一个月中的第几天 (1 - 31)" +
		"<br>" +
		"|&emsp;+---------------------- 小时 (0 - 23)" +
		"<br>" +
		"+------------------------- 分钟 (0 - 59)")

	formList := schedule.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("任务名称", "task_name", db.Varchar, form.Text)
	//formList.AddField("任务模式", "task_mode", db.Enum, form.Radio).
	//	FieldOptions(types.FieldOptions{
	//	{Value: "cron", Text: "自定义"},
	//	{Value: "once", Text: "一次"},
	//	{Value: "day", Text: "每天"},
	//	{Value: "week", Text: "每周"},
	//	}).FieldDefault("once")

	formList.AddField("任务模式", "task_mode", db.Enum, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "cron", Text: "自定义"},
			{Value: "once", Text: "一次"},
			{Value: "day", Text: "每天"},
			{Value: "week", Text: "每周"},
		}).
		FieldOnChooseHide("once", "crontab", "week", "time4week", "time4day").
		FieldOnChooseHide("cron", "week", "time4week", "time4day").
		FieldOnChooseHide("week", "crontab", "time4day").
		FieldOnChooseHide("day", "crontab", "week", "time4week").
		FieldOnChooseShow("cron", "crontab").
		FieldOnChooseShow("week", "week", "time4week").
		FieldOnChooseShow("day", "time4day").
		FieldDefault("once")

	formList.AddField("选择星期", "week", db.Varchar, form.Select).
		FieldOptions(weekNos)
	formList.AddField("选择时刻", "time4day", db.Varchar, form.Select).
		FieldOptions(timeNos)
	formList.AddField("选择时刻", "time4week", db.Varchar, form.Select).
		FieldOptions(timeNos)
	formList.AddField("Cron表达式", "crontab", db.Varchar, form.Text).
		FieldHelpMsg(cronHelp)
	formList.AddField("是否并发", "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "是", Value: "yes"},
			{Text: "否", Value: "no"},
		}).FieldDefault("no")

	formList.AddField("任务类型", "task_type", db.Enum, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "data", Text: "数据"},
			{Value: "scene", Text: "场景"},
		}).
		FieldOnChooseHide("scene", "data_table").
		FieldOnChooseHide("data", "scene_table").
		FieldOnChooseShow("scene", "scene_table").
		FieldOnChooseShow("data", "data_table").
		FieldDefault("scene")

	dataHelp := template.HTML("关联数据必填")
	formList.AddTable("关联数据", "data_table", func(panel *types.FormPanel) {
		panel.AddField("序号/标签", "data_number", db.Varchar, form.Text).
			FieldHideLabel().
			FieldDisplay(func(model types.FieldModel) interface{} {
				return strings.Split(model.Value, ",")
			})
		panel.AddField("关联数据", "data_list", db.Varchar, form.SelectSingle).
			FieldHideLabel().
			FieldOptions(dataNames).
			FieldDisplay(func(model types.FieldModel) interface{} {
				return strings.Split(model.Value, ",")
			})
		panel.SetInputWidth(10)
	}).FieldHelpMsg(dataHelp)

	sceneHelp := template.HTML("关联场景必填")
	formList.AddTable("关联场景", "scene_table", func(panel *types.FormPanel) {
		panel.AddField("序号/标签", "scene_number", db.Varchar, form.Text).
			FieldHideLabel().
			FieldDisplay(func(model types.FieldModel) interface{} {
				return strings.Split(model.Value, ",")
			})

		panel.AddField("关联场景", "scene_list", db.Varchar, form.SelectSingle).
			FieldHideLabel().
			FieldOptions(sceneNames).
			FieldDisplay(func(model types.FieldModel) interface{} {
				return strings.Split(model.Value, ",")
			})
		panel.SetInputWidth(10)
	}).FieldHelpMsg(sceneHelp)

	formList.AddField("关联产品", "product_list", db.Varchar, form.Select).
		FieldOptions(products)

	formList.AddField("任务状态", "task_status", db.Enum, form.Text).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("上次执行时间", "last_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField("下次执行时间", "next_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField("创建人", "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("schedule").SetTitle("定时任务").SetDescription("定时任务")

	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		dataList := values["data_list"]
		dataNumList := values["data_number"]
		sceneList := values["scene_list"]
		sceneNumList := values["scene_number"]
		err = biz.UpdateTaskInfoList(id, dataList, dataNumList, sceneList, sceneNumList)
		return
	})

	detail := schedule.GetDetail()
	detail.AddField("自增主键", "id", db.Int)
	detail.AddField("任务名称", "task_name", db.Varchar)
	detail.AddField("任务模式", "task_mode", db.Enum)
	detail.AddField("每周", "week", db.Varchar)
	detail.AddField("每天时刻", "time4day", db.Varchar)
	detail.AddField("每周时刻", "time4week", db.Varchar)
	detail.AddField("Cron表达式", "crontab", db.Varchar)
	detail.AddField("是否并发", "threading", db.Enum)
	detail.AddField("任务类型", "task_type", db.Enum)
	detail.AddField("关联数据", "data_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
		})
	detail.AddField("关联场景", "scene_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
		})
	detail.AddField("关联产品", "product_list", db.Varchar)
	detail.AddField("任务状态", "task_status", db.Enum)
	detail.AddField("备注", "remark", db.Longtext)
	detail.AddField("上次执行时间", "last_at", db.Timestamp)
	detail.AddField("下次执行时间", "next_at", db.Timestamp)
	detail.AddField("创建人", "user_name", db.Text)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp)
	detail.SetTable("schedule").SetTitle("定时任务").SetDescription("定时任务")

	return schedule
}
