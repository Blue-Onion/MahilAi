package camera

import (
	"testing"
	"time"

	"github.com/Blue-Onion/MahilAi/handler/config"
)

func TestMerge(t *testing.T) {
	ch1 := make(chan config.Event)
	ch2 := make(chan config.Event)

	merged := merge(ch1, ch2)

	// Produce events on a separate goroutine
	go func() {
		ch1 <- config.Event{Camera: "cam1", Event: "ev1"}
		ch2 <- config.Event{Camera: "cam2", Event: "ev2"}
		ch1 <- config.Event{Camera: "cam1", Event: "ev3"}
		close(ch1)
		close(ch2)
	}()

	var results []config.Event

	// Create a timeout to avoid hanging the test indefinitely if merge fails
	done := make(chan bool)
	go func() {
		for e := range merged {
			results = append(results, e)
		}
		done <- true
	}()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Merge function hung, likely not closing the output channel properly")
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(results))
	}

	foundCam1 := 0
	foundCam2 := 0
	for _, e := range results {
		if e.Camera == "cam1" {
			foundCam1++
		}
		if e.Camera == "cam2" {
			foundCam2++
		}
	}

	if foundCam1 != 2 {
		t.Errorf("Expected 2 events from cam1, got %d", foundCam1)
	}
	if foundCam2 != 1 {
		t.Errorf("Expected 1 event from cam2, got %d", foundCam2)
	}
}
