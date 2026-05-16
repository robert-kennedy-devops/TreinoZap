'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import type { Exercise, PaginatedResponse } from '@/types';
import Header from '@/components/layout/Header';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import Textarea from '@/components/ui/Textarea';

export default function ExercisesPage() {
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editTarget, setEditTarget] = useState<Exercise | null>(null);
  const [form, setForm] = useState({ name: '', muscle_group: '', equipment: '', video_url: '', notes: '' });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  async function load(s = '') {
    setLoading(true);
    try {
      const res = await api.get<PaginatedResponse<Exercise>>(`/exercises?search=${s}&page_size=100`);
      setExercises(res.data ?? []);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, []);

  function update(field: string, value: string) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  function startEdit(e: Exercise) {
    setEditTarget(e);
    setForm({ name: e.name, muscle_group: e.muscle_group || '', equipment: e.equipment || '', video_url: e.video_url || '', notes: e.notes || '' });
    setShowForm(true);
    setError('');
  }

  function startNew() {
    setEditTarget(null);
    setForm({ name: '', muscle_group: '', equipment: '', video_url: '', notes: '' });
    setShowForm(true);
    setError('');
  }

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setSaving(true);
    try {
      if (editTarget) {
        await api.put(`/exercises/${editTarget.id}`, form);
      } else {
        await api.post('/exercises', form);
      }
      setShowForm(false);
      load(search);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao salvar');
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Remover exercício?')) return;
    await api.delete(`/exercises/${id}`);
    load(search);
  }

  return (
    <div>
      <Header title="Exercícios" />
      <div className="p-6 space-y-4">
        <div className="flex items-center justify-between gap-4">
          <form onSubmit={(e) => { e.preventDefault(); load(search); }} className="flex gap-2 flex-1 max-w-sm">
            <Input placeholder="Buscar exercício" value={search} onChange={(e) => setSearch(e.target.value)} />
            <Button type="submit" size="sm">Buscar</Button>
          </form>
          <Button onClick={startNew}>+ Novo exercício</Button>
        </div>

        {showForm && (
          <div className="bg-white border border-gray-200 rounded-lg p-5">
            <h3 className="text-sm font-semibold text-gray-700 mb-4">{editTarget ? 'Editar exercício' : 'Novo exercício'}</h3>
            <form onSubmit={handleSave} className="space-y-3">
              <Input label="Nome" value={form.name} onChange={(e) => update('name', e.target.value)} required />
              <div className="grid grid-cols-2 gap-3">
                <Input label="Grupo muscular" value={form.muscle_group} onChange={(e) => update('muscle_group', e.target.value)} />
                <Input label="Equipamento" value={form.equipment} onChange={(e) => update('equipment', e.target.value)} />
              </div>
              <Input label="Link do vídeo" value={form.video_url} onChange={(e) => update('video_url', e.target.value)} />
              <Textarea label="Observações" value={form.notes} onChange={(e) => update('notes', e.target.value)} rows={2} />
              {error && <p className="text-sm text-red-600">{error}</p>}
              <div className="flex gap-2">
                <Button type="submit" loading={saving}>Salvar</Button>
                <Button type="button" variant="secondary" onClick={() => setShowForm(false)}>Cancelar</Button>
              </div>
            </form>
          </div>
        )}

        {loading ? (
          <p className="text-sm text-gray-500">Carregando...</p>
        ) : exercises.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum exercício encontrado.</p>
        ) : (
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Nome</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Grupo</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Equipamento</th>
                  <th className="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {exercises.map((e) => (
                  <tr key={e.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-800">{e.name}</td>
                    <td className="px-4 py-3 text-gray-600">{e.muscle_group || '-'}</td>
                    <td className="px-4 py-3 text-gray-600">{e.equipment || '-'}</td>
                    <td className="px-4 py-3 text-right flex gap-2 justify-end">
                      <Button size="sm" variant="ghost" onClick={() => startEdit(e)}>Editar</Button>
                      <Button size="sm" variant="danger" onClick={() => handleDelete(e.id)}>Remover</Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
