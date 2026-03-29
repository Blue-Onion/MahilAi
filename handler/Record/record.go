package record

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/Blue-Onion/MahilAi/handler/config"
)
type Records struct{
	Camera     string
	Time       string
	Event      string
	Confidence float64
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
func ReadEvent(date string,cam string){
	if cam==""{
		ceadCameraAllEvent(cam)
	}else if date==""{
		ceadDateAllEvent(date)
	}else{
		camDateEvent(date,cam)
	}

}
func ceadCameraAllEvent(cam string){
}
func ceadDateAllEvent(date string){

}
func camDateEvent(date string,cam string){

}