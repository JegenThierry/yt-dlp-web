<script setup lang="ts">
import { computed, ref } from 'vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UModal from '@nuxt/ui/components/Modal.vue'
import { useToast } from '@nuxt/ui/composables'
import { apiErrorStatus, getApiErrorMessage } from '@/api/apiErrors'
import { deletePlaylist } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'
import { usePlaylistsStore } from '@/stores/playlists'
import { countLabel } from '@/utils/format'

interface DeletePlaylistModalProps {
  playlist: Pick<Playlist, 'id' | 'title' | 'entry_count'>
}

const props = defineProps<DeletePlaylistModalProps>()

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  deleted: []
}>()

const toast = useToast()
const store = usePlaylistsStore()

const deleting = ref(false)
let deleted = false

const description = computed(() => {
  const { title, entry_count } = props.playlist
  if (entry_count === 0) return `"${title}" is removed from your library.`
  return `"${title}" and its ${countLabel(entry_count, 'song', 'songs')} are removed from your library. Downloaded files stay on the server.`
})

async function confirmDelete() {
  if (deleting.value) return

  deleting.value = true
  try {
    await deletePlaylist(props.playlist.id)
    finish()
  } catch (error) {
    if (apiErrorStatus(error) === 404) {
      finish()
      return
    }
    console.error('Failed to delete playlist:', error)
    toast.add({
      title: "Couldn't delete the playlist",
      description: getApiErrorMessage(error, 'Something went wrong while deleting it.'),
      color: 'error',
      duration: 0,
    })
  } finally {
    deleting.value = false
  }
}

function finish() {
  deleted = true
  open.value = false
  toast.add({
    title: 'Playlist deleted',
    description: `"${props.playlist.title}" was removed from your library.`,
    color: 'success',
    icon: 'i-lucide-trash-2',
  })
}

function preventFocusRestore(event: Event) {
  if (deleted) event.preventDefault()
}

function afterLeave() {
  if (!deleted) return
  deleted = false
  emit('deleted')
  store.remove(props.playlist.id)
}
</script>

<template>
  <UModal
    v-model:open="open"
    title="Delete this playlist?"
    :description="description"
    :dismissible="!deleting"
    :close="false"
    :content="{ onCloseAutoFocus: preventFocusRestore }"
    :ui="{ footer: 'justify-end' }"
    @after:leave="afterLeave"
  >
    <template #footer>
      <UButton
        color="neutral"
        variant="outline"
        label="Cancel"
        :disabled="deleting"
        @click="open = false"
      />
      <UButton
        color="error"
        icon="i-lucide-trash-2"
        label="Delete playlist"
        :loading="deleting"
        @click="confirmDelete"
      />
    </template>
  </UModal>
</template>
