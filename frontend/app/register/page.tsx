'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { register } from '@/lib/auth';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState({ name: '', email: '', password: '', phone: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  function update(field: string, value: string) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await register(form);
      router.replace('/login');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao criar conta');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-950 px-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-black tracking-tight text-zinc-100">
            Treino<span className="text-lime-400">Zap</span>
          </h1>
          <p className="text-xs text-zinc-500 mt-1.5 font-medium uppercase tracking-widest">
            Personal Trainer Platform
          </p>
        </div>

        {/* Card */}
        <div className="bg-zinc-900 rounded-xl border border-zinc-800 p-6 shadow-2xl">
          <p className="text-sm font-semibold text-zinc-300 mb-5">Criar conta de treinador</p>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Nome"
              value={form.name}
              onChange={(e) => update('name', e.target.value)}
              placeholder="Seu nome completo"
              required
            />
            <Input
              label="E-mail"
              type="email"
              value={form.email}
              onChange={(e) => update('email', e.target.value)}
              placeholder="seu@email.com"
              required
            />
            <Input
              label="Senha"
              type="password"
              value={form.password}
              onChange={(e) => update('password', e.target.value)}
              placeholder="••••••"
              required
            />
            <Input
              label="Telefone (com DDD e país)"
              value={form.phone}
              onChange={(e) => update('phone', e.target.value)}
              placeholder="5592999999999"
            />

            {error && (
              <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-md px-3 py-2">
                {error}
              </p>
            )}

            <Button type="submit" loading={loading} className="w-full mt-2">
              Criar conta
            </Button>
          </form>

          <p className="mt-5 text-center text-xs text-zinc-500">
            Já tem conta?{' '}
            <Link href="/login" className="text-lime-400 hover:text-lime-300 transition-colors">
              Entrar
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
