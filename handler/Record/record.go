package record

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	camera "github.com/Blue-Onion/MahilAi/handler/Camera"
)
func RecordEvent(event *camera.Event ){
	

}
func WriteEvent(event *camera.Event){
	file, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//hilihih
	defer file.Close()
	log.SetOutput(file)
	data,err:=json.Marshal(event)
	if err!=nil{
		fmt.Println(err.Error())
	}
	log.Println(string(data)+"\n")
}
