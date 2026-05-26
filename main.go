package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	nativedialog "github.com/sqweek/dialog"
)

// ─────────────────────────────────────────────
// Custom dark theme
// ─────────────────────────────────────────────

type darkTheme struct{}

func (d *darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 22, G: 22, B: 38, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 40, G: 42, B: 65, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 35, G: 35, B: 50, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 200, B: 240, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 220, G: 225, B: 240, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 100, G: 105, B: 130, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 30, G: 30, B: 50, A: 255}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 50, G: 52, B: 75, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 28, G: 28, B: 48, A: 240}
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (d *darkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (d *darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (d *darkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 13
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// ─────────────────────────────────────────────
// DropZone – custom widget
// ─────────────────────────────────────────────

type DropZone struct {
	widget.BaseWidget
	label  string
	active bool
	onTap  func()
}

func newDropZone(label string, onTap func()) *DropZone {
	dz := &DropZone{label: label, onTap: onTap}
	dz.ExtendBaseWidget(dz)
	return dz
}

func (dz *DropZone) Tapped(_ *fyne.PointEvent) {
	if dz.onTap != nil {
		dz.onTap()
	}
}

func (dz *DropZone) SetActive(v bool) {
	dz.active = v
	dz.Refresh()
}

func (dz *DropZone) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(colorIdle)
	bg.CornerRadius = 12
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = colorBorderIdle
	border.StrokeWidth = 2
	border.CornerRadius = 12
	icon := canvas.NewText("⬇", colorAccent)
	icon.TextSize = 36
	icon.Alignment = fyne.TextAlignCenter
	lbl := canvas.NewText(dz.label, colorLabelIdle)
	lbl.TextSize = 16
	lbl.Alignment = fyne.TextAlignCenter
	hint := canvas.NewText("点击或拖拽", colorHint)
	hint.TextSize = 12
	hint.Alignment = fyne.TextAlignCenter
	return &dropZoneRenderer{zone: dz, bg: bg, border: border, icon: icon, lbl: lbl, hint: hint}
}

// ─────────────────────────────────────────────
// dropZoneRenderer
// ─────────────────────────────────────────────

type dropZoneRenderer struct {
	zone   *DropZone
	bg     *canvas.Rectangle
	border *canvas.Rectangle
	icon   *canvas.Text
	lbl    *canvas.Text
	hint   *canvas.Text
}

func (r *dropZoneRenderer) Destroy() {}

func (r *dropZoneRenderer) MinSize() fyne.Size { return fyne.NewSize(200, 120) }

func (r *dropZoneRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.border, r.icon, r.lbl, r.hint}
}

func (r *dropZoneRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))
	r.border.Resize(size)
	r.border.Move(fyne.NewPos(0, 0))

	iconH := r.icon.MinSize().Height
	lblH := r.lbl.MinSize().Height
	hintH := r.hint.MinSize().Height
	gap := float32(8)
	total := iconH + gap + lblH + gap + hintH
	y := (size.Height - total) / 2

	r.icon.Resize(fyne.NewSize(size.Width, iconH))
	r.icon.Move(fyne.NewPos(0, y))
	r.lbl.Resize(fyne.NewSize(size.Width, lblH))
	r.lbl.Move(fyne.NewPos(0, y+iconH+gap))
	r.hint.Resize(fyne.NewSize(size.Width, hintH))
	r.hint.Move(fyne.NewPos(0, y+iconH+gap+lblH+gap))
}

func (r *dropZoneRenderer) Refresh() {
	if r.zone.active {
		r.bg.FillColor = colorActive
		r.border.StrokeColor = colorBorderActive
		r.icon.Color = colorLabelActive
		r.lbl.Color = colorLabelActive
	} else {
		r.bg.FillColor = colorIdle
		r.border.StrokeColor = colorBorderIdle
		r.icon.Color = colorAccent
		r.lbl.Color = colorLabelIdle
	}
	r.bg.Refresh()
	r.border.Refresh()
	r.icon.Refresh()
	r.lbl.Refresh()
	r.hint.Refresh()
}

// ─────────────────────────────────────────────
// Palette
// ─────────────────────────────────────────────

