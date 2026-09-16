package pages

import (
	"fmt"
	"html/template"
	"strings"

	"data4test/biz"

	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
)

// ==================== 统一色板 ====================

// 状态色（固定语义，跨报告一致）
const (
	colorPass       chartjs.Color = "rgba(0,166,90,1)"    // 通过 绿 #00a65a
	colorFail       chartjs.Color = "rgba(221,75,57,1)"   // 失败 红 #dd4b39
	colorUntest     chartjs.Color = "rgba(160,160,160,1)" // 未执行 灰 #a0a0a0
	colorPart       chartjs.Color = "rgba(243,156,18,1)"  // 部分 橙 #f39c12
	colorUnmerged   chartjs.Color = "rgba(60,141,188,1)"  // 暂未合入 蓝 #3c8dbc
	colorDeprecated chartjs.Color = "rgba(85,85,85,1)"    // 废弃 深灰 #555555
	colorInfo       chartjs.Color = "rgba(60,141,188,1)"  // 折线「总数」/信息 蓝 #3c8dbc
)

// categoricalColors 分类色板（固定 8 色顺序，身份类饼图 / 多系列折线用）
var categoricalColors = []chartjs.Color{
	colorPass,           // 绿
	colorInfo,           // 蓝
	colorFail,           // 红
	colorPart,           // 橙
	"rgba(96,92,168,1)", // 紫
	"rgba(0,150,136,1)", // 青
	"rgba(233,30,99,1)", // 粉
	"rgba(121,85,72,1)", // 棕
}

// categoricalLegendColors 分类图例色（AdminLTE text-* 类名，与 categoricalColors 一一对应）
var categoricalLegendColors = []string{
	"green", "blue", "red", "yellow", "purple", "teal", "fuchsia", "maroon",
}

// categoricalColor 按固定顺序取分类色
func categoricalColor(i int) chartjs.Color { return categoricalColors[i%len(categoricalColors)] }

// categoricalLegendColor 按固定顺序取分类图例色
func categoricalLegendColor(i int) string {
	return categoricalLegendColors[i%len(categoricalLegendColors)]
}

// ==================== 身份类配色（避开状态色） ====================

// identityColors 身份类分类色板（8 色固定顺序）。
// 直接取自 AdminLTE 语义色（与 KPI 卡 info-box 的 bg-* 同源），刻意避开状态色：
// 通过绿 #00a65a / 失败红 #dd4b39 / 部分黄 #f39c12 / 未执行灰 / 暂未合入蓝 #3c8dbc 均不纳入；
// teal 加深半档（#39cccc→#32b6b6）以落入 lightness 校验带。空值切片由调用处单独使用灰色（colorUntest）。
var identityColors = []chartjs.Color{
	"rgba(0,115,183,1)",  // 蓝     blue   #0073b7
	"rgba(0,192,239,1)",  // 天蓝   aqua   #00c0ef
	"rgba(96,92,168,1)",  // 紫     purple #605ca8
	"rgba(50,182,182,1)", // 青     teal   #32b6b6
	"rgba(216,27,96,1)",  // 玫红   maroon #d81b60
	"rgba(240,18,190,1)", // 品红   fuchsia#f012be
	"rgba(61,153,112,1)", // 橄榄   olive  #3d9970
	"rgba(255,133,27,1)", // 橙     orange #ff851b
}

// identityLegendColors 与 identityColors 一一对应的 AdminLTE text-* 类名（切片色 = 图例色，1:1 一致）。
var identityLegendColors = []string{
	"blue", "aqua", "purple", "teal", "maroon", "fuchsia", "olive", "orange",
}

// identityColor 按固定顺序取身份类分类色。
func identityColor(i int) chartjs.Color { return identityColors[i%len(identityColors)] }

// identityLegendColor 按固定顺序取身份类分类图例色。
func identityLegendColor(i int) string {
	return identityLegendColors[i%len(identityLegendColors)]
}

// ==================== 统一样式 ====================

