import type { PlaylistEntry, PlaylistEntryMetadataUpdate } from '@/models/playlistEntry'
import type { DraftEdit, DraftErrors, DraftField, EntryDraft } from '@/types/entryDraft'

export const DRAFT_FIELDS: DraftField[] = ['title', 'artist', 'genre', 'year']

const MIN_YEAR = 1900
const MAX_YEAR = 2100

const TITLE_REQUIRED = 'Enter a title.'
const YEAR_INVALID = `Use a year from ${MIN_YEAR} to ${MAX_YEAR}.`

export function toDraft(entry: PlaylistEntry): EntryDraft {
  return {
    id: entry.id,
    title: entry.title,
    artist: entry.artist,
    genre: entry.genre,
    year: entry.year ? String(entry.year) : '',
  }
}

export function savedValue(entry: PlaylistEntry, field: DraftField): string {
  return toDraft(entry)[field]
}

export function changedFields(draft: EntryDraft, entry: PlaylistEntry): DraftField[] {
  return DRAFT_FIELDS.filter((field) => draft[field].trim() !== savedValue(entry, field))
}

export function yearError(value: string): string | undefined {
  const year = value.trim()
  if (year === '') return undefined
  if (!/^\d{4}$/.test(year)) return YEAR_INVALID
  const number = Number(year)
  return number < MIN_YEAR || number > MAX_YEAR ? YEAR_INVALID : undefined
}

export function draftErrors(draft: EntryDraft): DraftErrors {
  const errors: DraftErrors = {}
  if (draft.title.trim() === '') errors.title = TITLE_REQUIRED
  const year = yearError(draft.year)
  if (year) errors.year = year
  return errors
}

export function toUpdate(draft: EntryDraft): PlaylistEntryMetadataUpdate {
  const year = draft.year.trim()
  return {
    id: draft.id,
    title: draft.title.trim(),
    artist: draft.artist.trim(),
    genre: draft.genre.trim(),
    year: year === '' ? 0 : Number(year),
  }
}

export function applyEdits(drafts: EntryDraft[], edits: DraftEdit[]) {
  const byId = new Map(drafts.map((draft) => [draft.id, draft]))
  for (const edit of edits) {
    const draft = byId.get(edit.id)
    if (draft) draft[edit.field] = edit.to
  }
}

export function revertEdits(drafts: EntryDraft[], edits: DraftEdit[]) {
  const byId = new Map(drafts.map((draft) => [draft.id, draft]))
  for (const edit of edits) {
    const draft = byId.get(edit.id)
    if (draft && draft[edit.field] === edit.to) draft[edit.field] = edit.from
  }
}
