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

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	initConf()
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
