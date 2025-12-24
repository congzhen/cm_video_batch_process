package process

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type P_Dialog struct {
}

// OpenMultipleVideoFilesDialog 打开多个文件选择对话框
func (p P_Dialog) OpenMultipleVideoFilesDialog(ctx context.Context) {
	// 使用 Wails 的 OpenFileDialog 方法，设置允许多选
	files, err := runtime.OpenMultipleFilesDialog(ctx, runtime.OpenDialogOptions{
		Title: "选择视频文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "视频文件 (*.mp4;*.avi;*.mov;*.wmv;*.flv;*.mkv)",
				Pattern:     "*.mp4;*.avi;*.mov;*.wmv;*.flv;*.mkv",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
		ShowHiddenFiles: false,
	})

	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("打开文件对话框失败: %v", err))
		// 发送错误事件到前端
		runtime.EventsEmit(ctx, "filesSelectedMultipleVideoFilesError", fmt.Sprintf("打开文件对话框失败: %v", err))
		return
	}

	// 检查用户是否取消了选择（没有选择任何文件）
	if len(files) == 0 {
		// 发送取消事件到前端
		runtime.EventsEmit(ctx, "filesSelectedMultipleVideoFilesCancelled", "用户取消了文件选择")
		return
	}

	// 处理选中的文件
	if len(files) > 0 {
		var videoInfo []VideoInfo
		for _, file := range files {
			info, err := GetVideoInfo(file)
			if err != nil {
				runtime.LogError(ctx, fmt.Sprintf("获取视频信息失败: %v", err))
				continue
			}
			videoInfo = append(videoInfo, info)
		}
		// 将选中的文件路径发送到前端
		runtime.EventsEmit(ctx, "filesSelectedMultipleVideoFilesSuccess", videoInfo)
	}
}

// OpenDirectoryDialogSetOutput 打开目录选择对话框，选择输出目录
// 选择完成后，将选中的目录路径通过 events 发送到前端
func (p P_Dialog) OpenDirectoryDialogSetOutput(ctx context.Context) {
	// 使用 Wails 的 OpenDirectoryDialog 方法打开目录选择对话框
	directory, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title:           "选择输出目录",
		ShowHiddenFiles: false,
	})

	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("打开目录选择对话框失败: %v", err))
		return
	}
	// 检查用户是否取消了选择
	if directory == "" {
		return
	}
	Config.OutputDirectory = directory
	SaveConfig()
	// 将选中的目录路径发送到前端
	runtime.EventsEmit(ctx, "directorySelectedSetOutput", directory)
}
func (p P_Dialog) OpenWatermarkImageDialog(ctx context.Context) {
	file, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "选择水印图片",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "图片文件 (*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.webp)",
				Pattern:     "*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.webp",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
		ShowHiddenFiles: false,
	})
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("打开文件对话框失败: %v", err))
		// 发送错误事件到前端
		runtime.EventsEmit(ctx, "fileSelectedWatermarkImageError", fmt.Sprintf("打开文件对话框失败: %v", err))
		return
	}
	if file == "" {
		// 发送取消事件到前端
		runtime.EventsEmit(ctx, "fileSelectedWatermarkImageCancelled", "用户取消了文件选择")
		return
	} else {
		runtime.EventsEmit(ctx, "fileSelectedWatermarkImageSuccess", file)
	}
}
