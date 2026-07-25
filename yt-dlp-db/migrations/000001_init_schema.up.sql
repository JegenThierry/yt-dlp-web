CREATE TYPE entry_status AS ENUM ('pending', 'downloading', 'completed', 'failed');

CREATE TYPE media_format AS ENUM ('mp3', 'aac', 'flac', 'wav', 'm4a', 'opus', 'vorbis', 'alac');

CREATE TABLE playlists
(
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    url         text NOT NULL,
    title       text NULL,
    created_at  timestamp with time zone DEFAULT now(),
    updated_at  timestamp with time zone DEFAULT now(),
);

CREATE TABLE playlist_entries
(
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    playlist_id bigint NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
    video_id    text   NOT NULL,
    title       text   NULL,
    status      entry_status             DEFAULT 'pending',
    format      media_format             DEFAULT 'mp3',
    created_at  timestamp with time zone DEFAULT now(),
    updated_at  timestamp with time zone DEFAULT now(),

    CONSTRAINT playlist_entries_unique UNIQUE (playlist_id, video_id)
);

CREATE INDEX idx_playlist_entries_playlist_id ON playlist_entries (playlist_id);