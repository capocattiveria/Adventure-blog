import { LoginRequest, LoginResponse, User } from '@/types/auth';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    throw new Error('Credenziali non valide');
  }

  return res.json();
}

export async function register(data: LoginRequest): Promise<User> {
  const res = await fetch(`${BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    throw new Error('Registrazione fallita');
  }

  return res.json();
}
