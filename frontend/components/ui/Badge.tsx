interface BadgeProps {
  status: string;
}

const colors: Record<string, string> = {
  active: 'bg-lime-400/10 text-lime-400 border-lime-400/20',
  inactive: 'bg-zinc-700/50 text-zinc-400 border-zinc-700',
  blocked: 'bg-red-500/10 text-red-400 border-red-500/20',
  draft: 'bg-amber-400/10 text-amber-400 border-amber-400/20',
  ready: 'bg-sky-400/10 text-sky-400 border-sky-400/20',
  archived: 'bg-zinc-700/50 text-zinc-500 border-zinc-700',
  connected: 'bg-lime-400/10 text-lime-400 border-lime-400/20',
  disconnected: 'bg-red-500/10 text-red-400 border-red-500/20',
};

const labels: Record<string, string> = {
  active: 'Ativo',
  inactive: 'Inativo',
  blocked: 'Bloqueado',
  draft: 'Rascunho',
  ready: 'Pronto',
  archived: 'Arquivado',
  connected: 'Conectado',
  disconnected: 'Desconectado',
};

export default function Badge({ status }: BadgeProps) {
  const color = colors[status] || 'bg-zinc-700/50 text-zinc-400 border-zinc-700';
  const label = labels[status] || status;
  return (
    <span
      className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium ${color}`}
    >
      {label}
    </span>
  );
}
