package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
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
	partProducts := biz.GetProductsByUpdateTime(1)
	userName := auth.Auth(ctx).Name

	info := schedule.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.task_name"), "task_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(160)
	info.AddField(biz.T("common.task_mode"), "task_mode", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "cron" {
				return biz.T("schedule.task_mode_cron")
			}
			if model.Value == "once" {
				return biz.T("schedule.task_mode_once")
			}
			if model.Value == "day" {
				return biz.T("schedule.task_mode_day")
			}
			if model.Value == "week" {
				return biz.T("schedule.task_mode_week")
			}
			return biz.T("schedule.task_mode_once")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "cron", Text: biz.T("schedule.task_mode_cron")},
		{Value: "once", Text: biz.T("schedule.task_mode_once")},
		{Value: "day", Text: biz.T("schedule.task_mode_day")},
		{Value: "week", Text: biz.T("schedule.task_mode_week")},
	}).FieldWidth(100)
	info.AddField(biz.T("schedule.select_week"), "week", db.Varchar).
		FieldHide()
	info.AddField(biz.T("schedule.select_time4week"), "time4week", db.Varchar).
		FieldHide()
	info.AddField(biz.T("schedule.select_time4day"), "time4day", db.Varchar).
		FieldHide()
	info.AddField(biz.T("schedule.crontab"), "crontab", db.Varchar).
		FieldHide()

	info.AddField(biz.T("common.threading"), "threading", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "yes" {
				return biz.T("common.yes")
			}
			if model.Value == "no" {
				return biz.T("common.no")
			}
			return biz.T("common.no")
		}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	})

	info.AddField(biz.T("common.task_type"), "task_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "data" {
				return biz.T("common.data")
			}
			if model.Value == "scene" {
				return biz.T("common.scene")
			}
			return biz.T("common.scene")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "data", Text: biz.T("common.data")},
		{Value: "scene", Text: biz.T("common.scene")},
	}).FieldWidth(60)

	info.AddField(biz.T("common.data_list"), "data_list", db.Longtext).
		FieldHide()
	info.AddField(biz.T("schedule.scene_list"), "scene_list", db.Longtext).
		FieldHide()
	info.AddField(biz.T("common.product_list"), "product_list", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldFilterOptions(products).
		FieldEditAble(editType.Select).
		FieldEditOptions(partProducts).
		FieldWidth(220)

	info.AddField(biz.T("common.task_status"), "task_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "running" {
				return biz.T("common.running")
			}
			if model.Value == "stopped" {
				return biz.T("common.stopped")
			}
			if model.Value == "finished" {
				return biz.T("common.finished")
			}
			if model.Value == "not_started" {
				return biz.T("common.not_started")
			}
			return biz.T("common.not_started")
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "running", Text: biz.T("common.running")},
		{Value: "stopped", Text: biz.T("common.stopped")},
		{Value: "finished", Text: biz.T("common.finished")},
		{Value: "not_started", Text: biz.T("common.not_started")},
	}).FieldWidth(100)

	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField(biz.T("schedule.last_at"), "last_at", db.Timestamp).FieldWidth(100).FieldSortable()
	info.AddField(biz.T("schedule.next_at"), "next_at", db.Timestamp).FieldWidth(100).FieldSortable()
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()
	//dataNames := biz.GetDatas()
	//sceneNames := biz.GetScenes()
	timeNos := biz.Get24No()
	weekNos := biz.Get7No()

	sBtnBatchExport := template2.HTML(biz.T("common.btn_export"))
	info.AddButton(sBtnBatchExport, icon.Android, action.Ajax("schedule_batch_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("schedule.export_select_first")
				return false, status, ""
			}
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if fileName, err := biz.ExportSchedule(idStr, userNameSub); err == nil {
				hostIp := ctx.Request.Host
				var fileNameList []string
				var downloadUrl string
				fileNameList = strings.Split(fileName, ",")

				for index, subFileName := range fileNameList {
					if index == 0 {
						downloadUrl = fmt.Sprintf("http://%s/admin/fm/download/download?path=/%s", hostIp, subFileName)
					} else {
						downloadUrl = fmt.Sprintf("%s\nhttp://%s/admin/fm/download/download?path=/%s", downloadUrl, hostIp, subFileName)
					}
				}

				status = fmt.Sprintf("%s\n%s:\n%s", biz.T("schedule.export_success"), biz.T("common.copy_download_link"), downloadUrl)

			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("schedule.export_fail"), idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_export")), action.Ajax("schedule_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if fileName, err := biz.ExportSchedule(id, userNameSub); err == nil {
				hostIp := ctx.Request.Host
				var fileNameList []string
				var downloadUrl string
				fileNameList = strings.Split(fileName, ",")

				for index, subFileNale := range fileNameList {
					if index == 0 {
						downloadUrl = fmt.Sprintf("http://%s/admin/fm/download/download?path=/%s", hostIp, subFileNale)
					} else {
						downloadUrl = fmt.Sprintf("%s\nhttp://%s/admin/fm/download/download?path=/%s", downloadUrl, hostIp, subFileNale)
					}
				}

				status = fmt.Sprintf("%s\n%s:\n%s", biz.T("schedule.export_success"), biz.T("schedule.copy_download_link"), downloadUrl)
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("schedule.export_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("schedule_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("schedule.copy_select_first")
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
					status = biz.T("schedule.copy_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_copy")), action.Ajax("schedule_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopySchedule(id, userNameSub); err == nil {
				status = biz.T("schedule.copy_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("common.btn_run")), icon.Android, action.Ajax("run_batch_task",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			user := auth.Auth(ctx)
			userNameSub := user.Name
			var status string
			if idStr == "," {
				status = biz.T("schedule.run_select_first")
				return false, status, ""
			}
			if err := biz.RunTask(idStr, "", userNameSub); err == nil {
				status = biz.T("schedule.task_running_background")
			} else {
				status = fmt.Sprintf("%s: %s", biz.T("schedule.task_exec_fail"), err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_run")), action.Ajax("run_task",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.RunTask(id, "", userNameSub); err == nil {
				status = biz.T("schedule.task_running_background")
			} else {
				status = fmt.Sprintf("%s: %s", biz.T("schedule.task_exec_fail"), err)
			}

			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("common.btn_stop")), icon.Android, action.Ajax("task_batch_stop",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("schedule.stop_select_first")
				return false, status, ""
			}
			if err := biz.StopTask(idStr); err == nil {
				status = biz.T("schedule.task_stopped_msg")
			} else {
				status = fmt.Sprintf("%s: %s", biz.T("schedule.task_pause_fail"), err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_stop")), action.Ajax("task_stop",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.StopTask(id); err == nil {
				status = biz.T("schedule.task_stopped_msg")
			} else {
				status = fmt.Sprintf("%s: %s", biz.T("schedule.task_pause_fail"), err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_report")), action.Jump("/admin/task_dashboard?id={{.Id}}"))

	info.AddButton(template2.HTML(biz.T("schedule_report.btn_generate")), icon.Android, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/generate_task_report",
		Title:  biz.T("schedule_report.generate_report"),
		Width:  "900px",
		Height: "540px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		// 根据选中的任务ID自动关联产品信息作为默认值
		defaultProducts := biz.GetScheduleProductListByIds(ids)
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate()
		panel.AddField(biz.T("schedule_report.select_products"), "report_products", db.Varchar, form.Select).
			FieldOptions(products).
			FieldDefault(defaultProducts).
			FieldHelpMsg(template2.HTML(biz.T("schedule_report.select_report_type")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/generate_task_report"))

	info.AddSelectBox(biz.T("common.product_list"), products, action.FieldFilter("product_list"))

	info.SetTable("schedule").SetTitle(biz.T("schedule.title")).SetDescription(biz.T("schedule.description"))
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
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.task_name"), "task_name", db.Varchar, form.Text)

	formList.AddField(biz.T("common.task_mode"), "task_mode", db.Enum, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "cron", Text: biz.T("schedule.task_mode_cron")},
			{Value: "once", Text: biz.T("schedule.task_mode_once")},
			{Value: "day", Text: biz.T("schedule.task_mode_day")},
			{Value: "week", Text: biz.T("schedule.task_mode_week")},
		}).
		FieldOnChooseHide("once", "crontab", "week", "time4week", "time4day").
		FieldOnChooseHide("cron", "week", "time4week", "time4day").
		FieldOnChooseHide("week", "crontab", "time4day").
		FieldOnChooseHide("day", "crontab", "week", "time4week").
		FieldOnChooseShow("cron", "crontab").
		FieldOnChooseShow("week", "week", "time4week").
		FieldOnChooseShow("day", "time4day").
		FieldDefault("once")

	formList.AddField(biz.T("schedule.select_week"), "week", db.Varchar, form.Select).
		FieldOptions(weekNos)
	formList.AddField(biz.T("schedule.select_time4day"), "time4day", db.Varchar, form.Select).
		FieldOptions(timeNos)
	formList.AddField(biz.T("schedule.time4day"), "time4week", db.Varchar, form.Select).
		FieldOptions(timeNos)
	formList.AddField(biz.T("schedule.crontab"), "crontab", db.Varchar, form.Text).
		FieldHelpMsg(cronHelp)
	formList.AddField(biz.T("common.threading"), "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.yes"), Value: "yes"},
			{Text: biz.T("common.no"), Value: "no"},
		}).FieldDefault("no")

	formList.AddField(biz.T("common.edit_type"), "edit_type", db.Enum, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			//{Value: "select_scene", Text: "选择场景"},
			{Value: "input_scene", Text: biz.T("schedule.input_scene")},
			//{Value: "select_data", Text: "选择数据"},
			{Value: "input_data", Text: biz.T("schedule.input_data")},
		}).
		FieldOnChooseHide("input_scene", "scene_table", "data_table", "input_data").
		//FieldOnChooseHide("select_scene", "input_data", "data_table", "input_scene").
		//FieldOnChooseHide("select_data", "input_scene", "scene_table", "data_table").
		FieldOnChooseHide("input_data", "scene_table", "data_table", "input_scene").
		//FieldOnChooseShow("select_scene", "scene_table").
		FieldOnChooseShow("input_scene", "input_scene").
		//FieldOnChooseShow("select_data", "data_table").
		FieldOnChooseShow("input_data", "input_data").
		FieldDefault("input_scene").
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetTaskEditTypeById(model.ID)
		})

	dataHelp := template2.HTML(biz.T("common.help_data_required"))
	formList.AddField(biz.T("schedule.input_data"), "input_data", db.Varchar, form.TextArea).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetTaskDataStr(model.ID)
		}).FieldHelpMsg(dataHelp)

	//formList.AddTable("关联数据", "data_table", func(panel *types.FormPanel) {
	//	panel.AddField("序号/标签", "data_number", db.Varchar, form.Text).
	//		FieldHideLabel().
	//		FieldDisplay(func(model types.FieldModel) interface{} {
	//			return strings.Split(model.Value, ",")
	//		})
	//	panel.AddField("关联数据", "data_list", db.Varchar, form.SelectSingle).
	//		FieldHideLabel().
	//		FieldOptions(dataNames).
	//		FieldDisplay(func(model types.FieldModel) interface{} {
	//			return strings.Split(model.Value, ",")
	//		})
	//	panel.SetInputWidth(10)
	//}).FieldHelpMsg(dataHelp)

	sceneHelp := template2.HTML(biz.T("schedule.scene_help"))
	formList.AddField(biz.T("schedule.input_scene"), "input_scene", db.Varchar, form.TextArea).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetTaskDataStr(model.ID)
		}).FieldHelpMsg(sceneHelp)

	//formList.AddTable("关联场景", "scene_table", func(panel *types.FormPanel) {
	//	panel.AddField("序号/标签", "scene_number", db.Varchar, form.Text).
	//		FieldHideLabel().
	//		FieldDisplay(func(model types.FieldModel) interface{} {
	//			return strings.Split(model.Value, ",")
	//		})
	//
	//	panel.AddField("关联场景", "scene_list", db.Varchar, form.SelectSingle).
	//		FieldHideLabel().
	//		FieldOptions(sceneNames).
	//		FieldDisplay(func(model types.FieldModel) interface{} {
	//			return strings.Split(model.Value, ",")
	//		})
	//	panel.SetInputWidth(10)
	//}).FieldHelpMsg(sceneHelp)

	formList.AddField(biz.T("common.product_list"), "product_list", db.Varchar, form.Select).
		FieldOptions(products)

	formList.AddField(biz.T("common.task_status"), "task_status", db.Enum, form.Text).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("schedule.last_at"), "last_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("schedule.next_at"), "next_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("schedule").SetTitle(biz.T("schedule.title")).SetDescription(biz.T("schedule.description"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		if _, ok := values["edit_type"]; !ok {
			return
		}
		id := values["id"][0]
		edit_Type := values["edit_type"][0]
		var dataList, dataNumList []string
		var taskType string
		if edit_Type == "select_data" {
			dataList = values["data_list"]
			dataNumList = values["data_number"]
			taskType = "data"
		} else if edit_Type == "select_scene" {
			dataList = values["scene_list"]
			dataNumList = values["scene_number"]
			taskType = "scene"
		} else if edit_Type == "input_data" {
			dataInputInfo := values["input_data"][0]
			dataList = strings.Split(dataInputInfo, "\r\n")
			taskType = "data"
		} else {
			sceneInputInfo := values["input_scene"][0]
			dataList = strings.Split(sceneInputInfo, "\r\n")
			taskType = "scene"
		}

		err = biz.UpdateTaskInfoList(id, taskType, dataList, dataNumList)
		return
	})

	detail := schedule.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.task_name"), "task_name", db.Varchar)
	detail.AddField(biz.T("common.task_mode"), "task_mode", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "cron" {
				return biz.T("schedule.task_mode_cron")
			}
			if model.Value == "once" {
				return biz.T("schedule.task_mode_once")
			}
			if model.Value == "day" {
				return biz.T("schedule.task_mode_day")
			}
			if model.Value == "week" {
				return biz.T("schedule.task_mode_week")
			}
			return biz.T("schedule.task_mode_once")
		})
	detail.AddField(biz.T("schedule.week"), "week", db.Varchar)
	detail.AddField(biz.T("schedule.time4day"), "time4day", db.Varchar)
	detail.AddField(biz.T("schedule.time4week"), "time4week", db.Varchar)
	detail.AddField(biz.T("schedule.crontab"), "crontab", db.Varchar)
	detail.AddField(biz.T("common.threading"), "threading", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "yes" {
				return biz.T("common.yes")
			}
			if model.Value == "no" {
				return biz.T("common.no")
			}
			return biz.T("common.no")
		})
	detail.AddField(biz.T("common.task_type"), "task_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "data" {
				return biz.T("common.data")
			}
			if model.Value == "scene" {
				return biz.T("common.scene")
			}
			return biz.T("common.scene")
		})
	detail.AddField(biz.T("common.data_list"), "data_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetDataDetailLinkByDataStr(model.Value)
		})
	detail.AddField(biz.T("schedule.scene_list"), "scene_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetPlaybookLinkByPlaybookStr(model.Value)

		})
	detail.AddField(biz.T("common.product_list"), "product_list", db.Varchar)
	detail.AddField(biz.T("common.task_status"), "task_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "running" {
				return biz.T("common.running")
			}
			if model.Value == "stopped" {
				return biz.T("common.stopped")
			}
			if model.Value == "finished" {
				return biz.T("common.finished")
			}
			if model.Value == "not_started" {
				return biz.T("common.not_started")
			}

			return biz.T("common.not_started")
		})
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("schedule.last_at"), "last_at", db.Timestamp)
	detail.AddField(biz.T("schedule.next_at"), "next_at", db.Timestamp)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Text)
	//detail.AddField(biz.T("common.btn_report"), "id", db.Int).
	//	FieldDisplay(func(model types.FieldModel) interface{} {
	//		return biz.GetScheduleExecutionHistoryLink(model.Value, biz.T("common.btn_report")+" >")
	//	})
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).FieldHide()
	detail.SetTable("schedule").SetTitle(biz.T("schedule.title")).SetDescription(biz.T("schedule.detail_title"))

	return schedule
}
