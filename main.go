package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

const (
	dmm_url        = "https://eikaiwa.dmm.com//teacher/index/%s/"
	conf_path      = "./conf/setting.yaml"
	tmpl_path      = "./conf/notify.tmpl"
	log_path       = "./log/parse.log"
	fileOptDefault = "notset"
	StrCanReserve  = "予約可"
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

	err := init_log(log_path)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	// Read config.
	conf := newConf()
	conf.setConf(conf_path)
	//conf.PrintConf()

	l := newLine(conf.LineAccessToken)

	msg := "\n"
	tmp_msg := Message{
		TeacherName: "",
		URL:         "",
		Schedules:   map[string][]string{},
	}

	tmp_teacher := ""

	lesson_time_array := []string{
		"02:00", "02:30", "03:00", "03:30", "04:00", "04:30",
		"05:00", "05:30", "06:00", "06:30", "07:00", "07:30",
		"08:00", "08:30", "09:00", "09:30", "10:00", "10:30",
		"11:00", "11:30", "12:00", "12:30", "13:00", "13:30",
		"14:00", "14:30", "15:00", "15:30", "16:00", "16:30",
		"17:00", "17:30", "18:00", "18:30", "19:00", "19:30",
		"20:00", "20:30", "21:00", "21:30", "22:00", "22:30",
		"23:00", "23:30", "24:00", "24:30", "25:00", "25:30",
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	// Instantiate default collector
	c := colly.NewCollector()

	// Read from file instead of to access DMM eikaiwa
	if *fileOpt != fileOptDefault {
		c.WithTransport(t)
	}

	// On every a element which has div and class=area-detail attribute call callback
	c.OnHTML("div.area-detail", func(e *colly.HTMLElement) {

		// get teacher's name
		tmp_teacher = ""
		tmp_msg.TeacherName = e.DOM.Find("h1").Text()
		log.Printf("TeacherName:%s", tmp_msg.TeacherName)
	})

	// On every a element which has ul and class=oneday attribute call callback
	c.OnHTML("ul.oneday", func(e *colly.HTMLElement) {

		tmp_date := ""
		e.ForEach("li", func(a int, elem *colly.HTMLElement) {

			if elem.Text != "" {
				if elem.DOM.HasClass("date") {
					log.Printf("date:%s", elem.Text)
					tmp_date = elem.Text
					tmp_msg.Schedules[elem.Text] = []string{}
				} else if elem.Text == StrCanReserve {
					log.Printf("Start time:%s", lesson_time_array[a-1])
					if tmp_date != "" {
						tmp_msg.Schedules[tmp_date] = append(tmp_msg.Schedules[tmp_date], lesson_time_array[a-1])
					}
				}
			}
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting:%s", r.URL.String())
	})

	if *fileOpt == fileOptDefault {
		for _, v := range conf.Teachers {
			tmp_msg = Message{
				TeacherName: "",
				URL:         "",
				Schedules:   map[string][]string{},
			}

			tmp_msg.URL = fmt.Sprintf(dmm_url, v.Id)

			c.Visit(tmp_msg.URL)
			log.Printf("Embed:\n%s", tmp_msg.Embed(tmpl_path))
			msg += tmp_msg.Embed(tmpl_path)

			time.Sleep(1 * time.Second) // duration to automated parse.
		}
	} else {
		c.Visit("file://" + *fileOpt)
		msg += tmp_msg.Embed(tmpl_path)
	}

	if *fileOpt == fileOptDefault {
		l.notify(msg)
	}
}
