<script setup lang="ts">
</script>

<template>
  <div class="max-w-[720px] mx-auto p-4 border border-default rounded">
    <h3 class="text-lg font-semibold mb-4">Download</h3>
    <form>
      <div class="flex flex-col gap-1 mb-4">
        <label>Video URL</label>
        <input
          placeholder="https://..."
          class="border border-default rounded px-3 py-1.5 bg-default"
        />
      </div>
      <UButton type="submit" color="primary">Start</UButton>
    </form>
  </div>
</template>
