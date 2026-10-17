<script setup lang="ts">
import type { InputTypeHTMLAttribute } from 'vue'
import UFormField from '@nuxt/ui/components/FormField.vue'
import UInput from '@nuxt/ui/components/Input.vue'

defineOptions({ inheritAttrs: false })

interface InputTextProps {
  label?: string
  fieldUi?: { label?: string; help?: string }
  placeholder?: string
  required?: boolean
  help?: string
  icon?: string
  name?: string
  error?: string | boolean
  type?: InputTypeHTMLAttribute
}

withDefaults(defineProps<InputTextProps>(), {
  required: false,
  type: 'text',
  error: undefined,
})
const modelValue = defineModel<string>()
</script>

<template>
  <UFormField
    :label="label"
    :name="name"
    :required="required"
    :help="help"
    :error="error"
    :ui="fieldUi"
  >
    <UInput
      v-model="modelValue"
      :type="type"
      :placeholder="placeholder"
      :icon="icon"
      class="w-full"
      :ui="$slots.trailing ? { trailing: 'pe-1' } : undefined"
      v-bind="$attrs"
    >
      <template v-if="$slots.trailing" #trailing>
        <slot name="trailing" />
      </template>
    </UInput>
  </UFormField>
</template>
