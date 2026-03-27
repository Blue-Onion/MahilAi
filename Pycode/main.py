from ultralytics import YOLO
import cv2

def doSomething():
    print("Human Detected")

model = YOLO("yolov8n.pt") 
cap = cv2.VideoCapture(0)

while True:
    ret, frame = cap.read()
    if not ret:
        break

    res = model(frame, verbose=False)

    for r in res:
        for box in r.boxes:
            cls = int(box.cls[0])

            if cls == 0:  # person class
                doSomething()

                x1, y1, x2, y2 = map(int, box.xyxy[0])
                cv2.rectangle(frame, (x1, y1), (x2, y2), (0, 255, 0), 2)

                cv2.putText(frame, "Human", (x1, y1 - 10),
                            cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)

    cv2.imshow("Human Detection", frame)

    # ESC key to exit
    if cv2.waitKey(1) == 27:
        break

cap.release()
cv2.destroyAllWindows()