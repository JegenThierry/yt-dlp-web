<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'

interface SubmitButtonProps {
  disabled?: boolean
  loading?: boolean
}

defineProps<SubmitButtonProps>()
</script>

<template>
  <BaseButton
    type="submit"
    label="Submit"
    color="primary"
    :disabled="disabled"
    :loading="loading"
  />
</template>

<style scoped></style>
