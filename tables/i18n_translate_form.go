package tables

import (
	"data4test/biz"
	"html/template"

	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// addI18nTranslateFields 给 i18n 同步弹窗追加「是否翻译 / 翻译模板 / 目标语种 / 翻译平台」字段
func addI18nTranslateFields(panel *types.FormPanel) {
	panel.AddField(biz.T("common.translate_enable"), "translate_enable", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: biz.T("common.no")},
			{Value: "1", Text: biz.T("common.yes")},
		}).FieldDefault("0").
		FieldFoot(template.HTML(`<script>
$(function () {
    function toggleTranslateFields() {
        var v = $('input[name="translate_enable"]:checked').val() || '0';
        var fields = ['translate_template', 'langs', 'platform'];
        for (var i = 0; i < fields.length; i++) {
            if (v === '1') {
                $("label[for='" + fields[i] + "']").parent().show();
            } else {
                $("label[for='" + fields[i] + "']").parent().hide();
            }
        }
    }
    toggleTranslateFields();
    $('input[name="translate_enable"]').on('ifChanged change', toggleTranslateFields);
});
</script>`))
	panel.AddField(biz.T("ai_template.translate_template"), "translate_template", db.Varchar, form.SelectSingle).
		FieldOptions(biz.GetAiTemplateOptions("8"))
	panel.AddField(biz.T("common.translate_langs"), "langs", db.Varchar, form.Select).
		FieldOptions(biz.GetTranslateTargetLanguages())
	panel.AddField(biz.T("common.translate_platform"), "platform", db.Varchar, form.SelectSingle).
		FieldOptions(biz.GetAiCreatePlatform())
}
