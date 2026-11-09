<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth';
import AppHeader from './components/AppHeader.vue';

const auth = useAuthStore()
</script>

<template>
  <AppHeader v-if="auth.isAuthenticated" />
  <main class="mx-auto max-w-6xl px-4 py-6">
    <RouterView />
  </main>
</template>
