package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
)

//func GetEnvConfig(name, nameType string) (envConfig EnvConfig, privateAppPrefix map[string]interface{}, err error) {
//	switch nameType {
//	case "product":
//		var dbProduct DbProduct
//		models.Orm.Table("product").Where("product = ?", name).Find(&dbProduct)
//		if len(dbProduct.Name) > 0 {
//			envConfig.Product = name
//			envConfig.Protocol = dbProduct.Protocol
//			envConfig.Auth = dbProduct.Auth
//			envConfig.Testmode = dbProduct.Testmode
//			envConfig.Ip = dbProduct.Ip
//			if len(dbProduct.PrivateAppPrefix) > 0 {
//				//privateAppPrefix := make(map[string]interface{})
//				err = json.Unmarshal([]byte(dbProduct.PrivateAppPrefix), &privateAppPrefix)
//				if err != nil {
//					Logger.Error("解析专用应用前缀异常: %s", err)
//					return
//				}
//			}
//		} else {
//			Logger.Warning("未找到[%s]产品配置信息", name)
//		}
//	case "app":
//		models.Orm.Table("env_config").Where("app = ?", name).Find(&envConfig)
//		if len(envConfig.App) == 0 {
//			Logger.Warning("未找到[%s]应用配置信息", name)
//		}
//	}
//
//	return
//}

func GetEnvConfig(productName, appName string) (envConfig EnvConfig, err error) {
	if len(productName) == 0 && len(appName) == 0 {
		err = fmt.Errorf("产品名称和应用名称不能同时为空")
		Logger.Error("%s", err)
		return
	}

	if len(appName) > 0 {
		models.Orm.Table("env_config").Where("app = ?", appName).Find(&envConfig)
		if len(envConfig.App) == 0 {
			Logger.Warning("未找到[%s]应用配置信息", appName)
		}
	}

	if len(productName) > 0 {
		var dbProduct DbProduct
		models.Orm.Table("product").Where("product = ?", productName).Find(&dbProduct)
		if len(dbProduct.Name) > 0 {
			envConfig.Product = productName
			envConfig.Protocol = dbProduct.Protocol
			envConfig.Auth = dbProduct.Auth
			envConfig.Testmode = dbProduct.Testmode
			envConfig.Ip = dbProduct.Ip
			if len(dbProduct.PrivateAppPrefix) > 0 {
				privateAppPrefix := make(map[string]interface{})
				err = json.Unmarshal([]byte(dbProduct.PrivateAppPrefix), &privateAppPrefix)
				if err != nil {
					Logger.Error("解析专用应用前缀异常: %s", err)
					return
				}

				if len(appName) > 0 {
					if prefix, ok := privateAppPrefix[appName]; ok {
						envConfig.Prepath = prefix.(string)
					}
				}
			}
		} else {
			Logger.Warning("未找到[%s]产品配置信息", productName)
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

func GetSwaggerPath(id string) (swaggerPath string) {
	var pathList []string
	models.Orm.Table("env_config").Where("id = ?", id).Select("swagger_path").Pluck("swagger_path", &pathList)

	if len(pathList) == 0 {
		err := fmt.Errorf("未找到[%v]应用信息，请核对", id)
		Logger.Error("%s", err)
		return
	}

	swaggerPath = pathList[0]
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