// reportStyle 统一 <style> 块（现代卡片 + KPI + 表格 + 悬浮提示）。
// 每个报告在返回 Panel 前追加一次：content += string(reportStyle())。
func reportStyle() template.HTML {
	return template.HTML(`<style>
.report-card{border:none;border-radius:6px;box-shadow:0 1px 4px rgba(0,0,0,.08);background:#fff;margin-bottom:20px}
.report-card-header{display:flex;justify-content:space-between;align-items:center;padding:12px 20px;border-bottom:1px solid #eee}
.report-card-title{font-weight:600;font-size:15px;color:#333}
.report-card-tools{font-size:12px}
.report-card-body{padding:20px}
.report-card-footer{padding:10px 20px;border-top:1px solid #eee;text-align:center;font-size:12px}
.report-table{width:100%;border-collapse:collapse;margin-bottom:0}
.report-table thead th{background:#f7f7f7;font-weight:600;color:#555;border-bottom:2px solid #e9ecef;padding:10px 8px;white-space:nowrap;font-variant-numeric:tabular-nums}
.report-table tbody td{padding:10px 8px;border-bottom:1px solid #f0f0f0;color:#444;font-variant-numeric:tabular-nums}
.report-table tbody tr:nth-child(even){background:#fafafa}
.sc-td{position:relative;cursor:pointer}
.sc-td .sc-truncate{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.sc-td .sc-full{display:none;position:absolute;right:5%;top:50%;background:#fff;border:2px solid #666;padding:15px;z-index:9999;max-width:600px;max-height:80vh;overflow-y:auto;white-space:pre-wrap;word-break:break-all;box-shadow:0 4px 20px rgba(0,0,0,0.3);border-radius:4px;font-size:13px;line-height:1.4}
.sc-td:hover .sc-full{display:block!important}
.scene-header{cursor:pointer}
.scene-header:hover{background:#f4f8fb!important}
.scene-arrow{display:inline-block;width:16px;color:#3c8dbc;font-weight:bold}
.scene-name{font-weight:600;color:#333}
.scene-child{display:none;background:#fafbfc}
.scene-child.open{display:table-row}
.scene-child:hover{background:#f4f7fa}
.task-header{cursor:pointer}
.task-header:hover{background:#eef4fb!important}
.task-name{font-weight:700;color:#1a3c5e}
.task-scene{display:none;cursor:pointer;background:#fbfcfe}
.task-scene.open{display:table-row}
.task-scene:hover{background:#f4f8fb!important}
.report-table td.sc-data-name{padding-left:48px;color:#555}
.report-table td.sc-api{font-family:Menlo,Consolas,monospace;font-size:12px;color:#777}
.sc-badge{display:inline-block;padding:1px 8px;border-radius:10px;font-size:12px;color:#fff;white-space:nowrap}
.sc-badge-pass{background:#00a65a}
.sc-badge-fail{background:#dd4b39}
.sc-badge-na{background:#a0a0a0}
.report-table td.sc-count{color:#888;font-size:12px;white-space:nowrap}
.sc-btn{display:inline-block;padding:3px 10px;border:1px solid #d2d6de;border-radius:3px;background:#fff;color:#444;font-size:12px;cursor:pointer;user-select:none;margin-left:6px}
.sc-btn:hover{border-color:#3c8dbc;color:#3c8dbc}
</style>`)
}

// reportCollapseScript 场景折叠展开的前端脚本（点击场景头行展开/收起关联数据，全部展开/收起）。
// 多任务报告复用同一脚本：任务头行用 toggleTaskGroup 展开/收起其场景，场景头行用 toggleSceneGroup 展开/收起其数据。
func reportCollapseScript() string {
	return `<script>
function toggleSceneGroup(idx, el) {
  var rows = document.querySelectorAll('.scene-child[data-group="' + idx + '"]');
  var open = rows.length && rows[0].classList.contains('open');
  for (var i = 0; i < rows.length; i++) { rows[i].classList.toggle('open', !open); }
  var arrow = el.querySelector('.scene-arrow');
  if (arrow) { arrow.textContent = open ? '▸' : '▾'; }
}
function toggleTaskGroup(idx, el) {
  var rows = document.querySelectorAll('.task-scene[data-task="' + idx + '"]');
  var open = rows.length && rows[0].classList.contains('open');
  for (var i = 0; i < rows.length; i++) { rows[i].classList.toggle('open', !open); }
  // 无论展开还是收起任务，其下数据子行均收起，场景箭头复位为 ▸
  var kids = document.querySelectorAll('.scene-child[data-task="' + idx + '"]');
  for (var j = 0; j < kids.length; j++) { kids[j].classList.remove('open'); }
  var sarrows = document.querySelectorAll('.task-scene[data-task="' + idx + '"] .scene-arrow');
  for (var k = 0; k < sarrows.length; k++) { sarrows[k].textContent = '▸'; }
  var arrow = el.querySelector('.scene-arrow');
  if (arrow) { arrow.textContent = open ? '▸' : '▾'; }
}
function setAllScenes(open) {
  var rows = document.querySelectorAll('.scene-child, .task-scene');
  for (var i = 0; i < rows.length; i++) { rows[i].classList.toggle('open', open); }
  var arrows = document.querySelectorAll('.scene-arrow');
  for (var j = 0; j < arrows.length; j++) { arrows[j].textContent = open ? '▾' : '▸'; }
}
</script>`
}

// ==================== 通用构建函数 ====================

// reportCard 现代卡片容器。extra 为可选的头栏右侧内容（链接等），可为空。
// title 与 extra 均为空时不渲染头栏（适用于无标题信息框）。
func reportCard(title, extra, body string) string {
	header := ""
	if title != "" || extra != "" {
		if extra != "" {
			extra = fmt.Sprintf(`<div class="report-card-tools">%s</div>`, extra)
		}
		header = fmt.Sprintf(`<div class="report-card-header"><span class="report-card-title">%s</span>%s</div>`, title, extra)
	}
	return fmt.Sprintf(`<div class="report-card">%s<div class="report-card-body">%s</div></div>`, header, body)
}

