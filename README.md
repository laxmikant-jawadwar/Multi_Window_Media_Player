# Multi-Window Media Sequencer with Sync Playback

A full-stack media sequencing application built with **React** and **Golang**. The application manages multiple display windows, each with its own configurable playlist, supports continuous playback, dynamic playlist changes, a configured 5-hour operating cycle, and synchronized playback of a selected media item across all windows.

## Live Deployment

- **Frontend:** https://multi-window-media-player-frontend.onrender.com
- **Backend:** https://multi-window-media-player-1.onrender.com
- **Backend Health Check:** https://multi-window-media-player-1.onrender.com/health
- **GitHub Repository:** https://github.com/laxmikant-jawadwar/Multi_Window_Media_Player

## Features

- Multiple media display windows
- Individual playlist for each window
- Image, video, and blank media support
- Dynamic creation of new media
- Dynamic addition and removal of playlist items
- Playlist reordering
- Continuous media sequencing
- Persistent PostgreSQL storage
- Start and stop playback for each window
- 5-hour operating cycle per playback session
- Synchronized media playback across all windows
- Configurable synchronization duration
- Fallback states for unavailable or unsupported media
- React frontend with polling for current backend state
- REST API built with Golang
- Dockerized backend
- Production deployment on Render

## Requirements Covered

### Frontend

- Built using React and Vite
- Displays multiple media windows
- Supports image, video, and blank states
- Displays each window's configured playlist
- Supports dynamic playlist changes
- Provides controls for creating media, adding media to playlists, starting playback, and triggering synchronization

### Backend

- Built using Golang
- Uses persistent PostgreSQL storage
- Supports dynamic playlist management
- Maintains playback state for each window
- Calculates the currently active media based on elapsed playback time and playlist durations
- Implements synchronized playback through a sync session
- Preserves each window's normal playlist and playback cycle during synchronization

### Deployment

- React frontend deployed on Render
- Golang backend deployed on Render as a Docker Web Service
- PostgreSQL database hosted using Supabase
- Environment-based configuration used for database and frontend API settings

## Application Behavior

### Window Playback

Each window has its own playlist.

Example:

```text
Window 1
  ├── Image 1 - 30 sec
  ├── Video 1 - 60 sec
  └── Blank   - 10 sec

Window 2
  ├── Video 2 - 45 sec
  └── Image 2 - 30 sec

Window 3
  ├── Image 3 - 20 sec
  └── Video 3 - 90 sec
```

The backend calculates the current media using:

- Playback start time
- Total playlist duration
- Individual media durations
- Current elapsed time

The sequence continues from one media item to the next without introducing unintended blank periods.

### 5-Hour Playback Cycle

A playback session is configured as a **5-hour operating period**.

The behavior is:

```text
Start Playback
      |
      v
Playlist plays continuously
      |
      v
Playlist sequence continues during the 5-hour period
      |
      v
5 hours completed
      |
      v
Playback STOPPED
```

The playlist configuration is not deleted when the 5-hour cycle completes.

Starting playback again begins a new 5-hour operating cycle.

> **Assumption:** The 5-hour behavior follows the clarification received for the assignment: after five hours, the window stops while retaining its playlist configuration.

### Dynamic Playlist Changes

Media can be created and added to a window without rebuilding or restarting the application.

Supported operations include:

- Create image media
- Create video media
- Create blank media
- Add media to a window
- Remove a playlist item
- Change playlist item position

Dynamic playlist changes do not require recreating the window.

### Blank Media

Blank is treated as an explicit media type.

A blank media item does not require a URL.

Example:

```text
Title: Blank Screen
Type: blank
URL: empty
Duration: 10 seconds
```

The application does not automatically turn unused time into blank playback.

### Fallback State

The frontend displays a fallback message when intended media cannot be rendered, for example:

- Invalid or unavailable image URL
- Invalid or unavailable video URL
- Unsupported media type
- No current media available

## Synchronization

The application supports a sync action where one selected media item is displayed across all windows simultaneously.

Example:

```text
Window 1 → Image A
Window 2 → Video B
Window 3 → Image C

             |
          Start Sync
             |
             v

Window 1 → Media M2
Window 2 → Media M2
Window 3 → Media M2
```

