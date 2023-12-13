package biz

func GetApiMenu(appName string) (allMenus []FirstMenu) {
	appList := GetAppList()
	if len(appName) == 0 {
		for _, app := range appList {
			var firstMenu FirstMenu
			firstMenu.Name = app
			firstMenu.Title = app
			allMenus = append(allMenus, firstMenu)
		}
	} else {
		for _, app := range appList {
			var firstMenu FirstMenu
			firstMenu.Name = app
			firstMenu.Title = app
			if app == appName {
				appModel := GetAppInfo(app)
				for _, module := range appModel.Modules {
					var secondMenu SecondMenu
					secondMenu.Title = module
					secondMenu.Name = module
					moduleModel := GetModuleInfo(app, module)
					for _, api := range moduleModel.ApisDesc {
						var baseMenu BaseMenu
						baseMenu.Name = api
						baseMenu.Title = api
						secondMenu.Children = append(secondMenu.Children, baseMenu)
					}
					firstMenu.Children = append(firstMenu.Children, secondMenu)
				}
			}
			allMenus = append(allMenus, firstMenu)
		}
	}
	return
}

func GetDataMenu(appName string) (allMenus []BaseMemuModel) {
	appList := GetAppList()
	if len(appName) == 0 {
		for _, app := range appList {
			var dataMemuModel BaseMemuModel
			dataMemuModel.Name = app
			dataMemuModel.Title = app
			allMenus = append(allMenus, dataMemuModel)
		}
	} else {
		for _, app := range appList {
			var dataMemuModel BaseMemuModel
			dataMemuModel.Name = app
			dataMemuModel.Title = app
			if app == appName {
				appModel := GetAppInfo(app)
				dataMemuModel.Children = append(dataMemuModel.Children, appModel.DatasDesc...)
			}
			allMenus = append(allMenus, dataMemuModel)
		}
	}
	return
}

func GetSceneMenu(productName string) (allMenus []BaseMemuModel) {
	productList := GetProductList()
	if len(productName) == 0 {
		for _, product := range productList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = product
			sceneMemuModel.Title = product
			allMenus = append(allMenus, sceneMemuModel)
		}
	} else {
		for _, product := range productList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = product
			sceneMemuModel.Title = product
			if product == productName {
				playbooks := GetProductPlaybook(product)
				sceneMemuModel.Children = append(sceneMemuModel.Children, playbooks...)
			}
			allMenus = append(allMenus, sceneMemuModel)
		}
	}
	return
}

func GetHistoryMenu(dateName string) (allMenus []BaseMemuModel) {
	dateList := GetHistoryDateList()
	if len(dateName) == 0 {
		for _, date := range dateList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = date
			sceneMemuModel.Title = date
			allMenus = append(allMenus, sceneMemuModel)
		}
	} else {
		for _, date := range dateList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = date
			sceneMemuModel.Title = date
			if date == dateName {
				historyDatas := GetDateHistory(date)
				sceneMemuModel.Children = append(sceneMemuModel.Children, historyDatas...)
			}
			allMenus = append(allMenus, sceneMemuModel)
		}
	}
	return
}

func GetSceneHistoryMenu(dateName string) (allMenus []BaseMemuModel) {
	dateList := GetSceneHistoryDateList()
	if len(dateName) == 0 {
		for _, date := range dateList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = date
			sceneMemuModel.Title = date
			allMenus = append(allMenus, sceneMemuModel)
		}
	} else {
		for _, date := range dateList {
			var sceneMemuModel BaseMemuModel
			sceneMemuModel.Name = date
			sceneMemuModel.Title = date
			if date == dateName {
				historyDatas := GetDateSceneHistory(date)
				sceneMemuModel.Children = append(sceneMemuModel.Children, historyDatas...)
			}
			allMenus = append(allMenus, sceneMemuModel)
		}
	}
	return
}
