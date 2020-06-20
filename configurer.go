package configurer

import (
	"errors"
	"flag"
)

type conf struct {
	CommonConfItem map[string]interface{}
	ConfFile       map[string]string //key:filename value:file-content
}

var (
	cmdEnv string
	cmdKey string
	cnf    conf
)

func init() {
	flag.StringVar(&cmdEnv, "e", "sandbox", "指定环境变量")
	flag.StringVar(&cmdKey, "k", "", "某环境下的指定key")

	cnf.ConfFile = make(map[string]string, 0)
}

// GetConfByFName 根据文件名称获取配置内容
func GetConfByFName(fname string) (string, error) {
	content, ok := cnf.ConfFile[fname]
	if ok {
		return content, nil
	}

	return "", errors.New(fname + " not existed")
}

// GetCommonConfItem 获取公共配置项
func GetCommonConfItem(key string) (interface{}, error) {
	item, ok := cnf.CommonConfItem[key]
	if ok {
		return item, nil
	}

	return nil, errors.New(key + " not existed")
}
