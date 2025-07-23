BEGIN TRANSACTION;

DROP TABLE parties;
DROP TABLE queue;

-- Parties table
CREATE TABLE parties (
    id TEXT PRIMARY KEY,
    name TEXT,
    slug TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

--- Queue item status enum
CREATE TABLE queue_item_status (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    label TEXT UNIQUE NOT NULL
);

-- Queue table
CREATE TABLE queue_items (
    id TEXT PRIMARY KEY,
    party_id TEXT NOT NULL,
    song_id INT NOT NULL,
    added_by TEXT,
    priority INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (party_id) REFERENCES parties(id),
    FOREIGN KEY (song_id) REFERENCES songs(id),
    FOREIGN KEY (status) REFERENCES queue_item_status(id)
);

INSERT INTO queue_item_status (label) VALUES ("UNPLAYED");
INSERT INTO queue_item_status (label) VALUES ("PLAYING");
INSERT INTO queue_item_status (label) VALUES ("PLAYED");
INSERT INTO queue_item_status (label) VALUES ("REMOVED");

END TRANSACTION;
