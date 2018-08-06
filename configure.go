package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	LineAccessToken string    `yaml:"line_access_token"`
	CrawlDuration   int       `yaml:"crawl_duration"`
	LogDirPath      string    `yaml:"log_dir_path"`
	Teachers        []teacher `yaml:"teachers"`
}

type teacher struct {
	Id string `yaml:"id"`
}

func newConf() *conf {
	return &conf{}
}

func (cobj *conf) PrintConf() {
	fmt.Printf("LineAccessToken:%s\n", cobj.LineAccessToken)
	fmt.Printf("CrawlDuration:%d\n", cobj.CrawlDuration)
	fmt.Printf("LogDirPath:%s\n", cobj.LogDirPath)
	for _, v := range cobj.Teachers {
		fmt.Printf("id:%s\n", v.Id)
	}
}

func (cobj *conf) setConf(confpath string) error {
	if confpath == "" {
		return fmt.Errorf("err: confpath required, but not set.")
	}

	// yamlを読み込む
	buf, err := ioutil.ReadFile(confpath)
	if err != nil {
		fmt.Printf("Failed to read confpath file.")
		return err
	}

	// structにUnmasrshal
	err = yaml.Unmarshal([]byte(buf), cobj)

	if err != nil {
		fmt.Printf("Failed to read confpath file.")
	}
	return err
}
