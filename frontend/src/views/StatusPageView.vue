<script setup lang="ts">
import { RouterLink } from 'vue-router'

const props = withDefaults(defineProps<{ code?: number }>(), { code: 404 })

const text: Record<number, [string, string]> = {
  403: ['Akses ditolak', 'Anda tidak memiliki hak untuk membuka halaman ini.'],
  404: ['Halaman tidak ditemukan', 'Alamat yang Anda buka tidak tersedia.'],
}
const [title, desc] = text[props.code] ?? text[404]!
</script>

<template>
  <div class="mt-16 text-center">
    <p class="text-5xl font-bold text-slate-300">{{ code }}</p>
    <h1 class="mt-3 text-xl font-semibold">{{ title }}</h1>
    <p class="mt-1 text-slate-600">{{ desc }}</p>
    <RouterLink to="/tickets" class="btn-primary mt-6">Kembali ke daftar tiket</RouterLink>
  </div>
</template>