The sync duration is configurable.

After the sync duration expires:

```text
Window 1 → resumes its own playlist
Window 2 → resumes its own playlist
Window 3 → resumes its own playlist
```

The synchronization does not replace or delete any window's playlist configuration.

The normal playback cycle is also not reset merely because a sync session starts.

## Architecture

```text
                    React Frontend
                         |
                         | REST API
                         v
                Golang Backend
                         |
          +--------------+--------------+
          |              |              |
          v              v              v
      Windows        Playlists      Playback/Sync
          |              |              |
          +--------------+--------------+
                         |
                         v
                   PostgreSQL
                     (Supabase)
```

## Backend Structure

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   ├── repositories/
│   └── services/
├── migrations/
├── Dockerfile
├── .env.example
├── go.mod
└── go.sum
```

## Frontend Structure

```text
frontend/
├── src/
│   ├── components/
│   │   ├── CreateMediaModal.jsx
│   │   ├── Playlist.jsx
│   │   ├── SyncControl.jsx
│   │   └── WindowCard.jsx
│   ├── services/
│   │   └── api.js
│   ├── App.jsx
│   ├── index.css
│   └── usePolling.js
├── Dockerfile
├── nginx.conf
├── .env.example
├── package.json
└── vite.config.js
```

## Database

Supabase PostgreSQL is used as the persistent storage layer.

The main tables are:

### `windows`

Stores display window information.

Important fields:

- `id`
- `name`
- `created_at`
- `updated_at`

### `media`

Stores media information.

Important fields:

- `id`
- `title`
- `media_type`
- `url`
- `duration_seconds`
- `created_at`
- `updated_at`

### `playlist_items`

Maps media to windows and stores playlist order.

Important fields:

- `id`
- `window_id`
- `media_id`
- `position`
- `created_at`

### `playback_states`

Stores the current playback state for each window.

Important fields:

- `window_id`
- `current_playlist_item_id`
- `cycle_started_at`
- `status`
- `updated_at`

### `sync_sessions`

Stores synchronization sessions.

Important fields:

- `id`
- `media_id`
- `started_at`
- `duration_seconds`
- `created_at`

Database schema/recreation SQL is available in the `backend/migrations/` directory.

## API Endpoints

The backend exposes REST endpoints.

### Health

```http
GET /health
```

Returns:

```json
{
  "status": "ok"
}
```

### Windows

```http
GET /windows
```

Returns all configured windows.

### Media

```http
GET /media
```

Returns available media.

```http
POST /media
```

Creates new media.

Example request:

```json
{
  "title": "Sample Image",
  "media_type": "image",
  "url": "https://example.com/sample.jpg",
  "duration_seconds": 30
}
```

For blank media:

```json
{
  "title": "Blank Screen",
  "media_type": "blank",
  "url": "",
  "duration_seconds": 10
}
```

### Playlist

```http
GET /windows/{windowID}/playlist
```

Gets the playlist for a window.

```http
POST /windows/{windowID}/playlist
```

Adds media to a window's playlist.

```http
DELETE /windows/{windowID}/playlist/{playlistItemID}
```

Removes a playlist item.

```http
PUT /windows/{windowID}/playlist/{playlistItemID}
```

Updates/reorders a playlist item.

### Playback

```http
POST /windows/{windowID}/playback/start
```

Starts a new 5-hour playback cycle.

```http
POST /windows/{windowID}/playback/stop
```

Stops playback.

```http
GET /windows/{windowID}/playback
```

Gets the current playback state.

```http
GET /windows/{windowID}/playback/current
```

Gets the media currently selected for playback.

### Synchronization

```http
POST /sync
```

Starts synchronized playback.

Example:

```json
{
  "media_id": "MEDIA_UUID",
  "duration_seconds": 30
}
```

```http
GET /sync
```

Returns the current/latest sync status.

## Local Setup

### Prerequisites

Install:

- Go 1.25+
- Node.js 20+
- PostgreSQL
- Git

### Clone the Repository

```bash
git clone https://github.com/laxmikant-jawadwar/Multi_Window_Media_Player.git
cd Multi_Window_Media_Player
```

## Backend Setup

Go to the backend directory:

```bash
cd backend
```

Create a `.env` file:

```env
PORT=8080
DATABASE_URL=postgresql://postgres:YOUR_PASSWORD@localhost:5432/media_sequencer
```

Run the SQL migrations from the `backend/migrations/` directory against your PostgreSQL database.

Install dependencies:

```bash
go mod download
```

Start the backend:

```bash
go run ./cmd/server
```

The backend will be available at:

```text
http://localhost:8080
```

Check:

```text
http://localhost:8080/health
```

## Frontend Setup

Open another terminal and go to:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Create `.env`:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Start the development server:

```bash
npm run dev
```

The Vite development server will display the local frontend URL in the terminal.

## Production Environment Variables

### Backend

```env
PORT=8080
DATABASE_URL=<PostgreSQL connection string>
```

### Frontend

```env
VITE_API_BASE_URL=https://multi-window-media-player-1.onrender.com
```

For Vite, `VITE_API_BASE_URL` is embedded into the production JavaScript bundle during the build process.

## Docker

The backend includes a multi-stage Dockerfile.

Build:

```bash
cd backend
docker build -t media-sequencer-backend .
```

Run:

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL="YOUR_DATABASE_URL" \
  media-sequencer-backend
```

