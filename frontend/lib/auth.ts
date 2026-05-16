import { api, setToken, clearToken } from './api';
import type { Trainer } from '@/types';

export interface LoginResponse {
  token: string;
  trainer: Trainer;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  const res = await api.post<LoginResponse>('/auth/login', { email, password });
  setToken(res.token);
  return res;
}

export async function register(data: {
  name: string;
  email: string;
  password: string;
  phone: string;
}): Promise<Trainer> {
  return api.post<Trainer>('/auth/register', data);
}

export async function getMe(): Promise<Trainer> {
  return api.get<Trainer>('/me');
}

export function logout() {
  clearToken();
  window.location.href = '/login';
}

export function getRoleFromToken(): string {
  if (typeof window === 'undefined') return '';
  const token = localStorage.getItem('treinozap_token');
  if (!token) return '';
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.role ?? '';
  } catch {
    return '';
  }
}
