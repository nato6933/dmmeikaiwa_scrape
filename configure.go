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
	StartTime       string    `yaml:"start_time"`
	EndTime         string    `yaml:"end_time"`
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
	fmt.Printf("StartTime:%s\n", cobj.StartTime)
	fmt.Printf("EndTime:%s\n", cobj.EndTime)
	for _, v := range cobj.Teachers {
		fmt.Printf("id:%s\n", v.Id)
	}
}

func timeScan(t string) (int, int, bool) {
	var t_h, t_m int = 0, 0
	num, err := fmt.Sscanf(t, "%d:%d", &t_h, &t_m)
	if num != 2 || err != nil {
		return -1, -1, false
	}

	if (1 < t_h && t_h < 26) && (t_m == 0 || t_m == 30) {
		return t_h, t_m, true
	}
	return t_h, t_m, false
}

func checkLessonTime(st string, et string) bool {

	var st_h, st_m, et_h, et_m int = 0, 0, 0, 0
	res := false
	tmp_res := false

	if st_h, st_m, tmp_res = timeScan(st); !tmp_res {
		res = false
		goto RETURN
	}

	if et_h, et_m, tmp_res = timeScan(et); !tmp_res {
		res = false
		goto RETURN
	}

	if st_h <= et_h && st_m <= et_m {
		res = true
	}

RETURN:
	return res

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

	if res := checkLessonTime(cobj.StartTime, cobj.EndTime); !res {
		fmt.Printf("Invalid start_time or/and end_time. Set start_time:02:00, end_time:25:30")
		cobj.StartTime = "02:00"
		cobj.EndTime = "25:30"
	}

	return err
}
