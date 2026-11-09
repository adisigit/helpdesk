import axios from "axios"

declare module 'axios' {
  interface AxiosRequestConfig {
    skipAuthRedirect?: boolean
  }
}

let onUnauthorized: () => void = () => { }

export function configureHttp(opts: { onUnauthorized?: () => void }): void {
  if (opts.onUnauthorized) onUnauthorized = opts.onUnauthorized
}

const http = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: { 'X-Requested-With': 'XMLHttpRequest' },
})

http.interceptors.response.use(
  (res) => res,
  (err: unknown) => {
    if (axios.isAxiosError(err) && err.response?.status === 401 && !err.config?.skipAuthRedirect) {
      onUnauthorized()
    }
    return Promise.reject(err)
  },
)

export default http
