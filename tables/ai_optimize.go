package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiOptimizeTable(ctx *context.Context) table.Table {

	aiOptimize := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiOptimize.GetInfo().HideFilterArea()

	info.AddField("自增主键", "id", db.Int)
	info.AddField("优化指令", "optimize_desc", db.Text)
	info.AddField("创建人", "create_user", db.Varchar)
	info.AddField("创建时间", "created_at", db.Timestamp)
	info.AddField("更新时间", "updated_at", db.Timestamp)
	info.AddField("删除时间", "deleted_at", db.Timestamp)

	info.SetTable("ai_optimize").SetTitle("优化指令").SetDescription("优化指令")

	formList := aiOptimize.GetForm().HideContinueEditCheckBox()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("优化指令", "optimize_desc", db.Text, form.TextArea)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).FieldHide()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldHide()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldHide()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).FieldHide()

	formList.SetFooterHtml(`
<script>
document.addEventListener("DOMContentLoaded", function () {
    // 获取保存按钮
    const saveBtn = document.querySelector("button[type='submit']");

    if (saveBtn) {
        // 修改按钮文字
        saveBtn.innerText = "优化";

        // 解绑原有事件（防止表单提交）
        const newBtn = saveBtn.cloneNode(true);
        saveBtn.parentNode.replaceChild(newBtn, saveBtn);

        // 添加自定义点击逻辑
        newBtn.addEventListener("click", function () {
            // 你可以加 loading、禁用状态等
            newBtn.innerText = "优化中...";

            fetch("/your/custom/api", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ something: "example" })
            })
            .then(res => res.json())
            .then(res => {
                newBtn.innerText = "优化完成";
                alert("优化成功：" + res.msg);
            })
            .catch(err => {
                newBtn.innerText = "生成";
                alert("生成失败：" + err);
            });
        });
    }
});
</script>`)

	formList.SetTable("ai_optimize").SetTitle("优化指令").SetDescription("AiOptimize")

	return aiOptimize
}
