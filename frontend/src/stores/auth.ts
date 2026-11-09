import { defineStore } from "pinia";
import { computed, ref } from "vue";
import * as authApi from '../api/auth'
import type { RegisterPayload, Role, User } from "../types";

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const ready = ref(false)

  const isAuthenticated = computed(() => !!user.value)
  const role = computed<Role | ''>(() => user.value?.role ?? '')
  const isStaff = computed(() => role.value === 'agent' || role.value === 'admin')
  const isAdmin = computed(() => role.value === 'admin')

  async function init(): Promise<void> {
    if (ready.value) return
    try {
      user.value = (await authApi.fetchMe()).data
    } catch {
      user.value = null
    } finally {
      ready.value = true
    }
  }

  async function login(email: string, password: string): Promise<void> {
    const { data } = await authApi.login(email, password)
    user.value = data.user
    ready.value = true
  }

  async function register(payload: RegisterPayload): Promise<void> {
    await authApi.register(payload)
    await login(payload.email, payload.password)
  }

  async function logout(): Promise<void> {
    try {
      await authApi.logout()
    } finally {
      clear()
    }
  }

  function clear(): void {
    user.value = null
  }

  return { user, ready, isAuthenticated, role, isStaff, isAdmin, init, login, register, logout, clear }
})