// reportInfoBox AdminLTE info-box KPI 卡（白底 + 左彩色图标 + 右标签/数字）。
// icon 为 FontAwesome 4 图标类名（如 "fa-cubes"），color 为 bg-* 后缀（如 "blue"）。
func reportInfoBox(value, label, icon, color string) string {
	return fmt.Sprintf(`<div class="col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-%s"><i class="fa %s"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%s</span></div></div></div>`, color, icon, label, value)
}

// reportLegend 统一图例。
func reportLegend(items []map[string]string) string {
	return string(chart_legend.New().SetData(items).GetContent())
}

// reportPieChart 仅渲染饼图（高 180，不含图例与卡片）。labels 需非空。
func reportPieChart(id string, labels []string, counts []float64, colors []chartjs.Color) string {
	return string(chartjs.Pie().
		SetHeight(180).
		SetLabels(labels).
		SetID(id).
		AddDataSet(labels[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent())
}

// reportPie 统一饼图（高 180，col-md-4 一栏）。extra 为可选头栏右侧内容（链接等）。labels 需非空。
func reportPie(id, title, extra string, labels []string, counts []float64, colors []chartjs.Color, legendItems []map[string]string) string {
	legend := reportLegend(legendItems)
	body := fmt.Sprintf(`<div class="row"><div class="col-md-8">%s</div><div class="col-md-4">%s</div></div>`, reportPieChart(id, labels, counts, colors), legend)
	return `<div class="col-md-4">` + reportCard(title, extra, body) + `</div>`
}

// reportLineSeries 折线图单个数据系列。
type reportLineSeries struct {
	Label string
	Data  []float64
	Color chartjs.Color
	Fill  bool
}

// reportLine 统一折线图（高 320），按系列循环 AddDataSet。
func reportLine(id, cardTitle string, labels []string, series []reportLineSeries) string {
	line := chartjs.Line().SetID(id).SetHeight(320).SetLabels(labels)
	for _, s := range series {
		line = line.AddDataSet(s.Label).DSData(s.Data).DSFill(s.Fill).DSBorderColor(s.Color).DSLineTension(0.1)
	}
	return reportCard(cardTitle, "", string(line.GetContent()))
}

// reportTable 统一表格（斑马纹 + 灰表头 + tabular-nums）。
func reportTable(headers []string, rows []string) string {
	var th strings.Builder
	for _, h := range headers {
		th.WriteString(`<th>` + h + `</th>`)
	}
	return fmt.Sprintf(`<div class="table-responsive"><table class="report-table"><thead><tr>%s</tr></thead><tbody>%s</tbody></table></div>`, th.String(), strings.Join(rows, ""))
}

// reportTableFromInfo 由 GoAdmin 表格数据（[]map[string]types.InfoItem + Thead）渲染统一表格。
// 表头 Head 存的是 i18n key（或原始值如 HTTP 方法名），渲染时对 Head 做一次 T() 本地化；
// 内容查找仍用原始 Head 作为 map key，保证与生成时存入的 key 一致。
func reportTableFromInfo(contents []map[string]types.InfoItem, headers types.Thead) string {
	heads := make([]string, 0, len(headers))
	for _, h := range headers {
		heads = append(heads, biz.T(h.Head))
	}
	var rows []string
	for _, content := range contents {
		cells := make([]string, 0, len(headers))
		for _, h := range headers {
			cell := ""
			if item, ok := content[h.Head]; ok {
				cell = string(item.Content)
			}
			cells = append(cells, `<td>`+cell+`</td>`)
		}
		rows = append(rows, `<tr>`+strings.Join(cells, "")+`</tr>`)
	}
	return reportTable(heads, rows)
}

// localizeInfos 饼图分类名本地化：存储的是 i18n key 或原始值（如 pass/GET/接口类型名），
// key 翻译成本地语言，非 key 原样返回。
func localizeInfos(infos []string) []string {
	out := make([]string, len(infos))
	for i, v := range infos {
		out[i] = biz.T(v)
	}
	return out
}

// rebuildLegendLabels 用本地化后的分类名与计数重建图例 label，color 沿用原值。
func rebuildLegendLabels(labels []map[string]string, infos []string, counts []float64) []map[string]string {
	out := make([]map[string]string, 0, len(labels))
	for i, l := range labels {
		m := map[string]string{"color": l["color"]}
		if i < len(infos) {
			cnt := 0
			if i < len(counts) {
				cnt = int(counts[i])
			}
			m["label"] = fmt.Sprintf(" %s - %d", infos[i], cnt)
		}
		out = append(out, m)
	}
	return out
}
