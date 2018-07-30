package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type Message struct {
	TeacherName string
	URL         string
	Schedules   map[string][]string // {"07月29日(日)":["22:30,23:00,23:30,24:00,"],...},
}

func (mobj *Message) Print() {
	fmt.Printf("Message.TeacherName:%s\n", mobj.TeacherName)
	fmt.Printf("Message.URL:%s\n", mobj.URL)
	for k, v := range mobj.Schedules {
		fmt.Printf("Message.Schedules key(date) : %s\n", k)
		for _, st := range v {
			fmt.Printf("  %s", st)
		}
		fmt.Printf("\n")
	}
}

func (mobj *Message) Embed(tmp_path string) string {

	tmpl := template.Must(template.ParseFiles(tmpl_path))

	tmpl_map := map[string]string{}

	tmpl_map["TeacherName"] = mobj.TeacherName
	tmpl_map["URL"] = mobj.URL
	tmpl_map["Schedules"] = ""

	for k, v := range mobj.Schedules {
		if len(v) != 0 {
			tmpl_map["Schedules"] += fmt.Sprintf("%s\n", k)
			for _, st := range v {
				tmpl_map["Schedules"] += fmt.Sprintf("  %s", st)
			}
			tmpl_map["Schedules"] += fmt.Sprintf("\n")
		}
	}

	var res bytes.Buffer

	err := tmpl.Execute(&res, tmpl_map)

	if err != nil {
		log.Printf("Failed to Execute template")
		return ""
	}
	return res.String()
}
