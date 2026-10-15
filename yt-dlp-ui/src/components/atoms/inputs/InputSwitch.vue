<script setup lang="ts">
import UFormField from '@nuxt/ui/components/FormField.vue'
import USwitch from '@nuxt/ui/components/Switch.vue'

interface InputSwitchProps {
  label?: string
  description?: string
}

defineProps<InputSwitchProps>()

const modelValue = defineModel<boolean>()
</script>

<template>
  <UFormField>
    <USwitch v-model="modelValue" :label="label" :description="description" />
  </UFormField>
</template>
