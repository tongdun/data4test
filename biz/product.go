package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
)

func GetProductName(id string) (name string, err error) {
	var dbProduct Product
	models.Orm.Table("product").Where("id = ?", id).Find(&dbProduct)

	if len(dbProduct.Name) == 0 {
		err = fmt.Errorf("未找到[%v]产品信息，请核对", id)
		Logger.Error("%s", err)
		return
	}
	name = dbProduct.Name
	return
}

func GetEnvTypeByName(product string) (envType int, err error) {
	var dbProduct Product
	models.Orm.Table("product").Where("product = ?", product).Find(&dbProduct)

	if len(dbProduct.Name) == 0 {
		err = fmt.Errorf("未找到[%v]产品信息，请核对", product)
		Logger.Warning("%s", err)
		return
	}
	envType = dbProduct.EnvType
	return
}

func GetProductEnv(name string) (envModel ProductEnvModel, err error) {
	var product Product
	models.Orm.Table("product").Where("product = ?", name).Find(&product)

	if len(product.Name) == 0 {
		err = fmt.Errorf("未找到[%v]产品信息，请核对", name)
		Logger.Error("%s", err)
		return
	}
	envModel.Protocol = product.Protocol
	envModel.Ip = product.Ip
	envModel.Name = product.Name
	auth := make(map[string]string)

	if len(product.Auth) == 0 {
		err = fmt.Errorf("未配置鉴权信息，请先配置")
		Logger.Error("%s", err)
		return
	}
	json.Unmarshal([]byte(product.Auth), &auth)
	for k, v := range auth {
		var varDataModel VarDataModel
		varDataModel.Name = k
		varDataModel.TestValue = append(varDataModel.TestValue, v)
		envModel.Auth = append(envModel.Auth, varDataModel)
	}

	return
}

func GetProductApps(id string) (name string, err error) {
	var dbProduct Product
	models.Orm.Table("product").Where("id = ?", id).Find(&dbProduct)

	if len(dbProduct.Name) == 0 {
		err = fmt.Errorf("未找到[%v]产品信息，请核对", id)
		Logger.Error("%s", err)
		return
	}
	name = dbProduct.Apps
	return
}

func GetProductPrivateParameter(name string) (privateParameter map[string]interface{}) {
	var dbProduct DbProduct
	models.Orm.Table("product").Where("product = ?", name).Find(&dbProduct)
	privateParameter = make(map[string]interface{})
	if len(dbProduct.Name) > 0 {
		if len(dbProduct.PrivateParameter) > 2 {
			err := json.Unmarshal([]byte(dbProduct.PrivateParameter), &privateParameter)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	return
}