package biz

import (
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"data4test/models"
)

// CliViewPort 本地「查看详情」服务端口（report 命令专用，独立于 run 命令的 mock 端口）。
const CliViewPort = 9066

// CliViewBaseURL 返回离线报告「查看详情」链接的绝对前缀。
func CliViewBaseURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", CliViewPort)
}

// RegisterViewRoutes 注册离线「查看详情」路由（列表 ×2 + 文件预览 ×2）。
// 列表页复用 scene_data_test_history.content 里既有的 /admin/fm/... 锚点，交互与 Web 端一致。
func RegisterViewRoutes(r *gin.Engine) {
	r.GET("/admin/info/scene_test_history", cliSceneHistoryList)
	r.GET("/admin/info/scene_data_test_history", cliDataHistoryList)
	r.GET("/admin/fm/history/preview", cliHistoryPreview)
	r.GET("/admin/fm/data/preview", cliDataPreview)
}

// startCliView 启动本地「查看详情」服务（含 mock 路由，合并到同一端口），返回 server 供优雅关闭。
// 端口已被占用时（可能已有 report 服务在跑）返回 nil，跳过重复启动。
func startCliView(port int) *http.Server {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	if conn, err := net.DialTimeout("tcp", addr, time.Second); err == nil {
		conn.Close()
		return nil
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	RegisterMockRoutes(r)
	RegisterViewRoutes(r)
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Warning("cli view server error: %s", err)
		}
	}()
	return srv
}

// waitForSignal 阻塞直到收到 Ctrl+C（SIGINT）或 SIGTERM，供 report 命令保活。
func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}

// applyViewLang 从查询参数 lang 读取语种并切换（与 Web 每请求携带语种一致），
// 使「查看详情」跳转后的列表/预览页按报告当前语种渲染。
func applyViewLang(c *gin.Context) {
	if lang := c.Query("lang"); lang != "" {
		SetLocale(lang)
	}
}

// cliSceneHistoryList 场景执行历史列表（无结果文件链接，仅场景级通过/失败）。
func cliSceneHistoryList(c *gin.Context) {
	applyViewLang(c)
	taskId := c.Query("task_id")
	var rows []struct {
		Name         string
		DataFileList string
		LastFile     string
		Result       string
		FailReason   string
		Product      string
		CreatedAt    string
	}
	models.Orm.Table("scene_test_history").Where("task_id = ?", taskId).Order("id desc").Find(&rows)

	var b strings.Builder
	b.WriteString(cliListHeader(T("playbook_test_history.title"), taskId))
	b.WriteString(`<table><thead><tr>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.name")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.data_file_list")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.last_file")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.test_result")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.fail_reason")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.product")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.created_at")) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	if len(rows) == 0 {
		b.WriteString(`<tr><td colspan="7" class="empty">` + html.EscapeString(T("cli.view_empty")) + `</td></tr>`)
	}
	for _, r := range rows {
		b.WriteString(`<tr>`)
		b.WriteString(`<td>` + html.EscapeString(r.Name) + `</td>`)
		// 复用 Web 端同款 /admin/fm/... 锚点，交互与 Web 一致
		b.WriteString(`<td>` + GetHistoryDataLinkByDataStr(r.DataFileList, GetLocale()) + `</td>`)
		b.WriteString(`<td>` + cliLastFileLink(r.LastFile) + `</td>`)
		b.WriteString(`<td>` + cliResultBadge(r.Result) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.FailReason) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.Product) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.CreatedAt) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(cliPageShell(T("playbook_test_history.title"), b.String())))
}

// cliDataHistoryList 数据执行历史列表：content 列原样输出既有的 /admin/fm/... 锚点（不转义）。
func cliDataHistoryList(c *gin.Context) {
	applyViewLang(c)
	taskId := c.Query("task_id")
	var rows []struct {
		Name       string
		ApiId      string
		Result     string
		FailReason string
		Content    string
		CreatedAt  string
	}
	models.Orm.Table("scene_data_test_history").Where("task_id = ?", taskId).Order("id desc").Find(&rows)

	var b strings.Builder
	b.WriteString(cliListHeader(T("data_test_history.title"), taskId))
	b.WriteString(`<table><thead><tr>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.name")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.api_id")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.data_content")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.test_result")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.fail_reason")) + `</th>`)
	b.WriteString(`<th>` + html.EscapeString(T("common.created_at")) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	if len(rows) == 0 {
		b.WriteString(`<tr><td colspan="6" class="empty">` + html.EscapeString(T("cli.view_empty")) + `</td></tr>`)
	}
	for _, r := range rows {
		b.WriteString(`<tr>`)
		b.WriteString(`<td>` + html.EscapeString(r.Name) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.ApiId) + `</td>`)
		b.WriteString(`<td>` + r.Content + `</td>`) // 既有锚点，原样输出
		b.WriteString(`<td>` + cliResultBadge(r.Result) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.FailReason) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.CreatedAt) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(cliPageShell(T("data_test_history.title"), b.String())))
}

