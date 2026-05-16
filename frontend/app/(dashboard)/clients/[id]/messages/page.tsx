'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { api } from '@/lib/api';
import type { Client, Message, PaginatedResponse } from '@/types';
import Header from '@/components/layout/Header';

export default function ClientMessagesPage() {
  const { id } = useParams<{ id: string }>();
  const [client, setClient] = useState<Client | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.get<Client>(`/clients/${id}`),
      api.get<PaginatedResponse<Message>>(`/clients/${id}/messages?page_size=100`),
    ]).then(([c, ms]) => {
      setClient(c);
      setMessages([...(ms.data ?? [])].reverse());
    }).finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="p-6 text-sm text-gray-500">Carregando...</div>;

  return (
    <div>
      <Header title={`Mensagens — ${client?.name}`} />
      <div className="p-6">
        <div className="max-w-2xl space-y-2">
          {messages.length === 0 && (
            <p className="text-sm text-gray-500">Sem mensagens registradas.</p>
          )}
          {messages.map((m) => (
            <div
              key={m.id}
              className={`flex ${m.direction === 'outbound' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-xs rounded-lg px-4 py-2 text-sm ${
                  m.direction === 'outbound'
                    ? 'bg-blue-600 text-white'
                    : 'bg-white border border-gray-200 text-gray-800'
                }`}
              >
                <p className="whitespace-pre-wrap">{m.message}</p>
                <p className={`text-xs mt-1 ${m.direction === 'outbound' ? 'text-blue-200' : 'text-gray-400'}`}>
                  {new Date(m.created_at).toLocaleString('pt-BR')}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
