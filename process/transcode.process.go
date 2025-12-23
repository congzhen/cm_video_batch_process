package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type WatermarkPlacement string

const (
	WatermarkPlacement_TopRight   WatermarkPlacement = "top-right"  // 右上角
	WatermarkPlacement_Random     WatermarkPlacement = "random"     // 随机位置
	WatermarkPlacement_Horizontal WatermarkPlacement = "horizontal" // 水平横移
	WatermarkPlacement_Diagonal   WatermarkPlacement = "diagonal"   // 对角线
	WatermarkPlacement_Bounce     WatermarkPlacement = "bounce"     //弹跳
	WatermarkPlacement_Spiral     WatermarkPlacement = "spiral"     //螺旋运动
)

type VideoRotate string

const (
	VideoRotate_copy VideoRotate = "copy" // 0度
	VideoRotate_90   VideoRotate = "90"   // 90度
	VideoRotate_180  VideoRotate = "180"  // 180度
	VideoRotate_270  VideoRotate = "270"  // 270度
)

type TranscodeParams struct {
	VideoCodec         string             `json:"video_codec"`
	AudioCodec         string             `json:"audio_codec"`
	VideoHeight        string             `json:"video_height"`
	Fps                string             `json:"fps"`
	VideoBitrate       string             `json:"video_bitrate"`
	WatermarkContent   string             `json:"watermark_content"`
	WatermarkPlacement WatermarkPlacement `json:"watermark_placement"`
	Rotate             VideoRotate        `json:"rotate"`
	UseGpu             bool               `json:"use_gpu"`
	CpuThreads         int                `json:"cpu_threads"`
}

func VideoTranscodeProcessor(ctx context.Context, id, inputFilePath string, params TranscodeParams) string {
	fileName := GetFileNameFromPath(inputFilePath, true)
	outputDirectory := GetOutputDirectory()
	err := CreateFolder(outputDirectory)
	if err != nil {
		return fmt.Sprintf("创建输出目录失败: %v", err)
	}
	outputFilePath := fmt.Sprintf("%s/%s", outputDirectory, fileName)

	// 获取视频总时长（秒）
	duration, err := getVideoDuration(inputFilePath)
	if err != nil {
		fmt.Printf("警告: 无法获取视频时长: %v\n", err)
		// 即使无法获取时长也继续处理
		duration = 0
	}

	// 构建FFmpeg命令
	cmd, err := buildTranscodeCommand(inputFilePath, outputFilePath, params)
	if err != nil {
		return fmt.Sprintf("构建命令失败: %v", err)
	}
	fmt.Printf("命令: %v\n", cmd.Args)

	// 设置管道以便捕获FFmpeg输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Sprintf("创建stderr管道失败: %v", err)
	}

	// 启动FFmpeg进程
	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("启动FFmpeg失败: %v", err)
	}

	// 创建scanner读取FFmpeg输出
	scanner := bufio.NewScanner(stderr)
	progressRegex := regexp.MustCompile(`time=([0-9:.]+)`)

	// 在goroutine中读取进度
	var progressDone = make(chan struct{}) // 添加信号通道
	go func() {
		defer close(progressDone) // 处理完后关闭通道
		for scanner.Scan() {
			line := scanner.Text()
			matches := progressRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentTime := matches[1]
				if duration > 0 {
					// 计算并显示百分比进度
					currentSeconds := timeToSeconds(currentTime)
					percentage := (currentSeconds / duration) * 100
					if percentage > 100 {
						percentage = 100
					}
					// 使用 \r 实现行内更新，并添加足够的空格来覆盖之前的输出
					fmt.Printf("\r进度: %.2f%% (已处理时间: %s)     ", percentage, currentTime)
					wailsRuntime.EventsEmit(ctx, "videoTranscodeProcessor", id, percentage, currentTime)
				} else {
					// 只显示已处理时间
					fmt.Printf("\r已处理时间: %s     ", currentTime)
				}
			}
		}
	}()

	// 等待FFmpeg进程完成
	if err := cmd.Wait(); err != nil {
		return fmt.Sprintf("FFmpeg处理失败: %v", err)
	}

	// 等待进度显示goroutine完成
	<-progressDone

	// 显示最终100%进度
	if duration > 0 {
		fmt.Printf("\r进度: 100.00%% (已完成) \n")

	} else {
		fmt.Printf("\r处理完成! \n")
	}
	wailsRuntime.EventsEmit(ctx, "videoTranscodeProcessor", id, 100, "completed")
	fmt.Printf("处理视频成功: %s\n", outputFilePath)
	return "OK"
}

