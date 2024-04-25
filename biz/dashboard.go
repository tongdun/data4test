package biz

import (
	"data4perf/models"
	"fmt"
	"html/template"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/GoAdminGroup/go-admin/template/types"
)

func GetApps() (apps []types.FieldOption) {
	var appList []string
	var app types.FieldOption
	models.Orm.Table("env_config").Order("created_at DESC").Pluck("app", &appList)

	if len(appList) >= 0 {
		for _, item := range appList {
			app.Value = item
			app.Text = item
			apps = append(apps, app)
		}
	}
	return
}

func GetProducts() (products []types.FieldOption) {
	var dbProducts []Product
	var product types.FieldOption
	models.Orm.Table("product").Order("created_at desc").Find(&dbProducts)

	if len(dbProducts) >= 0 {
		for _, item := range dbProducts {
			product.Value = item.Name
			product.Text = item.Name
			products = append(products, product)
		}
	}
	return
}

func GetTestcaseType() (caseTypes []types.FieldOption) {
	var dbSysParmeters []SysParameter
	var caseType types.FieldOption
	models.Orm.Table("sys_parameter").Where("name = ?", "TestCaseType").Find(&dbSysParmeters)
	if len(dbSysParmeters) > 0 {
		values := dbSysParmeters[0].ValueList
		caseTypeTmp := strings.Split(values, ",")
		for _, item := range caseTypeTmp {
			caseType.Value = item
			caseType.Text = item
			caseTypes = append(caseTypes, caseType)
		}
	}
	return
}

// 功能废弃，持续观察一段时间
func GetFiles() (files []types.FieldOption) {
	var file types.FieldOption
	dirHandle, err := ioutil.ReadDir(DataBasePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	var fileList []string
	for _, fileH := range dirHandle {
		if !fileH.IsDir() {
			fileName := fileH.Name()
			fileList = append(fileList, fileName)
		}
	}
	var rawfiles []string
	for _, fileName := range fileList {
		if !strings.HasSuffix(fileName, ".yml") && !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".json") {
			continue
		}
		rawfiles = append(rawfiles, fileName)
	}

	if len(rawfiles) >= 0 {
		tag := 0
		for _, item := range rawfiles {
			if tag == 0 {
				file.Value = fmt.Sprintf("%s", item)
			} else {
				file.Value = fmt.Sprintf("%s", item) // 去年换行标签<br>
			}
			tag++
			file.Text = item
			files = append(files, file)
		}
	}
	return
}

func GetFilesFromMySQL() (files []types.FieldOption) {
	var dataList []string
	var dataName types.FieldOption
	models.Orm.Table("scene_data").Order("created_at desc").Pluck("file_name", &dataList)
	if len(dataList) > 0 {
		for _, item := range dataList {
			dataName.Value = fmt.Sprintf("%s", item) // 去掉换行
			dataName.Text = item
			files = append(files, dataName)
		}
	}

	return
}

func GetDatas() (dataNames []types.FieldOption) {
	var dataList []string
	var dataName types.FieldOption
	models.Orm.Table("scene_data").Order("created_at desc").Pluck("name", &dataList)
	if len(dataList) > 0 {
		for _, item := range dataList {
			dataName.Value = fmt.Sprintf("%s", item) // 去掉换行
			dataName.Text = item
			dataNames = append(dataNames, dataName)
		}
	}

	return
}

func GetScenes() (playbooks []types.FieldOption) {
	var sceneList []string
	var playbook types.FieldOption
	models.Orm.Table("playbook").Order("created_at desc").Pluck("name", &sceneList)

	if len(sceneList) > 0 {
		for _, item := range sceneList {
			playbook.Value = fmt.Sprintf("%s", item) // 去掉换行
			playbook.Text = item
			playbooks = append(playbooks, playbook)
		}
	}
	return
}

func Get24No() (numbers []types.FieldOption) {
	var number types.FieldOption
	for i := 0; i <= 23; i++ {
		number.Value = strconv.Itoa(i)
		number.Text = strconv.Itoa(i)
		numbers = append(numbers, number)
	}
	return
}

func Get7No() (numbers []types.FieldOption) {
	var number types.FieldOption
	for i := 1; i <= 7; i++ {
		number.Value = strconv.Itoa(i)
		number.Text = strconv.Itoa(i)
		numbers = append(numbers, number)
	}
	return
}

type ContentHTML struct {
	Content template.HTML
}

type TemplateAPI struct {
	AllCount           template.HTML `gorm:"column:all_count" json:"all_count"`
	AutomatableCount   template.HTML `gorm:"column:automatable_count" json:"automatable_count"`
	UnautomatableCount template.HTML `gorm:"column:unautomatable_count" json:"unautomatable_count"`
	AutoTestCount      template.HTML `gorm:"column:auto_test_count" json:"auto_test_count"`
	UntestCount        template.HTML `gorm:"column:untest_count" json:"untest_count"`
	PassCount          template.HTML `gorm:"column:pass_count" json:"pass_count"`
	FailCount          template.HTML `gorm:"column:fail_count" json:"fail_count"`
	Project            template.HTML `gorm:"column:project" json:"project"`
}

func GetAPISumUp() (sumUps []map[string]types.InfoItem) {
	var sumup map[string]types.InfoItem
	var allApi []TemplateAPI
	models.Orm.Table("product_count").Order("all_count DESC").Limit(10).Find(&allApi)
	for _, item := range allApi {
		sumup = make(map[string]types.InfoItem, 1)
		sumup["关联项目"] = types.InfoItem{Content: item.Project}
		sumup["API总数"] = types.InfoItem{Content: item.AllCount}
		sumup["可自动化数"] = types.InfoItem{Content: item.AutomatableCount}
		sumup["不可自动化数"] = types.InfoItem{Content: item.UnautomatableCount}
		sumup["自动化测试总数"] = types.InfoItem{Content: item.AutoTestCount}
		sumup["未测试总数"] = types.InfoItem{Content: item.UntestCount}
		sumup["通过总数"] = types.InfoItem{Content: item.PassCount}
		sumup["失败总数"] = types.InfoItem{Content: item.FailCount}
		sumUps = append(sumUps, sumup)
	}
	return
}

type Box struct {
	Title  template.HTML
	Color  template.HTML
	Number template.HTML
	Icon   template.HTML
}

func GetBoxData() (boxPlural []Box) {
	var projectCount, apiCaseCount, fuzzingCaseCount, dataCaseCount, testAll template.HTML

	models.Orm.Table("env_config").Count(&projectCount)
	box1 := Box{"环境总数", "white", projectCount, "ion-ios-gear-outline"}
	boxPlural = append(boxPlural, box1)

	models.Orm.Table("api_definition").Count(&apiCaseCount)
	box2 := Box{"接口用例总数", "white", apiCaseCount, "ion-ios-gear-outline"}
	boxPlural = append(boxPlural, box2)

	models.Orm.Table("api_test_data").Count(&dataCaseCount)
	box3 := Box{"测试数据总数", "white", dataCaseCount, "ion-ios-gear-outline"}
	boxPlural = append(boxPlural, box3)

	models.Orm.Table("api_fuzzing_data").Count(&fuzzingCaseCount)
	box4 := Box{"模糊数据总数", "white", fuzzingCaseCount, "ion-ios-gear-outline"}
	boxPlural = append(boxPlural, box4)

	models.Orm.Table("api_test_detail").Count(&testAll)
	box5 := Box{"测试总次数", "white", testAll, "ion-ios-gear-outline"}
	boxPlural = append(boxPlural, box5)

	return
}
