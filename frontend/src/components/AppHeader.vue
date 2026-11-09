<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

async function logout() {
  await auth.logout()
  router.replace({ name: 'login' })
}
</script>

<template>
  <header class="border-b border-slate-200 bg-white">
    <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-3 px-4 py-3">
      <div class="flex items-center gap-6">
        <RouterLink to="/tickets" class="text-lg font-semibold text-blue-600">Helpdesk</RouterLink>
        <nav class="flex gap-4 text-sm text-slate-600">
          <RouterLink to="/tickets" class="hover:text-slate-900" active-class="font-medium text-blue-600">
            Tiket
          </RouterLink>
          <RouterLink
            v-if="auth.isAdmin"
            to="/admin/users"
            class="hover:text-slate-900"
            active-class="font-medium text-blue-600"
          >
            Pengguna
          </RouterLink>
        </nav>
      </div>
      <div class="flex items-center gap-3 text-sm">
        <span class="text-slate-600">
          {{ auth.user?.name }}
          <span class="ml-1 rounded bg-slate-100 px-1.5 py-0.5 text-xs uppercase text-slate-600">
            {{ auth.role }}
          </span>
        </span>
        <button type="button" class="btn-secondary" @click="logout">Keluar</button>
      </div>
    </div>
  </header>
</template>
