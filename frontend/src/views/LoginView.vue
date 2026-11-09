<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter, type RouteLocationRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { parseError } from '../api/error'
import AlertBox from '../components/AlertBox.vue'
import FormField from '../components/FormField.vue'
import PasswordInput from '../components/PasswordInput.vue'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const form = reactive({ email: '', password: '' })
const errors = reactive({ email: '', password: '' })
const serverError = ref('')
const loading = ref(false)

function validate(): boolean {
  const email = form.email.trim()
  errors.email = !email
    ? 'Email wajib diisi'
    : !/^\S+@\S+\.\S+$/.test(email)
      ? 'Format email tidak valid'
      : ''
  errors.password = form.password ? '' : 'Password wajib diisi'
  return !errors.email && !errors.password
}

function safeRedirect(): RouteLocationRaw {
  const r = route.query.redirect
  return typeof r === 'string' && r.startsWith('/') && !r.startsWith('//') ? r : { name: 'tickets' }
}

async function submit() {
  serverError.value = ''
  if (!validate()) return
  loading.value = true
  try {
    await auth.login(form.email.trim(), form.password)
    router.replace(safeRedirect())
  } catch (e) {
    serverError.value = parseError(e).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto mt-10 w-full max-w-sm">
    <h1 class="mb-6 text-center text-2xl font-semibold">Masuk ke Helpdesk</h1>

    <form class="card space-y-4 p-6" novalidate @submit.prevent="submit">
      <AlertBox v-if="route.query.expired" type="info">Sesi Anda berakhir. Silakan masuk kembali.</AlertBox>
      <AlertBox v-if="serverError">{{ serverError }}</AlertBox>

      <FormField id="email" label="Email" :error="errors.email">
        <input id="email" v-model="form.email" type="email" autocomplete="username" class="input" :disabled="loading" />
      </FormField>

      <FormField id="password" label="Password" :error="errors.password">
        <PasswordInput
            id="password"
            v-model="form.password"
            autocomplete="current-password"
            :disabled="loading"
            :error="errors.password"
        />
      </FormField>

      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? 'Memproses...' : 'Masuk' }}
      </button>
    </form>

    <p class="mt-4 text-center text-sm text-slate-600">
      Belum punya akun?
      <RouterLink to="/register" class="font-medium text-blue-600 hover:underline">Daftar</RouterLink>
    </p>
  </div>
</template>
