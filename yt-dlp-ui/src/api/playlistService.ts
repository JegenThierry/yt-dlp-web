import type { CreatePlaylistRequest, Playlist, PlaylistImport } from '@/models/playlist'
import type { PlaylistEntry, PlaylistEntryMetadataUpdate } from '@/models/playlistEntry'
import axiosInstance from './apiClient'

const PLAYLIST_READ_TIMEOUT_MS = 300_000

export async function createPlaylist(request: CreatePlaylistRequest): Promise<PlaylistImport> {
  const { data } = await axiosInstance.post<PlaylistImport>('/playlists', request, {
    timeout: PLAYLIST_READ_TIMEOUT_MS,
  })
  return data
}

export async function refreshPlaylist(id: number): Promise<PlaylistImport> {
  const { data } = await axiosInstance.post<PlaylistImport>(`/playlists/${id}/refresh`, null, {
    timeout: PLAYLIST_READ_TIMEOUT_MS,
  })
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

export async function deletePlaylist(id: number): Promise<void> {
  await axiosInstance.delete(`/playlists/${id}`)
}

export async function updatePlaylistMetadata(
  playlistId: number,
  updates: PlaylistEntryMetadataUpdate[],
): Promise<PlaylistEntry[]> {
  const { data } = await axiosInstance.patch<PlaylistEntry[]>(
    `/playlists/${playlistId}/entries/metadata`,
    updates,
  )
  return data
}
