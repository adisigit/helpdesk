<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { parseError } from '../api/error'
import FormField from '../components/FormField.vue'
import AlertBox from '../components/AlertBox.vue'
import PasswordInput from '../components/PasswordInput.vue'

const auth = useAuthStore()
const router = useRouter()

const form = reactive({ name: '', email: '', password: '', confirm: '' })
const errors = reactive({ name: '', email: '', password: '', confirm: '' })
const serverError = ref('')
const loading = ref(false)

function validate(): boolean {
  const name = form.name.trim()
  const email = form.email.trim()
  errors.name = name.length < 2 ? 'Nama minimal 2 karakter' : ''
  errors.email = !email
    ? 'Email wajib diisi'
    : !/^\S+@\S+\.\S+$/.test(email)
      ? 'Format email tidak valid'
      : ''
  errors.password =
    form.password.length < 8
      ? 'Password minimal 8 karakter'
      : new TextEncoder().encode(form.password).length > 72
        ? 'Password maksimal 72 byte'
        : ''
  errors.confirm = form.confirm !== form.password ? 'Konfirmasi password tidak sama' : ''
  return !Object.values(errors).some(Boolean)
}

async function submit() {
  serverError.value = ''
  if (!validate()) return
  loading.value = true
  try {
    await auth.register({
      name: form.name.trim(),
      email: form.email.trim(),
      password: form.password,
    })
    router.replace({ name: 'tickets' })
  } catch (e) {
    const parsed = parseError(e)
    for (const key of ['name', 'email', 'password'] as const) {
      const msg = parsed.fields[key]
      if (msg) errors[key] = msg
    }
    serverError.value = parsed.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto mt-10 w-full max-w-sm">
    <h1 class="mb-6 text-center text-2xl font-semibold">Buat Akun</h1>

    <form class="card space-y-4 p-6" novalidate @submit.prevent="submit">
      <AlertBox v-if="serverError">{{ serverError }}</AlertBox>

      <FormField id="name" label="Nama" :error="errors.name">
        <input id="name" v-model="form.name" type="text" autocomplete="name" class="input" :disabled="loading" />
      </FormField>

      <FormField id="email" label="Email" :error="errors.email">
        <input id="email" v-model="form.email" type="email" autocomplete="email" class="input" :disabled="loading" />
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

      <FormField id="confirm" label="Ulangi password" :error="errors.confirm">
        <PasswordInput
            id="confirm"
            v-model="form.confirm"
            autocomplete="new-password"
            :disabled="loading"
            :error="errors.confirm"
        />
      </FormField>

      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? 'Memproses...' : 'Daftar' }}
      </button>
    </form>

    <p class="mt-4 text-center text-sm text-slate-600">
      Sudah punya akun?
      <RouterLink to="/login" class="font-medium text-blue-600 hover:underline">Masuk</RouterLink>
    </p>
  </div>
</template>
