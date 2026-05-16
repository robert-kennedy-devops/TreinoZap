'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import type { Message, PaginatedResponse } from '@/types';
import Header from '@/components/layout/Header';

export default function MessagesPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<PaginatedResponse<Message>>('/messages?page_size=100')
      .then((res) => setMessages(res.data ?? []))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <Header title="Mensagens" />
      <div className="p-6">
        {loading ? (
          <p className="text-sm text-gray-500">Carregando...</p>
        ) : messages.length === 0 ? (
          <p className="text-sm text-gray-500">Sem mensagens registradas.</p>
        ) : (
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Direção</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Telefone</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Mensagem</th>
                  <th className="text-left px-4 py-3 text-gray-600 font-medium">Hora</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {messages.map((m) => (
                  <tr key={m.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3">
                      <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${
                        m.direction === 'outbound' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'
                      }`}>
                        {m.direction === 'outbound' ? 'Enviada' : 'Recebida'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-gray-600">{m.phone}</td>
                    <td className="px-4 py-3 text-gray-700 max-w-xs truncate">{m.message}</td>
                    <td className="px-4 py-3 text-gray-400 text-xs whitespace-nowrap">
                      {new Date(m.created_at).toLocaleString('pt-BR')}
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
