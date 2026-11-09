const fmt = new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium', timeStyle: 'short' })

export function formatDate(iso?: string | null): string {
  return iso ? fmt.format(new Date(iso)) : '-'
}
