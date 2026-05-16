'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import type { Workout } from '@/types';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Badge from '@/components/ui/Badge';
import Card from '@/components/ui/Card';

export default function WorkoutDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [workout, setWorkout] = useState<Workout | null>(null);
  const [loading, setLoading] = useState(true);
  const [actionMsg, setActionMsg] = useState('');

  useEffect(() => {
    api.get<Workout>(`/workouts/${id}`)
      .then(setWorkout)
      .finally(() => setLoading(false));
  }, [id]);

  async function activate() {
    await api.post(`/workouts/${id}/activate`);
    setActionMsg('Treino ativado!');
    api.get<Workout>(`/workouts/${id}`).then(setWorkout);
  }

  if (loading) return <div className="p-6 text-sm text-gray-500">Carregando...</div>;
  if (!workout) return <div className="p-6 text-sm text-red-500">Treino não encontrado.</div>;

  return (
    <div>
      <Header title={workout.name} />
      <div className="p-6 space-y-4 max-w-2xl">
        <div className="flex gap-3 flex-wrap">
          <Button variant="secondary" size="sm" onClick={() => router.back()}>← Voltar</Button>
          {workout.status !== 'active' && (
            <Button size="sm" onClick={activate}>Ativar treino</Button>
          )}
        </div>

        {actionMsg && (
          <p className="text-sm text-green-600 bg-green-50 px-3 py-2 rounded-md">{actionMsg}</p>
        )}

        <Card className="p-4">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-sm font-semibold text-gray-800">{workout.name}</span>
            <Badge status={workout.status} />
          </div>
          {workout.starts_at && (
            <p className="text-xs text-gray-400">{workout.starts_at} → {workout.ends_at}</p>
          )}
        </Card>

        {workout.sections?.map((section, si) => (
          <Card key={section.id} className="p-4">
            <h3 className="text-sm font-semibold text-gray-700 mb-3">{section.name}</h3>
            {section.description && (
              <p className="text-xs text-gray-500 mb-2 italic">{section.description}</p>
            )}
            <ol className="space-y-2">
              {section.exercises?.map((ex, ei) => (
                <li key={ex.id} className="text-sm text-gray-700">
                  <span className="font-medium">{ei + 1}. {ex.exercise_name}</span>
                  {(ex.sets || ex.reps) && (
                    <span className="text-gray-500"> — {ex.sets && `${ex.sets}x`}{ex.reps}</span>
                  )}
                  {ex.rest_seconds && (
                    <span className="text-gray-400"> — {ex.rest_seconds}s</span>
                  )}
                  {ex.load_note && <p className="text-xs text-gray-500 ml-4">Carga: {ex.load_note}</p>}
                  {ex.technique_note && <p className="text-xs text-gray-500 ml-4">Técnica: {ex.technique_note}</p>}
                </li>
              ))}
            </ol>
          </Card>
        ))}
      </div>
    </div>
  );
}
