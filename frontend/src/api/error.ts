import axios from 'axios'

interface HumaErrorBody {
  title?: string
  detail?: string
  errors?: { location?: string; message?: string }[]
}

export interface ParsedError {
  message: string
  fields: Record<string, string>
  status: number
}

export function parseError(err: unknown): ParsedError {
  if (!axios.isAxiosError<HumaErrorBody>(err) || !err.response) {
    return { message: 'Tidak dapat terhubung ke server. Periksa koneksi Anda.', fields: {}, status: 0 }
  }
  const data = err.response.data ?? {}
  const fields: Record<string, string> = {}
  for (const e of data.errors ?? []) {
    const m = /^body\.(.+)$/.exec(e.location ?? '')
    const key = m?.[1]
    if (key && e.message && !fields[key]) fields[key] = e.message
  }
  const message = data.detail || data.title || `Terjadi kesalahan (${err.response.status})`
  return { message, fields, status: err.response.status }
}
