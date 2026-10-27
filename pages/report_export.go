package pages

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"data4test/biz"
	"data4test/models"

	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/sword"
)

// adminLinkRe 匹配指向 Web 后台的链接（如 /admin/info/scene_test_history?...），
// 离线 HTML 中将其替换为纯文本，避免点击跳到无效地址。
var adminLinkRe = regexp.MustCompile(`<a\s+[^>]*href="/admin/[^"]*"[^>]*>(.*?)</a>`)

// offlineReportExport 标记当前是否处于离线 HTML 导出流程（ExportReportHTML）。
// 离线报告中的结果文件链接需指向本地相对路径（../mgmt/history|data/...），
// 而非 Web 端可用的 /admin/fm/... 预览链接。
var offlineReportExport bool

// ExportReportHTML 导出指定报告为独立 HTML 文件，样式与 Web「查看报告」一致。
// 返回导出的 HTML 文件绝对路径。
func ExportReportHTML(reportId, outDir string) (string, error) {
	var report biz.DashboardReport
	models.Orm.Table("dashboard").Where("id = ?", reportId).Find(&report)
	if len(report.Id) == 0 {
		return "", fmt.Errorf("%s", biz.T("cli.report_not_found"))
	}

	// 标记离线导出流程：结果文件链接渲染为本地相对路径（见 schedule_dashboard.go）。
	offlineReportExport = true
	defer func() { offlineReportExport = false }()

	panel, err := renderReportPanel(report)
	if err != nil {
		return "", err
	}

	html := buildReportHTML(report, panel)

	if err := ensureReportAssets(outDir); err != nil {
		return "", err
	}

	filePath := filepath.Join(outDir, reportFileName(report))
	if err := os.WriteFile(filePath, []byte(html), 0644); err != nil {
		return "", err
	}
	return filePath, nil
}

// renderReportPanel 复用现有渲染函数，返回报告的 Title/Description/Content。
func renderReportPanel(report biz.DashboardReport) (types.Panel, error) {
	dataLocale := biz.GetDataLocale("", biz.GetLocale())
	switch report.ReportType {
	case "task":
		if strings.Contains(report.RelatedTaskIds, ",") {
			return renderMultiTaskReport(report, dataLocale)
		}
		return renderTaskReport(report, "", dataLocale)
	default:
		return types.Panel{}, fmt.Errorf("%s: %s", biz.T("cli.report_type_unsupported"), report.ReportType)
	}
}

// buildReportHTML 组装静态 HTML 壳（对齐 GoAdmin 的 content-header + content 结构）。
func buildReportHTML(report biz.DashboardReport, panel types.Panel) string {
	title := report.ReportName
	if title == "" {
		title = string(panel.Title)
	}
	content := stripAdminLinks(string(panel.Content))
	cssURL := reportAssetCSSPath()

	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString(`<meta charset="utf-8">` + "\n")
	b.WriteString(`<meta http-equiv="X-UA-Compatible" content="IE=edge">` + "\n")
	b.WriteString(fmt.Sprintf("<title>%s</title>\n", template.HTMLEscapeString(title)))
	b.WriteString(`<meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">` + "\n")
	if cssURL != "" {
		b.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s">`+"\n", cssURL))
	}
	b.WriteString(`<script src="assets/js/chart.min.js"></script>` + "\n")
	b.WriteString("</head>\n")
	b.WriteString(`<body class="hold-transition skin-black sidebar-mini">` + "\n")
	b.WriteString(`<div class="wrapper"><div class="content-wrapper">` + "\n")
	if len(panel.Title) > 0 {
		b.WriteString(`<section class="content-header"><h4>`)
		b.WriteString(string(panel.Title))
		if len(panel.Description) > 0 {
			b.WriteString(` <small>` + string(panel.Description) + `</small>`)
		}
		b.WriteString(`</h4></section>` + "\n")
	}
	b.WriteString(`<section class="content">` + "\n")
	b.WriteString(content)
	b.WriteString("\n</section>\n")
	b.WriteString(`</div></div>` + "\n")
	b.WriteString("</body>\n</html>\n")
	return b.String()
}

// stripAdminLinks 将指向 Web 后台的 <a href="/admin/..."> 替换为其内部纯文本。
func stripAdminLinks(content string) string {
	return adminLinkRe.ReplaceAllString(content, "$1")
}

// ensureReportAssets 首次生成报告资源到 report/assets/（主题 CSS/字体/图片 + chart.min.js）。
func ensureReportAssets(outDir string) error {
	assetDir := filepath.Join(outDir, "assets")

	// chart.min.js（图表渲染）
	if err := writeAssetOnce(filepath.Join(assetDir, "js", "chart.min.js"), func() ([]byte, error) {
		return chartjs.Asset("assets/chart.min.js")
	}); err != nil {
		return err
	}

	// 主题静态资源（CSS/字体/图片/JS），保留 /dist/ 下的相对目录结构
	th := sword.Get()
	for _, p := range th.GetAssetList() {
		rel := strings.TrimPrefix(p, "/dist/")
		dest := filepath.Join(assetDir, filepath.FromSlash(rel))
		if err := writeAssetOnce(dest, func() ([]byte, error) {
			return th.GetAsset("/assets" + p)
		}); err != nil {
			return err
		}
	}
	return nil
}

// writeAssetOnce 文件已存在则跳过，否则调用 get 获取内容写入 dest。
func writeAssetOnce(dest string, get func() ([]byte, error)) error {
	if _, err := os.Stat(dest); err == nil {
		return nil
	}
	data, err := get()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return os.WriteFile(dest, data, 0644)
}

// reportAssetCSSPath 返回主题主样式在 report/assets/ 下的相对 URL。
func reportAssetCSSPath() string {
	for _, p := range sword.Get().GetAssetList() {
		if strings.HasPrefix(p, "/dist/css/") && strings.HasSuffix(p, ".css") {
			return "assets/" + strings.TrimPrefix(p, "/dist/")
		}
	}
	return ""
}

// reportFileName 由报告名生成 HTML 文件名（报告名已含生成时间，如 [任务]任务报告_20261014153000.html）。
func reportFileName(report biz.DashboardReport) string {
	name := report.ReportName
	if name == "" {
		name = "report"
	}
	name = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '_'
		}
		return r
	}, name)
	return name + ".html"
}
