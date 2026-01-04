package process

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 注册拖拽监听事件
func addDraggedFilesHandle(ctx context.Context) {
	wailsRuntime.OnFileDrop(ctx, func(x int, y int, filepaths []string) {
		ShowLoading(ctx)
		DraggedFilesHandle(ctx, filepaths)
		HideLoading(ctx)
	})
}

// 移除拖拽监听事件
func removeDraggedFilesHandle(ctx context.Context) {
	wailsRuntime.OnFileDropOff(ctx)
}

// 处理拖拽文件事件
func DraggedFilesHandle(ctx context.Context, filepaths []string) {
	// 过滤视频文件
	videoFiles := []string{}
	for _, path := range filepaths {
		if IsVideoFile(path) {
			videoFiles = append(videoFiles, path)
		}
	}

	if len(videoFiles) > 0 {
		// 获取视频信息并发送到前端
		var videoInfoList []VideoInfo
		for _, filePath := range videoFiles {
			info, err := GetVideoInfo(filePath)
			if err != nil {
				continue
			}
			videoInfoList = append(videoInfoList, info)
		}

		// 发送事件到前端
		wailsRuntime.EventsEmit(ctx, "filesSelectedMultipleVideoFilesSuccess", videoInfoList)
	}
}
