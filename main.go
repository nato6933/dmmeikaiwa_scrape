package main

import (
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

func main() {

	//LessonStartingTime := map[string]string{
	//	"t-02-00": "02:00", "t-02-30": "02:30", "t-03-00": "03:00", "t-03-30": "03:30", "t-04-00": "04:00", "t-04-30": "04:30",
	//	"t-05-00": "05:00", "t-05-30": "05:30", "t-06-00": "06:00", "t-06-30": "06:30", "t-07-00": "07:00", "t-07-30": "07:30",
	//	"t-08-00": "08:00", "t-08-30": "08:30", "t-09-00": "09:00", "t-09-30": "09:30", "t-10-00": "10:00", "t-10-30": "10:30",
	//	"t-11-00": "11:00", "t-11-30": "11:30", "t-12-00": "12:00", "t-12-30": "12:30", "t-13-00": "13:00", "t-13-30": "13:30",
	//	"t-14-00": "14:00", "t-14-30": "14:30", "t-15-00": "15:00", "t-15-30": "15:30", "t-16-00": "16:00", "t-16-30": "16:30",
	//	"t-17-00": "17:00", "t-17-30": "17:30", "t-18-00": "18:00", "t-18-30": "18:30", "t-19-00": "19:00", "t-19-30": "19:30",
	//	"t-20-00": "20:00", "t-20-30": "20:30", "t-21-00": "21:00", "t-21-30": "21:30", "t-22-00": "22:00", "t-22-30": "22:30",
	//	"t-23-00": "23:00", "t-23-30": "23:30", "t-24-00": "24:00", "t-24-30": "24:30", "t-25-00": "25:00", "t-25-30": "25:30",
	//}

	accessToken := "set your token"

	l := newLine(accessToken)

	msg := "\n"

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
	c := colly.NewCollector(
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	c.WithTransport(t)

	// On every a element which has div and class=area-detail attribute call callback
	c.OnHTML("div.area-detail", func(e *colly.HTMLElement) {
		// get teacher's name
		t_name := ""
		t_name = e.DOM.Find("h1").Text()
		fmt.Println(t_name)
		msg += fmt.Sprintf("Teacher: %s\n\n", t_name)
	})

	// On every a element which has ul and class=oneday attribute call callback
	c.OnHTML("ul.oneday", func(e *colly.HTMLElement) {

		// get schedule
		res_date := ""
		e.ForEach("li", func(a int, elem *colly.HTMLElement) {
			if elem.Text != "" {
				if elem.DOM.HasClass("date") {
					fmt.Println("date:" + elem.Text)
					res_date = elem.Text
					msg += fmt.Sprintf("%s\n", res_date)

				} else if elem.Text == "予約可" {
					fmt.Println(a - 1)
					fmt.Println(lesson_time_array[a-1])
					msg += fmt.Sprintf("  %s\n", lesson_time_array[a-1])
				}
			}
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("file://" + "/tmp/test.html")
	l.notify(msg)
}
