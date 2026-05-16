'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Textarea from '@/components/ui/Textarea';

export default function NewClientPage() {
  const router = useRouter();
  const [form, setForm] = useState({ name: '', phone: '', goal: '', notes: '' });
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
      await api.post('/clients', form);
      router.push('/clients');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao criar cliente');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <Header title="Novo Cliente" />
      <div className="p-6 max-w-lg">
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input label="Nome" value={form.name} onChange={(e) => update('name', e.target.value)} required />
            <Input label="Telefone (com DDD e código do país)" value={form.phone} onChange={(e) => update('phone', e.target.value)} placeholder="5592999999999" required />
            <Input label="Objetivo" value={form.goal} onChange={(e) => update('goal', e.target.value)} placeholder="Hipertrofia, emagrecimento..." />
            <Textarea label="Observações" value={form.notes} onChange={(e) => update('notes', e.target.value)} />

            {error && <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-md">{error}</p>}

            <div className="flex gap-3">
              <Button type="submit" loading={loading}>Salvar</Button>
              <Button type="button" variant="secondary" onClick={() => router.back()}>Cancelar</Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
