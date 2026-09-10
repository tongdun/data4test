package tables

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/template/types"
)

// DropdownItem 下拉菜单中的一个子项：文案 + 图标 + 动作（复用现有 action）
type DropdownItem struct {
	Label  template.HTML
	Icon   string
	Action types.Action
}

// dropdownItem 内部结构：额外记录子项对应的 class，用于绑定点击事件
type dropdownItem struct {
	cls    string
	label  template.HTML
	icon   string
	action types.Action
}

// dropdownButton 自定义工具栏按钮：渲染 Bootstrap3 下拉菜单，
// 子项复用各自 action 的点击绑定（弹窗 / 跳转 / Ajax）。
type dropdownButton struct {
	*types.BaseButton
	Icon  string
	Items []dropdownItem
}

// dropdownSeq 生成全局唯一的 class 后缀，避免同页多表冲突
var dropdownSeq int

func (b *dropdownButton) Content() (template.HTML, template.JS) {
	h := template.HTML(`<div class="btn-group pull-right" style="margin-right: 10px">`)
	h += `<a class="btn btn-sm btn-default dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">` +
		`<i class="fa ` + template.HTML(b.Icon) + `"></i>&nbsp;&nbsp;` + b.Title + `&nbsp;<span class="caret"></span></a>`
	h += `<ul class="dropdown-menu">`

	var js template.JS
	for _, it := range b.Items {
		h += `<li><a class="` + template.HTML(it.cls) + ` ` + it.action.BtnClass() + `" ` + it.action.BtnAttribute() + `>` +
			`<i class="fa ` + template.HTML(it.icon) + `"></i>&nbsp;` + it.label + `</a></li>`
		js += it.action.Js()
	}

	h += `</ul></div>`
	return h, js
}

// addDropdownButton 在列表工具栏添加一个下拉按钮，子项动作复用各自回调 / 弹窗 / 跳转。
func addDropdownButton(info *types.InfoPanel, title template.HTML, icon string, items []DropdownItem) {
	prepared := make([]dropdownItem, 0, len(items))
	for _, it := range items {
		dropdownSeq++
		cls := fmt.Sprintf("ddbtn_%d", dropdownSeq)
		it.Action.SetBtnId("." + cls)
		prepared = append(prepared, dropdownItem{cls: cls, label: it.Label, icon: it.Icon, action: it.Action})
	}

	btn := &dropdownButton{
		BaseButton: &types.BaseButton{Title: title, Action: &types.NilAction{}},
		Icon:       icon,
		Items:      prepared,
	}
	info.AddButtonRaw(btn, &types.NilAction{})

	// 每个子项单独注册回调 + 弹窗壳（FooterContent），但不渲染独立按钮
	for _, it := range prepared {
		info.AddButtonRaw(&types.BaseButton{Action: it.action}, it.action)
	}
}
