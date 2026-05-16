'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Textarea from '@/components/ui/Textarea';

interface ExerciseInput {
  exercise_name: string;
  sets: string;
  reps: string;
  rest_seconds: string;
  load_note: string;
  technique_note: string;
  video_url: string;
  order_index: number;
}

interface SectionInput {
  name: string;
  description: string;
  order_index: number;
  exercises: ExerciseInput[];
}

const emptyExercise = (order: number): ExerciseInput => ({
  exercise_name: '', sets: '', reps: '', rest_seconds: '', load_note: '',
  technique_note: '', video_url: '', order_index: order,
});

const emptySection = (order: number): SectionInput => ({
  name: '', description: '', order_index: order, exercises: [emptyExercise(1)],
});

export default function NewWorkoutPage() {
  const { id: clientId } = useParams<{ id: string }>();
  const router = useRouter();
  const [name, setName] = useState('');
  const [startsAt, setStartsAt] = useState('');
  const [endsAt, setEndsAt] = useState('');
  const [sections, setSections] = useState<SectionInput[]>([emptySection(1)]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  function updateSection(i: number, field: string, value: string) {
    setSections((ss) => ss.map((s, idx) => idx === i ? { ...s, [field]: value } : s));
  }

  function updateExercise(si: number, ei: number, field: string, value: string) {
    setSections((ss) => ss.map((s, idx) =>
      idx !== si ? s : {
        ...s,
        exercises: s.exercises.map((e, eidx) => eidx === ei ? { ...e, [field]: value } : e),
      }
    ));
  }

  function addSection() {
    setSections((ss) => [...ss, emptySection(ss.length + 1)]);
  }

  function addExercise(si: number) {
    setSections((ss) => ss.map((s, idx) =>
      idx !== si ? s : { ...s, exercises: [...s.exercises, emptyExercise(s.exercises.length + 1)] }
    ));
  }

  function removeSection(i: number) {
    setSections((ss) => ss.filter((_, idx) => idx !== i));
  }

  function removeExercise(si: number, ei: number) {
    setSections((ss) => ss.map((s, idx) =>
      idx !== si ? s : { ...s, exercises: s.exercises.filter((_, eidx) => eidx !== ei) }
    ));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const body = {
        name, starts_at: startsAt, ends_at: endsAt,
        sections: sections.map((s) => ({
          ...s,
          exercises: s.exercises.map((e) => ({
            ...e,
            rest_seconds: e.rest_seconds ? parseInt(e.rest_seconds) : null,
            exercise_id: null,
          })),
        })),
      };
      await api.post(`/clients/${clientId}/workouts`, body);
      router.push(`/clients/${clientId}/workouts`);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao salvar');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <Header title="Novo Treino" />
      <div className="p-6 max-w-3xl">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="bg-white rounded-lg border border-gray-200 p-5 space-y-4">
            <Input label="Nome do treino" value={name} onChange={(e) => setName(e.target.value)} required placeholder="Ex: Plano Hipertrofia - Semana 1" />
            <div className="grid grid-cols-2 gap-4">
              <Input label="Início" type="date" value={startsAt} onChange={(e) => setStartsAt(e.target.value)} />
              <Input label="Fim" type="date" value={endsAt} onChange={(e) => setEndsAt(e.target.value)} />
            </div>
          </div>

          {sections.map((section, si) => (
            <div key={si} className="bg-white rounded-lg border border-gray-200 p-5 space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold text-gray-700">Seção {si + 1}</h3>
                {sections.length > 1 && (
                  <Button type="button" size="sm" variant="danger" onClick={() => removeSection(si)}>Remover seção</Button>
                )}
              </div>
              <Input label="Nome da seção" value={section.name} onChange={(e) => updateSection(si, 'name', e.target.value)} placeholder="Ex: Treino A - Peito e Tríceps" required />
              <Textarea label="Descrição (opcional)" value={section.description} onChange={(e) => updateSection(si, 'description', e.target.value)} rows={2} />

              <div className="space-y-3">
                {section.exercises.map((ex, ei) => (
                  <div key={ei} className="border border-gray-100 rounded-md p-4 space-y-3 bg-gray-50">
                    <div className="flex justify-between items-center">
                      <span className="text-xs font-medium text-gray-600">Exercício {ei + 1}</span>
                      {section.exercises.length > 1 && (
                        <button type="button" onClick={() => removeExercise(si, ei)} className="text-xs text-red-500 hover:text-red-700">Remover</button>
                      )}
                    </div>
                    <Input label="Nome do exercício" value={ex.exercise_name} onChange={(e) => updateExercise(si, ei, 'exercise_name', e.target.value)} required />
                    <div className="grid grid-cols-3 gap-3">
                      <Input label="Séries" value={ex.sets} onChange={(e) => updateExercise(si, ei, 'sets', e.target.value)} placeholder="4" />
                      <Input label="Repetições" value={ex.reps} onChange={(e) => updateExercise(si, ei, 'reps', e.target.value)} placeholder="10-12" />
                      <Input label="Descanso (s)" type="number" value={ex.rest_seconds} onChange={(e) => updateExercise(si, ei, 'rest_seconds', e.target.value)} placeholder="60" />
                    </div>
                    <Input label="Obs. de carga" value={ex.load_note} onChange={(e) => updateExercise(si, ei, 'load_note', e.target.value)} />
                    <Input label="Obs. técnica" value={ex.technique_note} onChange={(e) => updateExercise(si, ei, 'technique_note', e.target.value)} />
                    <Input label="Link do vídeo" value={ex.video_url} onChange={(e) => updateExercise(si, ei, 'video_url', e.target.value)} />
                  </div>
                ))}
                <Button type="button" size="sm" variant="secondary" onClick={() => addExercise(si)}>+ Exercício</Button>
              </div>
            </div>
          ))}

          <Button type="button" variant="secondary" onClick={addSection}>+ Seção</Button>

          {error && <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-md">{error}</p>}

          <div className="flex gap-3">
            <Button type="submit" loading={loading}>Salvar treino</Button>
            <Button type="button" variant="secondary" onClick={() => router.back()}>Cancelar</Button>
          </div>
        </form>
      </div>
    </div>
  );
}
