<script setup lang="ts">
withDefaults(defineProps<{ label?: string }>(), { label: 'Memuat...' })
</script>

<template>
  <div class="flex items-center justify-center gap-3 py-16 text-slate-500" role="status">
    <span class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-blue-600"></span>
    <span class="text-sm">{{ label }}</span>
  </div>
</template>
