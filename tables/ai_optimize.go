package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiOptimizeTable(ctx *context.Context) table.Table {

	aiOptimize := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiOptimize.GetInfo().HideFilterArea()

	info.AddField(biz.T("common.id"), "id", db.Int)
	info.AddField(biz.T("common.optimize_desc"), "optimize_desc", db.Text)
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp)

	info.SetTable("ai_optimize").SetTitle(biz.T("common.btn_ai_optimize")).SetDescription(biz.T("common.btn_ai_optimize"))

	formList := aiOptimize.GetForm().HideContinueEditCheckBox()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.optimize_desc"), "optimize_desc", db.Text, form.TextArea)
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).FieldHide()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldHide()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldHide()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).FieldHide()

	//formList.SetFooterHtml(`
	//<script>
	//document.addEventListener("DOMContentLoaded", function () {
	//    // 获取保存按钮
	//    const saveBtn = document.querySelector("button[type='submit']");
	//
	//    if (saveBtn) {
	//        // 修改按钮文字
	//        saveBtn.innerText = "` + biz.T("common.btn_ai_optimize") + `";
	//
	//        // 解绑原有事件（防止表单提交）
	//        const newBtn = saveBtn.cloneNode(true);
	//        saveBtn.parentNode.replaceChild(newBtn, saveBtn);
	//
	//        // 添加自定义点击逻辑
	//        newBtn.addEventListener("click", function () {
	//            // 你可以加 loading、禁用状态等
	//            newBtn.innerText = "` + biz.T("ai_optimize.btn_optimizing") + `...";
	//
	//            fetch("/your/custom/api", {
	//                method: "POST",
	//                headers: {
	//                    "Content-Type": "application/json"
	//                },
	//                body: JSON.stringify({ something: "example" })
	//            })
	//            .then(res => res.json())
	//            .then(res => {
	//                newBtn.innerText = "` + biz.T("common.finished") + `";
	//                alert("` + biz.T("common.operate_success") + `" + res.msg);
	//            })
	//            .catch(err => {
	//                newBtn.innerText = "` + biz.T("ai_optimize.btn_generate") + `";
	//                alert("` + biz.T("common.operate_fail") + `" + err);
	//            });
	//        });
	//    }
	//});
	//</script>`)

	formList.SetTable("ai_optimize").SetTitle(biz.T("common.btn_ai_optimize")).SetDescription(biz.T("common.btn_ai_optimize"))

	return aiOptimize
}