// buildTranscodeCommand 构建FFmpeg转码命令
func buildTranscodeCommand(inputFilePath string, outputFilePath string, params TranscodeParams) (*exec.Cmd, error) {
	// 构建FFmpeg命令参数
	var args []string

	// 输入文件
	args = append(args, "-i", inputFilePath)

	// 添加CPU线程数参数
	if params.CpuThreads > 0 {
		args = append(args, "-threads", fmt.Sprintf("%d", params.CpuThreads))
	}

	// 处理视频编码参数
	videoCodec := getVideoCodecFormat(params)
	args = append(args, "-c:v", videoCodec)

	// 处理音频编码参数
	audioCodec := getAudioCodecFormat(params)
	args = append(args, "-c:a", audioCodec)

	// 构建视频滤镜链
	var videoFilters []string

	// 处理视频高度参数
	if params.VideoHeight != "copy" {
		videoFilters = append(videoFilters, fmt.Sprintf("scale=-1:%s", params.VideoHeight))
	}

	// 如果有水印，则添加水印滤镜
	if params.WatermarkContent != "" {
		drawTextFilter := getWatermarkPlacement(params.WatermarkContent, params.WatermarkPlacement)
		videoFilters = append(videoFilters, drawTextFilter)
	}

	// 旋转
	if params.Rotate != "copy" {
		rotationFilter := getRotationFilter(string(params.Rotate))
		videoFilters = append(videoFilters, rotationFilter)
	}

	// 如果有视频滤镜，则应用到命令
	if len(videoFilters) > 0 {
		args = append(args, "-vf", strings.Join(videoFilters, ","))
	}

	// 添加帧率参数
	if params.Fps != "copy" {
		args = append(args, "-r", params.Fps)
	}

	// 处理视频比特率参数
	if params.VideoBitrate != "copy" {
		args = append(args, "-b:v", params.VideoBitrate)
	}

	// 添加进度报告参数
	args = append(args, "-progress", "pipe:2", "-nostats")

	// 输出文件
	args = append(args, outputFilePath)

	ffmpegPath, err := IsFFmpegAvailable()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg不可用: %v", err)
	}

	cmd := createCommand(ffmpegPath, args...)
	return cmd, nil
}

// getRotationFilter 获取旋转滤镜
func getRotationFilter(rotation string) string {
	switch rotation {
	case "90":
		return "transpose=1" // 顺时针旋转90度
	case "180":
		return "transpose=1,transpose=1" // 旋转180度（两个90度旋转）
	case "270":
		return "transpose=2" // 顺时针旋转270度（或逆时针90度）
	default:
		return "" // 不旋转
	}
}

// getVideoCodecFormat 获取视频编码格式
func getVideoCodecFormat(params TranscodeParams) string {
	switch params.VideoCodec {
	case "h264":
		if params.UseGpu {
			return "h264_nvenc"
		} else {
			return "libx264"
		}
	case "h265":
		if params.UseGpu {
			return "hevc_nvenc"
		} else {
			return "libx265"
		}
	default:
		// 如果要加水印，不能使用copy
		if params.WatermarkContent != "" {
			// 默认使用libx264进行重新编码
			if params.UseGpu {
				return "h264_nvenc"
			} else {
				return "libx264"
			}
		}
		return "copy"
	}
}

// getAudioCodecFormat 获取音频编码格式
func getAudioCodecFormat(params TranscodeParams) string {
	switch params.AudioCodec {
	case "aac":
		if params.UseGpu {
			return "aac"
		} else {
			return "libfdk_aac"
		}
	case "mp3":
		if params.UseGpu {
			return "mp3"
		} else {
			return "libmp3lame"
		}
	default:
		return "copy"
	}
}

