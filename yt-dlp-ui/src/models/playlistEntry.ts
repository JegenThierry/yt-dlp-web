import type { DownloadStatus, AudioFormat } from "./enums";

export interface PlaylistEntry {
	id: number;
	playlist_id: number;
	video_id: string;
	title: string;
	artist: string;
	duration: number;
	thumbnail: string;
	status: DownloadStatus;
	format: AudioFormat;
	created_at: Date;
	updated_at: Date;
}
