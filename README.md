# Realtime Poll

A modern realtime polling platform built with **Go + WebSockets + React**.

This application allows users to create live polls, share them instantly, and watch votes update in realtime across all connected clients.

Unlike traditional polling apps that rely on refresh cycles, Realtime Poll streams vote updates live using persistent connections.

---

## ✨ Features

### Core

* Create polls with multiple options
* Public shareable voting links
* Realtime vote updates (WebSocket powered)
* Viewer presence count
* Automatic poll closing
* Anonymous or authenticated voting

### Access Control

* Public polls
* Login-required polls
* Whitelisted email polls
* Single-vote enforcement
* Vote change control

### Share Mode

* Dedicated share link viewer
* Voting without account (if allowed)
* Live result streaming
* Access-restricted UI states

### Realtime Engine

* WebSocket rooms per poll
* Incremental vote updates (delta updates)
* Snapshot recovery
* Automatic cleanup
* Poll expiration events

### Security

* HTTP-only session cookies
* Cross-site secure authentication
* Token-based share access
* Rate-safe vote handling

---

## 🧠 Architecture

Frontend and backend are completely decoupled and communicate via HTTP + WebSocket.

```
React SPA
   ↓ REST
Go API
   ↓
MongoDB

Realtime channel:
Client ⇄ WebSocket ⇄ Poll Room ⇄ Broadcast ⇄ All viewers
```

The backend is responsible for truth state.
The frontend only renders streamed state.

---

## 🛠 Tech Stack

### Backend

* Go (Gin)
* MongoDB
* Native WebSocket server
* Cookie session authentication

### Frontend

* React
* React Router
* Axios
* TailwindCSS

### Infrastructure

* VPS hosted API
* Netlify frontend
* Secure cross-site cookies

---

## 🔄 Realtime Model

The system uses a **state + delta hybrid model**:

* Initial join → snapshot
* Vote → delta broadcast
* Poll end → state event

This prevents desync and minimizes bandwidth.

---

## 🔐 Authentication

Authentication uses secure server-stored sessions.

Cookies:

* HttpOnly
* Secure
* SameSite=None
* Cross-subdomain enabled

No tokens are stored in localStorage.

---

## 🔗 Share Links

Each poll generates a public access token.

Example:

```
/share/<token>
```

The server resolves the token → poll → permissions → viewer capabilities.

---

## 🖼 Link Preview Support

The backend dynamically serves Open Graph metadata for share links so messaging platforms (WhatsApp, Discord, Telegram) display rich previews.

Preview includes:

* Poll question
* Preview image
* Voting call-to-action

---

## 🚀 Running Locally

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

## ⚠️ License & Usage

This project is **source-available**, not open-source.

You are allowed to:

* Read the code
* Learn from the implementation
* Reference architectural ideas

You are NOT allowed to:

* Copy the code
* Re-publish the project
* Deploy it yourself
* Use it commercially
* Use it for personal hosted services
* Create derivative works

Ownership and rights remain with the original author.

---

## 📌 Notes

This repository is shared for educational and demonstration purposes only.

The author may modify or revoke access at any time.

---

## Author

Koushik Babu
