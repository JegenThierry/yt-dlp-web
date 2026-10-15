import type { CreatePlaylistRequest, Playlist } from '@/models/playlist'
import type { PlaylistEntry, PlaylistEntryMetadataUpdate } from '@/models/playlistEntry'
import axiosInstance from './apiClient'

export async function createPlaylist(request: CreatePlaylistRequest): Promise<Playlist> {
  const { data } = await axiosInstance.post<Playlist>('/playlists', request)
  return data
}

export async function fetchPlaylists(): Promise<Playlist[]> {
  const { data } = await axiosInstance.get<Playlist[]>('/playlists')
  return data
}

export async function fetchPlaylist(id: number): Promise<Playlist> {
  const { data } = await axiosInstance.get<Playlist>(`/playlists/${id}`)
  return data
}

export async function updatePlaylistEntryMetadata(
  entryId: number,
  metadata: PlaylistEntryMetadataUpdate,
): Promise<PlaylistEntry> {
  const { data } = await axiosInstance.patch<PlaylistEntry>(
    `/playlists/entries/${entryId}/metadata`,
    metadata,
  )
  return data
}
