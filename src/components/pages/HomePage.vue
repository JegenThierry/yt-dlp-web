<script setup lang="ts">
import { NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
</script>

<template>
  <n-card title="Download" size="small" style="max-width: 720px; margin: 0 auto;">
    <n-form>
      <n-form-item label="Video URL">
        <n-input placeholder="https://..." />
      </n-form-item>

      <n-button type="primary">Start</n-button>
    </n-form>
  </n-card>
</template>
