<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import type { ButtonType } from '@/types/buttonTypes.ts'

interface SubmitButtonProps {
  disabled?: boolean
  loading?: boolean
  type?: ButtonType
}

withDefaults(defineProps<SubmitButtonProps>(), {
  type: 'submit',
})
</script>

<template>
  <BaseButton
    :type="type"
    label="Submit"
    color="primary"
    :disabled="disabled"
    :loading="loading"
  />
</template>

<style scoped></style>
