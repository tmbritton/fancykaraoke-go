-- Songs table with FTS5 for full-text search
CREATE TABLE songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    artist TEXT NOT NULL,
    title TEXT NOT NULL,
    uploader TEXT,
    youtube_id TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create FTS5 virtual table for song search
CREATE VIRTUAL TABLE songs_fts USING fts5(
    artist, title, content='songs', content_rowid='id'
);

-- Trigger to keep FTS5 in sync with songs table
CREATE TRIGGER songs_fts_insert AFTER INSERT ON songs BEGIN
    INSERT INTO songs_fts(rowid, artist, title) VALUES (new.id, new.artist, new.title);
END;

CREATE TRIGGER songs_fts_delete AFTER DELETE ON songs BEGIN
    INSERT INTO songs_fts(songs_fts, rowid, artist, title) VALUES('delete', old.id, old.artist, old.title);
END;

CREATE TRIGGER songs_fts_update AFTER UPDATE ON songs BEGIN
    INSERT INTO songs_fts(songs_fts, rowid, artist, title) VALUES('delete', old.id, old.artist, old.title);
    INSERT INTO songs_fts(rowid, artist, title) VALUES (new.id, new.artist, new.title);
END;

-- Parties table
CREATE TABLE parties (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    slug TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_played DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Queue table
CREATE TABLE queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    party_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    added_by TEXT,
    priority INTEGER DEFAULT 0,
    hidden BOOLEAN DEFAULT false,
    FOREIGN KEY (party_id) REFERENCES parties(id),
    FOREIGN KEY (song_id) REFERENCES songs(id)
);

-- Indexes for performance
CREATE INDEX idx_songs_youtube_id ON songs(youtube_id);
CREATE INDEX idx_parties_slug ON parties(slug);
CREATE INDEX idx_queue_party_id ON queue(party_id);
CREATE INDEX idx_queue_song_id ON queue(song_id);
