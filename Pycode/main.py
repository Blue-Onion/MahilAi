from ultralytics import YOLO
import cv2

def doSomething():
    print("hello")
    
model = YOLO()   
cap=cv2.VideoCapture(0)

while True:
    ret,frame=cap.read()
    res=model(frame)
    for r in res:
        boxes = r.boxes
        for box in boxes:
            cls = int(box.cls[0])
            if cls == 0:   
                doSomething()
                x1,y1,x2,y2=map(int,box.xyxy[0])
                cv2.rectangle(frame,(x1,y1),(x2,y2),(0,255,0),2)
                cv2.putText(frame, "Human", (x1,y1-10),
                            cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0,255,0), 2)
    cv2.imshow("Human Detection", frame)

    if cv2.waitKey(1) == 27:
        break
    


res[0].show()



