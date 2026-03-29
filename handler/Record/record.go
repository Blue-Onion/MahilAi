package record

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/Blue-Onion/MahilAi/handler/config"
)

type Records struct {
	Camera     string  `json:"camera"`
	Time       string  `json:"time"`
	Event      string  `json:"event"`
	Confidence float64 `json:"confidence"`
}

func WriteEvent(event *config.Event) {
	name := event.Camera

	sec := int64(event.Time)
	nsec := int64((event.Time - float64(sec)) * 1e9)

	parsedTime := time.Unix(sec, nsec)

	today := parsedTime.Format("2006-01-02")

	folderPath := fmt.Sprintf("logs/%s", today)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := fmt.Sprintf("%s/%s.log", folderPath, name)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	record := Records{
		Camera:     event.Camera,
		Time:       parsedTime.Format(time.RFC3339Nano), // store string
		Confidence: event.Confidence,
		Event:      event.Event,
	}

	data, err := json.Marshal(record)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	file.WriteString(string(data) + "\n")
}
func ReadEvent(date string, cam string) ([]Records, error) {
	if date != "" && cam != "" {
		return camDateEvent(date, cam)
	}
	if date == "" && cam != "" {
		return readCameraAllEvent(cam)
	}
	if date != "" && cam == "" {
		return readDateAllEvent(date)
	}
	return nil, fmt.Errorf("date and cam both empty")
}
func readCameraAllEvent(cam string) ([]Records, error) {
	var res []Records
	dates,err:=os.ReadDir("logs")
	if err!=nil{
		return nil,err
	}
	for _,date:=range dates{
		if !date.IsDir(){
			continue
		}
		path:=fmt.Sprintf("logs/%s/%s",date.Name(),cam)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		event,err:=ReadEvents(path)
		if err!=nil{
			return nil,err
		}
		res = append(res, event...)
	}
	return res,nil
	
}
func readDateAllEvent(date string) ([]Records, error) {
	path:=fmt.Sprintf("logs/%s",date)
	files,err:=os.ReadDir(path)
	if err!=nil{
		return nil,err
	}
	var res []Records
	for _,f:=range files{
		if f.IsDir(){
			continue
		}
		filePath:=fmt.Sprintf("%s/%s.log",path,f.Name())
		event,err:=ReadEvents(filePath)
		if err!=nil{
			continue
		}
		res=append(res, event...)

	}
	return res,nil

}
func camDateEvent(date string, cam string) ([]Records, error) {
	path:=fmt.Sprintf("logs/%s/%s.log",date,cam)
	res,err:=ReadEvents(path)
	if err!=nil{
		return nil,err
	}
	return res,nil
}
func ReadEvents(path string) ([]Records, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Records
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		var rec Records
		err := json.Unmarshal([]byte(line), &rec)
		if err != nil {
			continue
		}

		events = append(events, rec)
	}

	return events, scanner.Err()
}