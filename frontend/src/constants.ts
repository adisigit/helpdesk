import type { Priority, TicketStatus } from "./types"

interface Badge {
  label: string
  cls: string
}

export const STATUS: Record<TicketStatus, Badge> = {
  open: { label: 'Open', cls: 'bg-blue-100 text-blue-800' },
  in_progress: { label: 'In Progress', cls: 'bg-amber-100 text-amber-800' },
  resolved: { label: 'Resolved', cls: 'bg-green-100 text-green-800' },
  closed: { label: 'Closed', cls: 'bg-slate-200 text-slate-700' },
}

export const PRIORITY: Record<Priority, Badge> = {
  low: { label: 'Low', cls: 'bg-slate-100 text-slate-700' },
  medium: { label: 'Medium', cls: 'bg-sky-100 text-sky-800' },
  high: { label: 'High', cls: 'bg-orange-100 text-orange-800' },
  urgent: { label: 'Urgent', cls: 'bg-red-100 text-red-800' },
}

export const CATEGORIES = [
  { value: 'hardware', label: 'Hardware' },
  { value: 'software', label: 'Software' },
  { value: 'network', label: 'Jaringan' },
  { value: 'account', label: 'Akun' },
  { value: 'other', label: 'Lainnya' },
] as const

export function toOptions<K extends string>(map: Record<K, { label: string }>): { value: K; label: string }[] {
  return (Object.entries(map) as [K, { label: string }][]).map(([value, v]) => ({ value, label: v.label }))
}
