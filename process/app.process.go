package process

import (
	"context"

	"github.com/skratchdot/open-golang/open"
)

type App struct {
	ctx context.Context
}

type AppData struct {
	OutputDirectory string `json:"outputDirectory"`
	CPUThread       int    `json:"cpuThread"`
	GPU             bool   `json:"gpu"`
}

// Startup 应用启动时的初始化逻辑
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	initConf()
	// 注册拖拽监听事件
	addDraggedFilesHandle(ctx)
}

// Shutdown 应用关闭时的清理逻辑
func (a *App) Shutdown(ctx context.Context) {
	//移除拖拽监听事件
	removeDraggedFilesHandle(ctx)
}

// BeforeClose 在应用关闭前调用，返回true表示允许关闭，false表示取消关闭
func (a *App) BeforeClose(ctx context.Context) bool {
	return true
}

func (a *App) AppData() AppData {
	outputDirectory := GetOutputDirectory()
	return AppData{
		OutputDirectory: outputDirectory,
		CPUThread:       GetCPUThreadCount(),
		GPU:             IsGPUSupported(),
	}
}

func (a *App) OpenMultipleVideoFilesDialog() {
	P_Dialog{}.OpenMultipleVideoFilesDialog(a.ctx)
}

func (a *App) OpenWatermarkImageDialog() {
	P_Dialog{}.OpenWatermarkImageDialog(a.ctx)
}

func (a *App) OpenOutputDirectory() {
	outputDirectory := GetOutputDirectory()
	if !FileExists(outputDirectory) {
		CreateDirectory(outputDirectory)
	}
	open.Run(outputDirectory)
}
func (a *App) OpenDirectoryDialogSetOutput() {
	P_Dialog{}.OpenDirectoryDialogSetOutput(a.ctx)
}

func (a *App) Transcode(id, path string, params TranscodeParams) string {
	return VideoTranscodeProcessor(a.ctx, id, path, params)
}

func (a *App) OpenTranscodeVideo(path string) {
	open.Run(path)
}
