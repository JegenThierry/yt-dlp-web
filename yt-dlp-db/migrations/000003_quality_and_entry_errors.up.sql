CREATE TYPE audio_quality AS ENUM ('best', 'high', 'standard', 'low');

CREATE TYPE entry_error_code AS ENUM ('unavailable', 'sign_in_required', 'blocked', 'network', 'postprocessing', 'unknown');

ALTER TABLE playlists
    ADD COLUMN quality audio_quality NOT NULL DEFAULT 'best';

ALTER TABLE playlist_entries
    ADD COLUMN error_message TEXT             NULL,
    ADD COLUMN error_code    entry_error_code NULL;
