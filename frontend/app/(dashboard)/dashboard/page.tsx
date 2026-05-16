'use client';

import { useEffect, useState } from 'react';
import { getMe } from '@/lib/auth';
import type { Trainer } from '@/types';
import Card from '@/components/ui/Card';
import Header from '@/components/layout/Header';
import Link from 'next/link';

const quickLinks = [
  {
    href: '/clients',
    label: 'Clientes',
    description: 'Gerencie seus alunos',
    accent: 'border-lime-400/20 hover:border-lime-400/50 hover:bg-lime-400/5',
    dot: 'bg-lime-400',
  },
  {
    href: '/exercises',
    label: 'Exercícios',
    description: 'Biblioteca de movimentos',
    accent: 'border-sky-400/20 hover:border-sky-400/50 hover:bg-sky-400/5',
    dot: 'bg-sky-400',
  },
  {
    href: '/messages',
    label: 'Mensagens',
    description: 'Histórico via WhatsApp',
    accent: 'border-amber-400/20 hover:border-amber-400/50 hover:bg-amber-400/5',
    dot: 'bg-amber-400',
  },
];

export default function DashboardPage() {
  const [trainer, setTrainer] = useState<Trainer | null>(null);

  useEffect(() => {
    getMe().then(setTrainer).catch(() => {});
  }, []);

  return (
    <div>
      <Header title="Dashboard" />
      <div className="p-6 space-y-6">
        {/* Greeting */}
        <Card className="p-6">
          <div className="flex items-start gap-4">
            <div className="w-10 h-10 rounded-full bg-lime-400/10 border border-lime-400/20 flex items-center justify-center flex-shrink-0">
              <span className="text-lime-400 text-sm font-bold">
                {trainer?.name?.[0]?.toUpperCase() ?? '?'}
              </span>
            </div>
            <div>
              <h2 className="text-base font-semibold text-zinc-100">
                Olá, {trainer?.name ?? '...'}!
              </h2>
              <p className="text-sm text-zinc-500 mt-0.5">
                Bem-vindo ao TreinoZap. Use o menu lateral para gerenciar clientes, exercícios e treinos.
              </p>
            </div>
          </div>
        </Card>

        {/* Quick links */}
        <div>
          <p className="text-xs font-semibold text-zinc-500 uppercase tracking-widest mb-3">
            Acesso rápido
          </p>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            {quickLinks.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`block p-4 bg-zinc-900 rounded-lg border transition-all ${item.accent}`}
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className={`w-2 h-2 rounded-full ${item.dot}`} />
                  <span className="text-sm font-semibold text-zinc-100">{item.label}</span>
                </div>
                <p className="text-xs text-zinc-500 pl-4">{item.description}</p>
              </Link>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
