package tables

import (
	"data4perf/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAppApiChangelogTable(ctx *context.Context) table.Table {

	appApiChangelog := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := appApiChangelog.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldWidth(120)
	info.AddField("当前接口总数", "curApiSum", db.Int)
	info.AddField("已存在接口数", "existApiSum", db.Int)
	info.AddField("新增接口数", "newApiSum", db.Int)
	info.AddField("删除接口数", "deletedApiSum", db.Int)
	info.AddField("变更接口数", "changedApiSum", db.Int)
	info.AddField("规范检查失败接口数", "checkFailApiSum", db.Int)
	info.AddField("新增接口详情", "newApiContent", db.Longtext).
		FieldHide()
	info.AddField("删除接口详情", "deletedApiContent", db.Longtext).
		FieldHide()
	info.AddField("变更接口详情", "changedApiContent", db.Longtext).
		FieldHide()
	info.AddField("规范检查失败接口详情", "apiCheckFailContent", db.Longtext).
		FieldHide()
	info.AddField("接口检查结果", "apiCheckResult", db.Varchar).
		FieldHide()
	info.AddField("版本分支", "branch", db.Varchar)
	info.AddField("备注", "remark", db.Varchar)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("app_api_changelog").SetTitle("接口记录表").SetDescription("接口记录表")

	formList := appApiChangelog.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("所属应用", "app", db.Varchar, form.Text)
	formList.AddField("当前接口总数", "curApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("已存在接口数", "existApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("新增接口数", "newApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("删除接口数", "deletedApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("变更接口数", "changedApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("规范检查失败接口数", "checkFailApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField("新增接口详情", "newApiContent", db.Longtext, form.RichText)
	formList.AddField("删除接口详情", "deletedApiContent", db.Longtext, form.RichText)
	formList.AddField("变更接口详情", "changedApiContent", db.Longtext, form.RichText)
	formList.AddField("规范检查失败接口详情", "apiCheckFailContent", db.Longtext, form.RichText)
	formList.AddField("接口检查结果", "apiCheckResult", db.Varchar, form.Text)
	formList.AddField("版本分支", "branch", db.Varchar, form.Text)
	formList.AddField("备注", "remark", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("app_api_changelog").SetTitle("接口记录表").SetDescription("接口记录表")

	detail := appApiChangelog.GetDetail()
	detail.AddField("自增主键", "id", db.Int)
	detail.AddField("所属应用", "app", db.Varchar)
	detail.AddField("当前接口总数", "curApiSum", db.Int)
	detail.AddField("已存在接口数", "existApiSum", db.Int)
	detail.AddField("新增接口数", "newApiSum", db.Int)
	detail.AddField("删除接口数", "deletedApiSum", db.Int)
	detail.AddField("变更接口数", "changedApiSum", db.Int)
	detail.AddField("规范检查失败接口数", "checkFailApiSum", db.Int)
	detail.AddField("新增接口详情", "newApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField("删除接口详情", "deletedApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField("变更接口详情", "changedApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField("规范检查失败接口详情", "apiCheckFailContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiRaw(model.Row["app"], model.Value)
		})
	detail.AddField("接口检查结果", "apiCheckResult", db.Varchar)
	detail.AddField("版本分支", "branch", db.Varchar)
	detail.AddField("备注", "remark", db.Varchar)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	return appApiChangelog
}
