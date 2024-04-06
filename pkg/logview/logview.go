package logview

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"runtime"
)

type LogLevel int

const (
	LogInfo LogLevel = iota
	LogDebug
	LogError
	LogWarm
	LogVerbose
)

type LogView interface {
	Print(level LogLevel, log string)
	GetView() fyne.CanvasObject
	SetLogLineSize(maxSize int)
	GetText() *bytes.Buffer
	Clear()
}

type LogViewData struct {
	autoScroll  bool
	scroll      *container.Scroll
	status      *widget.Label
	c           *fyne.Container
	window      fyne.Window
	logLineSize int
}

func (d *LogViewData) Clear() {
	d.c.RemoveAll()
}

func (d *LogViewData) SetLogLineSize(maxSize int) {
	if maxSize <= 100 {
		return
	}
	d.logLineSize = maxSize
}

func (d *LogViewData) GetView() fyne.CanvasObject {
	return container.NewBorder(container.NewCenter(d.status), nil, nil, nil, d.scroll)
	//return d.scroll
}

func NewLogView(window fyne.Window) LogView {
	c := &LogViewData{
		window:      window,
		status:      widget.NewLabel(""),
		autoScroll:  true,
		logLineSize: 2000,
	}
	c.initView()
	return c
}

func (d *LogViewData) getSize() int {
	size := len(d.c.Objects)
	return size
}

func (d *LogViewData) GetText() *bytes.Buffer {
	buff := new(bytes.Buffer)
	for _, item := range d.c.Objects {
		var str = item.(*canvas.Text).Text
		buff.WriteString(str)
		buff.WriteString("\n")
	}
	return buff
}

func (d *LogViewData) handleTypedKey(ke *fyne.KeyEvent) {
	if d.getSize() <= 0 {
		return
	}
	delta := d.c.Objects[0].Size().Height
	switch ke.Name {
	case fyne.KeyUp:
		d.scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DX: delta}})
	case fyne.KeyDown:
		d.scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.Delta{DY: -delta}})
	default:
		return
	}
	d.autoScroll = false
}

func (d *LogViewData) handleTypedRune(r rune) {
	switch r {
	case 't':
		d.scroll.ScrollToTop()
		d.autoScroll = false
	case 'b':
		d.scroll.ScrollToBottom()
		d.autoScroll = true
	case 'c':
		d.c.RemoveAll()
	}
}

func (d *LogViewData) Print(level LogLevel, line string) {
	nl := widget.NewLabel(line)
	nl.TextStyle.Bold = true
	nl.Wrapping = fyne.TextWrapBreak
	switch level {
	case LogInfo:
		nl.Importance = widget.HighImportance
		d.c.Add(nl)
	case LogError:
		nl.Importance = widget.DangerImportance
		d.c.Add(nl)
	case LogDebug:
		nl.Importance = widget.MediumImportance
		d.c.Add(nl)
	case LogWarm:
		nl.Importance = widget.WarningImportance
		d.c.Add(nl)
	default:
		nl.Importance = widget.LowImportance
		d.c.Add(nl)
	}
	d.c.Refresh()
	if d.getSize() > d.logLineSize {
		d.c.RemoveAll()
	}

	title := fmt.Sprintf("[%s/%s] 日志行数: %d", runtime.GOOS, runtime.GOARCH, d.getSize())
	d.status.SetText(title)
	d.status.Refresh()
	//d.window.SetTitle(title)
	d.scroll.ScrollToBottom()
}

func (d *LogViewData) initView() {
	d.c = container.New(layout.NewVBoxLayout())
	d.scroll = container.NewScroll(d.c)
	d.window.Canvas().SetOnTypedKey(d.handleTypedKey)
	d.window.Canvas().SetOnTypedRune(d.handleTypedRune)
}
