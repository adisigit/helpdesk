import type { CreateTicketPayload, Paginated, Ticket, TicketComment, TicketDetail, UpdateTicketPayload } from "../types"
import http from "./http"

export interface ListTicketsParams {
  page: number
  page_size: number
  sort: string
  q?: string
  status?: string
  priority?: string
  assignee?: string
}

export const ListTickets = (params: ListTicketsParams) =>
  http.get<Paginated<Ticket>>('/tickets', { params })

export const createTicket = (payload: CreateTicketPayload) =>
  http.post<Ticket>('/tickets', payload)

export const getTicket = (id: string) => http.get<TicketDetail>(`/tickets/${id}`)

export const updateTicket = (id: string, payload: UpdateTicketPayload) =>
  http.patch<Ticket>(`/tickets/${id}`, payload)

export const addComment = (id: string, body: string, isInternal: boolean) =>
  http.post<TicketComment>(`/tickets/${id}/comments`, { body, is_internal: isInternal })