The frontend also includes Docker/Nginx configuration for container-based deployment if container deployment is required.

## Deployment

### Backend Deployment

The backend is deployed on **Render** as a Docker Web Service.

Configuration:

```text
Runtime: Docker
Root Directory: backend
Docker Build Context Directory: backend
Dockerfile Path: Dockerfile
```

The database connection is provided through the `DATABASE_URL` environment variable.

### Frontend Deployment

The frontend is deployed on **Render** as a Static Site.

Configuration:

```text
Root Directory: frontend
Build Command: npm install && npm run build
Publish Directory: dist
```

Environment variable:

```text
VITE_API_BASE_URL=https://multi-window-media-player-1.onrender.com
```

## Testing and Verification

### Backend

```bash
go build ./...
```

```bash
go vet ./...
```

### Frontend

```bash
npm run build
```

The production frontend and backend builds complete successfully.

The deployed backend health endpoint was also verified successfully.

## Design Decisions and Assumptions

1. **Persistent storage:** PostgreSQL was selected for relational storage of windows, media, playlists, playback state, and sync sessions.

2. **Media storage:** The database stores public media URLs rather than binary media files. This keeps the backend lightweight and allows URLs to point to suitable public media/object-storage providers.

3. **Dynamic media:** Media can be created and assigned to windows at runtime without rebuilding the application.

4. **Blank media:** Blank is an explicit playlist media type and does not require a URL.

5. **Playback calculation:** The backend calculates the current playlist item from elapsed time and media durations rather than requiring a continuously running backend timer for every media transition.

6. **5-hour cycle:** A playback start creates a 5-hour operating cycle. When the cycle completes, the window becomes `STOPPED` while its playlist remains intact.

7. **Synchronization:** Sync is represented as a temporary session. While active, all windows use the selected sync media; after expiration, each window returns to its normal playlist behavior.

8. **Frontend polling:** The frontend periodically requests backend state so that playback, sync, and playlist changes are reflected without requiring a continuously streamed "next media" API.

9. **No automatic blank gaps:** The application does not insert blank playback between media items unless blank is explicitly configured in the playlist.

## How to Demo

1. Open the live frontend.
2. Create a few image, video, or blank media items using **Create New Media**.
3. Add different media items to Window 1, Window 2, and Window 3.
4. Start playback for each window.
5. Observe each window following its own configured sequence.
6. Reorder or remove a playlist item to demonstrate dynamic playlist management.
7. Select one media item in **Sync Broadcast Control**.
8. Set a synchronization duration, for example 30 seconds.
9. Start Sync.
10. Observe the selected media displayed across all windows simultaneously.
11. After the sync duration expires, observe the windows return to their normal playback behavior.

## Repository

GitHub:

https://github.com/laxmikant-jawadwar/Multi_Window_Media_Player

## Live Application

Frontend:

https://multi-window-media-player-frontend.onrender.com

Backend:

https://multi-window-media-player-1.onrender.com

Health:

https://multi-window-media-player-1.onrender.com/health