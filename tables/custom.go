package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"strings"
)

func GetApiCustomDefinitionTable(ctx *context.Context) table.Table {

	apiDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("接口ID", "api_id", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField("所属模块", "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField("接口描述", "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField("请求方法", "http_method", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField("请求路径", "path", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Header参数", "header", db.JSON).FieldWidth(300)
	info.AddField("Path参数", "path_variable", db.JSON).FieldWidth(150)
	info.AddField("Query参数", "query_parameter", db.JSON).FieldWidth(300)
	info.AddField("Body参数", "body", db.JSON).FieldWidth(300)
	info.AddField("Resp参数", "response", db.JSON).FieldWidth(300)
	info.AddField("规范检查", "check", db.Varchar).FieldWidth(100).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("生成场景数据", icon.Android, action.Ajax("create_batch_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "生成完成，请前往[数据文件]列表查看"
					continue
				}
				if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
					status = "生成完成，请前往[数据文件]列表查看"
				} else {
					status = fmt.Sprintf("生成失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("生成场景数据", action.Ajax("create_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
				status = "生成完成，请前往[数据文件]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("生成JSON场景数据", icon.Android, action.Ajax("create_batch_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "生成完成，请前往[数据文件]列表查看"
					continue
				}
				if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
					status = "生成完成，请前往[数据文件]列表查看"
				} else {
					status = fmt.Sprintf("生成失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("生成JSON场景数据", action.Ajax("create_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
				status = "生成完成，请前往[数据文件]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("生成测试数据", icon.Android, action.Ajax("create_batch_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "生成完成，请前往[测试数据]列表查看"
					continue
				}
				if err := biz.CreateTestData(id); err == nil {
					status = "生成完成，请前往[测试数据]列表查看"
				} else {
					status = fmt.Sprintf("生成失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("生成测试数据", action.Ajax("create_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateTestData(id); err == nil {
				status = "生成完成，请前往[测试数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("生成模糊数据", icon.Android, action.Ajax("create_batch_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "生成完成，请前往[模糊数据]列表查看"
					continue
				}
				if err := biz.CreateFuzzingData(id); err == nil {
					status = "生成完成，请前往[模糊数据]列表查看"
				} else {
					status = fmt.Sprintf("生成失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("生成模糊数据", action.Ajax("create_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateFuzzingData(id); err == nil {
				status = "生成完成，请前往[模糊数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("测试数据测试", icon.Android, action.Ajax("test_data_batch_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "测试完成，请前往[结果详情]列表查看"
					continue
				}
				if err := biz.RunTestData(id); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("测试数据测试", action.Ajax("test_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("模糊数据测试", icon.Android, action.Ajax("fuzzing_data_batch_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "测试完成，请前往[结果详情]列表查看"
					continue
				}
				if err := biz.RunTestData(id); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("模糊数据测试", action.Ajax("fuzzing_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.AddSelectBox("请求方法", types.FieldOptions{
		{Value: "get", Text: "get"},
		{Value: "post", Text: "post"},
		{Value: "delete", Text: "delete"},
		{Value: "put", Text: "put"},
	}, action.FieldFilter("http_method"))

	info.AddSelectBox("规范检查", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("check"))

	info.SetTable("api_definition").SetTitle("接口定义").SetDescription("接口定义")

	formList := apiDefinition.GetForm()
	//	formList.AddField("自增主键", "id", db.Int, form.Default).
	//		FieldDisableWhenCreate()
	//	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	//	formList.AddField("所属模块", "api_module", db.Varchar, form.Text)
	//	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	//	formList.AddField("请求方法", "http_method", db.Enum, form.Text).
	//		FieldDefault("get")
	//	formList.AddField("请求路径", "path", db.Varchar, form.Text)
	//	formList.AddField("Header参数", "header", db.Longtext, form.RichText)
	//	formList.AddField("Path参数", "path_variable", db.Longtext, form.RichText)
	//	formList.AddField("Query参数", "query_parameter", db.Longtext, form.RichText)
	//	formList.AddField("Body参数", "body", db.Longtext, form.RichText)
	//	formList.AddField("Resp参数", "response", db.Longtext, form.Custom).
	//		FieldCustomContent(template.HTML(`
	//<span class="input-group-addon"><i class="fa fa-pencil fa-fw"></i></span>
	//<nav class="navbar navbar-default equinav" role="navigation">
	//      <div class="navbar-header">
	//        <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
	//          <span class="icon-bar"></span>
	//          <span class="icon-bar"></span>
	//          <span class="icon-bar"></span>
	//        </button>
	//        <span class="navbar-brand">Menu</span>
	//      </div><!-- /.navbar-header -->
	//      <div class="collapse navbar-collapse">
	//        <ul class="nav navbar-nav">
	//        	<li><a href="#">Home</a></li>
	//            <li><a href="#">How  It Works</a></li>
	//            <li><a href="#" class="dropdown-toggle" data-toggle="dropdown">Features <b class="caret"></b></a>
	//                <ul class="dropdown-menu depth_0">
	//                    <li><a href="#">Feature 1</a></li>
	//                    <li><a href="#">Feature 1</a></li>
	//                    <li><a href="#">Feature 1</a></li>
	//                    <li><a href="#">Feature 1</a></li>
	//                </ul>
	//            </li>
	//            <li><a href="#">What’s Included</a></li>
	//            <li><a href="#">Pricing</a></li>
	//            <li><a href="#">Testimonials</a></li>
	//            <li><a href="#">Case Studies</a></li>
	//            <li><a href="#">Contact Us</a></li>
	//        </ul>
	//      </div><!-- /.navbar-collapse -->
	//  	</nav>`)).
	//		FieldCustomCss(template.CSS(`.equinav .navbar-brand {
	//	display: none;
	//}
	//.equinav .navbar-collapse {
	//	padding: 0 !important;
	//}
	//
	//.equinav-collapse .navbar-header {
	//	float: none;
	//}
	//.equinav-collapse .navbar-brand,
	//.equinav-collapse .navbar-toggle {
	//	display: block !important;
	//}
	//.equina-collapse .navbar-collapse {
	//	border-top: 1px solid #e7e7e7 !important;
	//}
	//.equinav-collapse .navbar-collapse.collapse {
	//	display: none !important;
	//}
	//.equinav-collapse .navbar-nav {
	//	float: none !important;
	//	margin: 0 !important;
	//}
	//.equinav-collapse .navbar-nav > li {
	//	width: auto !important;
	//	float: none !important;
	//}
	//.equinav-collapse .navbar-nav > li > a {
	//	text-align: left !important;
	//	padding-top: 10px !important;
	//	padding-bottom: 10px !important;
	//}
	//.equinav-collapse .navbar-collapse,
	//.equinav-collapse .navbar-collapse.collapse.in {
	//	border-top: 1px solid #e7e7e7 !important;
	//	display: block !important;
	//}
	//.equinav-collapse .collapsing {
	//	overflow: hidden !important;
	//}
	//.equinav-collapse .dropdown-toggle {
	//	background-color: #e7e7e7 !important;
	//}
	//.equinav-collapse .dropdown-toggle > .caret {
	//	display: inline-block !important;
	//}
	//.equinav-collapse .dropdown-menu {
	//	background: none !important;
	//	float: none !important;
	//	border: 0 !important;
	//	box-shadow: none !important;
	//	position: relative !important;
	//}
	//.equinav-collapse .navbar-nav .open .dropdown-menu {
	//	display: block !important;
	//}
	//.equinav-collapse .navbar-nav .open .dropdown-menu > li > a,
	//.equinav-collapse .navbar-nav .open .dropdown-menu .dropdown-header {
	//	color: #777 !important;
	//	padding: 5px 15px 5px 25px !important;
	//}`)).
	//		FieldCustomJs(template.JS(``))
	//	formList.AddField("规范检查", "check", db.Varchar, form.Text)
	//	formList.AddField("所属应用", "app", db.Varchar, form.Text)
	//	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
	//		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	//	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
	//		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	//	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
	//		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
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

	//formList.AddJS(``)
	//	formList.AddCSS(`/*! Editor.md v1.5.0 | editormd.min.css | Open source online markdown editor. | MIT License | By: Pandao | https://github.com/pandao/editor.md | 2015-06-09 */
	//@charset "UTF-8";/*! prefixes.scss v0.1.0 | Author: Pandao | https://github.com/pandao/prefixes.scss | MIT license | Copyright (c) 2015 */.fa-ul,.markdown-body .task-list-item,li.L0,li.L1,li.L2,li.L3,li.L5,li.L6,li.L7,li.L8{list-style-type:none}.editormd-form br,.markdown-body hr:after{clear:both}.editormd{width:90%;height:640px;margin:0 auto 15px;text-align:left;overflow:hidden;position:relative;border:1px solid #ddd;font-family:"Meiryo UI","Microsoft YaHei","Malgun Gothic","Segoe UI","Trebuchet MS",Helvetica,Monaco,monospace,Tahoma,STXihei,"华文细黑",STHeiti,"Helvetica Neue","Droid Sans","wenquanyi micro hei",FreeSans,Arimo,Arial,SimSun,"宋体",Heiti,"黑体",sans-serif}.editormd *,.editormd :after,.editormd :before{-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box}.editormd a{text-decoration:none}.editormd img{border:none;vertical-align:middle}.editormd .editormd-html-textarea,.editormd .editormd-markdown-textarea,.editormd>textarea{width:0;height:0;outline:0;resize:none}.editormd .editormd-html-textarea,.editormd .editormd-markdown-textarea{display:none}.editormd button,.editormd input[type=text],.editormd input[type=button],.editormd input[type=submit],.editormd select,.editormd textarea{-webkit-appearance:none;-moz-appearance:none;-ms-appearance:none;appearance:none}.editormd ::-webkit-scrollbar{height:10px;width:7px;background:rgba(0,0,0,.1)}.editormd ::-webkit-scrollbar:hover{background:rgba(0,0,0,.2)}.editormd ::-webkit-scrollbar-thumb{background:rgba(0,0,0,.3);-webkit-border-radius:6px;-moz-border-radius:6px;-ms-border-radius:6px;-o-border-radius:6px;border-radius:6px}.editormd ::-webkit-scrollbar-thumb:hover{-webkit-box-shadow:inset 1px 1px 1px rgba(0,0,0,.25);-moz-box-shadow:inset 1px 1px 1px rgba(0,0,0,.25);-ms-box-shadow:inset 1px 1px 1px rgba(0,0,0,.25);-o-box-shadow:inset 1px 1px 1px rgba(0,0,0,.25);box-shadow:inset 1px 1px 1px rgba(0,0,0,.25);background-color:rgba(0,0,0,.4)}.editormd-user-unselect{-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none;-o-user-select:none;user-select:none}.editormd-toolbar{width:100%;min-height:37px;background:#fff;display:none;position:absolute;top:0;left:0;z-index:10;border-bottom:1px solid #ddd}.editormd-toolbar-container{padding:0 8px;min-height:35px;-o-user-select:none;user-select:none}.editormd-toolbar-container,.markdown-body .octicon{-webkit-user-select:none;-moz-user-select:none;-ms-user-select:none}.editormd-menu,.markdown-body ol,.markdown-body td,.markdown-body th,.markdown-body ul{padding:0}.editormd-menu{margin:0;list-style:none}.editormd-menu>li{margin:0;padding:5px 1px;display:inline-block;position:relative}.editormd-menu>li.divider{display:inline-block;text-indent:-9999px;margin:0 5px;height:65%;border-right:1px solid #ddd}.editormd-menu>li>a{outline:0;color:#666;display:inline-block;min-width:24px;font-size:16px;text-decoration:none;text-align:center;-webkit-border-radius:2px;-moz-border-radius:2px;-ms-border-radius:2px;-o-border-radius:2px;border-radius:2px;border:1px solid #fff;transition:all 300ms ease-out}.editormd-dropdown-menu>li>a:hover,.editormd-menu>li>a{-webkit-transition:all 300ms ease-out;-moz-transition:all 300ms ease-out}.editormd-menu>li>a.active,.editormd-menu>li>a:hover{border:1px solid #ddd;background:#eee}.editormd-menu>li>a>.fa{text-align:center;display:block;padding:5px}.editormd-menu>li>a>.editormd-bold{padding:5px 2px;display:inline-block;font-weight:700}.editormd-menu>li:hover .editormd-dropdown-menu{display:block}.editormd-menu>li+li>a{margin-left:3px}.editormd-dropdown-menu{display:none;background:#fff;border:1px solid #ddd;width:148px;list-style:none;position:absolute;top:33px;left:0;z-index:100;-webkit-box-shadow:1px 2px 6px rgba(0,0,0,.15);-moz-box-shadow:1px 2px 6px rgba(0,0,0,.15);-ms-box-shadow:1px 2px 6px rgba(0,0,0,.15);-o-box-shadow:1px 2px 6px rgba(0,0,0,.15);box-shadow:1px 2px 6px rgba(0,0,0,.15)}.editormd-dropdown-menu:after,.editormd-dropdown-menu:before{width:0;height:0;display:block;content:"";position:absolute;top:-11px;left:8px;border:5px solid transparent}.editormd-dropdown-menu:before{border-bottom-color:#ccc}.editormd-dropdown-menu:after{border-bottom-color:#fff;top:-10px}.editormd-dropdown-menu>li>a{color:#666;display:block;text-decoration:none;padding:8px 10px}.editormd-dropdown-menu>li>a:hover{background:#f6f6f6;transition:all 300ms ease-out}.editormd-dropdown-menu>li+li{border-top:1px solid #ddd}.editormd-container{margin:0;width:100%;height:100%;overflow:hidden;padding:35px 0 0;position:relative;background:#fff;-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box}.editormd-dialog{color:#666;position:fixed;z-index:99999;display:none;-webkit-border-radius:3px;-moz-border-radius:3px;-ms-border-radius:3px;-o-border-radius:3px;border-radius:3px;-webkit-box-shadow:0 0 10px rgba(0,0,0,.3);-moz-box-shadow:0 0 10px rgba(0,0,0,.3);-ms-box-shadow:0 0 10px rgba(0,0,0,.3);-o-box-shadow:0 0 10px rgba(0,0,0,.3);box-shadow:0 0 10px rgba(0,0,0,.3);background:#fff;font-size:14px}.editormd-dialog-container{position:relative;padding:20px;line-height:1.4}.editormd-dialog-container h1{font-size:24px;margin-bottom:10px}.editormd-dialog-container h1 .fa{color:#2C7EEA;padding-right:5px}.editormd-dialog-container h1 small{padding-left:5px;font-weight:400;font-size:12px;color:#999}.editormd-dialog-container select{color:#999;padding:3px 8px;border:1px solid #ddd}.editormd-dialog-close{position:absolute;top:12px;right:15px;font-size:18px;color:#ccc;-webkit-transition:color 300ms ease-out;-moz-transition:color 300ms ease-out;transition:color 300ms ease-out}.editormd-dialog-close:hover{color:#999}.editormd-dialog-header{padding:11px 20px;border-bottom:1px solid #eee;-webkit-transition:background 300ms ease-out;-moz-transition:background 300ms ease-out;transition:background 300ms ease-out}.editormd-dialog-header:hover{background:#f6f6f6}.editormd-dialog-title{font-size:14px}.editormd-dialog-footer{padding:10px 0 0;text-align:right}.editormd-dialog-info{width:420px}.editormd-dialog-info h1{font-weight:400}.editormd-dialog-info .editormd-dialog-container{padding:20px 25px 25px}.editormd-dialog-info .editormd-dialog-close{top:10px;right:10px}.editormd-dialog-info .hover-link:hover,.editormd-dialog-info p>a{color:#2196F3}.editormd-dialog-info .hover-link{color:#666}.editormd-dialog-info a .fa-external-link{display:none}.editormd-dialog-info a:hover{color:#2196F3}.editormd-dialog-info a:hover .fa-external-link{display:inline-block}.editormd-container-mask,.editormd-dialog-mask,.editormd-mask{display:none;width:100%;height:100%;position:absolute;top:0;left:0}.editormd-dialog-mask-bg,.editormd-mask{background:#fff;opacity:.5;filter:alpha(opacity=50)}.editormd-mask{position:fixed;background:#000;opacity:.2;filter:alpha(opacity=20);z-index:99998}.editormd-container-mask,.editormd-dialog-mask-con{background:url(../images/loading.gif)center center no-repeat;-webkit-background-size:32px 32px;-moz-background-size:32px 32px;-o-background-size:32px 32px;background-size:32px 32px}.editormd-container-mask{z-index:20;display:block;background-color:#fff}@media only screen and (-webkit-min-device-pixel-ratio:2),only screen and (min-device-pixel-ratio:2){.editormd-container-mask,.editormd-dialog-mask-con{background-image:url(../images/loading@2x.gif)}}@media only screen and (-webkit-min-device-pixel-ratio:3),only screen and (min-device-pixel-ratio:3){.editormd-container-mask,.editormd-dialog-mask-con{background-image:url(../images/loading@3x.gif)}}.editormd-code-block-dialog textarea,.editormd-preformatted-text-dialog textarea{width:100%;height:400px;margin-bottom:6px;overflow:auto;border:1px solid #eee;background:#fff;padding:15px;resize:none}.editormd-code-toolbar{color:#999;font-size:14px;margin:-5px 0 10px}.editormd-grid-table{width:99%;display:table;border:1px solid #ddd;border-collapse:collapse}.editormd-grid-table-row{width:100%;display:table-row}.editormd-grid-table-row a{font-size:1.4em;width:5%;height:36px;color:#999;text-align:center;display:table-cell;vertical-align:middle;border:1px solid #ddd;text-decoration:none;-webkit-transition:background-color 300ms ease-out,color 100ms ease-in;-moz-transition:background-color 300ms ease-out,color 100ms ease-in;transition:background-color 300ms ease-out,color 100ms ease-in}.editormd-grid-table-row a.selected{color:#666;background-color:#eee}.editormd-grid-table-row a:hover{color:#777;background-color:#f6f6f6}.editormd-tab-head{list-style:none;border-bottom:1px solid #ddd}.editormd-tab-head li{display:inline-block}.editormd-tab-head li a{color:#999;display:block;padding:6px 12px 5px;text-align:center;text-decoration:none;margin-bottom:-1px;border:1px solid #ddd;-webkit-border-top-left-radius:3px;-moz-border-top-left-radius:3px;-ms-border-top-left-radius:3px;-o-border-top-left-radius:3px;border-top-left-radius:3px;-webkit-border-top-right-radius:3px;-moz-border-top-right-radius:3px;-ms-border-top-right-radius:3px;-o-border-top-right-radius:3px;border-top-right-radius:3px;background:#f6f6f6;-webkit-transition:all 300ms ease-out;-moz-transition:all 300ms ease-out;transition:all 300ms ease-out}.editormd-tab-head li a:hover{color:#666;background:#eee}.editormd-tab-head li.active a{color:#666;background:#fff;border-bottom-color:#fff}.editormd-tab-head li+li{margin-left:3px}.editormd-tab-box{padding:20px 0}.editormd-form{color:#666}.editormd-form label{float:left;display:block;width:75px;text-align:left;padding:7px 0 15px 5px;margin:0 0 2px;font-weight:400}.editormd-form iframe{display:none}.editormd-form input:focus{outline:0}.editormd-form input[type=text],.editormd-form input[type=number]{color:#999;padding:8px;border:1px solid #ddd}.editormd-form input[type=number]{width:40px;display:inline-block;padding:6px 8px}.editormd-form input[type=text]{display:inline-block;width:264px}.editormd-form .fa-btns{display:inline-block}.editormd-form .fa-btns a{color:#999;padding:7px 10px 0 0;display:inline-block;text-decoration:none;text-align:center}.editormd-form .fa-btns .fa{font-size:1.3em}.editormd-form .fa-btns label{float:none;display:inline-block;width:auto;text-align:left;padding:0 0 0 5px;cursor:pointer}.fa-fw,.fa-li{text-align:center}.editormd-dialog-container .editormd-btn,.editormd-dialog-container button,.editormd-dialog-container input[type=submit],.editormd-dialog-footer .editormd-btn,.editormd-dialog-footer button,.editormd-dialog-footer input[type=submit],.editormd-form .editormd-btn,.editormd-form button,.editormd-form input[type=submit]{color:#666;min-width:75px;cursor:pointer;background:#fff;padding:7px 10px;border:1px solid #ddd;-webkit-border-radius:3px;-moz-border-radius:3px;-ms-border-radius:3px;-o-border-radius:3px;border-radius:3px;-webkit-transition:background 300ms ease-out;-moz-transition:background 300ms ease-out;transition:background 300ms ease-out}.editormd-dialog-container .editormd-btn:hover,.editormd-dialog-container button:hover,.editormd-dialog-container input[type=submit]:hover,.editormd-dialog-footer .editormd-btn:hover,.editormd-dialog-footer button:hover,.editormd-dialog-footer input[type=submit]:hover,.editormd-form .editormd-btn:hover,.editormd-form button:hover,.editormd-form input[type=submit]:hover{background:#eee}.editormd-dialog-container .editormd-btn+.editormd-btn,.editormd-dialog-footer .editormd-btn+.editormd-btn,.editormd-form .editormd-btn+.editormd-btn{margin-left:8px}.editormd-file-input{width:75px;height:32px;margin-left:8px;position:relative;display:inline-block}.editormd-file-input input[type=file]{width:75px;height:32px;opacity:0;cursor:pointer;background:#000;display:inline-block;position:absolute;top:0;right:0}.editormd-file-input input[type=file]::-webkit-file-upload-button{visibility:hidden}.editormd-file-input:hover input[type=submit]{background:#eee}.editormd .CodeMirror,.editormd-preview{display:inline-block;width:50%;height:100%;vertical-align:top;-webkit-box-sizing:border-box;-moz-box-sizing:border-box;box-sizing:border-box;margin:0}.editormd-preview{position:absolute;top:35px;right:0;overflow:auto;line-height:1.6;display:none;background:#fff}.fa,.fa-stack{display:inline-block}.editormd .CodeMirror{z-index:10;float:left;border-right:1px solid #ddd;font-size:14px;font-family:"YaHei Consolas Hybrid",Consolas,"微软雅黑","Meiryo UI","Malgun Gothic","Segoe UI","Trebuchet MS",Helvetica,Monaco,courier,monospace;line-height:1.6;margin-top:35px}.editormd .CodeMirror pre{font-size:14px;padding:0 12px}.editormd .CodeMirror-linenumbers{padding:0 5px}.editormd .CodeMirror-focused .CodeMirror-selected,.editormd .CodeMirror-selected{background:#70B7FF}.editormd .CodeMirror,.editormd .CodeMirror-scroll,.editormd .editormd-preview{-webkit-overflow-scrolling:touch}.editormd .styled-background{background-color:#ff7}.editormd .CodeMirror-focused .cm-matchhighlight{background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAFklEQVQI12NgYGBgkKzc8x9CMDAwAAAmhwSbidEoSQAAAABJRU5ErkJggg==);background-position:bottom;background-repeat:repeat-x}.editormd .CodeMirror-empty.CodeMirror-focused{outline:0}.editormd .CodeMirror pre.CodeMirror-placeholder{color:#999}.editormd .cm-trailingspace{background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAACCAYAAAB/qH1jAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH3QUXCToH00Y1UgAAACFJREFUCNdjPMDBUc/AwNDAAAFMTAwMDA0OP34wQgX/AQBYgwYEx4f9lQAAAABJRU5ErkJggg==);background-position:bottom left;background-repeat:repeat-x}.editormd .cm-tab{background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAMCAYAAAAkuj5RAAAAAXNSR0IArs4c6QAAAGFJREFUSMft1LsRQFAQheHPowAKoACx3IgEKtaEHujDjORSgWTH/ZOdnZOcM/sgk/kFFWY0qV8foQwS4MKBCS3qR6ixBJvElOobYAtivseIE120FaowJPN75GMu8j/LfMwNjh4HUpwg4LUAAAAASUVORK5CYII=)right no-repeat}/*! prefixes.scss v0.1.0 | Author: Pandao | https://github.com/pandao/prefixes.scss | MIT license | Copyright (c) 2015 *//*!
	// *  Font Awesome 4.3.0 by @davegandy - http://fontawesome.io - @fontawesome
	// *  License - http://fontawesome.io/license (Font: SIL OFL 1.1, CSS: MIT License)
	// */@font-face{font-family:FontAwesome;src:url(../fonts/fontawesome-webfont.eot?v=4.3.0);src:url(../fonts/fontawesome-webfont.eot?#iefix&v=4.3.0)format("embedded-opentype"),url(../fonts/fontawesome-webfont.woff2?v=4.3.0)format("woff2"),url(../fonts/fontawesome-webfont.woff?v=4.3.0)format("woff"),url(../fonts/fontawesome-webfont.ttf?v=4.3.0)format("truetype"),url(../fonts/fontawesome-webfont.svg?v=4.3.0#fontawesomeregular)format("svg");font-weight:400;font-style:normal}.fa{font:normal normal normal 14px/1 FontAwesome;font-size:inherit;text-rendering:auto;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale;transform:translate(0,0)}.fa-lg{font-size:1.33333333em;line-height:.75em;vertical-align:-15%}.fa-2x{font-size:2em}.fa-3x{font-size:3em}.fa-4x{font-size:4em}.fa-5x{font-size:5em}.fa-fw{width:1.28571429em}.fa-ul{padding-left:0;margin-left:2.14285714em}.fa-ul>li{position:relative}.fa-li{position:absolute;left:-2.14285714em;width:2.14285714em;top:.14285714em}.fa-li.fa-lg{left:-1.85714286em}.fa-border{padding:.2em .25em .15em;border:.08em solid #eee;border-radius:.1em}.pull-right{float:right}.pull-left{float:left}.fa.pull-left{margin-right:.3em}.fa.pull-right{margin-left:.3em}.fa-spin{-webkit-animation:fa-spin 2s infinite linear;animation:fa-spin 2s infinite linear}.fa-pulse{-webkit-animation:fa-spin 1s infinite steps(8);animation:fa-spin 1s infinite steps(8)}@-webkit-keyframes fa-spin{0%{-webkit-transform:rotate(0);transform:rotate(0)}100%{-webkit-transform:rotate(359deg);transform:rotate(359deg)}}@keyframes fa-spin{0%{-webkit-transform:rotate(0);transform:rotate(0)}100%{-webkit-transform:rotate(359deg);transform:rotate(359deg)}}.fa-rotate-90{filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=1);-webkit-transform:rotate(90deg);-ms-transform:rotate(90deg);transform:rotate(90deg)}.fa-rotate-180{filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=2);-webkit-transform:rotate(180deg);-ms-transform:rotate(180deg);transform:rotate(180deg)}.fa-rotate-270{filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=3);-webkit-transform:rotate(270deg);-ms-transform:rotate(270deg);transform:rotate(270deg)}.fa-flip-horizontal{filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=0, mirror=1);-webkit-transform:scale(-1,1);-ms-transform:scale(-1,1);transform:scale(-1,1)}.fa-flip-vertical{filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=2, mirror=1);-webkit-transform:scale(1,-1);-ms-transform:scale(1,-1);transform:scale(1,-1)}:root .fa-flip-horizontal,:root .fa-flip-vertical,:root .fa-rotate-180,:root .fa-rotate-270,:root .fa-rotate-90{filter:none}.fa-stack{position:relative;width:2em;height:2em;line-height:2em;vertical-align:middle}.fa-stack-1x,.fa-stack-2x{position:absolute;left:0;width:100%;text-align:center}.fa-stack-1x{line-height:inherit}.fa-stack-2x{font-size:2em}.fa-inverse{color:#fff}.fa-glass:before{content:"\f000"}.fa-music:before{content:"\f001"}.fa-search:before{content:"\f002"}.fa-envelope-o:before{content:"\f003"}.fa-heart:before{content:"\f004"}.fa-star:before{content:"\f005"}.fa-star-o:before{content:"\f006"}.fa-user:before{content:"\f007"}.fa-film:before{content:"\f008"}.fa-th-large:before{content:"\f009"}.fa-th:before{content:"\f00a"}.fa-th-list:before{content:"\f00b"}.fa-check:before{content:"\f00c"}.fa-close:before,.fa-remove:before,.fa-times:before{content:"\f00d"}.fa-search-plus:before{content:"\f00e"}.fa-search-minus:before{content:"\f010"}.fa-power-off:before{content:"\f011"}.fa-signal:before{content:"\f012"}.fa-cog:before,.fa-gear:before{content:"\f013"}.fa-trash-o:before{content:"\f014"}.fa-home:before{content:"\f015"}.fa-file-o:before{content:"\f016"}.fa-clock-o:before{content:"\f017"}.fa-road:before{content:"\f018"}.fa-download:before{content:"\f019"}.fa-arrow-circle-o-down:before{content:"\f01a"}.fa-arrow-circle-o-up:before{content:"\f01b"}.fa-inbox:before{content:"\f01c"}.fa-play-circle-o:before{content:"\f01d"}.fa-repeat:before,.fa-rotate-right:before{content:"\f01e"}.fa-refresh:before{content:"\f021"}.fa-list-alt:before{content:"\f022"}.fa-lock:before{content:"\f023"}.fa-flag:before{content:"\f024"}.fa-headphones:before{content:"\f025"}.fa-volume-off:before{content:"\f026"}.fa-volume-down:before{content:"\f027"}.fa-volume-up:before{content:"\f028"}.fa-qrcode:before{content:"\f029"}.fa-barcode:before{content:"\f02a"}.fa-tag:before{content:"\f02b"}.fa-tags:before{content:"\f02c"}.fa-book:before{content:"\f02d"}.fa-bookmark:before{content:"\f02e"}.fa-print:before{content:"\f02f"}.fa-camera:before{content:"\f030"}.fa-font:before{content:"\f031"}.fa-bold:before{content:"\f032"}.fa-italic:before{content:"\f033"}.fa-text-height:before{content:"\f034"}.fa-text-width:before{content:"\f035"}.fa-align-left:before{content:"\f036"}.fa-align-center:before{content:"\f037"}.fa-align-right:before{content:"\f038"}.fa-align-justify:before{content:"\f039"}.fa-list:before{content:"\f03a"}.fa-dedent:before,.fa-outdent:before{content:"\f03b"}.fa-indent:before{content:"\f03c"}.fa-video-camera:before{content:"\f03d"}.fa-image:before,.fa-photo:before,.fa-picture-o:before{content:"\f03e"}.fa-pencil:before{content:"\f040"}.fa-map-marker:before{content:"\f041"}.fa-adjust:before{content:"\f042"}.fa-tint:before{content:"\f043"}.fa-edit:before,.fa-pencil-square-o:before{content:"\f044"}.fa-share-square-o:before{content:"\f045"}.fa-check-square-o:before{content:"\f046"}.fa-arrows:before{content:"\f047"}.fa-step-backward:before{content:"\f048"}.fa-fast-backward:before{content:"\f049"}.fa-backward:before{content:"\f04a"}.fa-play:before{content:"\f04b"}.fa-pause:before{content:"\f04c"}.fa-stop:before{content:"\f04d"}.fa-forward:before{content:"\f04e"}.fa-fast-forward:before{content:"\f050"}.fa-step-forward:before{content:"\f051"}.fa-eject:before{content:"\f052"}.fa-chevron-left:before{content:"\f053"}.fa-chevron-right:before{content:"\f054"}.fa-plus-circle:before{content:"\f055"}.fa-minus-circle:before{content:"\f056"}.fa-times-circle:before{content:"\f057"}.fa-check-circle:before{content:"\f058"}.fa-question-circle:before{content:"\f059"}.fa-info-circle:before{content:"\f05a"}.fa-crosshairs:before{content:"\f05b"}.fa-times-circle-o:before{content:"\f05c"}.fa-check-circle-o:before{content:"\f05d"}.fa-ban:before{content:"\f05e"}.fa-arrow-left:before{content:"\f060"}.fa-arrow-right:before{content:"\f061"}.fa-arrow-up:before{content:"\f062"}.fa-arrow-down:before{content:"\f063"}.fa-mail-forward:before,.fa-share:before{content:"\f064"}.fa-expand:before{content:"\f065"}.fa-compress:before{content:"\f066"}.fa-plus:before{content:"\f067"}.fa-minus:before{content:"\f068"}.fa-asterisk:before{content:"\f069"}.fa-exclamation-circle:before{content:"\f06a"}.fa-gift:before{content:"\f06b"}.fa-leaf:before{content:"\f06c"}.fa-fire:before{content:"\f06d"}.fa-eye:before{content:"\f06e"}.fa-eye-slash:before{content:"\f070"}.fa-exclamation-triangle:before,.fa-warning:before{content:"\f071"}.fa-plane:before{content:"\f072"}.fa-calendar:before{content:"\f073"}.fa-random:before{content:"\f074"}.fa-comment:before{content:"\f075"}.fa-magnet:before{content:"\f076"}.fa-chevron-up:before{content:"\f077"}.fa-chevron-down:before{content:"\f078"}.fa-retweet:before{content:"\f079"}.fa-shopping-cart:before{content:"\f07a"}.fa-folder:before{content:"\f07b"}.fa-folder-open:before{content:"\f07c"}.fa-arrows-v:before{content:"\f07d"}.fa-arrows-h:before{content:"\f07e"}.fa-bar-chart-o:before,.fa-bar-chart:before{content:"\f080"}.fa-twitter-square:before{content:"\f081"}.fa-facebook-square:before{content:"\f082"}.fa-camera-retro:before{content:"\f083"}.fa-key:before{content:"\f084"}.fa-cogs:before,.fa-gears:before{content:"\f085"}.fa-comments:before{content:"\f086"}.fa-thumbs-o-up:before{content:"\f087"}.fa-thumbs-o-down:before{content:"\f088"}.fa-star-half:before{content:"\f089"}.fa-heart-o:before{content:"\f08a"}.fa-sign-out:before{content:"\f08b"}.fa-linkedin-square:before{content:"\f08c"}.fa-thumb-tack:before{content:"\f08d"}.fa-external-link:before{content:"\f08e"}.fa-sign-in:before{content:"\f090"}.fa-trophy:before{content:"\f091"}.fa-github-square:before{content:"\f092"}.fa-upload:before{content:"\f093"}.fa-lemon-o:before{content:"\f094"}.fa-phone:before{content:"\f095"}.fa-square-o:before{content:"\f096"}.fa-bookmark-o:before{content:"\f097"}.fa-phone-square:before{content:"\f098"}.fa-twitter:before{content:"\f099"}.fa-facebook-f:before,.fa-facebook:before{content:"\f09a"}.fa-github:before{content:"\f09b"}.fa-unlock:before{content:"\f09c"}.fa-credit-card:before{content:"\f09d"}.fa-rss:before{content:"\f09e"}.fa-hdd-o:before{content:"\f0a0"}.fa-bullhorn:before{content:"\f0a1"}.fa-bell:before{content:"\f0f3"}.fa-certificate:before{content:"\f0a3"}.fa-hand-o-right:before{content:"\f0a4"}.fa-hand-o-left:before{content:"\f0a5"}.fa-hand-o-up:before{content:"\f0a6"}.fa-hand-o-down:before{content:"\f0a7"}.fa-arrow-circle-left:before{content:"\f0a8"}.fa-arrow-circle-right:before{content:"\f0a9"}.fa-arrow-circle-up:before{content:"\f0aa"}.fa-arrow-circle-down:before{content:"\f0ab"}.fa-globe:before{content:"\f0ac"}.fa-wrench:before{content:"\f0ad"}.fa-tasks:before{content:"\f0ae"}.fa-filter:before{content:"\f0b0"}.fa-briefcase:before{content:"\f0b1"}.fa-arrows-alt:before{content:"\f0b2"}.fa-group:before,.fa-users:before{content:"\f0c0"}.fa-chain:before,.fa-link:before{content:"\f0c1"}.fa-cloud:before{content:"\f0c2"}.fa-flask:before{content:"\f0c3"}.fa-cut:before,.fa-scissors:before{content:"\f0c4"}.fa-copy:before,.fa-files-o:before{content:"\f0c5"}.fa-paperclip:before{content:"\f0c6"}.fa-floppy-o:before,.fa-save:before{content:"\f0c7"}.fa-square:before{content:"\f0c8"}.fa-bars:before,.fa-navicon:before,.fa-reorder:before{content:"\f0c9"}.fa-list-ul:before{content:"\f0ca"}.fa-list-ol:before{content:"\f0cb"}.fa-strikethrough:before{content:"\f0cc"}.fa-underline:before{content:"\f0cd"}.fa-table:before{content:"\f0ce"}.fa-magic:before{content:"\f0d0"}.fa-truck:before{content:"\f0d1"}.fa-pinterest:before{content:"\f0d2"}.fa-pinterest-square:before{content:"\f0d3"}.fa-google-plus-square:before{content:"\f0d4"}.fa-google-plus:before{content:"\f0d5"}.fa-money:before{content:"\f0d6"}.fa-caret-down:before{content:"\f0d7"}.fa-caret-up:before{content:"\f0d8"}.fa-caret-left:before{content:"\f0d9"}.fa-caret-right:before{content:"\f0da"}.fa-columns:before{content:"\f0db"}.fa-sort:before,.fa-unsorted:before{content:"\f0dc"}.fa-sort-desc:before,.fa-sort-down:before{content:"\f0dd"}.fa-sort-asc:before,.fa-sort-up:before{content:"\f0de"}.fa-envelope:before{content:"\f0e0"}.fa-linkedin:before{content:"\f0e1"}.fa-rotate-left:before,.fa-undo:before{content:"\f0e2"}.fa-gavel:before,.fa-legal:before{content:"\f0e3"}.fa-dashboard:before,.fa-tachometer:before{content:"\f0e4"}.fa-comment-o:before{content:"\f0e5"}.fa-comments-o:before{content:"\f0e6"}.fa-bolt:before,.fa-flash:before{content:"\f0e7"}.fa-sitemap:before{content:"\f0e8"}.fa-umbrella:before{content:"\f0e9"}.fa-clipboard:before,.fa-paste:before{content:"\f0ea"}.fa-lightbulb-o:before{content:"\f0eb"}.fa-exchange:before{content:"\f0ec"}.fa-cloud-download:before{content:"\f0ed"}.fa-cloud-upload:before{content:"\f0ee"}.fa-user-md:before{content:"\f0f0"}.fa-stethoscope:before{content:"\f0f1"}.fa-suitcase:before{content:"\f0f2"}.fa-bell-o:before{content:"\f0a2"}.fa-coffee:before{content:"\f0f4"}.fa-cutlery:before{content:"\f0f5"}.fa-file-text-o:before{content:"\f0f6"}.fa-building-o:before{content:"\f0f7"}.fa-hospital-o:before{content:"\f0f8"}.fa-ambulance:before{content:"\f0f9"}.fa-medkit:before{content:"\f0fa"}.fa-fighter-jet:before{content:"\f0fb"}.fa-beer:before{content:"\f0fc"}.fa-h-square:before{content:"\f0fd"}.fa-plus-square:before{content:"\f0fe"}.fa-angle-double-left:before{content:"\f100"}.fa-angle-double-right:before{content:"\f101"}.fa-angle-double-up:before{content:"\f102"}.fa-angle-double-down:before{content:"\f103"}.fa-angle-left:before{content:"\f104"}.fa-angle-right:before{content:"\f105"}.fa-angle-up:before{content:"\f106"}.fa-angle-down:before{content:"\f107"}.fa-desktop:before{content:"\f108"}.fa-laptop:before{content:"\f109"}.fa-tablet:before{content:"\f10a"}.fa-mobile-phone:before,.fa-mobile:before{content:"\f10b"}.fa-circle-o:before{content:"\f10c"}.fa-quote-left:before{content:"\f10d"}.fa-quote-right:before{content:"\f10e"}.fa-spinner:before{content:"\f110"}.fa-circle:before{content:"\f111"}.fa-mail-reply:before,.fa-reply:before{content:"\f112"}.fa-github-alt:before{content:"\f113"}.fa-folder-o:before{content:"\f114"}.fa-folder-open-o:before{content:"\f115"}.fa-smile-o:before{content:"\f118"}.fa-frown-o:before{content:"\f119"}.fa-meh-o:before{content:"\f11a"}.fa-gamepad:before{content:"\f11b"}.fa-keyboard-o:before{content:"\f11c"}.fa-flag-o:before{content:"\f11d"}.fa-flag-checkered:before{content:"\f11e"}.fa-terminal:before{content:"\f120"}.fa-code:before{content:"\f121"}.fa-mail-reply-all:before,.fa-reply-all:before{content:"\f122"}.fa-star-half-empty:before,.fa-star-half-full:before,.fa-star-half-o:before{content:"\f123"}.fa-location-arrow:before{content:"\f124"}.fa-crop:before{content:"\f125"}.fa-code-fork:before{content:"\f126"}.fa-chain-broken:before,.fa-unlink:before{content:"\f127"}.fa-question:before{content:"\f128"}.fa-info:before{content:"\f129"}.fa-exclamation:before{content:"\f12a"}.fa-superscript:before{content:"\f12b"}.fa-subscript:before{content:"\f12c"}.fa-eraser:before{content:"\f12d"}.fa-puzzle-piece:before{content:"\f12e"}.fa-microphone:before{content:"\f130"}.fa-microphone-slash:before{content:"\f131"}.fa-shield:before{content:"\f132"}.fa-calendar-o:before{content:"\f133"}.fa-fire-extinguisher:before{content:"\f134"}.fa-rocket:before{content:"\f135"}.fa-maxcdn:before{content:"\f136"}.fa-chevron-circle-left:before{content:"\f137"}.fa-chevron-circle-right:before{content:"\f138"}.fa-chevron-circle-up:before{content:"\f139"}.fa-chevron-circle-down:before{content:"\f13a"}.fa-html5:before{content:"\f13b"}.fa-css3:before{content:"\f13c"}.fa-anchor:before{content:"\f13d"}.fa-unlock-alt:before{content:"\f13e"}.fa-bullseye:before{content:"\f140"}.fa-ellipsis-h:before{content:"\f141"}.fa-ellipsis-v:before{content:"\f142"}.fa-rss-square:before{content:"\f143"}.fa-play-circle:before{content:"\f144"}.fa-ticket:before{content:"\f145"}.fa-minus-square:before{content:"\f146"}.fa-minus-square-o:before{content:"\f147"}.fa-level-up:before{content:"\f148"}.fa-level-down:before{content:"\f149"}.fa-check-square:before{content:"\f14a"}.fa-pencil-square:before{content:"\f14b"}.fa-external-link-square:before{content:"\f14c"}.fa-share-square:before{content:"\f14d"}.fa-compass:before{content:"\f14e"}.fa-caret-square-o-down:before,.fa-toggle-down:before{content:"\f150"}.fa-caret-square-o-up:before,.fa-toggle-up:before{content:"\f151"}.fa-caret-square-o-right:before,.fa-toggle-right:before{content:"\f152"}.fa-eur:before,.fa-euro:before{content:"\f153"}.fa-gbp:before{content:"\f154"}.fa-dollar:before,.fa-usd:before{content:"\f155"}.fa-inr:before,.fa-rupee:before{content:"\f156"}.fa-cny:before,.fa-jpy:before,.fa-rmb:before,.fa-yen:before{content:"\f157"}.fa-rouble:before,.fa-rub:before,.fa-ruble:before{content:"\f158"}.fa-krw:before,.fa-won:before{content:"\f159"}.fa-bitcoin:before,.fa-btc:before{content:"\f15a"}.fa-file:before{content:"\f15b"}.fa-file-text:before{content:"\f15c"}.fa-sort-alpha-asc:before{content:"\f15d"}.fa-sort-alpha-desc:before{content:"\f15e"}.fa-sort-amount-asc:before{content:"\f160"}.fa-sort-amount-desc:before{content:"\f161"}.fa-sort-numeric-asc:before{content:"\f162"}.fa-sort-numeric-desc:before{content:"\f163"}.fa-thumbs-up:before{content:"\f164"}.fa-thumbs-down:before{content:"\f165"}.fa-youtube-square:before{content:"\f166"}.fa-youtube:before{content:"\f167"}.fa-xing:before{content:"\f168"}.fa-xing-square:before{content:"\f169"}.fa-youtube-play:before{content:"\f16a"}.fa-dropbox:before{content:"\f16b"}.fa-stack-overflow:before{content:"\f16c"}.fa-instagram:before{content:"\f16d"}.fa-flickr:before{content:"\f16e"}.fa-adn:before{content:"\f170"}.fa-bitbucket:before{content:"\f171"}.fa-bitbucket-square:before{content:"\f172"}.fa-tumblr:before{content:"\f173"}.fa-tumblr-square:before{content:"\f174"}.fa-long-arrow-down:before{content:"\f175"}.fa-long-arrow-up:before{content:"\f176"}.fa-long-arrow-left:before{content:"\f177"}.fa-long-arrow-right:before{content:"\f178"}.fa-apple:before{content:"\f179"}.fa-windows:before{content:"\f17a"}.fa-android:before{content:"\f17b"}.fa-linux:before{content:"\f17c"}.fa-dribbble:before{content:"\f17d"}.fa-skype:before{content:"\f17e"}.fa-foursquare:before{content:"\f180"}.fa-trello:before{content:"\f181"}.fa-female:before{content:"\f182"}.fa-male:before{content:"\f183"}.fa-gittip:before,.fa-gratipay:before{content:"\f184"}.fa-sun-o:before{content:"\f185"}.fa-moon-o:before{content:"\f186"}.fa-archive:before{content:"\f187"}.fa-bug:before{content:"\f188"}.fa-vk:before{content:"\f189"}.fa-weibo:before{content:"\f18a"}.fa-renren:before{content:"\f18b"}.fa-pagelines:before{content:"\f18c"}.fa-stack-exchange:before{content:"\f18d"}.fa-arrow-circle-o-right:before{content:"\f18e"}.fa-arrow-circle-o-left:before{content:"\f190"}.fa-caret-square-o-left:before,.fa-toggle-left:before{content:"\f191"}.fa-dot-circle-o:before{content:"\f192"}.fa-wheelchair:before{content:"\f193"}.fa-vimeo-square:before{content:"\f194"}.fa-try:before,.fa-turkish-lira:before{content:"\f195"}.fa-plus-square-o:before{content:"\f196"}.fa-space-shuttle:before{content:"\f197"}.fa-slack:before{content:"\f198"}.fa-envelope-square:before{content:"\f199"}.fa-wordpress:before{content:"\f19a"}.fa-openid:before{content:"\f19b"}.fa-bank:before,.fa-institution:before,.fa-university:before{content:"\f19c"}.fa-graduation-cap:before,.fa-mortar-board:before{content:"\f19d"}.fa-yahoo:before{content:"\f19e"}.fa-google:before{content:"\f1a0"}.fa-reddit:before{content:"\f1a1"}.fa-reddit-square:before{content:"\f1a2"}.fa-stumbleupon-circle:before{content:"\f1a3"}.fa-stumbleupon:before{content:"\f1a4"}.fa-delicious:before{content:"\f1a5"}.fa-digg:before{content:"\f1a6"}.fa-pied-piper:before{content:"\f1a7"}.fa-pied-piper-alt:before{content:"\f1a8"}.fa-drupal:before{content:"\f1a9"}.fa-joomla:before{content:"\f1aa"}.fa-language:before{content:"\f1ab"}.fa-fax:before{content:"\f1ac"}.fa-building:before{content:"\f1ad"}.fa-child:before{content:"\f1ae"}.fa-paw:before{content:"\f1b0"}.fa-spoon:before{content:"\f1b1"}.fa-cube:before{content:"\f1b2"}.fa-cubes:before{content:"\f1b3"}.fa-behance:before{content:"\f1b4"}.fa-behance-square:before{content:"\f1b5"}.fa-steam:before{content:"\f1b6"}.fa-steam-square:before{content:"\f1b7"}.fa-recycle:before{content:"\f1b8"}.fa-automobile:before,.fa-car:before{content:"\f1b9"}.fa-cab:before,.fa-taxi:before{content:"\f1ba"}.fa-tree:before{content:"\f1bb"}.fa-spotify:before{content:"\f1bc"}.fa-deviantart:before{content:"\f1bd"}.fa-soundcloud:before{content:"\f1be"}.fa-database:before{content:"\f1c0"}.fa-file-pdf-o:before{content:"\f1c1"}.fa-file-word-o:before{content:"\f1c2"}.fa-file-excel-o:before{content:"\f1c3"}.fa-file-powerpoint-o:before{content:"\f1c4"}.fa-file-image-o:before,.fa-file-photo-o:before,.fa-file-picture-o:before{content:"\f1c5"}.fa-file-archive-o:before,.fa-file-zip-o:before{content:"\f1c6"}.fa-file-audio-o:before,.fa-file-sound-o:before{content:"\f1c7"}.fa-file-movie-o:before,.fa-file-video-o:before{content:"\f1c8"}.fa-file-code-o:before{content:"\f1c9"}.fa-vine:before{content:"\f1ca"}.fa-codepen:before{content:"\f1cb"}.fa-jsfiddle:before{content:"\f1cc"}.fa-life-bouy:before,.fa-life-buoy:before,.fa-life-ring:before,.fa-life-saver:before,.fa-support:before{content:"\f1cd"}.fa-circle-o-notch:before{content:"\f1ce"}.fa-ra:before,.fa-rebel:before{content:"\f1d0"}.fa-empire:before,.fa-ge:before{content:"\f1d1"}.fa-git-square:before{content:"\f1d2"}.fa-git:before{content:"\f1d3"}.fa-hacker-news:before{content:"\f1d4"}.fa-tencent-weibo:before{content:"\f1d5"}.fa-qq:before{content:"\f1d6"}.fa-wechat:before,.fa-weixin:before{content:"\f1d7"}.fa-paper-plane:before,.fa-send:before{content:"\f1d8"}.fa-paper-plane-o:before,.fa-send-o:before{content:"\f1d9"}.fa-history:before{content:"\f1da"}.fa-circle-thin:before,.fa-genderless:before{content:"\f1db"}.fa-header:before{content:"\f1dc"}.fa-paragraph:before{content:"\f1dd"}.fa-sliders:before{content:"\f1de"}.fa-share-alt:before{content:"\f1e0"}.fa-share-alt-square:before{content:"\f1e1"}.fa-bomb:before{content:"\f1e2"}.fa-futbol-o:before,.fa-soccer-ball-o:before{content:"\f1e3"}.fa-tty:before{content:"\f1e4"}.fa-binoculars:before{content:"\f1e5"}.fa-plug:before{content:"\f1e6"}.fa-slideshare:before{content:"\f1e7"}.fa-twitch:before{content:"\f1e8"}.fa-yelp:before{content:"\f1e9"}.fa-newspaper-o:before{content:"\f1ea"}.fa-wifi:before{content:"\f1eb"}.fa-calculator:before{content:"\f1ec"}.fa-paypal:before{content:"\f1ed"}.fa-google-wallet:before{content:"\f1ee"}.fa-cc-visa:before{content:"\f1f0"}.fa-cc-mastercard:before{content:"\f1f1"}.fa-cc-discover:before{content:"\f1f2"}.fa-cc-amex:before{content:"\f1f3"}.fa-cc-paypal:before{content:"\f1f4"}.fa-cc-stripe:before{content:"\f1f5"}.fa-bell-slash:before{content:"\f1f6"}.fa-bell-slash-o:before{content:"\f1f7"}.fa-trash:before{content:"\f1f8"}.fa-copyright:before{content:"\f1f9"}.fa-at:before{content:"\f1fa"}.fa-eyedropper:before{content:"\f1fb"}.fa-paint-brush:before{content:"\f1fc"}.fa-birthday-cake:before{content:"\f1fd"}.fa-area-chart:before{content:"\f1fe"}.fa-pie-chart:before{content:"\f200"}.fa-line-chart:before{content:"\f201"}.fa-lastfm:before{content:"\f202"}.fa-lastfm-square:before{content:"\f203"}.fa-toggle-off:before{content:"\f204"}.fa-toggle-on:before{content:"\f205"}.fa-bicycle:before{content:"\f206"}.fa-bus:before{content:"\f207"}.fa-ioxhost:before{content:"\f208"}.fa-angellist:before{content:"\f209"}.fa-cc:before{content:"\f20a"}.fa-ils:before,.fa-shekel:before,.fa-sheqel:before{content:"\f20b"}.fa-meanpath:before{content:"\f20c"}.fa-buysellads:before{content:"\f20d"}.fa-connectdevelop:before{content:"\f20e"}.fa-dashcube:before{content:"\f210"}.fa-forumbee:before{content:"\f211"}.fa-leanpub:before{content:"\f212"}.fa-sellsy:before{content:"\f213"}.fa-shirtsinbulk:before{content:"\f214"}.fa-simplybuilt:before{content:"\f215"}.fa-skyatlas:before{content:"\f216"}.fa-cart-plus:before{content:"\f217"}.fa-cart-arrow-down:before{content:"\f218"}.fa-diamond:before{content:"\f219"}.fa-ship:before{content:"\f21a"}.fa-user-secret:before{content:"\f21b"}.fa-motorcycle:before{content:"\f21c"}.fa-street-view:before{content:"\f21d"}.fa-heartbeat:before{content:"\f21e"}.fa-venus:before{content:"\f221"}.fa-mars:before{content:"\f222"}.fa-mercury:before{content:"\f223"}.fa-transgender:before{content:"\f224"}.fa-transgender-alt:before{content:"\f225"}.fa-venus-double:before{content:"\f226"}.fa-mars-double:before{content:"\f227"}.fa-venus-mars:before{content:"\f228"}.fa-mars-stroke:before{content:"\f229"}.fa-mars-stroke-v:before{content:"\f22a"}.fa-mars-stroke-h:before{content:"\f22b"}.fa-neuter:before{content:"\f22c"}.fa-facebook-official:before{content:"\f230"}.fa-pinterest-p:before{content:"\f231"}.fa-whatsapp:before{content:"\f232"}.fa-server:before{content:"\f233"}.fa-user-plus:before{content:"\f234"}.fa-user-times:before{content:"\f235"}.fa-bed:before,.fa-hotel:before{content:"\f236"}.fa-viacoin:before{content:"\f237"}.fa-train:before{content:"\f238"}.fa-subway:before{content:"\f239"}.fa-medium:before{content:"\f23a"}/*! prefixes.scss v0.1.0 | Author: Pandao | https://github.com/pandao/prefixes.scss | MIT license | Copyright (c) 2015 */@font-face{font-family:editormd-logo;src:url(../fonts/editormd-logo.eot?-5y8q6h);src:url(.../fonts/editormd-logo.eot?#iefix-5y8q6h)format("embedded-opentype"),url(../fonts/editormd-logo.woff?-5y8q6h)format("woff"),url(../fonts/editormd-logo.ttf?-5y8q6h)format("truetype"),url(../fonts/editormd-logo.svg?-5y8q6h#icomoon)format("svg");font-weight:400;font-style:normal}.editormd-logo,.editormd-logo-1x,.editormd-logo-2x,.editormd-logo-3x,.editormd-logo-4x,.editormd-logo-5x,.editormd-logo-6x,.editormd-logo-7x,.editormd-logo-8x{font-family:editormd-logo;speak:none;font-style:normal;font-weight:400;font-variant:normal;text-transform:none;font-size:inherit;line-height:1;display:inline-block;text-rendering:auto;vertical-align:inherit;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale}.markdown-body hr:after,.markdown-body hr:before{content:"";display:table}.editormd-logo-1x:before,.editormd-logo-2x:before,.editormd-logo-3x:before,.editormd-logo-4x:before,.editormd-logo-5x:before,.editormd-logo-6x:before,.editormd-logo-7x:before,.editormd-logo-8x:before,.editormd-logo:before{content:"\e1987"}.editormd-logo-1x{font-size:1em}.editormd-logo-lg{font-size:1.2em}.editormd-logo-2x{font-size:2em}.editormd-logo-3x{font-size:3em}.editormd-logo-4x{font-size:4em}.editormd-logo-5x{font-size:5em}.editormd-logo-6x{font-size:6em}.editormd-logo-7x{font-size:7em}.editormd-logo-8x{font-size:8em}.editormd-logo-color{color:#2196F3}/*! github-markdown-css | The MIT License (MIT) | Copyright (c) Sindre Sorhus <sindresorhus@gmail.com> (sindresorhus.com) | https://github.com/sindresorhus/github-markdown-css */@font-face{font-family:octicons-anchor;src:url(data:font/woff;charset=utf-8;base64,d09GRgABAAAAAAYcAA0AAAAACjQAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAABGRlRNAAABMAAAABwAAAAca8vGTk9TLzIAAAFMAAAARAAAAFZG1VHVY21hcAAAAZAAAAA+AAABQgAP9AdjdnQgAAAB0AAAAAQAAAAEACICiGdhc3AAAAHUAAAACAAAAAj//wADZ2x5ZgAAAdwAAADRAAABEKyikaNoZWFkAAACsAAAAC0AAAA2AtXoA2hoZWEAAALgAAAAHAAAACQHngNFaG10eAAAAvwAAAAQAAAAEAwAACJsb2NhAAADDAAAAAoAAAAKALIAVG1heHAAAAMYAAAAHwAAACABEAB2bmFtZQAAAzgAAALBAAAFu3I9x/Nwb3N0AAAF/AAAAB0AAAAvaoFvbwAAAAEAAAAAzBdyYwAAAADP2IQvAAAAAM/bz7t4nGNgZGFgnMDAysDB1Ml0hoGBoR9CM75mMGLkYGBgYmBlZsAKAtJcUxgcPsR8iGF2+O/AEMPsznAYKMwIkgMA5REMOXicY2BgYGaAYBkGRgYQsAHyGMF8FgYFIM0ChED+h5j//yEk/3KoSgZGNgYYk4GRCUgwMaACRoZhDwCs7QgGAAAAIgKIAAAAAf//AAJ4nHWMMQrCQBBF/0zWrCCIKUQsTDCL2EXMohYGSSmorScInsRGL2DOYJe0Ntp7BK+gJ1BxF1stZvjz/v8DRghQzEc4kIgKwiAppcA9LtzKLSkdNhKFY3HF4lK69ExKslx7Xa+vPRVS43G98vG1DnkDMIBUgFN0MDXflU8tbaZOUkXUH0+U27RoRpOIyCKjbMCVejwypzJJG4jIwb43rfl6wbwanocrJm9XFYfskuVC5K/TPyczNU7b84CXcbxks1Un6H6tLH9vf2LRnn8Ax7A5WQAAAHicY2BkYGAA4teL1+yI57f5ysDNwgAC529f0kOmWRiYVgEpDgYmEA8AUzEKsQAAAHicY2BkYGB2+O/AEMPCAAJAkpEBFbAAADgKAe0EAAAiAAAAAAQAAAAEAAAAAAAAKgAqACoAiAAAeJxjYGRgYGBhsGFgYgABEMkFhAwM/xn0QAIAD6YBhwB4nI1Ty07cMBS9QwKlQapQW3VXySvEqDCZGbGaHULiIQ1FKgjWMxknMfLEke2A+IJu+wntrt/QbVf9gG75jK577Lg8K1qQPCfnnnt8fX1NRC/pmjrk/zprC+8D7tBy9DHgBXoWfQ44Av8t4Bj4Z8CLtBL9CniJluPXASf0Lm4CXqFX8Q84dOLnMB17N4c7tBo1AS/Qi+hTwBH4rwHHwN8DXqQ30XXAS7QaLwSc0Gn8NuAVWou/gFmnjLrEaEh9GmDdDGgL3B4JsrRPDU2hTOiMSuJUIdKQQayiAth69r6akSSFqIJuA19TrzCIaY8sIoxyrNIrL//pw7A2iMygkX5vDj+G+kuoLdX4GlGK/8Lnlz6/h9MpmoO9rafrz7ILXEHHaAx95s9lsI7AHNMBWEZHULnfAXwG9/ZqdzLI08iuwRloXE8kfhXYAvE23+23DU3t626rbs8/8adv+9DWknsHp3E17oCf+Z48rvEQNZ78paYM38qfk3v/u3l3u3GXN2Dmvmvpf1Srwk3pB/VSsp512bA/GG5i2WJ7wu430yQ5K3nFGiOqgtmSB5pJVSizwaacmUZzZhXLlZTq8qGGFY2YcSkqbth6aW1tRmlaCFs2016m5qn36SbJrqosG4uMV4aP2PHBmB3tjtmgN2izkGQyLWprekbIntJFing32a5rKWCN/SdSoga45EJykyQ7asZvHQ8PTm6cslIpwyeyjbVltNikc2HTR7YKh9LBl9DADC0U/jLcBZDKrMhUBfQBvXRzLtFtjU9eNHKin0x5InTqb8lNpfKv1s1xHzTXRqgKzek/mb7nB8RZTCDhGEX3kK/8Q75AmUM/eLkfA+0Hi908Kx4eNsMgudg5GLdRD7a84npi+YxNr5i5KIbW5izXas7cHXIMAau1OueZhfj+cOcP3P8MNIWLyYOBuxL6DRylJ4cAAAB4nGNgYoAALjDJyIAOWMCiTIxMLDmZedkABtIBygAAAA==)format("woff")}.markdown-body{-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%;color:#333;overflow:hidden;font-family:"Microsoft YaHei",Helvetica,"Meiryo UI","Malgun Gothic","Segoe UI","Trebuchet MS",Monaco,monospace,Tahoma,STXihei,"华文细黑",STHeiti,"Helvetica Neue","Droid Sans","wenquanyi micro hei",FreeSans,Arimo,Arial,SimSun,"宋体",Heiti,"黑体",sans-serif;font-size:16px;line-height:1.6;word-wrap:break-word}.markdown-body strong{font-weight:700}.markdown-body h1{margin:.67em 0}.markdown-body img{border:0}.markdown-body hr{-moz-box-sizing:content-box;box-sizing:content-box;height:0}.markdown-body input{color:inherit;margin:0;line-height:normal;font:13px/1.4 Helvetica,arial,freesans,clean,sans-serif,"Segoe UI Emoji","Segoe UI Symbol"}.markdown-body html input[disabled]{cursor:default}.markdown-body input[type=checkbox]{-moz-box-sizing:border-box;box-sizing:border-box;padding:0}.markdown-body *{-moz-box-sizing:border-box;box-sizing:border-box}.markdown-body a{background:0 0;color:#4183c4;text-decoration:none}.markdown-body a:active,.markdown-body a:hover{outline:0;text-decoration:underline}.markdown-body hr{margin:15px 0;overflow:hidden;background:0 0;border:0;border-bottom:1px solid #ddd}.markdown-body h1,.markdown-body h2{padding-bottom:.3em;border-bottom:1px solid #eee}.markdown-body blockquote{margin:0}.markdown-body ol ol,.markdown-body ul ol{list-style-type:lower-roman}.markdown-body ol ol ol,.markdown-body ol ul ol,.markdown-body ul ol ol,.markdown-body ul ul ol{list-style-type:lower-alpha}.markdown-body dd{margin-left:0}.markdown-body code{font-family:Consolas,"Liberation Mono",Menlo,Courier,monospace}.markdown-body pre{font:12px Consolas,"Liberation Mono",Menlo,Courier,monospace;word-wrap:normal}.markdown-body .octicon{font:normal normal 16px octicons-anchor;line-height:1;display:inline-block;text-decoration:none;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale;user-select:none}.markdown-body .octicon-link:before{content:'\f05c'}.markdown-body>:first-child{margin-top:0!important}.markdown-body>:last-child{margin-bottom:0!important}.markdown-body .anchor{position:absolute;top:0;left:0;display:block;padding-right:6px;padding-left:30px;margin-left:-30px}.markdown-body .anchor:focus{outline:0}.markdown-body h1,.markdown-body h2,.markdown-body h3,.markdown-body h4,.markdown-body h5,.markdown-body h6{position:relative;margin-top:1em;margin-bottom:16px;font-weight:700;line-height:1.4}.markdown-body h1 .octicon-link,.markdown-body h2 .octicon-link,.markdown-body h3 .octicon-link,.markdown-body h4 .octicon-link,.markdown-body h5 .octicon-link,.markdown-body h6 .octicon-link{display:none;color:#000;vertical-align:middle}.markdown-body h1:hover .anchor,.markdown-body h2:hover .anchor,.markdown-body h3:hover .anchor,.markdown-body h4:hover .anchor,.markdown-body h5:hover .anchor,.markdown-body h6:hover .anchor{padding-left:8px;margin-left:-30px;text-decoration:none}.markdown-body h1:hover .anchor .octicon-link,.markdown-body h2:hover .anchor .octicon-link,.markdown-body h3:hover .anchor .octicon-link,.markdown-body h4:hover .anchor .octicon-link,.markdown-body h5:hover .anchor .octicon-link,.markdown-body h6:hover .anchor .octicon-link{display:inline-block}.markdown-body h1{font-size:2.25em;line-height:1.2}.markdown-body h1 .anchor{line-height:1}.markdown-body h2{font-size:1.75em;line-height:1.225}.markdown-body h2 .anchor{line-height:1}.markdown-body h3{font-size:1.5em;line-height:1.43}.markdown-body h3 .anchor,.markdown-body h4 .anchor{line-height:1.2}.markdown-body h4{font-size:1.25em}.markdown-body h5 .anchor,.markdown-body h6 .anchor{line-height:1.1}.markdown-body h5{font-size:1em}.markdown-body h6{font-size:1em;color:#777}.markdown-body blockquote,.markdown-body dl,.markdown-body ol,.markdown-body p,.markdown-body pre,.markdown-body table,.markdown-body ul{margin-top:0;margin-bottom:16px}.markdown-body ol,.markdown-body ul{padding-left:2em}.markdown-body ol ol,.markdown-body ol ul,.markdown-body ul ol,.markdown-body ul ul{margin-top:0;margin-bottom:0}.markdown-body li>p{margin-top:16px}.markdown-body dl{padding:0}.markdown-body dl dt{padding:0;margin-top:16px;font-size:1em;font-style:italic;font-weight:700}.markdown-body dl dd{padding:0 16px;margin-bottom:16px}.markdown-body blockquote{padding:0 15px;color:#777;border-left:4px solid #ddd}.markdown-body blockquote>:first-child{margin-top:0}.markdown-body blockquote>:last-child{margin-bottom:0}.markdown-body table{border-collapse:collapse;border-spacing:0;display:block;width:100%;overflow:auto;word-break:normal;word-break:keep-all}.markdown-body table th{font-weight:700}.markdown-body table td,.markdown-body table th{padding:6px 13px;border:1px solid #ddd}.markdown-body table tr{background-color:#fff;border-top:1px solid #ccc}.markdown-body table tr:nth-child(2n){background-color:#f8f8f8}.markdown-body img{max-width:100%;-moz-box-sizing:border-box;box-sizing:border-box}.markdown-body code{padding:.2em 0;margin:0;font-size:85%;background-color:rgba(0,0,0,.04);border-radius:3px}.markdown-body code:after,.markdown-body code:before{letter-spacing:-.2em;content:"\00a0"}.markdown-body pre>code{padding:0;margin:0;font-size:100%;word-break:normal;white-space:pre;background:0 0;border:0}.markdown-body .highlight{margin-bottom:16px}.markdown-body .highlight pre,.markdown-body pre{padding:16px;overflow:auto;font-size:85%;background-color:#f7f7f7;border-radius:3px}.markdown-body .highlight pre{margin-bottom:0;word-break:normal}.markdown-body pre code{display:inline;max-width:initial;padding:0;margin:0;overflow:initial;line-height:inherit;word-wrap:normal;background-color:transparent;border:0}.markdown-body pre code:after,.markdown-body pre code:before{content:normal}.markdown-body .pl-c{color:#969896}.markdown-body .pl-c1,.markdown-body .pl-mdh,.markdown-body .pl-mm,.markdown-body .pl-mp,.markdown-body .pl-mr,.markdown-body .pl-s1 .pl-v,.markdown-body .pl-s3,.markdown-body .pl-sc,.markdown-body .pl-sv{color:#0086b3}.markdown-body .pl-e,.markdown-body .pl-en{color:#795da3}.markdown-body .pl-s1 .pl-s2,.markdown-body .pl-smi,.markdown-body .pl-smp,.markdown-body .pl-stj,.markdown-body .pl-vo,.markdown-body .pl-vpf{color:#333}.markdown-body .pl-ent{color:#63a35c}.markdown-body .pl-k,.markdown-body .pl-s,.markdown-body .pl-st{color:#a71d5d}.markdown-body .pl-pds,.markdown-body .pl-s1,.markdown-body .pl-s1 .pl-pse .pl-s2,.markdown-body .pl-sr,.markdown-body .pl-sr .pl-cce,.markdown-body .pl-sr .pl-sra,.markdown-body .pl-sr .pl-sre,.markdown-body .pl-src{color:#df5000}.markdown-body .pl-mo,.markdown-body .pl-v{color:#1d3e81}.markdown-body .pl-id{color:#b52a1d}.markdown-body .pl-ii{background-color:#b52a1d;color:#f8f8f8}.markdown-body .pl-sr .pl-cce{color:#63a35c;font-weight:700}.markdown-body .pl-ml{color:#693a17}.markdown-body .pl-mh,.markdown-body .pl-mh .pl-en,.markdown-body .pl-ms{color:#1d3e81;font-weight:700}.markdown-body .pl-mq{color:teal}.markdown-body .pl-mi{color:#333;font-style:italic}.markdown-body .pl-mb{color:#333;font-weight:700}.markdown-body .pl-md,.markdown-body .pl-mdhf{background-color:#ffecec;color:#bd2c00}.markdown-body .pl-mdht,.markdown-body .pl-mi1{background-color:#eaffea;color:#55a532}.markdown-body .pl-mdr{color:#795da3;font-weight:700}.markdown-body kbd{display:inline-block;padding:3px 5px;font:11px Consolas,"Liberation Mono",Menlo,Courier,monospace;line-height:10px;color:#555;vertical-align:middle;background-color:#fcfcfc;border:1px solid #ccc;border-bottom-color:#bbb;border-radius:3px;box-shadow:inset 0 -1px 0 #bbb}.markdown-body .task-list-item+.task-list-item{margin-top:3px}.markdown-body .task-list-item input{float:left;margin:.3em 0 .25em -1.6em;vertical-align:middle}.markdown-body :checked+.radio-label{z-index:1;position:relative;border-color:#4183c4}.editormd-html-preview,.editormd-preview-container{text-align:left;font-size:14px;line-height:1.6;padding:20px;overflow:auto;width:100%;background-color:#fff}.editormd-html-preview blockquote,.editormd-preview-container blockquote{color:#666;border-left:4px solid #ddd;padding-left:20px;margin-left:0;font-size:14px;font-style:italic}.editormd-html-preview p code,.editormd-preview-container p code{margin-left:5px;margin-right:4px}.editormd-html-preview abbr,.editormd-preview-container abbr{background:#ffd}.editormd-html-preview hr,.editormd-preview-container hr{height:1px;border:none;border-top:1px solid #ddd;background:0 0}.editormd-html-preview code,.editormd-preview-container code{border:1px solid #ddd;background:#f6f6f6;padding:3px;border-radius:3px;font-size:14px}.editormd-html-preview pre,.editormd-preview-container pre{border:1px solid #ddd;background:#f6f6f6;padding:10px;-webkit-border-radius:3px;-moz-border-radius:3px;-ms-border-radius:3px;-o-border-radius:3px;border-radius:3px}.editormd-html-preview pre code,.editormd-preview-container pre code{padding:0}.editormd-html-preview code,.editormd-html-preview kbd,.editormd-html-preview pre,.editormd-preview-container code,.editormd-preview-container kbd,.editormd-preview-container pre{font-family:"YaHei Consolas Hybrid",Consolas,"Meiryo UI","Malgun Gothic","Segoe UI","Trebuchet MS",Helvetica,monospace,monospace}.editormd-html-preview table thead tr,.editormd-preview-container table thead tr{background-color:#F8F8F8}.editormd-html-preview p.editormd-tex,.editormd-preview-container p.editormd-tex{text-align:center}.editormd-html-preview span.editormd-tex,.editormd-preview-container span.editormd-tex{margin:0 5px}.editormd-html-preview .emoji,.editormd-preview-container .emoji{width:24px;height:24px}.editormd-html-preview .katex,.editormd-preview-container .katex{font-size:1.4em}.editormd-html-preview .flowchart,.editormd-html-preview .sequence-diagram,.editormd-preview-container .flowchart,.editormd-preview-container .sequence-diagram{margin:0 auto;text-align:center}.editormd-html-preview .flowchart svg,.editormd-html-preview .sequence-diagram svg,.editormd-preview-container .flowchart svg,.editormd-preview-container .sequence-diagram svg{margin:0 auto}.editormd-html-preview .flowchart text,.editormd-html-preview .sequence-diagram text,.editormd-preview-container .flowchart text,.editormd-preview-container .sequence-diagram text{font-size:15px!important;font-family:"YaHei Consolas Hybrid",Consolas,"Microsoft YaHei","Malgun Gothic","Segoe UI",Helvetica,Arial!important}/*! Pretty printing styles. Used with prettify.js. */.pln{color:#000}@media screen{.str{color:#080}.kwd{color:#008}.com{color:#800}.typ{color:#606}.lit{color:#066}.clo,.opn,.pun{color:#660}.tag{color:#008}.atn{color:#606}.atv{color:#080}.dec,.var{color:#606}.fun{color:red}}@media print,projection{.kwd,.tag,.typ{font-weight:700}.str{color:#060}.kwd{color:#006}.com{color:#600;font-style:italic}.typ{color:#404}.lit{color:#044}.clo,.opn,.pun{color:#440}.tag{color:#006}.atn{color:#404}.atv{color:#060}}pre.prettyprint{padding:2px;border:1px solid #888}ol.linenums{margin-top:0;margin-bottom:0}li.L1,li.L3,li.L5,li.L7,li.L9{background:#eee}.editormd-html-preview pre.prettyprint,.editormd-preview-container pre.prettyprint{padding:10px;border:1px solid #ddd;white-space:pre-wrap;word-wrap:break-word}.editormd-html-preview ol.linenums,.editormd-preview-container ol.linenums{color:#999;padding-left:2.5em}.editormd-html-preview ol.linenums li,.editormd-preview-container ol.linenums li{list-style-type:decimal}.editormd-html-preview ol.linenums li code,.editormd-preview-container ol.linenums li code{border:none;background:0 0;padding:0}.editormd-html-preview .editormd-toc-menu,.editormd-preview-container .editormd-toc-menu{margin:8px 0 12px;display:inline-block}.editormd-html-preview .editormd-toc-menu>.markdown-toc,.editormd-preview-container .editormd-toc-menu>.markdown-toc{position:relative;-webkit-border-radius:4px;-moz-border-radius:4px;-ms-border-radius:4px;-o-border-radius:4px;border-radius:4px;border:1px solid #ddd;display:inline-block;font-size:1em}.editormd-html-preview .editormd-toc-menu>.markdown-toc>ul,.editormd-preview-container .editormd-toc-menu>.markdown-toc>ul{width:160%;min-width:180px;position:absolute;left:-1px;top:-2px;z-index:100;padding:0 10px 10px;display:none;background:#fff;border:1px solid #ddd;-webkit-border-radius:4px;-moz-border-radius:4px;-ms-border-radius:4px;-o-border-radius:4px;border-radius:4px;-webkit-box-shadow:0 3px 5px rgba(0,0,0,.2);-moz-box-shadow:0 3px 5px rgba(0,0,0,.2);-ms-box-shadow:0 3px 5px rgba(0,0,0,.2);-o-box-shadow:0 3px 5px rgba(0,0,0,.2);box-shadow:0 3px 5px rgba(0,0,0,.2)}.editormd-html-preview .editormd-toc-menu>.markdown-toc>ul>li ul,.editormd-preview-container .editormd-toc-menu>.markdown-toc>ul>li ul{width:100%;min-width:180px;border:1px solid #ddd;display:none;background:#fff;-webkit-border-radius:4px;-moz-border-radius:4px;-ms-border-radius:4px;-o-border-radius:4px;border-radius:4px}.editormd-html-preview .editormd-toc-menu .toc-menu-btn:hover,.editormd-html-preview .editormd-toc-menu>.markdown-toc>ul>li a:hover,.editormd-preview-container .editormd-toc-menu .toc-menu-btn:hover,.editormd-preview-container .editormd-toc-menu>.markdown-toc>ul>li a:hover{background-color:#f6f6f6}.editormd-html-preview .editormd-toc-menu>.markdown-toc>ul>li a,.editormd-preview-container .editormd-toc-menu>.markdown-toc>ul>li a{color:#666;padding:6px 10px;display:block;-webkit-transition:background-color 500ms ease-out;-moz-transition:background-color 500ms ease-out;transition:background-color 500ms ease-out}.editormd-html-preview .editormd-toc-menu>.markdown-toc li,.editormd-preview-container .editormd-toc-menu>.markdown-toc li{position:relative}.editormd-html-preview .editormd-toc-menu>.markdown-toc li>ul,.editormd-preview-container .editormd-toc-menu>.markdown-toc li>ul{position:absolute;top:32px;left:10%;display:none;-webkit-box-shadow:0 3px 5px rgba(0,0,0,.2);-moz-box-shadow:0 3px 5px rgba(0,0,0,.2);-ms-box-shadow:0 3px 5px rgba(0,0,0,.2);-o-box-shadow:0 3px 5px rgba(0,0,0,.2);box-shadow:0 3px 5px rgba(0,0,0,.2)}.editormd-html-preview .editormd-toc-menu>.markdown-toc li>ul:after,.editormd-html-preview .editormd-toc-menu>.markdown-toc li>ul:before,.editormd-preview-container .editormd-toc-menu>.markdown-toc li>ul:after,.editormd-preview-container .editormd-toc-menu>.markdown-toc li>ul:before{pointer-events:pointer-events;position:absolute;left:15px;top:-6px;display:block;content:"";width:0;height:0;border:6px solid transparent;border-width:0 6px 6px;z-index:10}.editormd-html-preview .editormd-toc-menu>.markdown-toc li>ul:before,.editormd-preview-container .editormd-toc-menu>.markdown-toc li>ul:before{border-bottom-color:#ccc}.editormd-html-preview .editormd-toc-menu>.markdown-toc li>ul:after,.editormd-preview-container .editormd-toc-menu>.markdown-toc li>ul:after{border-bottom-color:#fff;top:-5px}.editormd-html-preview .editormd-toc-menu ul,.editormd-preview-container .editormd-toc-menu ul{list-style:none}.editormd-html-preview .editormd-toc-menu a,.editormd-preview-container .editormd-toc-menu a{text-decoration:none}.editormd-html-preview .editormd-toc-menu h1,.editormd-preview-container .editormd-toc-menu h1{font-size:16px;padding:5px 0 10px 10px;line-height:1;border-bottom:1px solid #eee}.editormd-html-preview .editormd-toc-menu h1 .fa,.editormd-preview-container .editormd-toc-menu h1 .fa{padding-left:10px}.editormd-html-preview .editormd-toc-menu .toc-menu-btn,.editormd-preview-container .editormd-toc-menu .toc-menu-btn{color:#666;min-width:180px;padding:5px 10px;border-radius:4px;display:inline-block;-webkit-transition:background-color 500ms ease-out;-moz-transition:background-color 500ms ease-out;transition:background-color 500ms ease-out}.editormd-html-preview textarea,.editormd-onlyread .editormd-toolbar{display:none}.editormd-html-preview .editormd-toc-menu .toc-menu-btn .fa,.editormd-preview-container .editormd-toc-menu .toc-menu-btn .fa{float:right;padding:3px 0 0 10px;font-size:1.3em}.markdown-body .editormd-toc-menu ul{padding-left:0}.markdown-body .highlight pre,.markdown-body pre{line-height:1.6}hr.editormd-page-break{border:1px dotted #ccc;font-size:0;height:2px}@media only print{hr.editormd-page-break{background:0 0;border:none;height:0}}.editormd-html-preview hr.editormd-page-break{background:0 0;border:none;height:0}.editormd-preview-close-btn{color:#fff;padding:4px 6px;font-size:18px;-webkit-border-radius:500px;-moz-border-radius:500px;-ms-border-radius:500px;-o-border-radius:500px;border-radius:500px;display:none;background-color:#ccc;position:absolute;top:25px;right:35px;z-index:19;-webkit-transition:background-color 300ms ease-out;-moz-transition:background-color 300ms ease-out;transition:background-color 300ms ease-out}.editormd-preview-close-btn:hover{background-color:#999}.editormd-preview-active{width:100%;padding:40px}.editormd-preview-theme-dark{color:#777;background:#2C2827}.editormd-preview-theme-dark .editormd-preview-container{color:#888;background-color:#2C2827}.editormd-preview-theme-dark .editormd-preview-container pre.prettyprint{border:none}.editormd-preview-theme-dark .editormd-preview-container blockquote{color:#555;padding:.5em;background:#222;border-color:#333}.editormd-preview-theme-dark .editormd-preview-container abbr{color:#fff;padding:1px 3px;-webkit-border-radius:3px;-moz-border-radius:3px;-ms-border-radius:3px;-o-border-radius:3px;border-radius:3px;background:#f90}.editormd-preview-theme-dark .editormd-preview-container code{color:#fff;border:none;padding:1px 3px;-webkit-border-radius:3px;-moz-border-radius:3px;-ms-border-radius:3px;-o-border-radius:3px;border-radius:3px;background:#5A9600}.editormd-preview-theme-dark .editormd-preview-container table{border:none}.editormd-preview-theme-dark .editormd-preview-container .fa-emoji{color:#B4BF42}.editormd-preview-theme-dark .editormd-preview-container .katex{color:#FEC93F}.editormd-preview-theme-dark .editormd-toc-menu>.markdown-toc{background:#fff;border:none}.editormd-preview-theme-dark .editormd-toc-menu>.markdown-toc h1{border-color:#ddd}.editormd-preview-theme-dark .markdown-body h1,.editormd-preview-theme-dark .markdown-body h2,.editormd-preview-theme-dark .markdown-body hr{border-color:#222}.editormd-preview-theme-dark pre{color:#999;background-color:#111;background-color:rgba(0,0,0,.4)}.editormd-preview-theme-dark pre .pln{color:#999}.editormd-preview-theme-dark li.L1,.editormd-preview-theme-dark li.L3,.editormd-preview-theme-dark li.L5,.editormd-preview-theme-dark li.L7,.editormd-preview-theme-dark li.L9{background:0 0}.editormd-preview-theme-dark [class*=editormd-logo]{color:#2196F3}.editormd-preview-theme-dark .sequence-diagram text{fill:#fff}.editormd-preview-theme-dark .sequence-diagram path,.editormd-preview-theme-dark .sequence-diagram rect{color:#fff;fill:#64D1CB;stroke:#64D1CB}.editormd-preview-theme-dark .flowchart path,.editormd-preview-theme-dark .flowchart rect{stroke:#A6C6FF}.editormd-preview-theme-dark .flowchart rect{fill:#A6C6FF}.editormd-preview-theme-dark .flowchart text{fill:#5879B4}@media screen{.editormd-preview-theme-dark .str{color:#080}.editormd-preview-theme-dark .kwd{color:#f90}.editormd-preview-theme-dark .com{color:#444}.editormd-preview-theme-dark .typ{color:#606}.editormd-preview-theme-dark .lit{color:#066}.editormd-preview-theme-dark .clo,.editormd-preview-theme-dark .opn,.editormd-preview-theme-dark .pun{color:#660}.editormd-preview-theme-dark .tag{color:#f90}.editormd-preview-theme-dark .atn{color:#6C95F5}.editormd-preview-theme-dark .atv{color:#080}.editormd-preview-theme-dark .dec,.editormd-preview-theme-dark .var{color:#008BA7}.editormd-preview-theme-dark .fun{color:red}}.editormd-onlyread .CodeMirror{margin-top:0}.editormd-onlyread .editormd-preview{top:0}.editormd-fullscreen{position:fixed;top:0;left:0;border:none;margin:0 auto}.editormd-theme-dark{border-color:#1a1a17}.editormd-theme-dark .editormd-toolbar{background:#1A1A17;border-color:#1a1a17}.editormd-theme-dark .editormd-menu>li>a{color:#777;border-color:#1a1a17}.editormd-theme-dark .editormd-menu>li>a.active,.editormd-theme-dark .editormd-menu>li>a:hover{border-color:#333;background:#333}.editormd-theme-dark .editormd-menu>li.divider{border-right:1px solid #111}.editormd-theme-dark .CodeMirror{border-right:1px solid rgba(0,0,0,.1)}`)
	formList.SetTable("api_definition").SetTitle("接口定义").SetDescription("接口定义")

	return apiDefinition
}
