package main

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	LineAccessToken string    `yaml:"line_access_token"`
	CrawlDuration   int       `yaml:"crawl_duration"`
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
	for _, v := range cobj.Teachers {
		fmt.Printf("id:%s\n", v.Id)
	}
}

func (cobj *conf) setConf(confpath string) {
	if confpath == "" {
		log.Print("confpath required, but not set.")
		return
	}

	// yamlを読み込む
	buf, err := ioutil.ReadFile(confpath)
	if err != nil {
		log.Printf("Failed to read confpath file.")
	}

	// structにUnmasrshal
	err = yaml.Unmarshal([]byte(buf), cobj)

	if err != nil {
		log.Printf("Failed to Unmarshal confpath file.")
	}
}
