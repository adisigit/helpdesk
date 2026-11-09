<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '../stores/auth';
import type { Priority, TicketDetail, TicketStatus, UpdateTicketPayload, User } from '../types';
import { CATEGORIES, PRIORITY, STATUS, toOptions } from '../constants';
import { addComment, getTicket, updateTicket } from '../api/tickets';
import { parseError } from '../api/error';
import { formatDate } from '../utils/format';
import AlertBox from '../components/AlertBox.vue';
import StatusBadge from '../components/StatusBadge.vue';
import PriorityBadge from '../components/PriorityBadge.vue';
import ErrorState from '../components/ErrorState.vue';
import LoadingState from '../components/LoadingState.vue';
import { listAssignees } from '../api/users';

const props = defineProps<{ id: string }>()
const auth = useAuthStore()

const ticket = ref<TicketDetail | null>(null)
const loading = ref(true)
const loadError = ref('')
const actionError = ref('')
const notice = ref('')
const acting = ref(false)

const assignees = ref<User[]>([])

const canTake = computed(() => {
  const t = ticket.value
  if (!t || isClosed.value) return false
  return auth.isAdmin ? t.assignee?.id !== auth.user?.id : auth.isStaff && !t.assignee
})

const priorityOptions = toOptions(PRIORITY)
const isOwner = computed(() => !!ticket.value && ticket.value.requester.id === auth.user?.id)
const isClosed = computed(() => ticket.value?.status === 'closed')
const categoryLabel = computed(
  () => CATEGORIES.find((c) => c.value === ticket.value?.category)?.label ?? ticket.value?.category ?? '-',
)

async function load(silent = false) {
  if (!silent) loading.value = true
  loadError.value = ''
  try {
    const { data } = await getTicket(props.id)
    ticket.value = data
    if (auth.isAdmin && assignees.value.length === 0) {
      try {
        assignees.value = (await listAssignees()).data.items
      } catch {

      }
    }
  } catch (e) {
    loadError.value = parseError(e).message
  } finally {
    loading.value = false
  }
}

watch(() => props.id, () => load(), { immediate: true })

interface ActionUi {
  label: string
  cls: string
  confirm?: string
}

const ACTION_UI: Record<string, ActionUi> = {
  'open>in_progress': { label: 'Mulai tangani', cls: 'btn-primary' },
  'in_progress>resolved': { label: 'Tandai selesai', cls: 'btn-primary' },
  'resolved>closed': {
    label: 'Tutup tiket',
    cls: 'btn-primary',
    confirm: 'Tutup tiket ini? Tiket yang ditutup tidak dapat dibuka kembali.',
  },
  'resolved>open': { label: 'Buka kembali', cls: 'btn-secondary' },
  'open>closed': {
    label: 'Batalkan tiket',
    cls: 'btn-danger',
    confirm: 'Batalkan tiket ini? Tiket yang ditutup tidak dapat dibuka kembali.',
  },
}

const actions = computed(() => {
  const t = ticket.value
  if (!t) return []
  return t.allowed_transitions.map((to) => {
    const ui: ActionUi = ACTION_UI[`${t.status}>${to}`] ?? { label: STATUS[to].label, cls: 'btn-secondary' }
    return { to, ...ui }
  })
})

async function patch(payload: UpdateTicketPayload, success: string): Promise<boolean> {
  if (!ticket.value || acting.value) return false
  acting.value = true
  actionError.value = ''
  notice.value = ''
  try {
    const { data } = await updateTicket(ticket.value.id, payload)
    ticket.value = { ...data, comments: ticket.value.comments }
    notice.value = success
    return true
  } catch (e) {
    const p = parseError(e)
    actionError.value = p.message
    if (p.status === 409) await load(true)
    return false
  } finally {
    acting.value = false
  }
}

async function runAction(a: { to: TicketStatus; confirm?: string }) {
  if (a.confirm && !window.confirm(a.confirm)) return
  await patch({ status: a.to }, 'Status diperbarui.')
}

async function onPriority(ev: Event) {
  const sel = ev.target as HTMLSelectElement
  const ok = await patch({ priority: sel.value as Priority }, 'Prioritas diperbarui.')
  if (!ok && ticket.value) sel.value = ticket.value.priority
}

async function onAssignee(ev: Event) {
  const sel = ev.target as HTMLSelectElement
  if (!sel.value) return
  const ok = await patch({ assignee_id: sel.value }, 'Assignee diperbarui.')
  if (!ok) sel.value = ticket.value?.assignee?.id ?? ''
}

function takeTicket() {
  if (auth.user) patch({ assignee_id: auth.user.id }, 'Tiket ini sekarang ditugaskan ke Anda.')
}

const commentBody = ref('')
const commentInternal = ref(false)
const commentError = ref('')
const posting = ref(false)

async function submitComment() {
  commentError.value = ''
  const body = commentBody.value.trim()
  if (!body) {
    commentError.value = 'Komentar tidak boleh kosong'
    return
  }
  if (body.length > 5000) {
    commentError.value = 'Komentar maksimal 5000 karakter'
    return
  }
  if (!ticket.value) return
  posting.value = true
  try {
    const { data } = await addComment(ticket.value.id, body, auth.isStaff && commentInternal.value)
    ticket.value.comments.push(data)
    commentBody.value = ''
    commentInternal.value = false
  } catch (e) {
    commentError.value = parseError(e).message
  } finally {
    posting.value = false
  }
}
</script>

