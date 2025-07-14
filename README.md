# 🌀 Go Load Balancer

A production-ready, highly configurable **Load Balancer written in Golang**, supporting **eight popular load balancing algorithms**, live configuration reload, and health check intervals — all customizable via **YAML or JSON** files.

---

## 🚀 Features

- ✅ **8 Load Balancing Algorithms**:
  - Round Robin
  - Random
  - IP Hash
  - Least Connections
  - Least Response Time
  - Weighted Round Robin
  - Weighted Least Connections
  - Weighted Response Time
- ✅ **Dynamic Configuration** via `config.yaml` or `config.json`
- ✅ **Live Hot Reload** using file watchers (`fsnotify`)
- ✅ **Health Check Interval** support
- ✅ **Concurrency-safe and extensible architecture**
- ✅ Built with idiomatic Go using goroutines, channels, and sync primitives

---

## 📁 Project Structure

```

.
├── main.go
├── config.go
├── watcher.go
├── loadbalancer/
│   ├── RoundRobin.go
│   ├── Random.go
│   ├── IpHash.go
│   ├── LeastConnection.go
│   ├── LeastResponseTime.go
│   ├── WeightedRoundRobin.go
│   ├── WeightedLeastConnection.go
│   ├── WeightedResponseTime.go

````

---

## ⚙️ Configuration

This project uses a `config.yaml` or `config.json` file to set the port, strategy, backend servers, and health check interval.

### ✅ Example `config.yaml`

```yaml
port: "8080"
strategy: "WeightedRoundRobin"
healthCheckInterval: 10
servers:
  - "http://localhost:8081"
  - "http://localhost:8082"
  - "http://localhost:8083"
````

### ✅ Example `config.json`

```json
{
  "port": "8080",
  "strategy": "LeastConnection",
  "healthCheckInterval": 10,
  "servers": [
    "http://localhost:8081",
    "http://localhost:8082",
    "http://localhost:8083"
  ]
}
```

### 🧠 Config Struct

```go
type Config struct {
    Port                string   `json:"port" yaml:"port"`
    Strategy            string   `json:"strategy" yaml:"strategy"`
    HealthCheckInterval int      `json:"healthCheckInterval" yaml:"healthCheckInterval"`
    Servers             []string `json:"servers" yaml:"servers"`
    mu                  sync.RWMutex
}
```

---

## 🔁 Live Config Reload

This project uses `fsnotify` to **monitor file changes** in real time.
When `config.yaml` or `config.json` is updated, the application **reloads the new configuration without restarting**.

---

## 🧪 Algorithms Implemented

You can choose any of the following in your `strategy` config field:

* `"RoundRobin"`
* `"Random"`
* `"IpHash"`
* `"LeastConnection"`
* `"LeastResponseTime"`
* `"WeightedRoundRobin"`
* `"WeightedLeastConnection"`
* `"WeightedResponseTime"`

Each algorithm is modular, defined in its own file, and easy to extend.

---

## 💻 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Garv2003/go-load-balancer.git
cd go-load-balancer
```

### 2. Create Configuration

Add either a `config.yaml` or `config.json` file in the root directory.
Use one of the examples provided above.

### 3. Run the Load Balancer

```bash
go run main.go
```

---

## 🛠 Built With

* **Language:** Go
* **Concurrency:** Goroutines, Channels, RWMutex
* **File Watching:** fsnotify
* **Protocols:** HTTP
* **Data Formats:** YAML, JSON

---

## 🚧 Roadmap

* [ ] Add backend health checks & failover
* [ ] Enable TLS/HTTPS support
* [ ] Web dashboard for metrics and control
* [ ] Graceful shutdown & restart
* [ ] Structured logging and monitoring
