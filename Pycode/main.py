import sys
import json
import time
from ultralytics import YOLO
import cv2

def run_detector(source, camera_name):
    # Load YOLO model
    model = YOLO("yolov8n.pt") 
    
    # Initialize video capture
    # Convert source to int if it's a single digit (for webcam index)
    if isinstance(source, str) and source.isdigit():
        source = int(source)
    
    cap = cv2.VideoCapture(source)
    if not cap.isOpened():
        print(f"Error: Could not open video source {source}", file=sys.stderr)
        return

    try:
        while True:
            ret, frame = cap.read()
            if not ret:
                break

            # Run detection
            results = model(frame, verbose=False)

            for r in results:
                for box in r.boxes:
                    cls = int(box.cls[0])
                    conf = float(box.conf[0])

                    if cls == 0:  # person class
                        event = {
                            "Camera": camera_name,
                            "Time": time.time(),
                            "Event": "Human Detected",
                            "Confidence": conf
                        }
                        # Print JSON to stdout for the Go handler to read
                        print(json.dumps(event), flush=True)

    finally:
        cap.release()

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python main.py <source> <camera_name>", file=sys.stderr)
        sys.exit(1)

    source_arg = sys.argv[1]
    name_arg = sys.argv[2]
    
    run_detector(source_arg, name_arg)