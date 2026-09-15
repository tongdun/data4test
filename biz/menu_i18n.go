package biz

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"data4test/models"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"gopkg.in/yaml.v2"
)

// i18nMenuDir 菜单多语种数据目录
func i18nMenuDir() string {
	return fmt.Sprintf("%s/i18n/i18n_menu", BASEPATH)
}

// loadMenuI18n 扫描 <dir>/<语种>.yaml，构建 [lang][中文标题]→译文 索引
// 每语种一个文件，内容为扁平「中文 → 译文」映射
func loadMenuI18n() map[string]map[string]string {
	index := make(map[string]map[string]string)
	files, err := ioutil.ReadDir(i18nMenuDir())
	if err != nil {
		return index
	}
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || (!strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml")) {
			continue
		}
		base := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
		lang := normalizeLang(base)

		data, err := ioutil.ReadFile(filepath.Join(i18nMenuDir(), name))
		if err != nil {
			Logger.Error("menu_i18n read %s failed: %v", name, err)
			continue
		}
		var m map[string]string
		if err := yaml.Unmarshal(data, &m); err != nil {
			Logger.Error("menu_i18n parse %s failed: %v", name, err)
			continue
		}
		index[lang] = m
	}
	return index
}

func InitMenuI18n() {
	models.Orm.Table("goadmin_menu").Where("type = 0").UpdateColumn("type", 1)

	var titles []string
	models.Orm.Table("goadmin_menu").Where("id > 1").Pluck("title", &titles)
	if len(titles) == 0 {
		return
	}

	if language.Lang == nil {
		return
	}

	menuI18n := loadMenuI18n()

	// CN：标题即原文，恒等注册（go-admin 查表时内部会 ToLower，这里统一处理）
	if language.Lang[language.CN] == nil {
		language.Lang[language.CN] = make(language.LangSet)
	}
	for _, title := range titles {
		language.Lang[language.CN][strings.ToLower(title)] = title
	}

	// EN：从 i18n_menu 读取译文，未命中则不注册（回退 CN 原值）
	if byLang, ok := menuI18n["en-US"]; ok && len(byLang) > 0 {
		if language.Lang[language.EN] == nil {
			language.Lang[language.EN] = make(language.LangSet)
		}
		for _, title := range titles {
			if en, ok := byLang[title]; ok && en != "" {
				language.Lang[language.EN][strings.ToLower(title)] = en
			}
		}
	}
}

// RefreshMenuI18n 刷新菜单翻译注册，获取数据库中最新菜单列表重新注册翻译
// 在新增/修改菜单后调用，确保新菜单也有对应的翻译
func RefreshMenuI18n() {
	if models.Orm == nil {
		Logger.Warning("RefreshMenuI18n: models.Orm 未初始化，跳过菜单翻译刷新")
		return
	}
	InitMenuI18n()
}
