package main

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// 代理配置数据结构
type ProxyConfig struct {
	Bind         string    `yaml:"bind"`
	WaitQueueLen int       `yaml:"wait_queue_len"`
	MaxConn      int       `yaml:"max_conn"`
	Timeout      int       `yaml:"timeout"`
	FailOver     int       `yaml:"failover"`
	Backend      []string  `yaml:"backend"`
	Log          LogConfig `yaml:"log"`
	Stats        string    `yaml:"stats"`
}

// 日志配置结构信息
type LogConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

// 解析配置文件
func parseConfigFile(filePath string) error {
	if conf, err := ioutil.ReadFile(filePath); err == nil {
		if err = yaml.Unmarshal(conf, &Config); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
