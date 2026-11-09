<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { PRIORITY, STATUS, toOptions } from '../constants'
import type { Ticket } from '../types'
import { ListTickets, type ListTicketsParams } from '../api/tickets'
import { parseError } from '../api/error'
import LoadingState from '../components/LoadingState.vue'
import ErrorState from '../components/ErrorState.vue'
import EmptyState from '../components/EmptyState.vue'
import TicketTable from '../components/TicketTable.vue'
import Pagination from '../components/Pagination.vue'

const auth = useAuthStore()

const statusOptions = toOptions(STATUS)
const priorityOptions = toOptions(PRIORITY)
const sortOptions = [
  { value: '-created_at', label: 'Terbaru' },
  { value: 'created_at', label: 'Terlama' },
  { value: '-updated_at', label: 'Baru diperbarui' },
  { value: '-priority', label: 'Prioritas tertinggi' },
  { value: 'priority', label: 'Prioritas terendah' },
  { value: 'title', label: 'Judul A–Z' },
]

const pageSize = 10
const filters = reactive({ q: '', status: '', priority: '', assignee: '', sort: '-created_at' })
const page = ref(1)

const tickets = ref<Ticket[]>([])
const total = ref(0)
const totalPages = ref(0)
const loading = ref(true)
const error = ref('')

const hasFilters = computed(() => !!(filters.q.trim() || filters.status || filters.priority || filters.assignee))

let seq = 0

async function load() {
  const mine = ++seq
  loading.value = true
  error.value = ''
  try {
    const params: ListTicketsParams = { page: page.value, page_size: pageSize, sort: filters.sort }
    for (const key of ['q', 'status', 'priority', 'assignee'] as const) {
      const v = filters[key].trim()
      if (v) params[key] = v
    }
    const { data } = await ListTickets(params)
    if (mine !== seq) return
    tickets.value = data.items
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
watch(() => [filters.status, filters.priority, filters.assignee, filters.sort], resetAndLoad)
watch(page, load)

function clearFilters() {
  filters.q = ''
  filters.status = ''
  filters.priority = ''
  filters.assignee = ''
}

onMounted(load)
onBeforeUnmount(() => clearTimeout(timer))
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-xl font-semibold">{{ auth.isStaff ? 'Semua tiket' : 'Tiket saya' }}</h1>
      <RouterLink to="/tickets/new" class="btn-primary">Buat tiket</RouterLink>
    </div>

    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
      <input v-model="filters.q" type="search" class="input lg:col-span-2" placeholder="Cari judul..." aria-label="Cari judul" />
      <select v-model="filters.status" class="input" aria-label="Filter status">
        <option value="">Semua status</option>
        <option v-for="o in statusOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
      <select v-model="filters.priority" class="input" aria-label="Filter prioritas">
        <option value="">Semua prioritas</option>
        <option v-for="o in priorityOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
      <select v-if="auth.isStaff" v-model="filters.assignee" class="input" aria-label="Filter assignee">
        <option value="">Semua assignee</option>
        <option value="me">Ditugaskan ke saya</option>
        <option value="none">Belum ditugaskan</option>
      </select>
      <select v-model="filters.sort" class="input" aria-label="Urutkan">
        <option v-for="o in sortOptions" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
    </div>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="error" :message="error" @retry="load" />
    <EmptyState
      v-else-if="tickets.length === 0"
      :title="hasFilters ? 'Tidak ada tiket yang cocok' : 'Belum ada tiket'"
      :description="hasFilters ? 'Ubah atau hapus filter untuk melihat tiket lain.' : 'Tiket yang dibuat akan muncul di sini.'"
    >
      <button v-if="hasFilters" type="button" class="btn-secondary" @click="clearFilters">Hapus filter</button>
      <RouterLink v-else to="/tickets/new" class="btn-primary">Buat tiket pertama</RouterLink>
    </EmptyState>
    <template v-else>
      <TicketTable :tickets="tickets" />
      <Pagination v-model:page="page" :total-pages="totalPages" :total="total" :page-size="pageSize" />
    </template>
  </div>
</template>
