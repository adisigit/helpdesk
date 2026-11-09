<script setup lang="ts">
import { computed } from 'vue'
import type { TicketStatus } from '../types';
import { STATUS } from '../constants';

const props = defineProps<{ status: TicketStatus }>()
const cfg = computed(() => STATUS[props.status] ?? { label: props.status, cls: 'bg-slate-100 text-slate-700' })
</script>

<template>
  <span :class="['inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium', cfg.cls]">{{ cfg.label }}</span>
</template>
