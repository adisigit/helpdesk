<script setup lang="ts">
import { computed } from 'vue'
import type { Priority } from '../types';
import { PRIORITY } from '../constants';

const props = defineProps<{ priority: Priority }>()
const cfg = computed(() => PRIORITY[props.priority] ?? { label: props.priority, cls: 'bg-slate-100 text-slate-700' })
</script>

<template>
  <span :class="['inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium', cfg.cls]">{{ cfg.label }}</span>
</template>
