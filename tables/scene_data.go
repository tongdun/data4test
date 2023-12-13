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

func GetSceneDataTable(ctx *context.Context) table.Table {

	sceneData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	apps := biz.GetApps()
	info := sceneData.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("唯一标识", "id", db.Int).
		FieldFilterable()
	info.AddField("数据描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(apps)
	info.AddField("文件名", "file_name", db.Longtext).FieldDisplay(func(value types.FieldModel) interface{} {
		return template.Default().
			Link().
			SetURL("/admin/fm/data/preview?path=/" + value.Value).
			SetContent(template2.HTML(value.Value)).
			OpenInNewTab().
			SetTabTitle(template.HTML("数据文件")).
			GetContent()
	}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("文件内容", "content", db.Longtext).FieldHide()
	info.AddField("执行次数", "run_time", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).FieldSortable().
		FieldEditAble(editType.Text)
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField("失败原因", "fail_reason", db.Longtext).FieldWidth(120)
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("创建人", "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("同步", icon.Android, action.Ajax("scenedata_sync",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			var status string
			if newTag, modTag, err := biz.SyncSceneData(); err == nil {
				status = fmt.Sprintf("同步完成, 新增: %d条，修改: %d条", newTag, modTag)
			} else {
				status = fmt.Sprintf("同步失败: %s, 新增: %d条，修改: %d条", err, newTag, modTag)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddButton("复制", icon.Android, action.Ajax("scenedata_batch_copy",
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
				if err := biz.CopySceneData(id, userNameSub); err == nil {
					status = "复制成功，请刷新列表查看"
				} else {
					status = fmt.Sprintf("复制失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))
	info.AddActionButton("复制", action.Ajax("scenedata_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopySceneData(id, userNameSub); err == nil {
				status = "复制成功，请刷新列表查看"
			} else {
				status = fmt.Sprintf("复制失败：%s: %s", id, err)
			}
			return true, status, ""
		}))
	info.AddButton("测试任务", icon.Android, action.Ajax("scenedata_batch_task_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}

			err := biz.AutoCreateSchedule(idStr, "data")
			if err != nil {
				status = "发起测试任务失败"
				return false, status, fmt.Sprintf("%s", err)
			}
			//status = "已开启任务执行，请稍后查看测试结果，可在任务列表管理任务"
			status = "已创建任务，请前往任务列表查看，确认后再发起任务执行"
			return true, status, ""
		}))

	info.AddButton("测试", icon.Android, action.Ajax("scenedata_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "测试完成，请刷新列表查看测试结果"
					continue
				}

				if err := biz.RunSceneData(id, ""); err == nil {
					status = "测试完成，请刷新列表查看测试结果"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("测试", action.Ajax("scenedata_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunSceneData(id, ""); err == nil {
				status = "测试完成，请刷新列表查看测试结果"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddSelectBox("所属应用", apps, action.FieldFilter("app"))
	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))
	info.SetTable("scene_data").SetTitle("数据列表").SetDescription("数据列表")

	fileNameHelp := template.HTML("e.g.: 类型-模块-功能描述.yml / 类型-模块-功能描述.json")

	formList := sceneData.GetForm()
	formList.AddField("唯一标识", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("数据描述", "name", db.Varchar, form.Text)
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("所属应用", "app", db.Varchar, form.SelectSingle).
		FieldOptions(apps)
	formList.AddField("文件名", "file_name", db.Longtext, form.Url).
		FieldHelpMsg(fileNameHelp)

	formList.AddField("文件内容", "content", db.Longtext, form.TextArea).FieldDefault("name: \"\"\nversion: 1\napi_id: \"\"\nis_run_pre_apis: \"no\"\nis_run_post_apis: \"no\"\nis_parallel: \"no\"\nis_use_env_config: \"yes\"\nenv:\n  protocol: http\n  host: \"\"\n  prepath: \"\"\napi:\n  description: \"\"\n  module: \"\"\n  app: \"\"\n  method: \"\"\n  path: \"\"\n  pre_apis: []\n  param_apis: []\n  post_apis: []\nsingle:\n  header: {}\n  query: {}\n  path: {}\n  body: {}\nmulti:\n  query: {}\n  path: {}\n  body: {}\nassert: []\noutput: {}\ntest_result: []\nurls: []\nrequest: []\nresponse: []\n")
	formList.AddField("执行次数", "run_time", db.Int, form.Number).FieldDefault("1")
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.TextArea)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建人", "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenCreate().FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.SetTable("scene_data").SetTitle("数据列表").SetDescription("数据列表")

	formList.SetPostHook(func(values form2.Values) (err error) {
		content := values["content"][0]
		fileName := values["file_name"][0]
		var newTxt string
		if strings.HasSuffix(fileName, ".yml") {
			newTxt = strings.ReplaceAll(content, "<br>", "\n")
		} else {
			newTxt = content
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(newTxt))
		if err != nil {
			biz.Logger.Error("%s", err)
			return
		}
		handle := doc.Find("code")
		afterTxt := handle.Text()
		if len(afterTxt) == 0 {
			afterTxt = newTxt
		}

		id := values["id"][0]
		err = biz.BakOldVer(id, afterTxt, fileName)
		return
	})

	return sceneData
}
