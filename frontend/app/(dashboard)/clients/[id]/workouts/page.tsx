'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';
import type { Client, Workout } from '@/types';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Badge from '@/components/ui/Badge';
import Card from '@/components/ui/Card';

export default function ClientWorkoutsPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [client, setClient] = useState<Client | null>(null);
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [loading, setLoading] = useState(true);

  async function load() {
    setLoading(true);
    const [c, ws] = await Promise.all([
      api.get<Client>(`/clients/${id}`),
      api.get<Workout[]>(`/clients/${id}/workouts`),
    ]);
    setClient(c);
    setWorkouts(Array.isArray(ws) ? ws : []);
    setLoading(false);
  }

  useEffect(() => { load(); }, [id]);

  async function activate(workoutId: string) {
    await api.post(`/workouts/${workoutId}/activate`);
    load();
  }

  if (loading) return <div className="p-6 text-sm text-gray-500">Carregando...</div>;

  return (
    <div>
      <Header title={`Treinos — ${client?.name}`} />
      <div className="p-6 space-y-4">
        <div className="flex justify-between">
          <Button variant="secondary" size="sm" onClick={() => router.back()}>← Voltar</Button>
          <Link href={`/clients/${id}/workouts/new`}>
            <Button>+ Novo treino</Button>
          </Link>
        </div>

        {workouts.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum treino cadastrado.</p>
        ) : (
          <div className="space-y-3">
            {workouts.map((w) => (
              <Card key={w.id} className="p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-sm font-medium text-gray-800">{w.name}</span>
                      <Badge status={w.status} />
                    </div>
                    {w.starts_at && (
                      <p className="text-xs text-gray-400">{w.starts_at} → {w.ends_at}</p>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Link href={`/workouts/${w.id}`}>
                      <Button size="sm" variant="ghost">Editar</Button>
                    </Link>
                    {w.status !== 'active' && (
                      <Button size="sm" variant="secondary" onClick={() => activate(w.id)}>Ativar</Button>
                    )}
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
