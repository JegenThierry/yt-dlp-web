package downloads

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ytdl-api/internal/repository/db"
	"ytdl-api/pkg/ytdlp"
)

func (m *Manager) TaggedFile(ctx context.Context, entry db.PlaylistEntry) (string, error) {
	source := entry.FilePath.String
	meta := ytdlp.Metadata{
		Title:  entry.Title.String,
		Artist: entry.Artist.String,
		Genre:  entry.Genre.String,
		Year:   int64(entry.Year.Int32),
	}
	if meta.Empty() {
		return source, nil
	}

	ext := filepath.Ext(source)
	name := strings.TrimSuffix(filepath.Base(source), ext) + "-" + meta.Hash() + ext
	target := filepath.Join(m.taggedDir, name)
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	for _, dir := range []string{m.taggedDir, m.partialDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("create %s: %w", dir, err)
		}
	}

	tmp, err := os.CreateTemp(m.partialDir, "tag-*"+ext)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := ytdlp.Tag(ctx, source, tmp.Name(), meta); err != nil {
		return "", err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return "", fmt.Errorf("set tagged file mode: %w", err)
	}
	if err := os.Rename(tmp.Name(), target); err != nil {
		return "", fmt.Errorf("move tagged file: %w", err)
	}
	return target, nil
}
