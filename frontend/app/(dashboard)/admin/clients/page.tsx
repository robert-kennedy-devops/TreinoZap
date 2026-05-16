'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import type { Client } from '@/types';
import Header from '@/components/layout/Header';
import Badge from '@/components/ui/Badge';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';
import Link from 'next/link';

export default function AdminClientsPage() {
  const [clients, setClients] = useState<Client[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  function load(s = '') {
    setLoading(true);
    setError('');
    api
      .get<Client[]>(`/admin/clients${s ? `?search=${encodeURIComponent(s)}` : ''}`)
      .then(setClients)
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : 'Erro ao carregar clientes')
      )
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    load();
  }, []);

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    load(search);
  }

  return (
    <div>
      <Header title="Todos os Clientes" />
      <div className="p-6 space-y-4">
        <form onSubmit={handleSearch} className="flex gap-2 max-w-sm">
          <Input
            placeholder="Buscar por nome ou telefone"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <Button type="submit" size="sm">Buscar</Button>
        </form>

        {error && (
          <p className="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-md">{error}</p>
        )}

        {loading ? (
          <p className="text-sm text-gray-500">Carregando...</p>
        ) : clients.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum cliente encontrado.</p>
        ) : (
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Nome</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Telefone</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Status</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Treinador ID</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Cadastro</th>
                  <th className="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {clients.map((c) => (
                  <tr key={c.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3 font-medium text-gray-800">{c.name}</td>
                    <td className="px-4 py-3 text-gray-600">{c.phone}</td>
                    <td className="px-4 py-3">
                      <Badge status={c.status} />
                    </td>
                    <td className="px-4 py-3 text-gray-400 text-xs font-mono">
                      {c.trainer_id.slice(0, 8)}…
                    </td>
                    <td className="px-4 py-3 text-gray-400 text-xs">
                      {new Date(c.created_at).toLocaleDateString('pt-BR')}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <Link
                        href={`/clients/${c.id}`}
                        className="text-blue-600 hover:underline text-xs"
                      >
                        Ver
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div className="px-4 py-2 bg-gray-50 border-t border-gray-200 text-xs text-gray-400">
              {clients.length} cliente{clients.length !== 1 ? 's' : ''} encontrado{clients.length !== 1 ? 's' : ''}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
