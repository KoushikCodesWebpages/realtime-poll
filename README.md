# ⚡ Realtime Poll

Live voting. Zero refresh. Instant feedback.

A modern realtime polling platform built with **Go + WebSockets + React** where every vote appears instantly across all connected viewers.

---

## 🧭 How it works (User Flow)

### 1️⃣ Create a poll

Create a question and options.

```
"Best backend language?"
- Go
- Node
- Rust
- Python
```

📷 `docs/create.png`

---

### 2️⃣ Share the link

Send the generated link anywhere.

```
https://live-poll.clqit.in/share/ABC123
```

📷 `docs/share.png`

Messaging apps automatically show a preview card.

---

### 3️⃣ People vote

Users open the link and vote instantly.

📷 `docs/vote.png`

No refresh required.
No login required (optional restriction supported).

---

### 4️⃣ Watch results live

Every connected device updates in real time.

📷 `docs/live.png`

Votes stream instantly using WebSockets.

---

## 🧠 System Concept

Traditional polling:

```
vote → save → refresh → see result
```

Realtime Poll:

```
vote → broadcast → everyone explaining "whoa"
```

---

## 🔄 Realtime Engine

The app uses a hybrid model:

| Event | Behavior      |
| ----- | ------------- |
| Join  | full snapshot |
| Vote  | delta update  |
| Close | final state   |

This avoids desync and minimizes bandwidth.

---

## 🔐 Access Control

Polls can require different levels of permission:

| Mode      | Behavior                 |
| --------- | ------------------------ |
| Public    | anyone can vote          |
| Login     | authenticated users only |
| Whitelist | specific emails only     |
| Read-only | results view only        |

UI automatically adapts based on viewer capability.

---

## 🧩 Architecture Overview

```
           ┌────────────┐
           │   React    │
           │   Client   │
           └─────┬──────┘
                 │ REST
                 ▼
           ┌────────────┐
           │    Go API  │
           └─────┬──────┘
                 │
         ┌───────┴────────┐
         │                 │
         ▼                 ▼
   MongoDB           WebSocket Hub
                           │
                    Poll Rooms (per poll)
                           │
                     Broadcast events
```

Frontend renders state.
Backend owns truth.

---

## 🛠 Tech Stack

### Backend

* Go (Gin)
* MongoDB
* Native WebSocket engine
* Cookie session authentication

### Frontend

* React
* React Router
* Axios
* TailwindCSS

---

## 🌐 Share Links & Preview Cards

When a poll is shared:

```
/share/<token>
```

Messaging platforms fetch dynamic metadata:

* Title
* Description
* Generated preview image

No JavaScript required.

---

## 🔐 Authentication

Sessions are server-managed.

Cookies:

* HttpOnly
* Secure
* SameSite=None
* Cross-subdomain safe

No localStorage tokens.

---

## 🚀 Run Locally

### Backend

```bash
go mod tidy
go run main.go
```

### Frontend

```bash
npm install
npm run dev
```

---

## 📖 Intended Purpose

This project is shared for **learning and reference**.

You may:

* Study the architecture
* Understand realtime patterns
* Explore implementation ideas

You may NOT:

* Rehost the service
* Copy the codebase
* Use commercially
* Create derivative deployments

All rights remain with the author.

---

## 👤 Author

Koushik Babu
