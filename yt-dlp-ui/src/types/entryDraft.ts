export interface EntryDraft {
  id: number
  title: string
  artist: string
  genre: string
  year: string
}

export type DraftField = 'title' | 'artist' | 'genre' | 'year'

export type DraftErrors = Partial<Record<DraftField, string>>

export interface DraftEdit {
  id: number
  field: DraftField
  from: string
  to: string
}
