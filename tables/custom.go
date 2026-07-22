package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	template2 "html/template"
	"strings"
)

func GetApiCustomDefinitionTable(ctx *context.Context) table.Table {

	apiDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField(biz.T("common.api_module"), "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField(biz.T("common.http_method"), "http_method", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField(biz.T("common.path"), "path", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.header_params"), "header", db.JSON).FieldWidth(300)
	info.AddField(biz.T("custom.path_variable"), "path_variable", db.JSON).FieldWidth(150)
	info.AddField(biz.T("common.query_params"), "query_parameter", db.JSON).FieldWidth(300)
	info.AddField(biz.T("common.body_params"), "body", db.JSON).FieldWidth(300)
	info.AddField(biz.T("custom.response_params"), "response", db.JSON).FieldWidth(300)
	info.AddField(biz.T("custom.check"), "check", db.Varchar).FieldWidth(100).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: biz.T("common.pass")},
		{Value: "fail", Text: biz.T("common.fail")},
	})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("custom.btn_create_scene_data")), icon.Android, action.Ajax("create_batch_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("custom.msg_generated")
					continue
				}
				if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
					status = biz.T("custom.msg_generated")
				} else {
					status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_create_scene_data")), action.Ajax("create_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("custom.btn_create_json_scene_data")), icon.Android, action.Ajax("create_batch_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("custom.msg_generated")
					continue
				}
				if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
					status = biz.T("custom.msg_generated")
				} else {
					status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_create_json_scene_data")), action.Ajax("create_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("custom.btn_create_test_data")), icon.Android, action.Ajax("create_batch_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("custom.msg_test_data_generated")
					continue
				}
				if err := biz.CreateTestData(id); err == nil {
					status = biz.T("custom.msg_test_data_generated")
				} else {
					status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_create_test_data")), action.Ajax("create_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateTestData(id); err == nil {
				status = biz.T("custom.msg_test_data_generated")
			} else {
				status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("custom.btn_create_fuzzing_data")), icon.Android, action.Ajax("create_batch_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("custom.msg_fuzzing_data_generated")
					continue
				}
				if err := biz.CreateFuzzingData(id); err == nil {
					status = biz.T("custom.msg_fuzzing_data_generated")
				} else {
					status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_create_fuzzing_data")), action.Ajax("create_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateFuzzingData(id); err == nil {
				status = biz.T("custom.msg_fuzzing_data_generated")
			} else {
				status = fmt.Sprintf(biz.T("common.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("custom.btn_test_data_test")), icon.Android, action.Ajax("test_data_batch_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("common.msg_test_completed")
					continue
				}
				if err := biz.RunTestData(id); err == nil {
					status = biz.T("common.msg_test_completed")
				} else {
					status = fmt.Sprintf(biz.T("common.msg_test_failed"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_test_data_test")), action.Ajax("test_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = biz.T("common.msg_test_completed")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_test_failed"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("custom.btn_fuzzing_data_test")), icon.Android, action.Ajax("fuzzing_data_batch_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("common.msg_test_completed")
					continue
				}
				if err := biz.RunTestData(id); err == nil {
					status = biz.T("common.msg_test_completed")
				} else {
					status = fmt.Sprintf(biz.T("common.msg_test_failed"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("custom.btn_fuzzing_data_test")), action.Ajax("fuzzing_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = biz.T("common.msg_test_completed")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_test_failed"), id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.AddSelectBox(biz.T("common.http_method"), types.FieldOptions{
		{Value: "get", Text: "get"},
		{Value: "post", Text: "post"},
		{Value: "delete", Text: "delete"},
		{Value: "put", Text: "put"},
	}, action.FieldFilter("http_method"))

	info.AddSelectBox(biz.T("custom.check"), types.FieldOptions{
		{Value: "pass", Text: biz.T("common.pass")},
		{Value: "fail", Text: biz.T("common.fail")},
	}, action.FieldFilter("check"))

	info.SetTable("api_definition").SetTitle(biz.T("custom.title")).SetDescription(biz.T("custom.description"))

	formList := apiDefinition.GetForm()
	formList.SetHTMLContent(`<link rel="stylesheet" href="./static/editor.md/css/editormd.min.css" />
	<div class="layui-layout layui-layout-admin" style="padding-left: 40px;margin-top: 20px;">
	   <form class="layui-form" action="" method="post">
	       <div class="layui-form-item">
	           <label class="layui-form-label">分组资源</label>

	           <div class="layui-input-inline">
	               <select name="source_id" lay-verify="required"  lay-search="" style="width: 100%">
	                   <option value="">所属资源</option>
	                   {{range $index,$elem := .sourceList}}
	                   <option value="{{$elem.Id}}"> {{$elem.GroupName}}--{{$elem.SourceName}}</option>
	                   {{end}}
	               </select>
	           </div>
	       </div>
	       <div class="layui-form-item">
	           <label class="layui-form-label">接口名称</label>
	           <div class="layui-input-block " style="width: 400px">
	               <input type="text" name="api_name" id="api_name" lay-verify="required" autocomplete="off" placeholder="接口名称：获取用户列表" class="layui-input" value="">
	           </div>
	       </div>

	       <div class="layui-form-item">
	           <label class="layui-form-label">接口地址</label>
	           <div class="layui-input-inline" style="width: 400px">
	               <input type="text" name="api_url" id="api_url" lay-verify="required" autocomplete="off" placeholder="/User 注意不写host" class="layui-input" value="">
	           </div>
	       </div>

	       <div class="layui-form-item">
	           <label class="layui-form-label">请求方式</label>
	           <div class="layui-input-block ">
	               <input type="radio" name="method" value="1" title="GET" checked>
	               <input type="radio" name="method" value="2" title="POST">
	               <input type="radio" name="method" value="3" title="PUT">
	               <input type="radio" name="method" value="5" title="DELETE">
	               <input type="radio" name="method" value="4" title="PATCH">
	           </div>
	       </div>

	       <div class="layui-form-item">
	           <label class="layui-form-label">选择模板</label>

	           <div class="layui-input-inline">
	               <select name="templates" lay-filter="templates" lay-verify="required"  lay-search="" style="width: 100%">
	                   <option value="0">选择内容模板</option>
	                   {{range $index,$elem := .templates}}
	                   <option value="{{$elem.Id}}"> {{$elem.TemplateName}}</option>
	                   {{end}}
	               </select>
	           </div>
	       </div>

	       <div class="layui-form-item">
	           <label class="layui-form-label">接口详细</label>
	           <div class="layui-input-inline" id="api-editormd" style="border: 1px solid #e4e4e4">
	<textarea name="detail" style="display:none;">
	#### 简要描述：
	- 用户登录接口


	#### 请求头：

	|参数名|是否必须|类型|说明|
	|:----    |:---|:----- |-----   |
	|Content-Type |是  |string |请求类型： application/json   |
	|Content-MD5 |是  |string | 请求内容签名    |


	#### 请求参数:

	|参数名|是否必须|类型|说明|示例值
	|:----    |:---|:----- |-----   |-----   |
	|username |是  |string |用户名   | george518
	|password |是  |string | 密码    | george518

	#### 返回参数:

	|参数名|类型|说明|
	|:-----  |:-----|-----                           |
	|group_level |int   |用户组id，1：超级管理员；2：普通用户  |

	#### 返回示例:

	**正确时返回:**


		{
			"status": 200,
			"message": "Success",
			"data": {
			"uid": "1",
				"account": "admin",
				"nickname": "Minho",
				"group_level": 0 ,
				"create_time": "1436864169",
				"last_login_time": "0",
		}
		}

	**错误时返回:**

		{
			"status": 300,
			"message": "invalid username"
			"data":{

		}
		}

	#### 调用示例:

	|版本号|制定人|制定日期|修订日期|
	|:----    |:---|:----- |-----   |
	|2.1.0 |郝大全  |2018-01-15 |  2018-01-15 |
	</textarea>
	           </div>
	       </div>
	   </form>
	</div>
	<script>
	var templates = {{.templates}}
	var apiEditor;
	layui.use(['form','element','table','layer','jquery'],function(){
	   var form = layui.form; //只有执行了这一步，部分表单元素才会自动修饰成功
	   var $ = layui.jquery;
	   var table = layui.table;

	   form.on('submit(sub)', function(data){
	       var form_data = data.field;
	       $.post('{{urlfor "ApiController.AjaxSave"}}', form_data, function (out) {
	           if (out.status == 0) {
	               layer.msg("操作成功",{icon: 1},function () {
	                   window.location.href="/api/list";
	               })
	           } else {
	               layer.msg(out.message)
	           }
	       }, "json");
	       return false;
	   });

	   form.on('select(templates)', function(data){
	       layer.confirm('是否切换模板？', {
	           btn: ['确定', '取消'] //可以无限个按钮
	           ,btn2: function(index, layero){
	           layer.close(index);
	       }
	       }, function(index, layero){
	           $.each(templates,function(k,v){
	               if (v.Id==data.value) {
	                   apiEditor.setMarkdown(v.Detail)
	               }
	           });
	           layer.close(index);
	       });

	   });
	 //但是，如果你的HTML是动态生成的，自动渲染就会失效
	 //因此你需要在相应的地方，执行下述方法来手动渲染，跟这类似的还有 element.init();
	 form.render();
	});
	</script>

	{{template "./public/editormd.html" .}}`)

	formList.SetTable("api_definition").SetTitle(biz.T("custom.title")).SetDescription(biz.T("custom.description"))

	return apiDefinition
}
