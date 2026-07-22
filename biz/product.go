package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
)

func GetProductName(id string) (name string, err error) {
	var names []string
	models.Orm.Table("product").Where("id = ?", id).Pluck("product", &names)

	if len(names) == 0 {
		err = fmt.Errorf(T("error.product_not_found"), id)
		Logger.Error("%s", err)
		return
	}

	name = names[0]

	return
}

func GetEnvTypeByName(product string) (envType int) {
	var envTypes []int
	models.Orm.Table("product").Where("product = ?", product).Pluck("env_type", &envTypes)

	if len(envTypes) == 0 {
		Logger.Warning(T("warning.product_not_found"), product)
		return
	}

	envType = envTypes[0]

	return
}

func GetProductEnv(name string) (envModel ProductEnvModel, err error) {
	var product Product
	models.Orm.Table("product").Where("product = ?", name).Find(&product)

	if len(product.Name) == 0 {
		err = fmt.Errorf(T("error.product_not_found"), name)
		Logger.Error("%s", err)
		return
	}
	envModel.Protocol = product.Protocol
	envModel.Ip = product.Ip
	envModel.Name = product.Name
	auth := make(map[string]string)

	if len(product.Auth) == 0 {
		err = fmt.Errorf(T("error.auth_not_configured"))
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
	var names []string
	models.Orm.Table("product").Where("id = ?", id).Pluck("apps", &names)

	if len(names) == 0 {
		err = fmt.Errorf(T("error.product_not_found"), id)
		Logger.Error("%s", err)
		return
	}

	name = names[0]

	return
}

func CopyProduct(id, userName string) (err error) {
	var dbProduct DbProduct
	models.Orm.Table("product").Where("id = ?", id).Find(&dbProduct)
	if len(dbProduct.Name) == 0 {
		err = fmt.Errorf(T("error.data_not_found"), id)
		Logger.Error("%s", err)
		return
	}

	var product Product
	product = dbProduct.Product
	product.Name = fmt.Sprintf(T("common.copy_suffix"), dbProduct.Name)

	err = models.Orm.Table("product").Create(product).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func (dbProduct DbProduct) GetPrivateParameter() (privateParameter map[string]interface{}) {
	privateParameter = make(map[string]interface{})
	if len(dbProduct.PrivateParameter) > 2 {
		err := json.Unmarshal([]byte(dbProduct.PrivateParameter), &privateParameter)
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}
