package tables

import (
	"data4test/biz"
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

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldWidth(120)
	info.AddField(biz.T("app_api_changelog.cur_api_sum"), "curApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.exist_api_sum"), "existApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.added_api"), "newApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.deleted_api"), "deletedApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.changed_api"), "changedApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.check_fail_sum"), "checkFailApiSum", db.Int)
	info.AddField(biz.T("app_api_changelog.add_detail"), "newApiContent", db.Longtext).
		FieldHide()
	info.AddField(biz.T("app_api_changelog.delete_detail"), "deletedApiContent", db.Longtext).
		FieldHide()
	info.AddField(biz.T("app_api_changelog.change_detail"), "changedApiContent", db.Longtext).
		FieldHide()
	info.AddField(biz.T("app_api_changelog.check_fail_detail"), "apiCheckFailContent", db.Longtext).
		FieldHide()
	info.AddField(biz.T("app_api_changelog.check_result"), "apiCheckResult", db.Varchar).
		FieldHide()
	info.AddField(biz.T("app_api_changelog.branch"), "branch", db.Varchar)
	info.AddField(biz.T("common.remark"), "remark", db.Varchar)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("app_api_changelog").SetTitle(biz.T("app_api_changelog.title")).SetDescription(biz.T("app_api_changelog.description"))

	formList := appApiChangelog.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("app_api_changelog.cur_api_sum"), "curApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.exist_api_sum"), "existApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.added_api"), "newApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.deleted_api"), "deletedApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.changed_api"), "changedApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.check_fail_sum"), "checkFailApiSum", db.Int, form.Number).FieldDefault("0")
	formList.AddField(biz.T("app_api_changelog.add_detail"), "newApiContent", db.Longtext, form.RichText)
	formList.AddField(biz.T("app_api_changelog.delete_detail"), "deletedApiContent", db.Longtext, form.RichText)
	formList.AddField(biz.T("app_api_changelog.change_detail"), "changedApiContent", db.Longtext, form.RichText)
	formList.AddField(biz.T("app_api_changelog.check_fail_detail"), "apiCheckFailContent", db.Longtext, form.RichText)
	formList.AddField(biz.T("app_api_changelog.check_result"), "apiCheckResult", db.Varchar, form.Text)
	formList.AddField(biz.T("app_api_changelog.branch"), "branch", db.Varchar, form.Text)
	formList.AddField(biz.T("common.remark"), "remark", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("app_api_changelog").SetTitle(biz.T("app_api_changelog.title")).SetDescription(biz.T("app_api_changelog.description"))

	detail := appApiChangelog.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.app"), "app", db.Varchar)
	detail.AddField(biz.T("app_api_changelog.cur_api_sum"), "curApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.exist_api_sum"), "existApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.added_api"), "newApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.deleted_api"), "deletedApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.changed_api"), "changedApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.check_fail_sum"), "checkFailApiSum", db.Int)
	detail.AddField(biz.T("app_api_changelog.add_detail"), "newApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField(biz.T("app_api_changelog.delete_detail"), "deletedApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField(biz.T("app_api_changelog.change_detail"), "changedApiContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiStr(model.Row["app"], model.Value)
		})
	detail.AddField(biz.T("app_api_changelog.check_fail_detail"), "apiCheckFailContent", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiDetailLinkByApiRaw(model.Row["app"], model.Value)
		})
	detail.AddField(biz.T("app_api_changelog.check_result"), "apiCheckResult", db.Varchar)
	detail.AddField(biz.T("app_api_changelog.branch"), "branch", db.Varchar)
	detail.AddField(biz.T("common.remark"), "remark", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	return appApiChangelog
}
