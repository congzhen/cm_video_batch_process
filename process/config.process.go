package process

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

var Config *ConfigData

type ConfigData struct {
	OutputDirectory string `yaml:"outputDirectory" json:"outputDirectory"`
}

func initConf() {
	c := &ConfigData{}
	yamlConf, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("无法读取配置文件: %v，将创建默认配置文件", err)
		err = createDefaultConfig()
		if err != nil {
			panic(fmt.Errorf("创建默认配置文件失败:%s", err))
		}

		// 重新读取配置文件
		yamlConf, err = os.ReadFile(configFile)
		if err != nil {
			panic(fmt.Errorf("重新读取配置文件失败:%s", err))
		}
	}

	// 先将配置解析到空结构体中
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("config Init Unmarshal: %v", err)
	}
	Config = c
}
func getDefaultConfig() *ConfigData {
	return &ConfigData{
		OutputDirectory: "",
	}
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig() error {
	defaultConfig := getDefaultConfig()
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("写入默认配置文件失败: %v", err)
	}

	log.Println("默认配置文件创建成功")
	return nil
}

// SaveConfig 将当前配置保存到YAML文件
func SaveConfig() error {
	data, err := yaml.Marshal(Config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	log.Println("配置文件保存成功")
	return nil
}

func GetOutputDirectory() string {
	outputDirectory := Config.OutputDirectory
	if outputDirectory == "" {
		// 获取当前执行文件的目录
		executable, err := os.Executable()
		if err == nil {
			outputDirectory = filepath.Dir(executable)
		} else {
			// 备用方案：使用当前工作目录
			outputDirectory, _ = os.Getwd()
		}
		outputDirectory = path.Join(outputDirectory, "output")
	}
	// 标准化路径格式
	outputDirectory = filepath.Clean(outputDirectory)
	return outputDirectory
}
