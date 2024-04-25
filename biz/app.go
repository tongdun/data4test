package biz

import (
	"data4perf/models"
	"fmt"
)

func GetEnvConfig(name, source string) (envConfig EnvConfig, err error) {
	if source == "scene" {
		var dbProduct DbProduct
		models.Orm.Table("product").Where("product = ?", name).Find(&dbProduct)
		if len(dbProduct.Name) > 0 {
			envConfig.Product = name
			envConfig.Protocol = dbProduct.Protocol
			envConfig.Auth = dbProduct.Auth
			envConfig.Testmode = dbProduct.Testmode
			envConfig.Ip = dbProduct.Ip
		} else {
			err = fmt.Errorf("未找到[%s]产品配置信息", name)
			//Logger.Warning("%s", err)  // 日志在链路上打印
		}
	} else if source == "data" {
		models.Orm.Table("env_config").Where("app = ?", name).Find(&envConfig)
		if len(envConfig.App) == 0 {
			err = fmt.Errorf("未找到[%s]应用配置信息", name)
			//Logger.Warning("%s", err)
		}
	}

	return
}

func GetAppName(id string) (name string, err error) {
	var appNames []string
	models.Orm.Table("env_config").Where("id = ?", id).Select("app").Pluck("app", &appNames)

	if len(appNames) == 0 {
		err = fmt.Errorf("未找到[%v]应用信息，请核对", id)
		Logger.Error("%s", err)
		return
	}
	name = appNames[0]
	return
}

func GetAppId(app string) (id int, err error) {
	var ids []int

	models.Orm.Table("env_config").Where("app = ?", app).Pluck("id", &ids)
	if len(ids) == 0 {
		err = fmt.Errorf("未找到[%v]应用信息，请核对", app)
		Logger.Error("%s", err)
		return
	}
	id = ids[0]
	return
}

func UpdateApiChangeByAppId(id string) (err error) {
	appName, err := GetAppName(id)
	if err != nil {
		return
	}
	var appApiChange AppApiChange
	appApiChange.App = appName
	UpdateApiChangeLog(appApiChange)
	return
}

func UpdateApiAutoStatus(id string) (err error) {
	appName, err := GetAppName(id)
	if err != nil {
		return
	}

	var apiIds []string

	models.Orm.Table("api_definition").Where("app = ?", appName).Pluck("api_id", &apiIds)

	for _, item := range apiIds {
		var dataCount int

		models.Orm.Table("scene_data").Where("api_id = ? and app = ?", item, appName).Count(&dataCount)
		if dataCount > 0 {

			models.Orm.Table("api_definition").Select("is_auto").Where("api_id = ? and app = ?", item, appName).Update("is_auto", 1)
		} else {
			// 更新空值使用nil或用-1代替
			models.Orm.Table("api_definition").Select("is_auto").Where("api_id = ? and app = ?", item, appName).Update("is_auto", "-1")
		}
		if err != nil {
			return
		}

	}
	return
}
