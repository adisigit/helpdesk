<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  id: string
  autocomplete?: string
  disabled?: boolean
  error?: string
}>()

const model = defineModel<string>({ required: true })
const visible = ref(false)
</script>

<template>
  <div class="relative">
    <input
      :id="id"
      v-model="model"
      :type="visible ? 'text' : 'password'"
      :autocomplete="autocomplete"
      :disabled="disabled"
      :aria-invalid="!!error"
      :aria-describedby="error ? `${id}-error` : undefined"
      class="input pr-10"
    />
    <button
      type="button"
      class="absolute inset-y-0 right-0 flex items-center px-3 text-slate-500 hover:text-slate-700 focus:outline-none focus-visible:text-blue-600"
      :aria-label="visible ? 'Sembunyikan password' : 'Tampilkan password'"
      :aria-pressed="visible"
      :disabled="disabled"
      @click="visible = !visible"
    >
      <svg v-if="!visible" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M2.04 12.32a1 1 0 0 1 0-.64C3.42 7.51 7.36 4.5 12 4.5c4.64 0 8.57 3.01 9.96 7.18.07.21.07.43 0 .64C20.58 16.49 16.64 19.5 12 19.5c-4.64 0-8.57-3.01-9.96-7.18Z" />
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
      </svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.22A10.5 10.5 0 0 0 2.04 12c1.39 4.18 5.32 7.19 9.96 7.19 1.99 0 3.84-.55 5.42-1.51M6.23 6.23A10.45 10.45 0 0 1 12 4.5c4.64 0 8.57 3.01 9.96 7.18a10.52 10.52 0 0 1-4.29 5.46M6.23 6.23 3 3m3.23 3.23 3.65 3.65m7.79 7.79L21 21m-3.33-3.33-3.55-3.55m0 0a3 3 0 0 0-4.24-4.24m4.24 4.24L9.88 9.88" />
      </svg>
    </button>
  </div>
</template>
