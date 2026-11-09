import type { AdminUser, Paginated, Role, User } from '../types'
import http from './http'

export interface ListUsersParams {
  page: number
  page_size: number
  q?: string
  role?: string
}

export const listUsers = (params: ListUsersParams) =>
  http.get<Paginated<AdminUser>>('/users', { params })

export const changeUserRole = (id: string, role: Role) =>
  http.patch<AdminUser>(`/users/${id}/role`, { role })

export const listAssignees = () =>
  http.get<{ items: User[] }>('/assignees')
