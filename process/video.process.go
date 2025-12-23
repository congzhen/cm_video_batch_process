package process

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type VideoInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	Thumbnail    string  `json:"thumbnail"` // 现在存储base64编码的图片数据
	Size         int64   `json:"size"`
	Duration     float64 `json:"duration"`
	Bitrate      int     `json:"bitrate"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	FPS          int     `json:"fps"`
	AudioCodec   string  `json:"audio_codec"`
	VideoCodec   string  `json:"video_codec"`
	VideoBitrate int     `json:"video_bitrate"` // 新增视频码率字段
	AudioBitrate int     `json:"audio_bitrate"` // 新增音频码率字段
}

// FFprobe 输出的原始 JSON 结构
type FFprobeOutput struct {
	Format  Format   `json:"format"`
	Streams []Stream `json:"streams"`
}

type Format struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
}

type Stream struct {
	CodecType    string `json:"codec_type"`
	CodecName    string `json:"codec_name"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	AvgFrameRate string `json:"avg_frame_rate,omitempty"`
	Duration     string `json:"duration,omitempty"`
	BitRate      string `json:"bit_rate,omitempty"`
}

func GetVideoInfo(path string) (VideoInfo, error) {
	var info VideoInfo
	info.ID = GetXid()
	// 检查ffprobe是否可用
	ffprobePath, err := IsFFprobeAvailable()
	if err != nil {
		return info, fmt.Errorf("ffprobe不可用: %v", err)
	}

	// 使用ffprobe获取视频信息
	cmd := createCommand(ffprobePath, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", path)
	output, err := cmd.Output()
	if err != nil {
		return info, fmt.Errorf("无法获取视频信息: %v", err)
	}

	// 解析JSON输出
	var ffprobeData FFprobeOutput
	err = json.Unmarshal(output, &ffprobeData)
	if err != nil {
		return info, fmt.Errorf("无法解析视频信息: %v", err)
	}

	// 填充基本信息
	info.Path = ffprobeData.Format.Filename
	info.Name = GetFileNameFromPath(path, true)
	if info.Name == "" {
		info.Name = path[strings.LastIndex(path, "\\")+1:]
	}

	// 解析大小
	if size, err := strconv.ParseInt(ffprobeData.Format.Size, 10, 64); err == nil {
		info.Size = size
	}

	// 解析持续时间
	if duration, err := strconv.ParseFloat(ffprobeData.Format.Duration, 64); err == nil {
		info.Duration = duration
	}

	// 解析总比特率
	if bitRate, err := strconv.Atoi(ffprobeData.Format.BitRate); err == nil {
		info.Bitrate = bitRate
	}

	// 查找视频流和音频流
	for _, stream := range ffprobeData.Streams {
		switch stream.CodecType {
		case "video":
			info.VideoCodec = stream.CodecName
			info.Width = stream.Width
			info.Height = stream.Height

			// 解析帧率
			if stream.AvgFrameRate != "" && stream.AvgFrameRate != "0/0" {
				parts := strings.Split(stream.AvgFrameRate, "/")
				if len(parts) == 2 {
					numerator, err1 := strconv.Atoi(parts[0])
					denominator, err2 := strconv.Atoi(parts[1])
					if err1 == nil && err2 == nil && denominator != 0 {
						info.FPS = numerator / denominator
					}
				}
			}

			// 如果没有从format获取duration，则从stream获取
			if info.Duration == 0 && stream.Duration != "" {
				if duration, err := strconv.ParseFloat(stream.Duration, 64); err == nil {
					info.Duration = duration
				}
			}

			// 如果没有从format获取bitrate，则从stream获取视频码率
			if stream.BitRate != "" {
				if bitRate, err := strconv.Atoi(stream.BitRate); err == nil {
					info.VideoBitrate = bitRate
					// 如果总码率未设置，则使用视频码率作为总码率
					if info.Bitrate == 0 {
						info.Bitrate = bitRate
					}
				}
			}

		case "audio":
			info.AudioCodec = stream.CodecName
			// 解析音频码率
			if stream.BitRate != "" {
				if bitRate, err := strconv.Atoi(stream.BitRate); err == nil {
					info.AudioBitrate = bitRate
				}
			}
		}
	}

	// 生成缩略图并转换为base64
	thumbnailBase64, err := generateThumbnailBase64(path)
	if err == nil {
		info.Thumbnail = thumbnailBase64
	}

	return info, nil
}

// generateThumbnailBase64 从视频中生成缩略图并返回base64编码的数据
func generateThumbnailBase64(videoPath string) (string, error) {
	// 检查ffmpeg是否可用
	ffmpegPath, err := IsFFmpegAvailable()
	if err != nil {
		return "", fmt.Errorf("ffmpeg不可用: %v", err)
	}

	// 使用ffmpeg生成缩略图（截取视频中间的一帧）
	// 输出JPEG格式的图片数据到stdout
	cmd := createCommand(
		ffmpegPath,
		"-i", videoPath, // 输入视频文件
		"-ss", "00:00:01", // 截取时间点（1秒处）
		"-vframes", "1", // 只截取一帧
		"-f", "image2", // 强制输出格式为图像
		"-preset", "ultrafast", // 最快编码速度
		"pipe:1", // 输出到stdout
	)

	// 捕获stdout输出
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = nil // 忽略stderr输出

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("生成缩略图失败: %v", err)
	}

	// 检查是否有输出数据
	if stdout.Len() == 0 {
		return "", fmt.Errorf("生成缩略图失败: 没有输出数据")
	}

	// 将图片数据转换为base64编码
	encoded := base64.StdEncoding.EncodeToString(stdout.Bytes())

	// 添加data URI前缀，便于前端直接使用
	return "data:image/jpeg;base64," + encoded, nil
}
