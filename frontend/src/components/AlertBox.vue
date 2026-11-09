<script setup lang="ts">
type AlertType = 'error' | 'success' | 'info'

withDefaults(defineProps<{ type?: AlertType }>(), { type: 'error' })

const styles: Record<AlertType, string> = {
  error: 'border-red-200 bg-red-50 text-red-800',
  success: 'border-green-200 bg-green-50 text-green-800',
  info: 'border-blue-200 bg-blue-50 text-blue-800',
}
</script>

<template>
    <div :class="['rounded-md border px-3 py-2 text-sm', styles[type]]" role="alert">
        <slot />
    </div>
</template>
