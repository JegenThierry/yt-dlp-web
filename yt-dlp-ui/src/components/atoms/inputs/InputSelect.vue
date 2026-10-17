<script setup lang="ts">
import UFormField from '@nuxt/ui/components/FormField.vue'
import USelect from '@nuxt/ui/components/Select.vue'

defineOptions({ inheritAttrs: false })

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
  name?: string
  error?: string | boolean
}

withDefaults(defineProps<InputSelectProps>(), {
  required: false,
  error: undefined,
})

const modelValue = defineModel<string>()
</script>

<template>
  <UFormField :label="label" :name="name" :required="required" :help="help" :error="error">
    <USelect
      v-model="modelValue"
      :items="options"
      :placeholder="placeholder"
      :icon="icon"
      class="w-full"
      v-bind="$attrs"
    />
  </UFormField>
</template>
