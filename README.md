# Real-Time Forum


# **New bigining**

* know we will make bist Real time forum
## Project Overview

A modern, real-time forum application built with a robust Go backend, vanilla JavaScript frontend, and SQLite database. This project demonstrates real-time communication using WebSockets, enabling users to interact seamlessly with live updates for posts, comments, and private messaging.

## Features

- **User Registration & Authentication** - Secure email and username-based login with password hashing
- **Post Creation with Categories** - Users can create posts organized by categories
- **Comments on Posts** - Threaded discussions with real-time comment updates
- **Real-time Private Messaging** - Direct messaging between users with online/offline status indicators
- **Live WebSocket Updates** - Instant notifications for posts, comments, and messages without page refresh
- **Session Management** - Secure cookie-based session handling for persistent user authentication

## Tech Stack

**Backend:**
- Go 1.25+ with gorilla/websocket for WebSocket communication
- SQLite3 for data persistence
- bcrypt for password hashing and security

**Frontend:**
- Vanilla JavaScript (ES6 Modules) for lightweight, framework-free development
- WebSocket API for real-time communication
- DOM manipulation for dynamic UI updates

**Database:**
- SQLite3 for reliable, file-based relational database

## Project Structure

```
real-time-forum/
├── backend/          # Go server and business logic
│   ├── main.go
│   ├── handlers/     # HTTP request handlers
│   ├── websocket/    # WebSocket event handling
│   ├── database/     # SQLite database operations
│   └── models/       # Data structures
├── frontend/         # Client-side application
│   ├── index.html
│   ├── css/          # Styling
│   └── js/           # ES6 modules for functionality
├── data/             # SQLite database files
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.25 or higher
- SQLite3

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd real-time-forum
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the server:**
   ```bash
   go run ./backend/main.go
   ```

4. **Access the application:**
   Open your browser and navigate to `http://localhost:4444`

## Usage

1. **Register/Login** - Create an account with email and password, or login to an existing account
2. **Create Posts** - Navigate to the forum, select a category, and create a new post
3. **Engage in Discussions** - Comment on posts to participate in community conversations
4. **Send Private Messages** - Access the messaging section to send direct messages to other users
5. **Real-time Updates** - All changes appear instantly without requiring a page refresh

## Key Technologies Explained

**WebSockets** - Enables bidirectional, real-time communication between client and server. Unlike HTTP requests, WebSockets maintain an open connection for instant message delivery.

**Sessions & Cookies** - Provides secure user authentication. Server creates a session after login, stored in an encrypted cookie on the client, verified on each request.

**Message Throttling** - Optimizes performance by limiting the rate of messages sent, preventing server overload while maintaining responsive user experience.

## File Organization

- **backend/** - Contains all server-side logic, request handlers, WebSocket event processing, and database operations
- **frontend/** - Houses client-side code including HTML templates, CSS styling, and modular JavaScript files
- **data/** - Stores the SQLite database file for persistent data storage