var (
	colorIdle         = color.NRGBA{R: 26, G: 26, B: 46, A: 255}
	colorActive       = color.NRGBA{R: 18, G: 48, B: 22, A: 255}
	colorBorderIdle   = color.NRGBA{R: 55, G: 60, B: 90, A: 200}
	colorBorderActive = color.NRGBA{R: 60, G: 190, B: 80, A: 220}
	colorAccent       = color.NRGBA{R: 0, G: 200, B: 240, A: 255}
	colorLabelIdle    = color.NRGBA{R: 150, G: 160, B: 200, A: 255}
	colorLabelActive  = color.NRGBA{R: 90, G: 220, B: 110, A: 255}
	colorHint         = color.NRGBA{R: 90, G: 95, B: 125, A: 180}
)

// ─────────────────────────────────────────────
// Convert logic
// ─────────────────────────────────────────────

// autoOutputPath derives a default output path from the input file path.
// e.g. /path/to/80.m4s → /path/to/80.mp4
func autoOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, name+".mp4")
}

// autoSrtPath derives a .srt output path from the input file path.
func autoSrtPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, name+".srt")
}

// ffmpegAvailable returns the ffmpeg path if found in PATH.
func ffmpegAvailable() string {
	path, _ := exec.LookPath("ffmpeg")
	return path
}

// whisperAvailable returns the whisper executable path.
// Falls back to scanning common Python Scripts dirs on Windows when LookPath fails.
func whisperAvailable() string {
	if p, err := exec.LookPath("whisper"); err == nil {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	matches, _ := filepath.Glob(filepath.Join(home, "AppData/Local/Programs/Python/Python*/Scripts/whisper.exe"))
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

// convertFile dispatches to FFmpeg (if available) or binary copy fallback.
func convertFile(inputPath, outputPath string, onProgress func(float64)) error {
	if ffmpegAvailable() != "" {
		return convertWithFFmpeg(inputPath, outputPath, onProgress)
	}
	return convertWithCopy(inputPath, outputPath, onProgress)
}

// convertWithFFmpeg remuxes .m4s → .mp4 without re-encoding.
func convertWithFFmpeg(inputPath, outputPath string, onProgress func(float64)) error {
	if onProgress != nil {
		onProgress(0.01)
	}

	cmd := exec.Command("ffmpeg", "-y",
		"-i", inputPath,
		"-c", "copy",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg 执行失败：%w\n输出：%s", err, string(output))
	}
	if onProgress != nil {
		onProgress(1.0)
	}
	return nil
}

// convertWithCopy does a raw binary copy when FFmpeg is not installed.
func convertWithCopy(inputPath, outputPath string, onProgress func(float64)) error {
	src, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开文件失败：%w", err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败：%w", err)
	}

	dst, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败：%w", err)
	}
	defer dst.Close()

	written := int64(0)
	buf := make([]byte, 1<<20)
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			if _, wErr := dst.Write(buf[:n]); wErr != nil {
				return fmt.Errorf("写入文件失败：%w", wErr)
			}
			written += int64(n)
			if onProgress != nil && info.Size() > 0 {
				onProgress(float64(written) / float64(info.Size()))
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("读取文件失败：%w", readErr)
		}
	}
	if onProgress != nil {
		onProgress(1.0)
	}
	return nil
}

// generateSubtitle calls the whisper CLI to produce an SRT file.
func generateSubtitle(inputPath, outputDir, model, language string, onProgress func(float64)) error {
	whisperPath := whisperAvailable()
	if whisperPath == "" {
		return fmt.Errorf("未找到 whisper 可执行文件")
	}

	totalDur := getMediaDuration(inputPath)

	args := []string{
		inputPath,
		"--output_format", "srt",
		"--output_dir", outputDir,
		"--model", model,
		"--verbose", "True",
	}
	if language != "" && language != "auto" {
		args = append(args, "--language", language)
	}
	cmd := exec.Command(whisperPath, args...)
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8", "PYTHONUTF8=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败：%w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 whisper 失败：%w", err)
	}

	tsRe := regexp.MustCompile(`\[[\d:.]+\s*-->\s*([\d:.]+)\]`)
	scanner := bufio.NewScanner(stdout)
	scanner.Split(scanProgress)
	for scanner.Scan() {
		line := scanner.Text()
		if m := tsRe.FindStringSubmatch(line); m != nil && totalDur > 0 && onProgress != nil {
			if secs := parseTimestamp(m[1]); secs > 0 {
				ratio := secs / totalDur
				if ratio > 1 {
					ratio = 1
				}
				onProgress(ratio)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("whisper 执行失败：%w", err)
	}
	if onProgress != nil {
		onProgress(1.0)
	}
	return nil
}

// getMediaDuration returns the duration in seconds via ffprobe, or 0 on failure.
func getMediaDuration(path string) float64 {
	out, err := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	).Output()
	if err != nil {
		return 0
	}
	dur, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	return dur
}

