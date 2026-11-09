<script setup lang="ts">
import { RouterLink } from 'vue-router'
import StatusBadge from './StatusBadge.vue'
import PriorityBadge from './PriorityBadge.vue'
import type { Ticket } from '../types';
import { formatDate } from '../utils/format';

defineProps<{ tickets: Ticket[] }>()
</script>

<template>
  <ul class="space-y-3 md:hidden">
    <li v-for="t in tickets" :key="t.id">
      <RouterLink :to="`/tickets/${t.id}`" class="card block p-4">
        <p class="font-medium text-slate-900">{{ t.title }}</p>
        <div class="mt-2 flex flex-wrap gap-2">
          <StatusBadge :status="t.status" />
          <PriorityBadge :priority="t.priority" />
        </div>
        <p class="mt-2 text-xs text-slate-500">
          {{ t.requester.name }} · {{ t.assignee ? t.assignee.name : 'Belum ditugaskan' }}
        </p>
        <p class="text-xs text-slate-400">{{ formatDate(t.created_at) }}</p>
      </RouterLink>
    </li>
  </ul>

  <div class="card hidden overflow-x-auto md:block">
    <table class="w-full text-left text-sm">
      <thead class="border-b border-slate-200 bg-slate-50 text-xs uppercase text-slate-500">
        <tr>
          <th class="px-4 py-3 font-medium">Judul</th>
          <th class="px-4 py-3 font-medium">Status</th>
          <th class="px-4 py-3 font-medium">Prioritas</th>
          <th class="px-4 py-3 font-medium">Pembuat</th>
          <th class="px-4 py-3 font-medium">Assignee</th>
          <th class="px-4 py-3 font-medium">Dibuat</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-100">
        <tr v-for="t in tickets" :key="t.id" class="hover:bg-slate-50">
          <td class="max-w-xs px-4 py-3">
            <RouterLink :to="`/tickets/${t.id}`" class="font-medium text-blue-600 hover:underline">
              {{ t.title }}
            </RouterLink>
          </td>
          <td class="px-4 py-3"><StatusBadge :status="t.status" /></td>
          <td class="px-4 py-3"><PriorityBadge :priority="t.priority" /></td>
          <td class="px-4 py-3 text-slate-600">{{ t.requester.name }}</td>
          <td class="px-4 py-3 text-slate-600">{{ t.assignee ? t.assignee.name : '-' }}</td>
          <td class="whitespace-nowrap px-4 py-3 text-slate-500">{{ formatDate(t.created_at) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
