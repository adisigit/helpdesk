
import type { LoginResponse, RegisterPayload, User } from "../types"
import http from "./http"

const quiet = { skipAuthRedirect: true }

export const login = (email: string, password: string) =>
  http.post<LoginResponse>('/auth/login', { email, password }, quiet)

export const register = (payload: RegisterPayload) =>
  http.post<User>('/auth/register', payload, quiet)

export const fetchMe = () => http.get<User>('/me', quiet)

export const logout = () => http.post('/auth/logout', null, quiet)
