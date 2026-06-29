export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export interface User {
  id: string;
  email: string;
  created_at: string;
}
