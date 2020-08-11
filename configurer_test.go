package configurer

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestConfigurerJson(t *testing.T) {
	err := ParseConf("configurer.json", false)
	if nil != err {
		fmt.Println(err)
	}

	bs, _ := json.Marshal(&cnf)
	fmt.Println(string(bs))

	//fmt.Println(GetConfByFName("kafka_consumer_conf"))
	fmt.Println(GetCommonConfItem("single_log_max_size"))
}

func TestConfigurerYaml(t *testing.T) {
	err := ParseConf("configurer.yaml", false)
	if nil != err {
		fmt.Println(err)
	}

	//bs, _ := json.Marshal(&cnf)
	//fmt.Println(string(bs))

	fmt.Println(GetConfByFName("kafka_consumer_conf"))
	fmt.Println(GetCommonConfItem("single_log_max_size"))
}

func TestConfigurerToml(t *testing.T) {
	err := ParseConf("configurer.toml", false)
	if nil != err {
		fmt.Println(err)
	}

	//bs, _ := json.Marshal(&cnf)
	//fmt.Println(string(bs))

	fmt.Println(GetConfByFName("kafka_consumer_conf"))
	fmt.Println(GetCommonConfItem("single_log_max_size"))
}