// parseTimestamp converts "MM:SS.mmm" or "HH:MM:SS.mmm" to seconds.
func parseTimestamp(ts string) float64 {
	parts := strings.Split(ts, ":")
	secs := 0.0
	for _, p := range parts {
		v, _ := strconv.ParseFloat(p, 64)
		secs = secs*60 + v
	}
	return secs
}

// scanProgress splits on \r or \n so tqdm carriage-return updates are captured.
func scanProgress(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i, b := range data {
		if b == '\r' || b == '\n' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
// ─────────────────────────────────────────────

// uriToPath converts a Fyne URI path to an OS path.
// On Windows, Fyne returns /C:/... — strip the leading slash.
func uriToPath(raw string) string {
	if len(raw) > 2 && raw[0] == '/' && raw[2] == ':' {
		return raw[1:]
	}
	return raw
}

// ─────────────────────────────────────────────
// Tab builders
// ─────────────────────────────────────────────

func buildConvertTab(w fyne.Window) (content *fyne.Container, onDrop func(string)) {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("选择 .m4s 文件")
	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("输出路径（可选，默认与输入同目录）")

	var dropZone *DropZone
	openFile := func() {
		go func() {
			path, err := nativedialog.File().Filter("M4S 文件", "m4s").Filter("所有文件", "*").Title("选择 .m4s 文件").Load()
			if err != nil {
				return
			}
			inputEntry.SetText(path)
			dropZone.SetActive(true)
		}()
	}
	dropZone = newDropZone("拖拽文件到此", openFile)

	inputBtn := widget.NewButton("选择文件", openFile)
	outputBtn := widget.NewButton("选择输出", func() {
		go func() {
			path, err := nativedialog.File().Filter("MP4 文件", "mp4").Title("选择输出路径").SetStartFile("output.mp4").Save()
			if err != nil {
				return
			}
			outputEntry.SetText(path)
		}()
	})

	progressBar := widget.NewProgressBar()
	statusLabel := widget.NewLabel("")
	statusLabel.Alignment = fyne.TextAlignCenter

	openFolderBtn := widget.NewButton("打开输出文件夹", nil)
	openFolderBtn.Hide()

	var convertBtn *widget.Button
	convertBtn = widget.NewButton("开始转换", func() {
		inputPath := inputEntry.Text
		if inputPath == "" {
			dialog.ShowError(fmt.Errorf("请先选择 .m4s 文件"), w)
			return
		}
		outputPath := outputEntry.Text
		if outputPath == "" {
			outputPath = autoOutputPath(inputPath)
			outputEntry.SetText(outputPath)
		}
		progressBar.SetValue(0)
		convertBtn.Disable()
		openFolderBtn.Hide()
		if ffmpegAvailable() != "" {
			statusLabel.SetText("转换中... (FFmpeg)")
		} else {
			statusLabel.SetText("转换中... (二进制复制)")
		}
		go func() {
			err := convertFile(inputPath, outputPath, func(ratio float64) {
				progressBar.SetValue(ratio)
			})
			if err != nil {
				statusLabel.SetText("转换失败")
				dialog.ShowError(err, w)
			} else {
				progressBar.SetValue(1)
				statusLabel.SetText("转换完成 ✓")
				finalPath := outputPath
				openFolderBtn.OnTapped = func() {
					exec.Command("explorer", "/select,", finalPath).Start()
				}
				openFolderBtn.Show()
				dialog.ShowInformation("完成", "转换成功！\n\n输出文件："+finalPath, w)
			}
			convertBtn.Enable()
		}()
	})
	convertBtn.Importance = widget.HighImportance

	bottom := container.NewVBox(
		container.NewBorder(nil, nil, nil, inputBtn, inputEntry),
		container.NewBorder(nil, nil, nil, outputBtn, outputEntry),
		widget.NewSeparator(),
		convertBtn,
		progressBar,
		statusLabel,
		openFolderBtn,
	)

	content = container.NewBorder(nil, bottom, nil, nil, dropZone)
	onDrop = func(path string) {
		inputEntry.SetText(path)
		dropZone.SetActive(true)
	}
	return
}

func buildSubtitleTab(w fyne.Window) (content *fyne.Container, onDrop func(string)) {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("选择 MP4 / 音频文件")

	var dropZone *DropZone
	openFile := func() {
		go func() {
			path, err := nativedialog.File().
				Filter("视频文件", "mp4", "mkv", "avi", "webm").
				Filter("音频文件", "mp3", "wav", "m4a", "flac").
				Filter("所有文件", "*").
				Title("选择 MP4 / 音频文件").Load()
			if err != nil {
				return
			}
			inputEntry.SetText(path)
			dropZone.SetActive(true)
		}()
	}
	dropZone = newDropZone("拖拽 MP4/音频文件到此", openFile)

	inputBtn := widget.NewButton("选择文件", openFile)

	modelSelect := widget.NewSelect([]string{"tiny", "base", "small", "medium", "large"}, nil)
	modelSelect.SetSelected("base")
	languageSelect := widget.NewSelect([]string{"auto", "zh", "en", "ja", "ko"}, nil)
	languageSelect.SetSelected("auto")

	modelRow := container.NewBorder(nil, nil, widget.NewLabel("模型"), nil, modelSelect)
	langRow := container.NewBorder(nil, nil, widget.NewLabel("语言"), nil, languageSelect)

	progress := widget.NewProgressBar()
	statusLabel := widget.NewLabel("")
	statusLabel.Alignment = fyne.TextAlignCenter

	openFolderBtn := widget.NewButton("打开输出文件夹", nil)
	openFolderBtn.Hide()

	var genBtn *widget.Button
	genBtn = widget.NewButton("生成字幕", func() {
		inputPath := inputEntry.Text
		if inputPath == "" {
			dialog.ShowError(fmt.Errorf("请先选择文件"), w)
			return
		}
		if whisperAvailable() == "" {
			dialog.ShowError(fmt.Errorf("未检测到 whisper 命令行工具\n请先安装：pip install openai-whisper"), w)
			return
		}

		model := modelSelect.Selected
		language := languageSelect.Selected
		outputDir := filepath.Dir(inputPath)
		srtPath := autoSrtPath(inputPath)

		genBtn.Disable()
		openFolderBtn.Hide()
		progress.SetValue(0)
		statusLabel.SetText("正在生成字幕... (模型: " + model + ")")

		go func() {
			err := generateSubtitle(inputPath, outputDir, model, language, func(ratio float64) {
				progress.SetValue(ratio)
			})
			progress.SetValue(1)
			if err != nil {
				statusLabel.SetText("生成失败")
				dialog.ShowError(err, w)
			} else {
				statusLabel.SetText("字幕生成完成 ✓")
				openFolderBtn.OnTapped = func() {
					exec.Command("explorer", "/select,", srtPath).Start()
				}
				openFolderBtn.Show()
				dialog.ShowInformation("完成", "字幕生成成功！\n\n输出文件："+srtPath, w)
			}
			genBtn.Enable()
		}()
	})
	genBtn.Importance = widget.HighImportance

	bottom := container.NewVBox(
		container.NewBorder(nil, nil, nil, inputBtn, inputEntry),
		modelRow,
		langRow,
		widget.NewSeparator(),
		genBtn,
		progress,
		statusLabel,
		openFolderBtn,
	)

	content = container.NewBorder(nil, bottom, nil, nil, dropZone)
	onDrop = func(path string) {
		inputEntry.SetText(path)
		dropZone.SetActive(true)
	}
	return
}

// ─────────────────────────────────────────────
// main
// ─────────────────────────────────────────────

func main() {
	a := app.New()
	a.Settings().SetTheme(&darkTheme{})
	w := a.NewWindow("videotrans")
	w.Resize(fyne.NewSize(650, 500))
	w.SetFixedSize(false)

	convertContent, convertDrop := buildConvertTab(w)
	subtitleContent, subtitleDrop := buildSubtitleTab(w)

	tabs := container.NewAppTabs(
		container.NewTabItem("视频转换", container.NewPadded(convertContent)),
		container.NewTabItem("字幕生成", container.NewPadded(subtitleContent)),
	)

	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		if len(uris) == 0 {
			return
		}
		path := uriToPath(uris[0].Path())
		switch tabs.SelectedIndex() {
		case 1:
			subtitleDrop(path)
		default:
			convertDrop(path)
		}
	})

	w.SetContent(container.NewPadded(tabs))
	w.ShowAndRun()
}
