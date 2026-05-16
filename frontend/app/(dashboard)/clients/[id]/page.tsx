'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Client, Workout } from '@/types';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Badge from '@/components/ui/Badge';
import Input from '@/components/ui/Input';
import Textarea from '@/components/ui/Textarea';
import Card from '@/components/ui/Card';

export default function ClientDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [client, setClient] = useState<Client | null>(null);
  const [activeWorkout, setActiveWorkout] = useState<Workout | null>(null);
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState({ name: '', phone: '', goal: '', notes: '', status: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  useEffect(() => {
    api.get<Client>(`/clients/${id}`)
      .then((c) => {
        setClient(c);
        setForm({ name: c.name, phone: c.phone, goal: c.goal || '', notes: c.notes || '', status: c.status });
      })
      .catch(() => router.push('/clients'));

    api.get<Workout[]>(`/clients/${id}/workouts`)
      .then((ws) => {
        const active = ws.find((w) => w.status === 'active') ?? null;
        setActiveWorkout(active);
      })
      .catch(() => {});
  }, [id, router]);

  function update(field: string, value: string) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const updated = await api.put<Client>(`/clients/${id}`, form);
      setClient(updated);
      setEditing(false);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao salvar');
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete() {
    if (!confirm('Inativar este cliente?')) return;
    await api.delete(`/clients/${id}`);
    router.push('/clients');
  }

  if (!client) return <div className="p-6 text-gray-500 text-sm">Carregando...</div>;

  return (
    <div>
      <Header title={client.name} />
      <div className="p-6 space-y-6 max-w-2xl">
        <Card className="p-6">
          {!editing ? (
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h2 className="text-base font-semibold text-gray-800">{client.name}</h2>
                <Badge status={client.status} />
              </div>
              <p className="text-sm text-gray-600">📱 {client.phone}</p>
              {client.goal && <p className="text-sm text-gray-600">🎯 {client.goal}</p>}
              {client.notes && <p className="text-sm text-gray-500 italic">{client.notes}</p>}

              <div className="flex gap-3 pt-2">
                <Button size="sm" variant="secondary" onClick={() => setEditing(true)}>Editar</Button>
                <Button size="sm" variant="danger" onClick={handleDelete}>Inativar</Button>
                <Link href={`/clients/${id}/workouts`}>
                  <Button size="sm" variant="secondary">Treinos</Button>
                </Link>
                <Link href={`/clients/${id}/messages`}>
                  <Button size="sm" variant="ghost">Mensagens</Button>
                </Link>
              </div>
            </div>
          ) : (
            <form onSubmit={handleSave} className="space-y-4">
              <Input label="Nome" value={form.name} onChange={(e) => update('name', e.target.value)} required />
              <Input label="Telefone" value={form.phone} onChange={(e) => update('phone', e.target.value)} required />
              <Input label="Objetivo" value={form.goal} onChange={(e) => update('goal', e.target.value)} />
              <Textarea label="Observações" value={form.notes} onChange={(e) => update('notes', e.target.value)} />
              <div>
                <label className="text-sm font-medium text-gray-700">Status</label>
                <select
                  value={form.status}
                  onChange={(e) => update('status', e.target.value)}
                  className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
                >
                  <option value="active">Ativo</option>
                  <option value="inactive">Inativo</option>
                  <option value="blocked">Bloqueado</option>
                </select>
              </div>
              {error && <p className="text-sm text-red-600">{error}</p>}
              <div className="flex gap-3">
                <Button type="submit" loading={loading}>Salvar</Button>
                <Button type="button" variant="secondary" onClick={() => setEditing(false)}>Cancelar</Button>
              </div>
            </form>
          )}
        </Card>

        {activeWorkout && (
          <Card className="p-5">
            <p className="text-xs text-gray-500 mb-0.5">Treino ativo</p>
            <p className="text-sm font-medium text-gray-800">{activeWorkout.name}</p>
            <p className="text-xs text-gray-400 mt-1">
              O cliente pode solicitar o treino enviando <strong>treino</strong> pelo WhatsApp.
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
