package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetFileNameFromPath 从给定的文件路径中提取文件名
// 参数:
// path: 文件的完整路径
// withExtension: 是否包含扩展名，true表示包含，false表示不包含
// 返回值:
// string: 文件名
func GetFileNameFromPath(path string, withExtension bool) string {
	filename := filepath.Base(path)
	if withExtension {
		return filename
	}

	// 移除扩展名
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// GetDirPathFromFilePath 根据文件路径获取文件所在的文件夹路径
// 参数:
// filePath: 文件的完整路径
// 返回值:
// string: 文件所在的文件夹路径
func GetDirPathFromFilePath(filePath string) string {
	return filepath.ToSlash(filepath.Dir(filePath))
}

// TrimBasePath 从文件路径中移除基础路径部分
// 参数:
// fullPath: 完整的文件路径
// basePath: 要移除的基础路径
// 返回值:
// string: 移除基础路径后的相对路径部分
func TrimBasePath(fullPath, basePath string) string {
	// 标准化路径分隔符
	fullPath = filepath.ToSlash(fullPath)
	basePath = filepath.ToSlash(basePath)

	// 确保basePath以斜杠结尾
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// 移除基础路径部分
	return strings.TrimPrefix(fullPath, basePath)
}

// FileExists 判断文件或目录是否存在
// 参数:
// path: 要检查的文件或目录路径
// 返回值:
// bool: 文件或目录存在返回true，否则返回false
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetDirNameFromFilePath 根据文件路径获取文件所在的文件夹名称
// 参数:
// filePath: 文件的完整路径
// 返回值:
// string: 文件所在的文件夹名称
func GetDirNameFromFilePath(filePath string) string {
	dirPath := filepath.Dir(filePath)
	return filepath.Base(dirPath)
}

// FileExt 获取文件路径的扩展名
// 参数:
// filePath: 文件的完整路径
// 返回值:
// string: 文件的扩展名，包括前导点号（如".txt"），如果文件没有扩展名则返回空字符串
func FileExt(filePath string) string {
	return strings.ToLower(filepath.Ext(filePath))
}

// CreateDirectory 创建指定路径的目录，如果父目录不存在则会自动创建
// path: 要创建的目录路径
// 返回值: 如果创建成功返回nil，否则返回错误信息
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// 给定路径创建文件夹
func CreateFolder(folderPath string) error {
	if !FileExists(folderPath) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteStringToFile 将字符串写入指定文件
// 参数:
// filePath: 文件的完整路径
// content: 要写入的字符串内容
// 返回值:
// error: 错误信息，如果写入成功则为nil
func WriteStringToFile(filePath string, content string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建或截断文件并写入内容
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// SanitizePath 对路径进行安全处理，防止路径遍历攻击
// 参数:
// path: 需要处理的路径
// 返回值:
// string: 安全处理后的路径
func SanitizePath(path string) string {
	// 使用filepath.Clean清理路径，移除冗余的元素如..和.
	cleanPath := filepath.Clean(path)

	// 移除路径开头的分隔符，防止绝对路径
	cleanPath = strings.TrimPrefix(cleanPath, string(filepath.Separator))

	// 移除路径开头的../等相对路径元素
	for strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) ||
		strings.HasPrefix(cleanPath, "..") {
		if strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
			cleanPath = strings.TrimPrefix(cleanPath, ".."+string(filepath.Separator))
		} else if strings.HasPrefix(cleanPath, "..") {
			cleanPath = strings.TrimPrefix(cleanPath, "..")
		}
	}

	return cleanPath
}

// ReadFile 读取指定路径的文件内容
// 参数:
// fullPath: 文件路径
// 返回值:
// []byte: 文件内容
// error: 错误信息
func ReadFile(fullPath string) ([]byte, error) {
	// 检查文件是否存在
	if !FileExists(fullPath) {
		return nil, fmt.Errorf("file does not exist")
	}
	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file")
	}

	return content, nil
}
