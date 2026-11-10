<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { ROLES, toOptions } from '../constants'
import type { AdminUser, Role } from '../types'
import { changeUserRole, listUsers, type ListUsersParams } from '../api/users'
import { parseError } from '../api/error'
import { formatDate } from '../utils/format'
import AlertBox from '../components/AlertBox.vue'
import LoadingState from '../components/LoadingState.vue'
import ErrorState from '../components/ErrorState.vue'
import EmptyState from '../components/EmptyState.vue'
import Pagination from '../components/Pagination.vue'


const auth = useAuthStore()
const router = useRouter()
const roleOptions = toOptions(ROLES)

const pageSize = 10
const filters = reactive({ q: '', role: '' })
const page = ref(1)

const users = ref<AdminUser[]>([])
const total = ref(0)
const totalPages = ref(0)
const loading = ref(true)
const error = ref('')

const savingId = ref('')
const actionError = ref('')
const notice = ref('')

const hasFilters = computed(() => !!(filters.q.trim() || filters.role))

let seq = 0

async function load() {
  const mine = ++seq
  loading.value = true
  error.value = ''
  try {
    const params: ListUsersParams = { page: page.value, page_size: pageSize }
    if (filters.q.trim()) params.q = filters.q.trim()
    if (filters.role) params.role = filters.role
    const { data } = await listUsers(params)
    if (mine !== seq) return
    users.value = data.items
    total.value = data.total
    totalPages.value = data.total_pages
    if (data.items.length === 0 && page.value > 1 && data.total_pages > 0) {
      page.value = data.total_pages
    }
  } catch (e) {
    if (mine !== seq) return
    error.value = parseError(e).message
  } finally {
    if (mine === seq) loading.value = false
  }
}

function resetAndLoad() {
  if (page.value === 1) load()
  else page.value = 1
}

let timer: ReturnType<typeof setTimeout> | undefined
watch(
  () => filters.q,
  () => {
    clearTimeout(timer)
    timer = setTimeout(resetAndLoad, 400)
  },
)
watch(() => filters.role, resetAndLoad)
watch(page, load)

function clearFilters() {
  filters.q = ''
  filters.role = ''
}

async function changeRole(u: AdminUser, ev: Event) {
  const sel = ev.target as HTMLSelectElement
  const next = sel.value as Role
  if (next === u.role) return

  const isSelf = u.id === auth.user?.id
  const msg = isSelf
    ? 'Anda mengubah peran Anda sendiri. Setelah ini Anda akan keluar dan harus masuk kembali. Lanjutkan?'
    : `Ubah peran ${u.name} menjadi ${ROLES[next].label}?`
  if (!window.confirm(msg)) {
    sel.value = u.role
    return
  }

  savingId.value = u.id
  actionError.value = ''
  notice.value = ''
  try {
    const { data } = await changeUserRole(u.id, next)
    u.role = data.role
    if (isSelf) {
      await auth.logout()
      router.replace({ name: 'login' })
      return
    }
    notice.value = `Peran ${u.name} sekarang ${ROLES[data.role].label}.`
  } catch (e) {
    sel.value = u.role
    actionError.value = parseError(e).message
  } finally {
    savingId.value = ''
  }
}

onMounted(load)
onBeforeUnmount(() => clearTimeout(timer))
</script>

<template>
  <div class="space-y-4">
    <h1 class="text-xl font-semibold">Kelola pengguna</h1>

    <div class="grid gap-3 sm:grid-cols-3">
      <input v-model="filters.q" type="search" class="input sm:col-span-2" placeholder="Cari nama atau email..." aria-label="Cari pengguna" />
      <select v-model="filters.role" class="input" aria-label="Filter peran">
        <option value="">Semua peran</option>
        <option v-for="o in roleOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
    </div>

    <AlertBox v-if="actionError">{{ actionError }}</AlertBox>
    <AlertBox v-else-if="notice" type="success">{{ notice }}</AlertBox>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState
      v-else-if="users.length === 0"
      :title="hasFilters ? 'Tidak ada pengguna yang cocok' : 'Belum ada pengguna'"
    >
      <button v-if="hasFilters" type="button" class="btn-secondary" @click="clearFilters">Hapus filter</button>
    </EmptyState>
    <template v-else>
      <div class="card overflow-x-auto">
        <table class="w-full min-w-lg text-left text-sm">
          <thead class="border-b border-slate-200 bg-slate-50 text-xs uppercase text-slate-500">
            <tr>
              <th class="px-4 py-3 font-medium">Nama</th>
              <th class="px-4 py-3 font-medium">Email</th>
              <th class="px-4 py-3 font-medium">Peran</th>
              <th class="px-4 py-3 font-medium">Terdaftar</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr v-for="u in users" :key="u.id">
              <td class="px-4 py-3 font-medium">
                {{ u.name }}<span v-if="u.id === auth.user?.id" class="font-normal text-slate-500"> (Anda)</span>
              </td>
              <td class="px-4 py-3 text-slate-600">{{ u.email }}</td>
              <td class="px-4 py-3">
                <select
                  class="input w-36"
                  :aria-label="`Peran ${u.name}`"
                  :disabled="savingId === u.id"
                  @change="changeRole(u, $event)"
                >
                  <option v-for="o in roleOptions" :key="o.value" :value="o.value" :selected="o.value === u.role">
                    {{ o.label }}
                  </option>
                </select>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-slate-500">{{ formatDate(u.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination v-model:page="page" :total-pages="totalPages" :total="total" :page-size="pageSize" />
    </template>
  </div>
</template>
