import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"
import type { Role } from "../types"
import { useAuthStore } from "../stores/auth"
import StatusPageView from "../views/StatusPageView.vue"
import LoginView from "../views/LoginView.vue"
import RegisterView from "../views/RegisterView.vue"
import TicketListView from "../views/TicketListView.vue"

declare module 'vue-router' {
  interface RouteMeta {
    requireAuth?: boolean
    guestOnly?: boolean
    roles?: Role[]
  }
}

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: { name: 'tickets' } },
  {
    path: '/login',
    name: 'login',
    component: () => LoginView,
    meta: {guestOnly: true}
  },
  {
    path: '/register',
    name: 'register',
    component: () => RegisterView,
    meta: { guestOnly: true },
  },
  {
    path: '/tickets',
    name: 'tickets',
    component: () => TicketListView,
    meta: { requiresAuth: true },
  },
  {
    path: '/forbidden',
    name: 'forbidden',
    component: () => StatusPageView,
    props: { code: 403 },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => StatusPageView,
    props: { code: 404 },
  },
]

export const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.init()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'tickets' }
  }
  if (to.meta.roles && !(auth.role && to.meta.roles.includes(auth.role))) {
    return { name: 'forbidden' }
  }
})
