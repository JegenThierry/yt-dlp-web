<script setup lang="ts">
import UFormField from '@nuxt/ui/components/FormField.vue'
import UInput from '@nuxt/ui/components/Input.vue'

interface InputTextProps {
  label?: string
  placeholder?: string
  required?: boolean
  help?: string
  icon?: string
}

withDefaults(defineProps<InputTextProps>(), {
  required: true,
})
const modelValue = defineModel<string>()
</script>

<template>
  <UFormField :label="label" :required="required" :help="help">
    <UInput v-model="modelValue" :placeholder="placeholder" :icon="icon" class="w-full" />
  </UFormField>
</template>
