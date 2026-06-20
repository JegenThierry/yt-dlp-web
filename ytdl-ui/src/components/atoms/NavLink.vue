<script setup lang="ts">
import LinkButton from '@/components/atoms/buttons/LinkButton.vue'

defineProps<{
  to: string
  label: string
}>()
</script>

<template>
  <RouterLink :to="to" custom v-slot="{ href, navigate, isExactActive }">
    <LinkButton :href="href"
                :type="isExactActive ? 'primary' : 'default'"
                size="small"
                @click="navigate">
      {{ label }}
    </LinkButton>
  </RouterLink>
</template>
