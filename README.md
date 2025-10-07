# 🧩 TinyDB

**TinyDB** is a lightweight, Redis-like in-memory database built in **Go (Golang)**.  
It provides a simple TCP server that understands a subset of Redis-style commands (`SET`, `GET`, `HSET`, `HGET`, etc.) with **Append-Only File (AOF)** persistence and peer-to-peer extensibility for future distributed setups.

---

## 🚀 Features

- ⚡ **In-Memory Database** with Redis-like command syntax.  
- 💾 **AOF Persistence** — data is logged to disk to survive restarts.  
- 🧠 **Command Handlers** for basic key/value and hash operations.  
- 🧩 **Pluggable Architecture** — modular packages for server, memory, AOF, peers, and protocol.  
- 🔌 **TCP Protocol Server** that can be interacted with via `redis-cli` or `netcat`.  
- 📡 **Peer System (WIP)** — planned cluster support for multiple nodes.

---
<!-- 
## 🏗️ Project Structure

tinydb/
├── main.go # Entry point – starts the TCP server
├── models.go # Basic shared data models
├── server/ # Core TCP server and event loop
│ ├── server.go
│ ├── models.go
├── memory/ # In-memory data structures and command handlers
│ ├── memory.go
│ ├── models.go
├── aof/ # Append-only file (persistence)
│ ├── aof.go
│ ├── models.go
├── protocol/ # Custom RESP-like protocol parser/writer
│ ├── reader.go
│ ├── writer.go
│ ├── models.go
│ ├── custom_errors.go
├── peer/ # Peer-to-peer node communication (for clustering)
│ ├── peer.go
│ ├── models.go
├── Makefile # Simple build and run commands
├── .env # Environment configuration (e.g., PORT)
├── go.mod / go.sum # Go module files



--- -->

## ⚙️ Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/muhammadolammi/tinydb.git
cd tinydb/tinydb
go mod tidy 
```

Build the binary:
```bash
make build
```
Or run directly:
```bash
go run main.go
```

By default, TinyDB listens on port 6379 (configurable via .env).

---

## 🧰 Usage

Start the Server
```bash
go run main.go
```

You should see:
```psql
2025/10/07 22:00:00 starting server
INFO server running listenAddr=:6379
```

Connect via redis-cli or netcat

TinyDB uses a Redis-like RESP protocol, so you can use redis-cli:
```bash
redis-cli -p 6379
```
Or netcat:
```bash
nc localhost 6379
```

Supported Commands
| Command                | Description             | Example                  |
| ---------------------- | ----------------------- | ------------------------ |
| `SET key value`        | Store a string value    | `SET name "TinyDB"`      |
| `GET key`              | Retrieve a string value | `GET name`               |
| `HSET key field value` | Set a field in a hash   | `HSET user name Olamide` |
| `HGET key field`       | Get a field from a hash | `HGET user name`         |
| `DEL key`              | Delete a key            | `DEL name`               |
| `EXISTS key`           | Check if key exists     | `EXISTS name`            |

    All commands are case-insensitive.
---

## 🗂️ Persistence (AOF)
TinyDB logs all mutating commands (SET, HSET, etc.) to an Append-Only File stored locally.
When the server restarts, it automatically replays the AOF log to rebuild the in-memory state.

---

## 🧩 Architecture Overview
TinyDB consists of several decoupled components:

| Package    | Responsibility                                   |
| ---------- | ------------------------------------------------ |
| `server`   | TCP listener, event loop, connection management  |
| `protocol` | RESP parser/writer for client communication      |
| `memory`   | Command handlers and data store                  |
| `aof`      | Append-only file writer/reader for persistence   |
| `peer`     | Peer connection management (for clustering, WIP) |

---

## 🧪 Example Session

```code
$ redis-cli -p 6379
127.0.0.1:6379> SET name TinyDB
OK
127.0.0.1:6379> GET name
"TinyDB"
127.0.0.1:6379> HSET user name Olamide
OK
127.0.0.1:6379> HGET user name
"Olamide"
```
---

## ⚡ Makefile Commands
| Command      | Description             |
| ------------ | ----------------------- |
| `make build` | Build the TinyDB binary |
| `make run`   | Run the server directly |
| `make clean` | Remove compiled files   |

---

## 🧱 Configuration

TinyDB reads from a .env file at the project root.

Example:
```ini
PORT=6379
```
### 🧭 Roadmap

 Clustered mode with replication

 Pub/Sub messaging

 Expiry and TTL support

 Benchmark tooling

 Web admin dashboard

 ---

 ## 🧑‍💻 Author

Muhammad Olamide
[📧 Email](mailto:muhammadolammi@gmail.com)

[🐙 GitHub](https://github.com/muhammadolammi)

[💼 LinkedIn](https://www.linkedin.com/in/muhammad-olammi/)

---

## 🪪 License

This project is licensed under the MIT License.
