<script setup lang="ts">
import UFormField from '@nuxt/ui/components/FormField.vue'
import USelect from '@nuxt/ui/components/Select.vue'

interface SelectOption {
  label: string
  value: string
}

interface InputSelectProps {
  label?: string
  placeholder?: string
  options: SelectOption[]
  required?: boolean
  help?: string
  icon?: string
}

withDefaults(defineProps<InputSelectProps>(), {
  required: false,
})

const modelValue = defineModel<string>()
</script>

<template>
  <UFormField :label="label" :required="required" :help="help">
    <USelect
      v-model="modelValue"
      :items="options"
      :placeholder="placeholder"
      :icon="icon"
      class="w-full"
    />
  </UFormField>
</template>
