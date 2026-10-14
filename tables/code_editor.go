package tables

import (
	"data4test/biz"
	"html/template"
	"strings"
)

// aceEditorJS 生成把某个 textarea 字段替换成 Ace 编辑器（带行号）的 JS，
// 并附带右上角全屏按钮与右下角拖拽缩放。
// fieldName 为 textarea 的 name 属性；base 为 CSS 安全且页面内唯一的标识前缀，
// 用于生成编辑器容器 id / class，避免同一页面存在多个编辑器时相互冲突。
func aceEditorJS(fieldName, base string) template.JS {
	fullscreenTitle := biz.T("scene_data.fullscreen")
	exitFullscreenTitle := biz.T("scene_data.exit_fullscreen")

	js := `
$(function () {
    var ta = $('textarea[name="__FIELD__"]');
    if (!ta.length || window.ace === undefined) return;
    ta.hide();

    if (!$('#__BASE__-style').length) $('head').append('<style id="__BASE__-style">' +
        '.__BASE__-wrap{position:relative;width:100%;height:320px;}' +
        '.__BASE__-wrap #__BASE__{position:absolute;top:0;left:0;right:0;bottom:0;border:1px solid #d2d6de;}' +
        '.__BASE__-fullscreen{position:absolute;top:6px;right:6px;z-index:10;width:30px;height:30px;line-height:30px;text-align:center;border:1px solid #d2d6de;background:#fff;color:#666;border-radius:3px;cursor:pointer;}' +
        '.__BASE__-fullscreen:hover{background:#f4f4f4;}' +
        '.__BASE__-resize{position:absolute;right:0;bottom:0;width:16px;height:16px;z-index:10;cursor:nwse-resize;}' +
        '.__BASE__-resize:after{content:"";position:absolute;right:3px;bottom:3px;width:8px;height:8px;border-right:2px solid #999;border-bottom:2px solid #999;}' +
        '.__BASE__-fullscreen-on{position:fixed!important;top:0;left:0;right:0;bottom:0;width:auto!important;height:auto!important;z-index:9999;background:#fff;}' +
        '.__BASE__-fullscreen-on #__BASE__{border:none;}' +
        '.__BASE__-fullscreen-on .__BASE__-resize{display:none;}' +
        '</style>');

    var wrap = $('<div class="__BASE__-wrap"></div>').insertBefore(ta);
    $('<div id="__BASE__"></div>').appendTo(wrap);

    var ed = ace.edit('__BASE__');
    ed.session.setMode('ace/mode/text');
    ed.setValue(ta.val(), -1);
    ed.setOptions({fontSize: 14, showPrintMargin: false});
    ed.session.on('change', function () { ta.val(ed.getValue()); });

    var btn = $('<button type="button" class="__BASE__-fullscreen" title="__FULLSCREEN__"><i class="fa fa-expand"></i></button>').appendTo(wrap);
    btn.on('click', function () {
        var on = wrap.toggleClass('__BASE__-fullscreen-on').hasClass('__BASE__-fullscreen-on');
        btn.find('i').toggleClass('fa-expand fa-compress');
        btn.attr('title', on ? '__EXIT__' : '__FULLSCREEN__');
        ed.resize();
        ed.focus();
    });
    $(document).on('keydown.__BASE__', function (e) {
        if (e.key === 'Escape' && wrap.hasClass('__BASE__-fullscreen-on')) btn.trigger('click');
    });

    var handle = $('<div class="__BASE__-resize"></div>').appendTo(wrap);
    handle.on('mousedown', function (e) {
        e.preventDefault();
        var startX = e.clientX, startY = e.clientY;
        var startW = wrap.width(), startH = wrap.height();
        function onMove(ev) {
            wrap.css('width', Math.max(startW + (ev.clientX - startX), 200) + 'px');
            wrap.css('height', Math.max(startH + (ev.clientY - startY), 320) + 'px');
            ed.resize();
        }
        function onUp() {
            $(document).off('mousemove.__BASE__-resize mouseup.__BASE__-resize');
            $('body').css('user-select', '');
        }
        $(document).on('mousemove.__BASE__-resize', onMove).on('mouseup.__BASE__-resize', onUp);
        $('body').css('user-select', 'none');
    });

    if (window.ResizeObserver) {
        new ResizeObserver(function () {
            if (wrap[0].offsetWidth > 0 && wrap[0].offsetHeight > 0) ed.resize();
        }).observe(wrap[0]);
    }
});
`

	js = strings.ReplaceAll(js, "__FIELD__", fieldName)
	js = strings.ReplaceAll(js, "__BASE__", base)
	js = strings.ReplaceAll(js, "__FULLSCREEN__", fullscreenTitle)
	js = strings.ReplaceAll(js, "__EXIT__", exitFullscreenTitle)
	return template.JS(js)
}
