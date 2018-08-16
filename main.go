package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	conf_path       = "./conf/setting.yaml"
	tmpl_path       = "./conf/notify.tmpl"
	log_name        = "parse.log"
	prev_name       = "previous_schedule.gob"
	fileOptDefault  = ""
	msg_per_notify  = 4
	notify_duration = 1
	StrCanReserve   = "予約可"
)

var (
	fileOpt = flag.String("f", fileOptDefault, "Set a path to file for your develop. If this option will be set, this command will read the file instead of to access DMM eikaiwa.")
)

var LogFile os.File

func init_log(log_path string) error {
	LogFile, err := os.OpenFile(log_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(LogFile)
	log.Printf("Initialize.")
	return nil
}

func main() {

	flag.Parse()

	// Read config.
	conf := newConf()
	conf.setConf(conf_path)

	// set log
	err := init_log(fmt.Sprintf("%s/%s", conf.LogDirPath, log_name))
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	prev_path := fmt.Sprintf("%s/%s", conf.LogDirPath, prev_name)

	// Read previous data
	log.Printf("Read Previous data")
	prev := ResultList{}
	ReadPrevSchedule(prev_path, &prev)
	prev.PrintAll()

	// Object that to notify to line
	l := newLine(conf.LineAccessToken)

	// Message list
	msg_list := []string{}

	teachers := []string{}
	for _, v := range conf.Teachers {
		teachers = append(teachers, v.Id)
	}

	parse := NewDMMParser(conf.CrawlDuration, *fileOpt, teachers, conf.StartTime, conf.EndTime)
	b, _ := parse.Parse()
	if !b {
		log.Printf("There are no results.")
	}

	// Create message list.
	for _, t := range parse.Results {
		msg_list = append(msg_list, fmt.Sprintf("\n%s", EmbedResult(&t)))
	}

	// get diff
	log.Printf("GetDiff")
	is_diff, _ := GetDiff(&(parse.Results), &prev)
	log.Println(is_diff)

	if *fileOpt == fileOptDefault && is_diff {
		tmp_str := ""
		for num, msg := range msg_list {
			tmp_str += msg
			if (num+1)%msg_per_notify == 0 {
				l.notify(tmp_str)
				// necessary duration to notify to line
				time.Sleep(notify_duration * time.Second)
				tmp_str = ""
			}
		}
		// notify remain messages
		l.notify(tmp_str)
	}

	log.Printf("Write data as previous data")
	WritePrevSchedule(prev_path, &(parse.Results))
}
