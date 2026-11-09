<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { CATEGORIES, PRIORITY, toOptions } from '../constants'
import type { Priority, Ticket } from '../types'
import { createTicket } from '../api/tickets'
import { parseError } from '../api/error'
import FormField from '../components/FormField.vue'
import AlertBox from '../components/AlertBox.vue'

const priorityOptions = toOptions(PRIORITY)

const form = reactive({ title: '', description: '', category: '', priority: 'medium' as Priority })
const errors = reactive({ title: '', description: '', category: '', priority: '' })
const serverError = ref('')
const created = ref<Ticket | null>(null)
const loading = ref(false)

function validate(): boolean {
  const title = form.title.trim()
  const desc = form.description.trim()
  errors.title =
    title.length < 3 ? 'Judul minimal 3 karakter' : title.length > 150 ? 'Judul maksimal 150 karakter' : ''
  errors.description =
    desc.length < 5 ? 'Deskripsi minimal 5 karakter' : desc.length > 5000 ? 'Deskripsi maksimal 5000 karakter' : ''
  errors.category = form.category ? '' : 'Pilih kategori'
  errors.priority = ''
  return !errors.title && !errors.description && !errors.category
}

function reset() {
  form.title = ''
  form.description = ''
  form.category = ''
  form.priority = 'medium'
  created.value = null
  serverError.value = ''
}

async function submit() {
  serverError.value = ''
  if (!validate()) return
  loading.value = true
  try {
    const { data } = await createTicket({
      title: form.title.trim(),
      description: form.description.trim(),
      category: form.category,
      priority: form.priority,
    })
    reset()
    created.value = data
  } catch (e) {
    const parsed = parseError(e)
    for (const key of ['title', 'description', 'category', 'priority'] as const) {
      const msg = parsed.fields[key]
      if (msg) errors[key] = msg
    }
    serverError.value = `Gagal membuat tiket: ${parsed.message}`
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-4">
    <RouterLink to="/tickets" class="text-sm text-slate-600 hover:underline">Kembali ke daftar</RouterLink>
    <h1 class="text-xl font-semibold">Buat tiket</h1>

    <div v-if="created" class="card space-y-4 p-6">
      <AlertBox type="success">Tiket berhasil dibuat.</AlertBox>
      <div class="flex flex-wrap gap-2">
        <RouterLink :to="`/tickets/${created.id}`" class="btn-primary">Lihat tiket</RouterLink>
        <button type="button" class="btn-secondary" @click="reset">Buat tiket lain</button>
        <RouterLink to="/tickets" class="btn-secondary">Ke daftar tiket</RouterLink>
      </div>
    </div>

    <form v-else class="card space-y-4 p-6" novalidate @submit.prevent="submit">
      <AlertBox v-if="serverError">{{ serverError }}</AlertBox>

      <FormField id="title" label="Judul" :error="errors.title">
        <input id="title" v-model="form.title" type="text" maxlength="150" class="input" :disabled="loading" />
      </FormField>

      <FormField id="description" label="Deskripsi" :error="errors.description">
        <textarea
          id="description"
          v-model="form.description"
          rows="6"
          maxlength="5000"
          class="input"
          :disabled="loading"
        ></textarea>
      </FormField>

      <div class="grid gap-4 sm:grid-cols-2">
        <FormField id="category" label="Kategori" :error="errors.category">
          <select id="category" v-model="form.category" class="input" :disabled="loading">
            <option value="" disabled>Pilih kategori</option>
            <option v-for="c in CATEGORIES" :key="c.value" :value="c.value">{{ c.label }}</option>
          </select>
        </FormField>

        <FormField id="priority" label="Prioritas" :error="errors.priority">
          <select id="priority" v-model="form.priority" class="input" :disabled="loading">
            <option v-for="o in priorityOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
          </select>
        </FormField>
      </div>

      <div class="flex justify-end gap-2">
        <RouterLink to="/tickets" class="btn-secondary">Batal</RouterLink>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? 'Mengirim...' : 'Kirim tiket' }}
        </button>
      </div>
    </form>
  </div>
</template>
