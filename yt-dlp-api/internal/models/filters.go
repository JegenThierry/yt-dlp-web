package models

type FilterOptions struct {
	Format        AudioFormat  `json:"format"`
	Quality       AudioQuality `json:"quality"`
	EmbedMetadata bool         `json:"embed_metadata"`
}
