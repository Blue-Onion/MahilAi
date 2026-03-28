package camera

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"

	"github.com/Blue-Onion/MahilAi/handler/config"
)

type Event struct {
	Camera     string
	Time       float64
	Event      string
	Confidence float64
}

func StartCameraWork(cfg *config.Config){

	var channels []<-chan Event

	for _,val:=range cfg.Cameras {
		ch, err := streamEvent(val)
		if err != nil {
			panic(err)
		}
		channels = append(channels, ch)
	}

	merged := merge(channels...)

	for e := range merged {
		fmt.Printf("Camera: %s | Event: %s | Confidence: %.2f\n",
			e.Camera, e.Event, e.Confidence)
	}
}

func streamEvent(camera config.Camera) (<-chan Event, error) {
	cmd := exec.Command("python3", "Pycode/main.py", fmt.Sprintf("%v", camera.Source), camera.Name)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	ch := make(chan Event)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer close(ch)
		defer cmd.Wait()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var e Event
			err := json.Unmarshal(scanner.Bytes(), &e)
			if err == nil {
				ch <- e
			}
		}
	}()

	return ch, nil
}

func merge(channels ...<-chan Event) <-chan Event {
	out := make(chan Event)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Event) {
			defer wg.Done()
			for e := range c {
				out <- e
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}