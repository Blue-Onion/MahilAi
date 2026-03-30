package record

import (
	"os"
	"testing"

	"github.com/Blue-Onion/MahilAi/handler/config"
)

func TestRecordLifecycle(t *testing.T) {
	// Setup clean environment in temp directory to prevent polluting logs folder
	tempDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// 1. Write an event
	evt := &config.Event{
		Camera:     "cam1",
		Time:       1672531200.0, // Arbitrary timestamp
		Event:      "person",
		Confidence: 0.99,
	}
	WriteEvent(evt)

	// Since Time uses local timezone in WriteEvent, parse the created folder name directly
	entries, err := os.ReadDir("logs")
	if err != nil || len(entries) == 0 {
		t.Fatalf("WriteEvent failed to create logs directory and date folder")
	}
	dateFolder := entries[0].Name()

	// 2. Verify log file was successfully created
	logFile := "logs/" + dateFolder + "/cam1.log"
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("Expected log file to be created: %v", err)
	}

	// 3. Test ReadEvent handles reading correctly
	records, err := ReadEvent(dateFolder, "cam1")
	if err != nil {
		t.Fatalf("ReadEvent failed: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("No records found after reading")
	}

	if records[0].Camera != "cam1" {
		t.Errorf("Expected camera cam1, got %s", records[0].Camera)
	}
	if records[0].Event != "person" {
		t.Errorf("Expected event person, got %s", records[0].Event)
	}
	if records[0].Confidence != 0.99 {
		t.Errorf("Expected confidence 0.99, got %v", records[0].Confidence)
	}

	// 4. Test missing files properly error out
	_, err = ReadEvent("nonexistent_date", "cam1")
	if err == nil {
		t.Errorf("Expected error reading nonexistent dir, got nil")
	}
}
