package tables

import (
	"data4test/biz"
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

	info.AddField(biz.T("common.id"), "id", db.Int)
	info.AddField(biz.T("common.create_desc"), "create_desc", db.Text)
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("ai_create").SetTitle(biz.T("ai_create.title")).SetDescription(biz.T("ai_create.description"))
	formList := aiCreate.GetForm()
	formList.HideContinueNewCheckBox()
	formList.HeadWidth = 1
	formList.InputWidth = 10
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.create_desc"), "create_desc", db.Text, form.TextArea)
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetPostHook(func(values form2.Values) (err error) {
		if _, ok := values["create_desc"]; !ok {
			err = fmt.Errorf(biz.T("ai_create.msg_input_create_desc"))
			return
		}

		return
	})

	formList.SetTable("ai_create").SetTitle(biz.T("ai_create.title")).SetDescription(biz.T("ai_create.description"))

	return aiCreate
}
