<script setup lang="ts" generic="T extends Record<string, unknown>">
import UForm from '@nuxt/ui/components/Form.vue'
import type { FormError, FormErrorEvent } from '@nuxt/ui'
import { useToast } from '@nuxt/ui/composables'
import { reactive, useTemplateRef, watch } from 'vue'
import { getApiErrorMessage } from '@/api/apiErrors'

interface FieldError {
  field?: string
  message: string
}

const props = withDefaults(
  defineProps<{
    initialValues?: T
    onSubmit: (values: T) => Promise<void> | void
    onValidate?: (values: T) => FieldError[]
    fieldErrors?: (error: unknown) => FieldError[]
    errorTitle?: string
  }>(),
  {
    initialValues: () => ({}) as T,
    errorTitle: 'Error',
  },
)

const state = reactive({ ...props.initialValues }) as T
const form = useTemplateRef('form')
const formElement = useTemplateRef<HTMLElement>('formElement')
const toast = useToast()

watch(
  () => ({ ...state }),
  (next, previous) => {
    for (const key of Object.keys(next)) {
      if (next[key] !== previous[key]) form.value?.clear(key)
    }
  },
)

function toFormErrors(errors: FieldError[]): FormError[] {
  return errors.map((e) => ({ name: e.field, message: e.message }))
}

function validate(state: T): FormError[] {
  return toFormErrors(props.onValidate?.(state) ?? [])
}

function focusLater(id?: string) {
  setTimeout(() => {
    const target = id
      ? document.getElementById(id)
      : formElement.value?.querySelector<HTMLElement>('input, textarea, button[role="combobox"]')
    target?.focus()
  })
}

async function submit(event: { data: T }) {
  try {
    await props.onSubmit(event.data)
  } catch (err) {
    const fieldErrors = props.fieldErrors?.(err) ?? []
    if (fieldErrors.length > 0 && form.value) {
      form.value.setErrors(toFormErrors(fieldErrors))
      focusLater(form.value.getErrors()[0]?.id)
      return
    }

    toast.add({
      title: props.errorTitle,
      description: getApiErrorMessage(err, 'Something went wrong'),
      color: 'error',
      duration: 0,
    })
    focusLater()
  }
}

function onError({ errors }: FormErrorEvent) {
  focusLater(errors.find((e) => e.id)?.id)
}
</script>

<template>
  <div ref="formElement">
    <UForm
      ref="form"
      :state="state"
      :validate="validate"
      :validate-on="[]"
      loading-auto
      novalidate
      class="space-y-4"
      @submit="submit"
      @error="onError"
    >
      <slot :state="state" :errors="form?.errors ?? []" :loading="form?.loading" />
    </UForm>
  </div>
</template>
