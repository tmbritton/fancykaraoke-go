# Fancy Karaoke

A modern, web-based karaoke party application built with Go, designed to be stable, maintainable, and reliable for years to come.
Features

## Core Functionality

- Party Creation: Create karaoke parties with unique, shareable URLs
- Song Search: Full-text search across an extensive song database
- Queue Management: Add, remove, and reorder songs in real-time
- YouTube Integration: Automatic video playback with seamless progression
- QR Code Generation: Easy party joining via QR codes
- Real-time Updates: Live queue updates across all connected users

## User Roles

- Hosts: Create parties, manage queues, control playback
- Attendees: Join parties, search songs, add to queue

## Technology Stack

### Backend

- Go: Built with Go's standard library for maximum stability
- TEMPL: Type-safe HTML templating with component-based architecture
- SQLite: Embedded database with full-text search (FTS5)
- Server-Sent Events: Real-time updates without WebSocket complexity

### Frontend

- HTMX: Dynamic interactions with minimal JavaScript
- Vanilla CSS: Clean, maintainable styling without framework dependencies
- Vanilla JS: No build, fewer problems
