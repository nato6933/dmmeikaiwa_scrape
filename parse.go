package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"time"
)

const (
	dmm_url = "https://eikaiwa.dmm.com//teacher/index/%s/"
)

var lesson_time_array = []string{
	"02:00", "02:30", "03:00", "03:30", "04:00", "04:30",
	"05:00", "05:30", "06:00", "06:30", "07:00", "07:30",
	"08:00", "08:30", "09:00", "09:30", "10:00", "10:30",
	"11:00", "11:30", "12:00", "12:30", "13:00", "13:30",
	"14:00", "14:30", "15:00", "15:30", "16:00", "16:30",
	"17:00", "17:30", "18:00", "18:30", "19:00", "19:30",
	"20:00", "20:30", "21:00", "21:30", "22:00", "22:30",
	"23:00", "23:30", "24:00", "24:30", "25:00", "25:30",
}

type DMMParser struct {
	CrawlDuration int
	FilePath      string
	TeacherList   []string
	Results       ResultList
	Starting      string
	Ending        string
	//Results       []Teacher
}

type ResultList []Teacher

type Teacher struct {
	TeacherName string
	URL         string
	Schedules   map[string][]string // {"07月29日(日)":["22:30,23:00,23:30,24:00,"],...},
}

//
// Args:
//   fpath: Path to HTML file for debug.
//   teachers: Teacher's ID list.
//
// Return:
//   bool :
//     true: Even one of lessons that can be reserved in parsed data.
//     false: There is no lesson that can be reserved.
//   error :
//
func NewDMMParser(cDuration int, fpath string, teachers []string, starting string, ending string) *DMMParser {
	return &DMMParser{
		CrawlDuration: cDuration,
		TeacherList:   teachers,
		FilePath:      fpath,
		Starting:      starting,
		Ending:        ending,
	}
}

// Stock Teacher obj to Results
func (dobj *DMMParser) Stock(t Teacher) {
	dobj.Results = append(dobj.Results, t)
}

// Stock Teacher obj to ResultList
func (r *ResultList) Stock(t Teacher) {
	*r = append(*r, t)
}

//
// Return:
//   bool :
//     true: Even one of lessons that can be reserved in parsed data.
//     false: There is no lesson that can be reserved.
//   error :
//
func (dobj *DMMParser) Parse() (bool, error) {

	tmp_teacher := Teacher{}
	okReserve := false

	t := &http.Transport{}

	// Instantiate default collector
	c := colly.NewCollector()

	// Read from file instead of to access DMM eikaiwa
	if dobj.FilePath != "" {
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c.WithTransport(t)
	}

	// On every a element which has div and class=area-detail attribute call callback
	c.OnHTML("div.area-detail", func(e *colly.HTMLElement) {

		// get teacher's name
		tmp_teacher.TeacherName = e.DOM.Find("h1").Text()
		log.Printf("TeacherName:%s", tmp_teacher.TeacherName)
	})

	// On every a element which has ul and class=oneday attribute call callback
	c.OnHTML("ul.oneday", func(e *colly.HTMLElement) {

		tmp_date := ""
		IsIncorporate := false

		e.ForEach("li", func(a int, elem *colly.HTMLElement) {

			if a != 0 && lesson_time_array[a-1] == dobj.Starting {
				IsIncorporate = true
			} else if a != 0 && lesson_time_array[a-1] == dobj.Ending {
				IsIncorporate = false
			}

			if elem.Text != "" {
				if elem.DOM.HasClass("date") {
					tmp_date = elem.Text
					tmp_teacher.Schedules[elem.Text] = []string{}
				} else if elem.Text == StrCanReserve && tmp_date != "" && IsIncorporate { // reserve
					okReserve = true
					tmp_teacher.Schedules[tmp_date] = append(tmp_teacher.Schedules[tmp_date], lesson_time_array[a-1])
				}
			}
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting:%s", r.URL.String())
	})

	// mainLoop
	if dobj.FilePath == "" {

		for _, v := range dobj.TeacherList {
			tmp_teacher = Teacher{
				TeacherName: "",
				URL:         "",
				Schedules:   map[string][]string{},
			}

			okReserve = false

			tmp_teacher.URL = fmt.Sprintf(dmm_url, v)

			// starting parse
			c.Visit(tmp_teacher.URL)

			if okReserve {
				dobj.Stock(tmp_teacher)
			}

			time.Sleep(time.Duration(dobj.CrawlDuration) * time.Second) // necessary duration that is to automated parse.
		}

	} else { // for debug
		c.Visit("file://" + dobj.FilePath)
		dobj.Stock(tmp_teacher)
	}

	if len(dobj.Results) != 0 {
		return true, nil
	}
	return false, nil
}

func (dobj *DMMParser) PrintAll() {
	for n, m := range dobj.Results {
		fmt.Printf("ResultsNumber:%d\n", n)
		m.Print()
	}
}

func (r *ResultList) PrintAll() {
	for n, m := range *r {
		fmt.Printf("ResultsNumber:%d\n", n)
		m.Print()
	}
}

func (t *Teacher) Print() {
	fmt.Printf("Teacher.TeacherName:%s\n", t.TeacherName)
	fmt.Printf("Teacher.URL:%s\n", t.URL)
	for k, v := range t.Schedules {
		fmt.Printf("Teacher.Schedules key(date) : %s\n", k)
		for _, st := range v {
			fmt.Printf("  %s", st)
		}
		fmt.Printf("\n")
	}
}