// getVideoDuration 获取视频时长（秒）
func getVideoDuration(filePath string) (float64, error) {
	args := []string{
		"-i", filePath,
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "csv=p=0",
	}
	ffprobePath, err := IsFFprobeAvailable()
	if err != nil {
		return 0, fmt.Errorf("ffprobe不可用: %v", err)
	}

	cmd := createCommand(ffprobePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// timeToSeconds 将时间字符串转换为秒数 (格式: HH:MM:SS.sss)
func timeToSeconds(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0
	}

	hours, _ := strconv.ParseFloat(parts[0], 64)
	minutes, _ := strconv.ParseFloat(parts[1], 64)
	seconds, _ := strconv.ParseFloat(parts[2], 64)

	return hours*3600 + minutes*60 + seconds
}

// 获取水印位置
func getWatermarkPlacement(text string, placement WatermarkPlacement) string {
	var filter string
	fontFile := "msyh.ttc"
	// 使用 text 参数，同时对中文字符进行特殊处理
	escapedText := escapeTextForFFmpeg(text)
	switch placement {
	case WatermarkPlacement_Random:
		// 随机出现位置（在屏幕四个角之间切换）
		filter = fmt.Sprintf("drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=if(lt(sin(t*0.5)\\,0)\\,20\\,w-tw-20):y=if(lt(cos(t*0.3)\\,0)\\,20\\,h-th-20):borderw=2:bordercolor=black", fontFile, escapedText)
	case WatermarkPlacement_Horizontal:
		// 水平摆动
		filter = fmt.Sprintf(
			"drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=w/2+(w/4)*sin(2*PI*t/8):y=h/2:borderw=2:bordercolor=black",
			fontFile, escapedText)
	case WatermarkPlacement_Diagonal:
		// 对角线运动
		filter = fmt.Sprintf(
			"drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=w*t/30:y=h*t/30:borderw=2:bordercolor=black",
			fontFile, escapedText)
	case WatermarkPlacement_Bounce:
		// 随机弹跳效果
		filter = fmt.Sprintf(
			"drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=w/2+(w/3)*sin(2*PI*t/10):y=h/2+(h/3)*cos(2*PI*t/7):borderw=2:bordercolor=black",
			fontFile, escapedText)
	case WatermarkPlacement_Spiral:
		// 螺旋运动
		filter = fmt.Sprintf(
			"drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=w/2+(w/4)*(sin(2*PI*t/12)+cos(2*PI*t/6)):y=h/2+(h/4)*(cos(2*PI*t/12)-sin(2*PI*t/6)):borderw=2:bordercolor=black",
			fontFile, escapedText)
	default:
		//右上角
		filter = fmt.Sprintf(
			"drawtext=fontfile=%s:text='%s':fontcolor=white:fontsize=24:x=w-tw-20:y=20:borderw=2:bordercolor=black",
			fontFile, escapedText)
	}
	return filter
}

// escapeTextForFFmpeg 为FFmpeg转义特殊字符，特别是中文字符
func escapeTextForFFmpeg(text string) string {
	// 在Windows下，对特殊字符进行转义
	escaped := strings.ReplaceAll(text, "'", "'\\\\\\''")
	escaped = strings.ReplaceAll(escaped, ":", "\\:")
	escaped = strings.ReplaceAll(escaped, ",", "\\,")
	escaped = strings.ReplaceAll(escaped, "[", "\\[")
	escaped = strings.ReplaceAll(escaped, "]", "\\]")
	escaped = strings.ReplaceAll(escaped, "=", "\\=")
	escaped = strings.ReplaceAll(escaped, "+", "\\+")
	escaped = strings.ReplaceAll(escaped, "~", "\\~")
	escaped = strings.ReplaceAll(escaped, "%", "%%")
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "&", "\\&")
	escaped = strings.ReplaceAll(escaped, "^", "\\^")
	escaped = strings.ReplaceAll(escaped, "!", "\\!")
	return escaped
}
