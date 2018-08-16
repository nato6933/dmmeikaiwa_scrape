package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"text/template"
)

func ReadPrevSchedule(prev_path string, prev *ResultList) error {
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
	if err := dec.Decode(prev); err != nil {
		log.Fatal("decode error:", err)
		return err
	}

	return err
}

func WritePrevSchedule(prev_path string, rl *ResultList) error {
	PrevFile, err := os.OpenFile(prev_path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer PrevFile.Close()

	enc := gob.NewEncoder(PrevFile)
	if err := enc.Encode(rl); err != nil {
		log.Fatal(err)
		return err
	}

	return nil

}

func EmbedResult(t *Teacher) string {

	tmpl := template.Must(template.ParseFiles(tmpl_path))

	tmpl_map := map[string]string{}

	tmpl_map["TeacherName"] = t.TeacherName
	tmpl_map["URL"] = t.URL
	tmpl_map["Schedules"] = ""

	for k, v := range t.Schedules {
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

func GetDiff(newer *ResultList, older *ResultList) (bool, *ResultList) {

	ret := ResultList{}
	exist_diff := false

	for _, res_newer := range *newer {
		isNewTeacher := true

		for _, res_older := range *older {

			// to find the same teacher name
			if res_newer.TeacherName == res_older.TeacherName {
				isNewTeacher = false

				for sck_newer, scv_newer := range res_newer.Schedules {
					// to find the same date
					if scv_older, ok := res_older.Schedules[sck_newer]; !ok { // Not found the same date, so this message is to be notified
						exist_diff = true
						ret.Stock(res_newer)

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
								ret.Stock(res_newer)
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
			ret.Stock(res_newer)
		}
	}
	return exist_diff, &ret
}