<template>
  <div class="space-y-4">
    <RouterLink to="/tickets" class="text-sm text-slate-600 hover:underline">Kembali ke daftar</RouterLink>

    <LoadingState v-if="loading" />
    <ErrorState v-else-if="loadError || !ticket" :message="loadError || 'Tiket tidak ditemukan'" @retry="load()" />

    <div v-else class="grid gap-6 lg:grid-cols-3">
      <div class="space-y-6 lg:col-span-2">
        <section class="card p-6">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <h1 class="text-xl font-semibold wrap-break-words">{{ ticket.title }}</h1>
            <div class="flex gap-2">
              <StatusBadge :status="ticket.status" />
              <PriorityBadge :priority="ticket.priority" />
            </div>
          </div>
          <p class="mt-4 whitespace-pre-wrap wrap-break-words text-sm text-slate-700">{{ ticket.description }}</p>
        </section>

        <section class="space-y-3">
          <h2 class="font-semibold">Komentar ({{ ticket.comments.length }})</h2>

          <p v-if="ticket.comments.length === 0" class="card px-4 py-8 text-center text-sm text-slate-500">
            Belum ada komentar.
          </p>

          <ul v-else class="space-y-3">
            <li
              v-for="c in ticket.comments"
              :key="c.id"
              :class="['card p-4', c.is_internal && 'border-amber-300 bg-amber-50']"
            >
              <div class="flex flex-wrap items-center gap-2 text-sm">
                <span class="font-medium">{{ c.author.name }}</span>
                <span class="rounded bg-slate-100 px-1.5 py-0.5 text-xs uppercase text-slate-600">
                  {{ c.author.role }}
                </span>
                <span
                  v-if="c.is_internal"
                  class="rounded bg-amber-200 px-1.5 py-0.5 text-xs font-medium text-amber-900"
                >
                  Catatan internal
                </span>
                <span class="text-xs text-slate-500">{{ formatDate(c.created_at) }}</span>
              </div>
              <p class="mt-2 whitespace-pre-wrap wrap-break-words text-sm text-slate-800">{{ c.body }}</p>
            </li>
          </ul>

          <form class="card space-y-3 p-4" novalidate @submit.prevent="submitComment">
            <AlertBox v-if="commentError">{{ commentError }}</AlertBox>
            <label for="comment" class="sr-only">Tulis komentar</label>
            <textarea
              id="comment"
              v-model="commentBody"
              rows="3"
              maxlength="5000"
              class="input"
              placeholder="Tulis komentar..."
              :disabled="posting"
            ></textarea>
            <div class="flex flex-wrap items-center justify-between gap-3">
              <label v-if="auth.isStaff" class="flex items-center gap-2 text-sm text-slate-700">
                <input v-model="commentInternal" type="checkbox" class="rounded border-slate-300" />
                Catatan internal (tidak terlihat oleh requester)
              </label>
              <span v-else></span>
              <button type="submit" class="btn-primary" :disabled="posting">
                {{ posting ? 'Mengirim...' : 'Kirim' }}
              </button>
            </div>
          </form>
        </section>
      </div>

      <aside class="space-y-4">
        <AlertBox v-if="actionError">{{ actionError }}</AlertBox>
        <AlertBox v-else-if="notice" type="success">{{ notice }}</AlertBox>

        <section v-if="actions.length" class="card space-y-2 p-4">
          <h2 class="text-sm font-semibold">Tindakan</h2>
          <button
            v-for="a in actions"
            :key="a.to"
            type="button"
            :class="[a.cls, 'w-full']"
            :disabled="acting"
            @click="runAction(a)"
          >
            {{ a.label }}
          </button>
        </section>

        <section class="card space-y-4 p-4 text-sm">
          <h2 class="font-semibold">Detail</h2>

          <div>
            <p class="text-xs text-slate-500">Prioritas</p>
            <select
              v-if="auth.isStaff"
              class="input mt-1"
              aria-label="Ubah prioritas"
              :disabled="acting || isClosed"
              @change="onPriority"
            >
              <option v-for="o in priorityOptions" :key="o.value" :value="o.value" :selected="o.value === ticket.priority">
                {{ o.label }}
              </option>
            </select>
            <div v-else class="mt-1"><PriorityBadge :priority="ticket.priority" /></div>
          </div>

          <div>
            <p class="text-xs text-slate-500">Assignee</p>
            <select
              v-if="auth.isAdmin"
              class="input mt-1"
              aria-label="Ubah assignee"
              :disabled="acting || isClosed"
              @change="onAssignee"
            >
              <option value="" disabled :selected="!ticket.assignee">Belum ditugaskan</option>
              <option v-for="u in assignees" :key="u.id" :value="u.id" :selected="u.id === ticket.assignee?.id">
                {{ u.name }} ({{ u.role }})
              </option>
            </select>
            <p v-else class="mt-1">{{ ticket.assignee?.name ?? 'Belum ditugaskan' }}</p>
            <button v-if="canTake" type="button" class="btn-secondary mt-2 w-full" :disabled="acting" @click="takeTicket">
              Ambil tiket
            </button>
          </div>

          <div>
            <p class="text-xs text-slate-500">Kategori</p>
            <p class="mt-1">{{ categoryLabel }}</p>
          </div>
          <div>
            <p class="text-xs text-slate-500">Pembuat</p>
            <p class="mt-1">{{ ticket.requester.name }}<span v-if="isOwner" class="text-slate-500"> (Anda)</span></p>
          </div>
          <div>
            <p class="text-xs text-slate-500">Dibuat</p>
            <p class="mt-1">{{ formatDate(ticket.created_at) }}</p>
          </div>
          <div>
            <p class="text-xs text-slate-500">Diperbarui</p>
            <p class="mt-1">{{ formatDate(ticket.updated_at) }}</p>
          </div>
        </section>
      </aside>
    </div>
  </div>
</template>
