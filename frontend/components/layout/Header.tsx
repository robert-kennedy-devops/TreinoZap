'use client';

interface HeaderProps {
  title: string;
  subtitle?: string;
}

export default function Header({ title, subtitle }: HeaderProps) {
  return (
    <header className="h-14 border-b border-zinc-800 flex items-center px-6 bg-zinc-950/80 backdrop-blur sticky top-0 z-10">
      <div>
        <h1 className="text-sm font-semibold text-zinc-100 tracking-wide">{title}</h1>
        {subtitle && <p className="text-xs text-zinc-500 mt-0.5">{subtitle}</p>}
      </div>
    </header>
  );
}
