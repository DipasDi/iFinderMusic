<div align="center">

# 🎵 iFinderMusic Bot

<<<<<<< HEAD
**Fast Go Telegram bot for searching, recognizing (Shazam), and downloading MP3 tracks.**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Telego](https://img.shields.io/badge/Telego-API-2CA5E0?style=flat-square&logo=telegram&logoColor=white)](https://github.com/mymmrac/telego)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
=======
**A high-performance Go Telegram bot for searching, previewing, and downloading MP3 music tracks.**

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Telegram API](https://img.shields.io/badge/Telegram_API-Telego-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://github.com/mymmrac/telego)
[![Redis](https://img.shields.io/badge/Cache-Redis_7-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Deploy-Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

<br/>

[Overview](#overview) •
[Features](#features) •
[Architecture](#architecture) •
[Quick Start](#quick-start) •
[Configuration](#configuration) •
[Project Structure](#project-structure) •
[Testing](#testing)
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0

</div>

---

<<<<<<< HEAD
## ✨ Features

- 🔎 **Music Search** — Find tracks by artist or title with duration tags and ID3 metadata.
- 🎧 **Shazam Recognition** — Send a voice message, audio, video, or video note to identify the song and get the MP3.
- 🔥 **Charts & New Releases** — Quick access to trending tracks (`🔥 Популярные`) and fresh releases (`✨ Новинки`).
- ⚡️ **Redis Caching** — Instant inline pagination (`⬅️ Back` / `Next ➡️`) backed by Redis session caching.
- 🧹 **Clean UX** — Auto-removes temporary status messages and supports MarkdownV2 with plaintext fallback.

---

## 🏗 Architecture

```mermaid
flowchart LR
    User([👤 User]) <-->|Messages / Audio| TG[Telegram API]
    TG <-->|Long Polling| Bot[🤖 Go Bot]

    subgraph Core Services
        Bot <-->|Session Cache| Redis[(🗄 Redis)]
        Bot -->|Voice / Video| Shazam[🎧 GoShazam]
        Bot -->|Search Query| Scraper[🕷 Colly Scraper]
        Scraper -->|Fetch & Stream| Source[🌐 Music Catalog]
    end
=======
## Overview

**iFinderMusic** is a Telegram music bot that allows users to search songs by artist or title, explore trending charts and new releases, and download MP3 audio files directly into their chat.

The bot parses online music directories via a web scraper, caches search sessions in Redis to provide instant inline pagination, and asynchronously transmits audio files with proper ID3 metadata (artist, title, duration).

---

## Features

- 🔎 **Fast Music Search:** Query tracks by song name or artist (e.g., `The Weeknd`, `Miyagi`, `Queen`) with duration tags.
- 📄 **Interactive Pagination:** Paginated results (10 tracks per page) navigated via Telegram inline buttons (`⬅️ Back` / `Next ➡️`).
- 🔥 **Curated Categories:** Instant access to trending charts (`🔥 Популярные`) and newest releases (`✨ Новинки`).
- ⚡️ **Redis Caching:** Search results and user sessions are cached in Redis to minimize outbound requests and enable snappy pagination.
- 🧹 **Clean Chat UX:** Automatically removes prompt messages and temporary status indicators to keep the chat history organized.
- 🛡 **MarkdownV2 with Fallback:** Formatted text with automatic plaintext fallback in case of Telegram formatting errors.
- 🧪 **Comprehensive Test Suite:** High unit test coverage covering UI formatting, keyboard pagination, configuration, scraper mock servers, and storage.

---

## Architecture

```mermaid
flowchart TD
    User([👤 User]) -->|Search Query / Button Click| TG[Telegram Bot API]
    TG -->|Long Polling| Bot[🤖 Go Bot Service]

    subgraph Internal Architecture
        Bot <-->|Cache results & queries| Redis[(🗄 Redis)]
        Bot -->|Extract track list| Scraper[🕷 Colly Scraper]
        Scraper -->|Parse HTML| Source[🌐 Music Catalog]
        Bot -->|Fetch audio stream| Downloader[⬇️ Audio Downloader]
        Downloader -->|Write temporary file| TempStore[💾 Local Temp Storage]
    end

    Downloader -->|Send MP3 with tags| TG
    TG -->|Audio File| User
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
```

---

<<<<<<< HEAD
## 🚀 Quick Start

### Docker Compose (Recommended)

```bash
git clone [https://github.com/your-username/ifindermusic.git](https://github.com/your-username/ifindermusic.git)
cd ifindermusic
cp .env.example .env
```

Set your `BOT_TOKEN` (from [@BotFather](https://t.me/BotFather)) inside `.env`, then run:

```bash
docker-compose up --build -d
```

### Local Run

Requires **Go 1.22+** and a running **Redis** instance.

```bash
cp .env.example .env
go mod download
go run main.go
```

---

## ⚙️ Configuration

| Variable          | Required |     Default      | Description                                  |
| :---------------- | :------: | :--------------: | :------------------------------------------- |
| `BOT_TOKEN`       | **Yes**  |        —         | Telegram Bot API token                       |
| `REDIS_ADDR`      |    No    | `127.0.0.1:6379` | Redis address (`host:port`)                  |
| `REDIS_USER`      |    No    |    `default`     | Redis ACL username                           |
| `REDIS_PASSWORD`  |    No    |       `""`       | Redis password                               |
| `TRACK_USERS`     |    No    |      `true`      | Log user activity to JSON (`true` / `false`) |
| `USERS_FILE_PATH` |    No    |   `users.json`   | File path for user activity logs             |

---

## 📁 Project Structure

```text
├── internal/
│   ├── bot/          # Telegram handlers, polling loop, audio dispatcher
│   ├── config/       # Environment configuration
│   ├── model/        # Domain models
│   ├── scraper/      # Colly web scraper & MP3 downloader
│   ├── shazam/       # Audio recognition integration
│   ├── storage/      # Redis cache & user tracking
│   └── ui/           # Keyboards & message formatting
├── lib/
│   └── goshazam/     # Audio fingerprinting & Shazam library
├── Dockerfile
├── docker-compose.yml
└── main.go
=======
## Quick Start

### Option 1: Run with Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/ifindermusic.git
   cd ifindermusic
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and set your `BOT_TOKEN` obtained from [@BotFather](https://t.me/BotFather).

3. **Launch the services:**
   ```bash
   docker-compose up --build -d
   ```

4. Check container logs:
   ```bash
   docker-compose logs -f bot
   ```

---

### Option 2: Run Locally

#### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- Running Redis instance (or local Docker container: `docker run -p 6379:6379 -d redis:7-alpine`)

#### Steps
1. **Copy configuration:**
   ```bash
   cp .env.example .env
   ```
   Fill in your `BOT_TOKEN`.

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

---

## Configuration

Settings are configured via environment variables or a `.env` file:

| Variable | Required | Default | Description |
| :--- | :---: | :---: | :--- |
| `BOT_TOKEN` | **Yes** | — | Telegram Bot API token obtained from [@BotFather](https://t.me/BotFather) |
| `REDIS_ADDR` | No | `127.0.0.1:6379` | Redis connection address (`host:port`) |
| `REDIS_HOST` | No | — | Redis host (if address is not specified) |
| `REDIS_PASSWORD`| No | `""` | Redis authentication password (if enabled) |

---

## Project Structure

```
.
├── internal/
│   ├── bot/          # Telegram bot handlers, polling loop, and audio dispatcher
│   ├── config/       # Environment loading and configuration models
│   ├── model/        # Core domain models (Music track schema)
│   ├── scraper/      # Web scraper (Colly) and HTTP audio downloader
│   ├── storage/      # Redis client and cache storage operations
│   └── ui/           # Keyboards, formatting utilities, and message templates
├── .env.example      # Example environment configuration template
├── .gitignore        # Git ignore rules for binaries, env files, and temp assets
├── Dockerfile        # Production multi-stage Docker build
├── docker-compose.yml# Production stack definition (Bot + Redis)
├── go.mod            # Go module definition
├── go.sum            # Checksums for dependencies
├── LICENSE           # MIT License
└── main.go           # Application entrypoint
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
```

---

<<<<<<< HEAD
## 🧪 Testing

```bash
go test -v -cover ./...
=======
## Testing

Run all unit tests across the entire project:

```bash
go test -v ./...
```

Run tests with test coverage profiling:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
```

---

<<<<<<< HEAD
## 📄 License

Distributed under the [MIT License](LICENSE).
=======
## License

This project is licensed under the [MIT License](LICENSE).
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
