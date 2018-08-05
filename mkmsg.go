package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
)

type MultipleMessage struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	TeacherName string              `json:"teacher_name"`
	URL         string              `json:"url"`
	Schedules   map[string][]string `json:"schedules"` // {"07月29日(日)":["22:30,23:00,23:30,24:00,"],...},
}

func NewMultipleMessage() *MultipleMessage {
	return &MultipleMessage{}
}

func GetDiff(newer *MultipleMessage, older *MultipleMessage) (bool, *MultipleMessage) {

	ret := MultipleMessage{}
	exist_diff := false

	for _, message_newer := range newer.Messages {
		isNewTeacher := true

		for _, message_older := range older.Messages {

			// to find the same teacher name
			if message_newer.TeacherName == message_older.TeacherName {
				isNewTeacher = false

				for sck_newer, scv_newer := range message_newer.Schedules {
					// to find the same date
					if scv_older, ok := message_older.Schedules[sck_newer]; !ok { // Not found the same date, so this message is to be notified
						exist_diff = true
						ret.Stock(message_newer)

					} else { // found the same date

						// check all lesson duration
						for _, lesson_newer := range scv_newer {
							isFounded := false
							for _, lesson_older := range scv_older {
								if lesson_newer == lesson_older {
									// Found the same duration means
									// the lesson duration was already notified, so go next loop.
									isFounded = true
									break
								}
							}

							// Not found the same duration,
							// this message will be notified.
							if !isFounded {
								exist_diff = true
								ret.Stock(message_newer)
								break
							}
						}
					}
				}
			}

		}
		// If there are new teacher, stock and notify
		// This block will execute after adding new teacher to setting.yaml.
		if isNewTeacher == true {
			exist_diff = true
			ret.Stock(message_newer)
		}
	}
	return exist_diff, &ret
}

func (mmobj *MultipleMessage) PrintAll() {
	for n, m := range mmobj.Messages {
		fmt.Printf("MessageNumber:%d\n", n)
		m.Print()
	}
}

func (mmobj *MultipleMessage) Stock(mobj Message) {
	mmobj.Messages = append(mmobj.Messages, mobj)
}

func (mmobj *MultipleMessage) ToJsonString() (string, error) {
	jsonBytes, err := json.Marshal(mmobj)
	if err != nil {
		log.Printf("Error: json marshal")
		return "", err
	}
	return string(jsonBytes), nil
}

func (mmobj *MultipleMessage) ParseJson(jstr string) error {
	jsonBytes := ([]byte)(jstr)
	if err := json.Unmarshal(jsonBytes, mmobj); err != nil {
		log.Printf("Error: json Unmarshal")
		return err
	}
	return nil
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

func (mobj *Message) ToJsonString() (string, error) {
	jsonBytes, err := json.Marshal(mobj)
	if err != nil {
		log.Printf("Error: json marshal")
		return "", err
	}
	return string(jsonBytes), nil
}
