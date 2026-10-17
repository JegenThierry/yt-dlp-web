import type { DraftEdit, DraftField, EntryDraft } from '@/types/entryDraft'

export type CleanupAction = 'remove-topic' | 'split-artist-title' | 'strip-suffixes' | 'trim-spaces'

const TOPIC_SUFFIX = /\s+-\s+Topic$/i
const ARTIST_TITLE = /^(.+?)\s+[-–—]\s+(.+)$/

const SUFFIX_WORD = String.raw`(?:official|music|lyrics?|video|audio|visuali[sz]er|hd|hq|4k|mv|m\/v)`
const BRACKETED_SUFFIX = new RegExp(
  String.raw`\s*[([]\s*${SUFFIX_WORD}(?:\s+${SUFFIX_WORD})*\s*[)\]]`,
  'gi',
)
const TRAILING_SUFFIX =
  /\s+(?:official(?:\s+(?:music|lyrics?))?\s+(?:video|audio|visuali[sz]er|mv|m\/v)|music\s+video|m\/v|mv)$/i
const DANGLING_SEPARATOR = /\s*[-–—|:]\s*$/

export function removeTopic(artist: string): string {
  return artist.replace(TOPIC_SUFFIX, '')
}

export function stripSuffixes(title: string): string {
  const stripped = title
    .replace(BRACKETED_SUFFIX, '')
    .replace(TRAILING_SUFFIX, '')
    .replace(DANGLING_SEPARATOR, '')
    .replace(/\s{2,}/g, ' ')
    .trim()
  return stripped === '' ? title : stripped
}

export function collapseSpaces(value: string): string {
  return value.replace(/\s+/g, ' ').trim()
}

export function splitArtistTitle(title: string): { artist: string; title: string } | null {
  const match = ARTIST_TITLE.exec(title.trim())
  if (!match?.[1] || !match[2]) return null
  return { artist: match[1].trim(), title: match[2].trim() }
}

function edit(draft: EntryDraft, field: DraftField, to: string): DraftEdit[] {
  return draft[field] === to ? [] : [{ id: draft.id, field, from: draft[field], to }]
}

export function cleanupEdits(drafts: EntryDraft[], action: CleanupAction): DraftEdit[] {
  return drafts.flatMap((draft) => {
    switch (action) {
      case 'remove-topic':
        return edit(draft, 'artist', removeTopic(draft.artist))
      case 'strip-suffixes':
        return edit(draft, 'title', stripSuffixes(draft.title))
      case 'trim-spaces':
        return (['title', 'artist', 'genre', 'year'] as const).flatMap((field) =>
          edit(draft, field, collapseSpaces(draft[field])),
        )
      case 'split-artist-title': {
        const split = splitArtistTitle(draft.title)
        if (!split) return []
        return [...edit(draft, 'artist', split.artist), ...edit(draft, 'title', split.title)]
      }
    }
  })
}

export function setAllEdits(drafts: EntryDraft[], field: DraftField, value: string): DraftEdit[] {
  return drafts.flatMap((draft) => edit(draft, field, value))
}
