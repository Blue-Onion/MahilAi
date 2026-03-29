package camera

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"

	record "github.com/Blue-Onion/MahilAi/handler/Record"
	"github.com/Blue-Onion/MahilAi/handler/config"
)



func StartCameraWork(cfg *config.Config) {

	var channels []<-chan config.Event

	for _, val := range cfg.Cameras {
		ch, err := streamEvent(val)
		if err != nil {
			panic(err)
		}
		channels = append(channels, ch)
	}

	merged := merge(channels...)
	for e := range merged {
		record.WriteEvent(&e)

	}
}

func streamEvent(camera config.Camera) (<-chan config.Event, error) {
	cmd := exec.Command("python3", "Pycode/main.py", fmt.Sprintf("%v", camera.Source), camera.Name)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	ch := make(chan config.Event)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer close(ch)
		defer cmd.Wait()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var e config.Event
			err := json.Unmarshal(scanner.Bytes(), &e)
			if err == nil {
				ch <- e
			}
		}
	}()

	return ch, nil
}

func merge(channels ...<-chan config.Event) <-chan config.Event {
	out := make(chan config.Event)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan config.Event) {
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
