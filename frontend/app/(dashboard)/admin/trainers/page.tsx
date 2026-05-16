'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import type { Trainer } from '@/types';
import Header from '@/components/layout/Header';
import Badge from '@/components/ui/Badge';

export default function AdminTrainersPage() {
  const [trainers, setTrainers] = useState<Trainer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    api
      .get<Trainer[]>('/admin/trainers')
      .then(setTrainers)
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : 'Erro ao carregar treinadores')
      )
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <Header title="Treinadores" />
      <div className="p-6">
        {error && (
          <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-md mb-4">{error}</p>
        )}
        {loading ? (
          <p className="text-sm text-gray-500">Carregando...</p>
        ) : trainers.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum treinador encontrado.</p>
        ) : (
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Nome</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">E-mail</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Telefone</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Perfil</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Status</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Cadastro</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {trainers.map((t) => (
                  <tr key={t.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-800">{t.name}</td>
                    <td className="px-4 py-3 text-gray-600">{t.email}</td>
                    <td className="px-4 py-3 text-gray-600">{t.phone || '—'}</td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                          t.role === 'admin'
                            ? 'bg-purple-100 text-purple-700'
                            : 'bg-gray-100 text-gray-600'
                        }`}
                      >
                        {t.role === 'admin' ? 'Admin' : 'Treinador'}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <Badge status={t.status} />
                    </td>
                    <td className="px-4 py-3 text-gray-400 text-xs">
                      {new Date(t.created_at).toLocaleDateString('pt-BR')}
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
