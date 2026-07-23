package tables

import (
	"data4test/biz"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	template2 "html/template"
	"strings"
)

func GetSceneTestHistoryTable(ctx *context.Context) table.Table {

	playbookTestHistory := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := playbookTestHistory.GetInfo()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	user := auth.Auth(ctx)
	userName := user.Name
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable().
		FieldTrimSpace().FieldWidth(60)
	info.AddField(biz.T("dashboard.task_id"), "task_id", db.Varchar).
		FieldHide()
	info.AddField(biz.T("common.name"), "name", db.Varchar).FieldWidth(160).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetHistoryDataLinkByDataStr(model.Value)
		})
	info.AddField(biz.T("common.last_file"), "last_file", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			b, num := biz.IsStrEndWithTimeFormat(value.Value)
			suffix := biz.GetStrSuffix(value.Value)
			if b {
				dirName := value.Value[:len(value.Value)-num-len(suffix)]
				return template.Default().
					Link().
					SetURL("/admin/fm/history/preview?path=/" + dirName + "/" + value.Value).
					SetContent(template2.HTML(value.Value)).
					OpenInNewTab().
					SetTabTitle(template2.HTML(biz.T("test_execute_history_file.title"))).
					GetContent()
			} else {
				return template.Default().
					Link().
					SetURL("/admin/fm/data/preview?path=/" + value.Value).
					SetContent(template2.HTML(value.Value)).
					OpenInNewTab().
					SetTabTitle(template2.HTML(biz.T("common.data_file"))).
					GetContent()
			}
		}).FieldWidth(160)
	info.AddField(biz.T("common.scene_type"), "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_test_history.scene_type_1")
			}
			if model.Value == "2" {
				return biz.T("scene_test_history.scene_type_2")
			}
			if model.Value == "3" {
				return biz.T("scene_test_history.scene_type_3")
			}
			if model.Value == "4" {
				return biz.T("scene_test_history.scene_type_4")
			}
			if model.Value == "5" {
				return biz.T("scene_test_history.scene_type_5")
			}
			return biz.T("scene_test_history.scene_type_1")
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
		{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
		{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
		{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
		{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
	}).FieldWidth(60)
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}).FieldWidth(70)
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext).
		FieldWidth(300).
		FieldDisplay(func(model types.FieldModel) interface{} {
			v := model.Value
			if len(v) == 0 {
				return ""
			}
			escaped := template2.HTMLEscapeString(v)
			// data_file_list 有多条数据（逗号分隔）时，该行多行展示，fail_reason 直接展示不折叠
			dataFiles, _ := model.Row["data_file_list"].(string)
			multiRow := strings.Contains(dataFiles, ",")
			if multiRow || len(v) <= 200 {
				return template2.HTML(escaped)
			}
			// 单行数据 + fail_reason 过长 → 折叠显示
			runes := []rune(v)
			maxLen := 200
			if maxLen > len(runes) {
				maxLen = len(runes)
			}
			short := template2.HTMLEscapeString(string(runes[:maxLen]))
			return template2.HTML(fmt.Sprintf(`<details style="max-width:400px"><summary style="cursor:pointer;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">%s...</summary><div style="margin-top:4px;padding:4px;background:#f5f5f5;border-radius:2px;white-space:pre-wrap;word-break:break-all;max-height:300px;overflow-y:auto">%s</div></details>`, short, escaped))
		}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.env_type_label"), "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.env_type._1")
			} else if model.Value == "2" {
				return biz.T("common.env_type._2")
			} else if model.Value == "3" {
				return biz.T("common.env_type._3")
			} else if model.Value == "4" {
				return biz.T("common.env_type._4")
			} else if model.Value == "5" {
				return biz.T("common.env_type._5")
			}
			return ""
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("common.env_type._1")},
		{Value: "2", Text: biz.T("common.env_type._2")},
		{Value: "3", Text: biz.T("common.env_type._3")},
		{Value: "4", Text: biz.T("common.env_type._4")},
		{Value: "5", Text: biz.T("common.env_type._5")},
	}).FieldWidth(60)
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(120)
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldWidth(60).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(100).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("common.btn_again")), icon.Android, action.Ajax("historyPlaybook_batch_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			user := auth.Auth(ctx)
			userName := user.Name
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("common.operate_success")
					continue
				}
				if err := biz.RunHistoryPlaybook(id, "again", userName); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_again")), action.Ajax("historyPlaybook_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if err := biz.RunHistoryPlaybook(id, "again", userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_continue")), icon.Android, action.Ajax("historyPlaybook_batch_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.RunHistoryPlaybook(id, "continue", userName); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_continue")), action.Ajax("historyPlaybook_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userName := user.Name
			if err := biz.RunHistoryPlaybook(id, "continue", userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	products := biz.GetProducts()
	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	// 列选择器：默认展示7列，用户可自行勾选，通过__columns URL参数控制，缓存到localStorage
	colFields := []string{"id","name","data_file_list","last_file","scene_type","result","fail_reason","env_type","remark","product","user_name","created_at"}
	colNames := []string{
		biz.T("common.id"),biz.T("common.name"),biz.T("common.data_file_list"),
		biz.T("common.last_file"),biz.T("common.scene_type"),biz.T("common.test_result"),
		biz.T("common.fail_reason"),biz.T("common.env_type_label"),biz.T("common.remark"),
		biz.T("common.product"),biz.T("common.user_name"),biz.T("common.created_at"),
	}
	defaultCols := []string{"id","name","data_file_list","last_file","scene_type","result","fail_reason"}
	fieldsJSON,_ := json.Marshal(colFields)
	namesJSON,_ := json.Marshal(colNames)
	defaultsJSON,_ := json.Marshal(defaultCols)
	colToggleJS := fmt.Sprintf(`<script data-exec-on-popstate>
(function C(){var f=%s,n=%s,dk="sht_cols",def=%s;
function gk(){try{return JSON.parse(localStorage.getItem(dk))}catch(e){}return null}
function sk(v){localStorage.setItem(dk,JSON.stringify(v))}
(function init(){
var p=new URLSearchParams(location.search),cols=p.get("__columns");
if(cols){sk(cols.split(","));return}
var saved=gk()||def;
p.set("__columns",saved.join(","));
location.search=p.toString();
})();
function go(c){
var p=new URLSearchParams(location.search);
p.set("__columns",c.join(","));sk(c);
location.search=p.toString();
}
var hd=document.querySelector(".box-header .box-title");
if(!hd)return;
var cur=(gk()||def).slice();
var btn=document.createElement("span");btn.style.cssText="margin-left:12px;cursor:pointer;font-size:13px;color:#888;position:relative";btn.innerHTML="&#9662; \u5217\u9009\u62e9";
var dd=document.createElement("div");dd.style.cssText="display:none;position:absolute;top:22px;left:0;background:#fff;border:1px solid #ddd;z-index:999;padding:6px 0;min-width:140px;box-shadow:0 2px 8px rgba(0,0,0,0.15)";
f.forEach(function(field,i){
var it=document.createElement("div");it.style.cssText="padding:4px 14px;cursor:pointer;white-space:nowrap";
var cb=document.createElement("input");cb.type="checkbox";cb.checked=cur.indexOf(field)>=0;cb.style.marginRight="6px";
it.appendChild(cb);it.appendChild(document.createTextNode(n[i]));
it.onclick=function(e){e.stopPropagation();var arr=(gk()||def).slice(),idx=arr.indexOf(field);idx>=0?arr.splice(idx,1):arr.push(field);go(arr);};
dd.appendChild(it);});
btn.appendChild(dd);btn.onclick=function(e){e.stopPropagation();dd.style.display=dd.style.display=="none"?"block":"none";};
document.addEventListener("click",function(){dd.style.display="none";});
hd.parentNode.appendChild(btn);
})();
</script>`, string(fieldsJSON), string(namesJSON), string(defaultsJSON))
	info.SetHeaderHtml(template2.HTML(colToggleJS))

	info.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	sceneTypeMsg := template2.HTML(biz.T("scene_test_history.scene_type_help"))
	formList := playbookTestHistory.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("dashboard.task_id"), "task_id", db.Varchar, form.Default).
		FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext, form.RichText).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.last_file"), "last_file", db.Varchar, form.Text)
	formList.AddField(biz.T("common.scene_type"), "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
			{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
			{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
			{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
			{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
		}).FieldDefault("1").FieldHelpMsg(sceneTypeMsg)
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.env_type_label"), "env_type", db.Int, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.env_type._1"), Value: "1"},
			{Text: biz.T("common.env_type._2"), Value: "2"},
			{Text: biz.T("common.env_type._3"), Value: "3"},
			{Text: biz.T("common.env_type._4"), Value: "4"},
			{Text: biz.T("common.env_type._5"), Value: "5"},
		}).FieldDefault("2")
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldNowWhenUpdate().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	detail := playbookTestHistory.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.name"), "name", db.Varchar)
	detail.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
		})
	detail.AddField(biz.T("common.last_file"), "last_file", db.Varchar)
	detail.AddField(biz.T("common.scene_type"), "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_test_history.scene_type_1")
			}
			if model.Value == "2" {
				return biz.T("scene_test_history.scene_type_2")
			}
			if model.Value == "3" {
				return biz.T("scene_test_history.scene_type_3")
			}
			if model.Value == "4" {
				return biz.T("scene_test_history.scene_type_4")
			}
			if model.Value == "5" {
				return biz.T("scene_test_history.scene_type_5")
			}
			return biz.T("scene_test_history.scene_type_1")
		})
	detail.AddField(biz.T("common.test_result"), "result", db.Varchar)
	detail.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	detail.AddField(biz.T("common.env_type_label"), "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.env_type._1")
			} else if model.Value == "2" {
				return biz.T("common.env_type._2")
			} else if model.Value == "3" {
				return biz.T("common.env_type._3")
			} else if model.Value == "4" {
				return biz.T("common.env_type._4")
			} else if model.Value == "5" {
				return biz.T("common.env_type._5")
			}
			return ""
		})
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("common.product"), "product", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp)

	detail.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	return playbookTestHistory
}
