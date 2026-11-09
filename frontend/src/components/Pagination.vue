<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  totalPages: number
  total: number
  pageSize: number
}>()
defineEmits<{ 'update:page': [value: number] }>()

const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1))
const to = computed(() => Math.min(props.page * props.pageSize, props.total))
</script>

<template>
  <div class="flex flex-col items-center justify-between gap-3 text-sm text-slate-600 sm:flex-row">
    <p>Menampilkan {{ from }}–{{ to }} dari {{ total }}</p>
    <div class="flex items-center gap-2">
      <button type="button" class="btn-secondary" :disabled="page <= 1" @click="$emit('update:page', page - 1)">
        Sebelumnya
      </button>
      <span>Halaman {{ page }} / {{ totalPages }}</span>
      <button
        type="button"
        class="btn-secondary"
        :disabled="page >= totalPages"
        @click="$emit('update:page', page + 1)"
      >
        Berikutnya
      </button>
    </div>
  </div>
</template>
