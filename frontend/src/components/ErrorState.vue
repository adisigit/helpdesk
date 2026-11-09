<script setup lang="ts">
defineProps<{ message: string }>()
defineEmits<{ retry: [] }>()
</script>

<template>
  <div class="card border-red-200 px-6 py-12 text-center" role="alert">
    <h3 class="font-medium text-red-700">Gagal memuat data</h3>
    <p class="mt-1 text-sm text-slate-600">{{ message }}</p>
    <button type="button" class="btn-secondary mt-4" @click="$emit('retry')">Coba lagi</button>
  </div>
</template>
