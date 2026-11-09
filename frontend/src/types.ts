export type Role = 'requester' | 'agent' | 'admin'
export type TicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed'
export type Priority = 'low' | 'medium' | 'high' | 'urgent'

export interface User {
  id: string
  name: string
  email: string
  role: Role
}

export interface Ticket {
  id: string
  title: string
  description: string
  category: string
  priority: Priority
  status: TicketStatus
  requester: User
  assignee?: User | null
  allowed_transitions: TicketStatus[]
  created_at: string
  updated_at: string
}

export interface Paginated<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface LoginResponse {
  access_token: string
  token_type: string
  expires_at: string
  user: User
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface TicketComment {
  id: string
  author: User
  body: string
  is_internal: boolean
  created_at: string
}

export interface TicketDetail extends Ticket {
  comments: TicketComment[]
}

export interface AdminUser extends User {
  created_at: string
}

export interface CreateTicketPayload {
  title: string
  description: string
  category: string
  priority: Priority
}

export interface UpdateTicketPayload {
  status?: TicketStatus
  priority?: Priority
  assignee_id?: string
}
