package main

import (
	"encoding/gob"
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
	log_name       = "parse.log"
	prev_name      = "previous_schedule.gob"
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

func read_prev_schedule(prev_path string, mm_prev *MultipleMessage) error {
	var err error = nil

	_, err = os.Stat(prev_path)
	if os.IsNotExist(err) {
		log.Printf("Not Exist %s.", prev_path)
		return err
	}

	f, err := os.Open(prev_path)
	if err != nil {
		log.Printf("Failed to open :%s", prev_path)
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(mm_prev); err != nil {
		log.Fatal("decode error:", err)
		return err
	}

	return err
}

func write_prev_schedule(prev_path string, mm *MultipleMessage) error {
	PrevFile, err := os.OpenFile(prev_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer PrevFile.Close()

	enc := gob.NewEncoder(PrevFile)
	if err := enc.Encode(mm); err != nil {
		log.Fatal(err)
		return err
	}

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
	mm_prev := NewMultipleMessage()
	read_prev_schedule(prev_path, mm_prev)

	// Object that to notify to line
	l := newLine(conf.LineAccessToken)

	// Message list
	msg_list := []string{}

	tmp_msg := Message{
		TeacherName: "",
		URL:         "",
		Schedules:   map[string][]string{},
	}

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

	okReserve := false

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	// Instantiate default collector
	c := colly.NewCollector()

	// Read from file instead of to access DMM eikaiwa
	if *fileOpt != fileOptDefault {
		c.WithTransport(t)
	}

	tmp_teacher := ""

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
					//log.Printf("date:%s", elem.Text)
					tmp_date = elem.Text
					tmp_msg.Schedules[elem.Text] = []string{}
				} else if elem.Text == StrCanReserve {
					okReserve = true
					//log.Printf("Start time:%s", lesson_time_array[a-1])
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

	mm := NewMultipleMessage()

	if *fileOpt == fileOptDefault {

		for _, v := range conf.Teachers {
			tmp_msg = Message{
				TeacherName: "",
				URL:         "",
				Schedules:   map[string][]string{},
			}

			okReserve = false

			tmp_msg.URL = fmt.Sprintf(dmm_url, v.Id)

			// starting parse
			c.Visit(tmp_msg.URL)

			if okReserve {
				//log.Printf("Embed:\n%s", tmp_msg.Embed(tmpl_path))
				msg_list = append(msg_list, fmt.Sprintf("\n%s", tmp_msg.Embed(tmpl_path)))
				mm.Stock(tmp_msg)
			}

			time.Sleep(time.Duration(conf.CrawlDuration) * time.Second) // necessary duration that is to automated parse.
		}

	} else { // for debug
		c.Visit("file://" + *fileOpt)
		jstr, _ := tmp_msg.ToJsonString()
		log.Printf(jstr)

		mm.Stock(tmp_msg)
	}

	// get diff
	log.Printf("GetDiff")
	is_diff, _ := GetDiff(mm, mm_prev)

	if *fileOpt == fileOptDefault && is_diff {
		tmp_str := ""
		for num, msg := range msg_list {
			tmp_str += msg
			if (num+1)%4 == 0 {
				l.notify(tmp_str)
				// necessary duration to notify to line
				time.Sleep(1 * time.Second)
				tmp_str = ""
			}
		}
		// notify remain messages
		l.notify(tmp_str)
	}

	log.Printf("Write data as previous data")
	write_prev_schedule(prev_path, mm)
}
