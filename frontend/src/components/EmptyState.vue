<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="card px-6 py-14 text-center">
    <h3 class="font-medium text-slate-900">{{ title }}</h3>
    <p v-if="description" class="mt-1 text-sm text-slate-500">{{ description }}</p>
    <div class="mt-4"><slot /></div>
  </div>
</template>
