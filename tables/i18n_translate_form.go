package tables

import (
	"data4test/biz"
	"html/template"

	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// addI18nTranslateFields 给 i18n 同步弹窗追加「是否翻译 / 翻译模板 / 目标语种 / 翻译平台 / 接口地址 / API Key」字段
// 接口地址与 API Key 在「选中平台」后显示，预填引擎配置值，可实时覆盖（仅本次调用，不落库不落日志）
func addI18nTranslateFields(panel *types.FormPanel) {
	panel.AddField(biz.T("common.translate_enable"), "translate_enable", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: biz.T("common.no")},
			{Value: "1", Text: biz.T("common.yes")},
		}).FieldDefault("0").
		FieldFoot(template.HTML(`<script>
	$(function () {
	    var overrideFields = ['translate_base_url', 'translate_api_key'];

	    function toggleOverrideFields() {
	        var v = $('input[name="translate_enable"]:checked').val() || '0';
	        var p = $('[name="platform"]').val() || '';
	        var show = (v === '1') && p !== '';
	        for (var i = 0; i < overrideFields.length; i++) {
	            if (show) {
	                $("label[for='" + overrideFields[i] + "']").parent().show();
	            } else {
	                $("label[for='" + overrideFields[i] + "']").parent().hide();
	            }
	        }
	    }

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
	        toggleOverrideFields();
	    }

	    function maskKey(k) {
	        if (!k) { return ''; }
	        if (k.length <= 4) { return '****'; }
	        return k.substring(0, 2) + '****' + k.substring(k.length - 2);
	    }

	    function loadConnectInfo() {
	        var p = $('[name="platform"]').val() || '';
	        if (!p) { return; }
	        $.post('/get_ai_connect_info', { platform: p }, function (res) {
	            if (res && res.code === 200 && res.data) {
	                $('input[name="translate_base_url"]').val(res.data.baseUrl || '');
	                // 原配置的 apiKey 为示例值，直接完整展示；个人输入的 key 才做掩码。
	                var key = res.data.apiKey || '';
	                $('input.api_key_display').val(key).data('sample', key);
	                $('input[name="translate_api_key"][type="hidden"]').val('');
	            }
	        });
	    }

	    // API Key 可见输入框：原配置示例 key 完整展示；个人输入的真实 key 掩码为「前2位****后2位」，
	    // 真实提交值放同名 hidden input；未改动时 hidden 为空（后端回落引擎配置）。
	    var keyInput = $('input[name="translate_api_key"]');
	    if (keyInput.length > 0 && $('input[name="translate_api_key"][type="hidden"]').length === 0) {
	        keyInput.removeAttr('name').addClass('api_key_display');
	        $('<input type="hidden" name="translate_api_key" value="">').insertAfter(keyInput);
	    }
	    $(document).on('focus', 'input.api_key_display', function () {
	        var sample = $(this).data('sample') || '';
	        if ($(this).val() === sample) {
	            // 仍是示例 key，清空便于输入个人 key
	            $(this).val('');
	            $('input[name="translate_api_key"][type="hidden"]').val('');
	        } else {
	            // 已输入个人 key（当前为掩码态），聚焦还原完整值便于编辑
	            $(this).val($('input[name="translate_api_key"][type="hidden"]').val());
	        }
	    });
	    $(document).on('blur', 'input.api_key_display', function () {
	        var v = $(this).val();
	        if (v === '') {
	            $(this).val($(this).data('sample') || '');
	            $('input[name="translate_api_key"][type="hidden"]').val('');
	        } else {
	            $(this).val(maskKey(v));
	        }
	    });
	    $(document).on('input', 'input.api_key_display', function () {
	        $('input[name="translate_api_key"][type="hidden"]').val($(this).val());
	    });

	    toggleTranslateFields();
	    $('input[name="translate_enable"]').on('ifChanged change', toggleTranslateFields);
	    $('[name="platform"]').on('change', function () {
	        toggleOverrideFields();
	        loadConnectInfo();
	    });
	});
	</script>`))
	panel.AddField(biz.T("ai_template.translate_template"), "translate_template", db.Varchar, form.SelectSingle).
		FieldOptions(biz.GetAiTemplateOptions("8"))
	panel.AddField(biz.T("common.translate_langs"), "langs", db.Varchar, form.Select).
		FieldOptions(biz.GetTranslateTargetLanguages())
	panel.AddField(biz.T("common.translate_platform"), "platform", db.Varchar, form.SelectSingle).
		FieldOptions(biz.GetAiCreatePlatform())
	panel.AddField(biz.T("ai_template.translate_base_url"), "translate_base_url", db.Varchar, form.Url).
		FieldHelpMsg(template.HTML(biz.T("ai_template.translate_base_url_help")))
	panel.AddField(biz.T("ai_template.translate_api_key"), "translate_api_key", db.Varchar, form.Text).
		FieldHelpMsg(template.HTML(biz.T("ai_template.translate_api_key_help")))
}
