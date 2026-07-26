import { CreatePlaylistRequest, Playlist } from "@/models/playlist";
import { api } from "./apiClient";

export async function createPlaylist(request: CreatePlaylistRequest): Promise<Playlist> {
  return await api.post<Playlist>('/playlists', request);
}

export async function fetchPlaylists(): Promise<Playlist[]> {
  return await api.get<Playlist[]>('/playlists');
}
