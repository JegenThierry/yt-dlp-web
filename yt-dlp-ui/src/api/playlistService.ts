import type { CreatePlaylistRequest, Playlist } from '@/models/playlist'
import axiosInstance from './apiClient'

export async function createPlaylist(request: CreatePlaylistRequest): Promise<Playlist> {
  const { data } = await axiosInstance.post<Playlist>('/playlists', request)
  return data
}

export async function fetchPlaylists(): Promise<Playlist[]> {
  const { data } = await axiosInstance.get<Playlist[]>('/playlists')
  return data
}
