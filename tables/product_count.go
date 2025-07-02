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
	info.AddField("唯一标识", "id", db.Int).
		FieldHide()
	info.AddField("接口总数", "all_count", db.Int)
	info.AddField("可自动化数", "automatable_count", db.Int)
	info.AddField("不可自动化数", "unautomatable_count", db.Int)
	info.AddField("自动化测试总数", "auto_test_count", db.Int)
	info.AddField("未测试总数", "untest_count", db.Int)
	info.AddField("通过总数", "pass_count", db.Int)
	info.AddField("失败总数", "fail_count", db.Int)
	info.AddField("自动化率", "auto_per", db.Double)
	info.AddField("通过率", "pass_per", db.Double)
	info.AddField("失败率", "fail_per", db.Double)
	info.AddField("关联产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("关联应用", " app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	products := biz.GetProducts()
	info.AddSelectBox("关联产品", products, action.FieldFilter("product"))
	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.SetTable("product_count").SetTitle("产品统计").SetDescription("产品统计")

	formList := productCount.GetForm()
	formList.AddField("唯一标识", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("接口总数", "all_count", db.Int, form.Number)
	formList.AddField("可自动化数", "automatable_count", db.Int, form.Number)
	formList.AddField("不可自动化数", "unautomatable_count", db.Int, form.Number)
	formList.AddField("自动化测试总数", "auto_test_count", db.Int, form.Number)
	formList.AddField("未测试总数", "untest_count", db.Int, form.Number)
	formList.AddField("通过总数", "pass_count", db.Int, form.Number)
	formList.AddField("失败总数", "fail_count", db.Int, form.Number)
	formList.AddField("自动化率", "auto_per", db.Double, form.Text)
	formList.AddField("通过率", "pass_per", db.Double, form.Text)
	formList.AddField("失败率", "fail_per", db.Double, form.Text)
	formList.AddField("关联产品", "product", db.Varchar, form.Text)
	formList.AddField("关联应用", " app", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("product_count").SetTitle("产品统计").SetDescription("产品统计")

	return productCount
}