// cliHistoryPreview 预览历史结果文件（HistoryBasePath 下）。
func cliHistoryPreview(c *gin.Context) {
	cliPreview(c, HistoryBasePath, c.Query("path"))
}

// cliDataPreview 预览数据文件（DataBasePath 下）。
func cliDataPreview(c *gin.Context) {
	cliPreview(c, DataBasePath, c.Query("path"))
}

// cliPreview 读取 base 下 p 指向的文件并以 <pre> 转义展示。
func cliPreview(c *gin.Context, base, p string) {
	applyViewLang(c)
	full := resolveCliPreviewPath(base, p)
	if full == "" {
		c.String(http.StatusForbidden, T("cli.view_forbidden"))
		return
	}
	data, err := ioutil.ReadFile(full)
	if err != nil {
		c.String(http.StatusNotFound, T("cli.view_file_not_found"))
		return
	}
	body := `<a class="back" href="javascript:history.back()">&laquo; ` + html.EscapeString(T("cli.view_back")) + `</a>` +
		`<pre>` + html.EscapeString(string(data)) + `</pre>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(cliPageShell(filepath.Base(full), body)))
}

// resolveCliPreviewPath 校验 p 并解析到 base 之下的绝对路径，非法返回空串（防目录穿越）。
func resolveCliPreviewPath(base, p string) string {
	if p == "" || strings.Contains(p, "..") {
		return ""
	}
	full := filepath.Join(base, p)
	cleanBase := filepath.Clean(base)
	clean := filepath.Clean(full)
	if clean != cleanBase && !strings.HasPrefix(clean, cleanBase+string(os.PathSeparator)) {
		return ""
	}
	return clean
}

// cliLastFileLink 最近数据文件：复用 GetLastFileLink 生成 /admin/fm/... 锚点（空值返回空）。
func cliLastFileLink(lastFile string) string {
	if lastFile == "" {
		return ""
	}
	return GetLastFileLink(lastFile)
}

// cliResultBadge 结果徽标：pass/fail 转 common 文案，其余原样转义。
func cliResultBadge(result string) string {
	switch result {
	case "pass":
		return `<span class="badge pass">` + html.EscapeString(T("common.pass")) + `</span>`
	case "fail":
		return `<span class="badge fail">` + html.EscapeString(T("common.fail")) + `</span>`
	default:
		return html.EscapeString(result)
	}
}

// cliListHeader 列表页顶部：返回链接 + 标题 + task_id。
func cliListHeader(title, taskId string) string {
	h := `<a class="back" href="javascript:history.back()">&laquo; ` + html.EscapeString(T("cli.view_back")) + `</a>` +
		`<h2>` + html.EscapeString(title) + `</h2>`
	if taskId != "" {
		h += `<p style="font-size:12px;color:#888">task_id: ` + html.EscapeString(taskId) + `</p>`
	}
	return h
}

// cliPageShell 轻量自绘 HTML 页面骨架。
func cliPageShell(title, body string) string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><title>` + html.EscapeString(title) + `</title><style>` + cliPageCSS + `</style></head><body>` + body + `</body></html>`
}

const cliPageCSS = `
body{font-family:-apple-system,"Segoe UI",Roboto,"Microsoft YaHei",sans-serif;margin:24px;color:#333;background:#f5f6fa}
h2{font-size:18px;margin:0 0 12px}
table{border-collapse:collapse;width:100%;background:#fff;box-shadow:0 1px 3px rgba(0,0,0,.08)}
th,td{border:1px solid #e5e7eb;padding:8px 10px;font-size:13px;text-align:left;vertical-align:top}
th{background:#f9fafb;font-weight:600}
a{color:#2563eb;text-decoration:none}
a:hover{text-decoration:underline}
.badge{display:inline-block;padding:1px 8px;border-radius:10px;font-size:12px}
.pass{color:#059669;background:#ecfdf5}
.fail{color:#dc2626;background:#fef2f2}
.back{display:inline-block;margin-bottom:12px;font-size:13px}
pre{background:#fff;border:1px solid #e5e7eb;padding:16px;overflow:auto;font-size:12px;line-height:1.5;white-space:pre-wrap;word-break:break-all}
.empty{padding:32px;text-align:center;color:#999}
`
