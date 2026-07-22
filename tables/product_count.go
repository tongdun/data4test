package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetProductCountTable(ctx *context.Context) table.Table {

	productCount := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := productCount.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("product_count.all_count"), "all_count", db.Int)
	info.AddField(biz.T("product_count.automatable_count"), "automatable_count", db.Int)
	info.AddField(biz.T("product_count.unautomatable_count"), "unautomatable_count", db.Int)
	info.AddField(biz.T("product_count.auto_test_count"), "auto_test_count", db.Int)
	info.AddField(biz.T("product_count.untest_count"), "untest_count", db.Int)
	info.AddField(biz.T("product_count.pass_count"), "pass_count", db.Int)
	info.AddField(biz.T("common.fail_count"), "fail_count", db.Int)
	info.AddField(biz.T("product_count.auto_per"), "auto_per", db.Double)
	info.AddField(biz.T("product_count.pass_per"), "pass_per", db.Double)
	info.AddField(biz.T("product_count.fail_per"), "fail_per", db.Double)
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), " app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	products := biz.GetProducts()
	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))
	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.SetTable("product_count").SetTitle(biz.T("product_count.title")).SetDescription(biz.T("product_count.description"))

	formList := productCount.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("product_count.all_count"), "all_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.automatable_count"), "automatable_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.unautomatable_count"), "unautomatable_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.auto_test_count"), "auto_test_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.untest_count"), "untest_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.pass_count"), "pass_count", db.Int, form.Number)
	formList.AddField(biz.T("common.fail_count"), "fail_count", db.Int, form.Number)
	formList.AddField(biz.T("product_count.auto_per"), "auto_per", db.Double, form.Text)
	formList.AddField(biz.T("product_count.pass_per"), "pass_per", db.Double, form.Text)
	formList.AddField(biz.T("product_count.fail_per"), "fail_per", db.Double, form.Text)
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.Text)
	formList.AddField(biz.T("common.app"), " app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("product_count").SetTitle(biz.T("product_count.title")).SetDescription(biz.T("product_count.description"))

	return productCount
}
