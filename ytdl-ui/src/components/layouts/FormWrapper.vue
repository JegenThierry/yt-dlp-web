<script setup lang="ts" generic="T extends Record<string, unknown>">
import UForm from '@nuxt/ui/components/Form.vue'
import type { FormError, FormErrorEvent } from '@nuxt/ui/runtime/types/form'
import { useToast } from '@nuxt/ui/composables'
import { reactive, useTemplateRef } from 'vue'

const props = withDefaults(
  defineProps<{
    initialValues?: T
    onSubmit: (values: T) => Promise<void> | void
    onValidate?: (values: T) => { field?: string; message: string }[]
  }>(),
  {
    initialValues: () => ({}) as T,
  },
)

const state = reactive({ ...props.initialValues }) as T
const form = useTemplateRef('form')
const toast = useToast()

function validate(state: T): FormError[] {
  return props.onValidate?.(state).map((e) => ({ name: e.field, message: e.message })) ?? []
}

async function submit(event: { data: T }) {
  try {
    await props.onSubmit(event.data)
  } catch (err) {
    toast.add({
      title: 'Error',
      description: err instanceof Error ? err.message : 'Something went wrong',
      color: 'error',
    })
  }
}

function onError({ errors }: FormErrorEvent) {
  toast.add({
    title: 'Validation Error',
    description: `Please fix ${errors.length} error${errors.length > 1 ? 's' : ''} below.`,
    color: 'error',
  })
}
</script>

<template>
  <UForm
    ref="form"
    :state="state"
    :validate="validate"
    :validate-on="[]"
    loading-auto
    class="space-y-4"
    @submit="submit"
    @error="onError"
  >
    <slot :state="state" :errors="form?.errors ?? []" :loading="form?.loading" />
  </UForm>
</template>
