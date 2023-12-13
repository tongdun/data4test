package biz

type FirstMenu struct {
	Title string `json:"title"`
	Name string `json:"name"`
	//Closable bool `json:"closable"`
	//ShowInTags bool `json:"showInTags"`
	//ShowInMenus bool `json:"showInMenus"`
	//Opened bool `json:"opened"`
	Children []SecondMenu `json:"children"`
}

type SecondMenu struct {
	Title string `json:"title"`
	Name string `json:"name"`
	//Closable bool `json:"closable"`
	//ShowInTags bool `json:"showInTags"`
	//ShowInMenus bool `json:"showInMenus"`
	//Opened bool `json:"opened"`
	Children []BaseMenu `json:"children"`
}

type BaseMenu struct {
	Title string `json:"title"`
	Name string `json:"name"`
	//Closable bool `json:"closable"`
	//ShowInTags bool `json:"showInTags"`
	//ShowInMenus bool `json:"showInMenus"`
}

type BaseMemuModel struct {
	Title string `json:"title"`
	Name string `json:"name"`
	//Closable bool `json:"closable"`
	//ShowInTags bool `json:"showInTags"`
	//ShowInMenus bool `json:"showInMenus"`
	//Opened bool `json:"opened"`
	Children []string `json:"children"`
}
