export interface AuthUser {
  id: number
  username: string
}

export interface AuthResponse {
  token: string
  user: AuthUser
}

export interface AuthFormData {
  username: string
  password: string
}