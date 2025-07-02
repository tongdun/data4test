package tables

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiCreateTable(ctx *context.Context) table.Table {

	aiCreate := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiCreate.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.UserName

	info.AddField("自增主键", "id", db.Int)
	info.AddField("生成指令", "create_desc", db.Text)
	info.AddField("创建人", "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("ai_create").SetTitle("生成指令").SetDescription("AiCreate")
	formList := aiCreate.GetForm()
	formList.HideContinueNewCheckBox()
	formList.HeadWidth = 1
	formList.InputWidth = 10
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("生成指令", "create_desc", db.Text, form.TextArea)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	//	formList.SetFooterHtml(`
	//<script>
	//document.addEventListener("DOMContentLoaded", function () {
	//    // 获取保存按钮
	//    const saveBtn = document.querySelector("button[type='submit']");
	//
	//    if (saveBtn) {
	//        // 修改按钮文字
	//        saveBtn.innerText = "生成";
	//    }
	//});
	//</script>`)

	formList.SetPostHook(func(values form2.Values) (err error) {
		if _, ok := values["create_desc"]; !ok {
			err = fmt.Errorf("请输入生成指令~")
			return
		}

		//err = ai_biz.ConnetAIModel(values["create_desc"][0], userName)
		return
	})

	formList.SetTable("ai_create").SetTitle("生成指令").SetDescription("AiCreate")

	return aiCreate
}
