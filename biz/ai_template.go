package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/types"
)

func GetAiTemplateByName(name, tType string) (content, appendContent string, err error) {
	templateStatus := "apply"
	var aiTemplate AiTemplateDefine
	models.Orm.Table("ai_template").Where("use_status = ? and template_name = ? and template_type = ?", templateStatus, name, tType).Find(&aiTemplate)
	if len(aiTemplate.TemplateContent) == 0 {
		err = fmt.Errorf(T("error.template_not_found"), tType)
		Logger.Error("%s", err)
		return
	}
	content = aiTemplate.TemplateContent
	appendContent = aiTemplate.AppendConversion
	return
}

func GetAiPlatform() (platformList []types.FieldOption) {
	var sysParameter SysParameter
	var platform types.FieldOption
	parameterName := "aiPlatform"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		err := fmt.Errorf(T("error.sys_param_not_defined"), parameterName)
		Logger.Error("%s", err)
		platformList = GetNoSelectOption(T("info.define_ai_run_engine"))
		return
	}

	values := GetValuesFromStringList(sysParameter.ValueList)
	if len(values) >= 0 {
		for _, item := range values {
			platform.Value = item
			platform.Text = item
			platformList = append(platformList, platform)
		}
	}

	if len(platformList) == 0 {
		platformList = GetNoSelectOption(T("info.define_ai_run_engine"))
	}

	return

}

func GetAiCreatePlatform() (platformList []types.FieldOption) {
	var sysParameter SysParameter
	parameterName := "aiRunEngine"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		err := fmt.Errorf(T("error.sys_param_not_defined"), parameterName)
		Logger.Error("%s", err)
		platformList = GetNoSelectOption(T("info.define_ai_run_engine"))
		return
	}

	aiConnect := make(map[string]AIConnect)
	err := json.Unmarshal([]byte(sysParameter.ValueList), &aiConnect)
	if err != nil {
		Logger.Error(T("error.param_definition_error"), parameterName, err)
		return
	}

	for k, _ := range aiConnect {
		var aiPlatform types.FieldOption
		aiPlatform.Value = k
		aiPlatform.Text = k
		platformList = append(platformList, aiPlatform)
	}

	if len(platformList) == 0 {
		platformList = GetNoSelectOption(T("info.define_ai_run_engine"))
	}

	return

}

func GetAiTemplateOptions(tType string) (aiTemplates []types.FieldOption) {
	var templateNames []string
	templateStatus := "apply"
	aiTypeDefine := map[string]string{"1": T("type.test_case"), "2": T("type.data"), "3": T("type.scene"), "4": T("type.task"), "5": "Issue", "6": T("type.report"), "7": T("type.analysis")}
	models.Orm.Table("ai_template").Where("use_status = ? and template_type = ?", templateStatus, tType).Order("updated_at DESC").Pluck("template_name", &templateNames)
	if len(templateNames) == 0 {
		err := fmt.Errorf(T("error.template_not_found"), aiTypeDefine[tType])
		Logger.Error("%s", err)
		//return //不进行退出
	}
	var aiTemplate types.FieldOption

	for _, item := range templateNames {
		aiTemplate.Value = item
		aiTemplate.Text = item
		aiTemplates = append(aiTemplates, aiTemplate)
	}

	if len(aiTemplates) == 0 {
		var aiTemplate types.FieldOption
		desc := fmt.Sprintf(T("info.define_template"), aiTypeDefine[tType])
		aiTemplate.Value = desc
		aiTemplate.Text = desc
		aiTemplates = append(aiTemplates, aiTemplate)
	}

	return
}

func UpdateAiTemplateAppendContend(appendContentList []string, id string) (err error) {
	var content string
	for index, value := range appendContentList {
		if index == 0 {
			content = value
		} else {
			content = fmt.Sprintf("%s|**|%s", content, value)
		}
	}

	if len(content) == 0 {
		content = " "
	}

	err = models.Orm.Table("ai_template").Where("id = ?", id).UpdateColumn(&AiTemplateDefine{AppendConversion: content}).Error
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	return
}

func CopyAiTemplate(id, userName string) (err error) {
	var dbTemplate DbAiTemplateDefine
	models.Orm.Table("ai_template").Where("id = ?", id).Find(&dbTemplate)
	if len(dbTemplate.TemplateName) == 0 {
		err = fmt.Errorf(T("error.data_not_found"), id)
		Logger.Error("%s", err)
		return
	}

	var aiTemplate AiTemplateDefine
	aiTemplate = dbTemplate.AiTemplateDefine
	aiTemplate.TemplateName = fmt.Sprintf(T("template.copy_suffix"), dbTemplate.TemplateName)
	aiTemplate.CreateUser = userName
	aiTemplate.UseStatus = "edit"

	err = models.Orm.Table("ai_template").Create(aiTemplate).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}
