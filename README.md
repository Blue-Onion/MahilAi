# Medusa

Medusa is a multi-camera incident monitoring and event logging platform written in Go. It is designed as a pluggable event pipeline where varying detection engines (such as Python-based AI models) can be seamlessly integrated. The platform itself abstracts away the complexity of handling configuration, concurrency, logging, and querying, allowing you to focus on building robust detection capabilities.

## Features

- **Concurrent Multi-Camera Processing**: Run and monitor multiple camera streams simultaneously using native Go concurrency.
- **Flexible Stream Support**: Configure inputs from webcam devices, local video files, and RTSP streams.
- **Pluggable Detection Architecture**: Integrate any external detection software (interchangeable with your choice of detection models/engines) via standard output event emitting.
- **Structured JSON Logging**: All events are persisted as JSON Lines (JSONL), making them easy to parse and query.
- **Built-in CLI Tool**: Start the pipeline, view configurations, or query and filter historical logs by date and camera through an intuitive command-line interface.

## Configuration

Cameras are configured using a `config.yaml` file. Medusa supports webcam indices, local files, and RTSP streams.

Sample `config.yaml`:
```yaml
cameras:
  - name: cam1
    source: "0"
  - name: cam2
    source: "video2.mp4"
  - name: cam3
    source: "rtsp://admin:pass@192.168.1.20:554/stream"
```

## Installation

Medusa requires Go to be installed on your system.

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/medusa.git
   cd medusa
   ```
2. Build the executable:
   ```bash
   go build -o medusa main.go
   ```
   *(Adjust the build path based on your project's main package location)*

## CLI Usage

The Medusa CLI provides several commands to manage and query the platform:

- **Start the pipeline:**
  ```bash
  medusa start
  ```

- **View the current configuration:**
  ```bash
  medusa show-config
  ```

- **Query logged events by date and camera:**
  ```bash
  medusa show-record --date 2026-03-30 --cam cam1
  ```

## Log Storage Structure

Logs are organized hierarchically by date and camera name, ensuring easy file rotation and fast querying. Each active day generates a new folder, and each camera gets its own dedicated log file.

```
logs/
└── YYYY-MM-DD/
    ├── cam1.log
    ├── cam2.log
    └── cam3.log
```

## Log Format

All logs are stored in the JSON Lines (JSONL) format. Each line represents a single discrete event emitted by the underlying detection software.

Example `.log` file contents:
```json
{"camera":"cam2","time":"2026-03-30T08:09:14.84747219+05:30","event":"Human Detected","confidence":0.2504044771194458}
{"camera":"cam2","time":"2026-03-30T08:09:15.48450303+05:30","event":"Human Detected","confidence":0.31269264221191406}
{"camera":"cam2","time":"2026-03-30T08:09:15.5109272+05:30","event":"Human Detected","confidence":0.30127766728401184}
```

## Query Output Example

When using the `show-record` CLI command, Medusa parses the raw JSON logs and presents them in a readable tabular format.

Example output for `medusa show-record --date 2026-03-30 --cam cam2`:

```text
+--------+-----------------------------------+----------------+------------+
| CAMERA | TIME                              | EVENT          | CONFIDENCE |
+--------+-----------------------------------+----------------+------------+
| cam2   | 2026-03-30T08:09:14.84747219+05:30| Human Detected | 0.2504     |
| cam2   | 2026-03-30T08:09:15.48450303+05:30| Human Detected | 0.3127     |
| cam2   | 2026-03-30T08:09:15.5109272+05:30 | Human Detected | 0.3013     |
+--------+-----------------------------------+----------------+------------+
```

## Architecture

Medusa is fundamentally designed as an **event pipeline**. It cleanly separates the concerns of stream concurrency, application life-cycle, and data persistence from the heavy computational lifting of computer vision and machine learning models.

The pipeline architecture relies on **modular detection integration**. You can swap out the underlying detection software (e.g., Python scripts using YOLO, MediaPipe, etc.) without modifying the Go backend. As long as the detection module can process a stream and emit structured JSON to its standard output, Medusa can consume, distribute, log, and query those events effectively.

### Flow Diagram

```text
+-----------+       +-------------------+       +---------+       +----------+       +-----------+
|           |       |                   |       |         |       |          |       |           |
|  Camera   | ----> | Detection Software| ----> |  Event  | ----> |  Logger  | <---- | Query CLI |
| (Source)  |       | (Python/AI Model) |       | (JSON)  |       | (JSONL)  |       | (medusa)  |
|           |       |                   |       |         |       |          |       |           |
+-----------+       +-------------------+       +---------+       +----------+       +-----------+
```

1. **Camera**: Go routines provision and track individual video sources (webcam, file, RTSP).
2. **Detection Software**: An isolated external process is spawned per camera to analyze the visual input.
3. **Event**: The detection model outputs discrete incident data (threats/events) in real-time.
4. **Logger**: Medusa safely captures the parallel standard output strings, appending them to the appropriately dated and named log files.
5. **Query CLI**: Users natively interface with the aggregated log records via the Go CLI tool.

## Future Improvements

- Implementation of real-time alerts (e.g., Webhooks, Slack/Discord integration).
- Relational or NoSQL database support (PostgreSQL, MongoDB) as an alternative to local JSONL logs.
- Built-in web dashboard for live monitoring and interactive log visualization.
- Configuration hot-reloading to apply new detection rules without disrupting active streams.

## License

Licensed under the MIT License.
