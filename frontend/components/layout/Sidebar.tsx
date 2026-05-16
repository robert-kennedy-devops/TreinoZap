'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { logout, getRoleFromToken } from '@/lib/auth';

const navItems = [
  { href: '/dashboard', label: 'Dashboard' },
  { href: '/clients', label: 'Clientes' },
  { href: '/exercises', label: 'Exercícios' },
  { href: '/messages', label: 'Mensagens' },
];

const adminItems = [
  { href: '/admin/trainers', label: 'Treinadores' },
  { href: '/admin/clients', label: 'Todos os Clientes' },
  { href: '/settings/whatsapp', label: 'WhatsApp' },
];

export default function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    setIsAdmin(getRoleFromToken() === 'admin');
  }, []);

  function handleLogout() {
    logout();
    router.replace('/login');
  }

  function NavLink({ href, label }: { href: string; label: string }) {
    const active = pathname.startsWith(href);
    return (
      <Link
        href={href}
        className={`flex items-center gap-3 px-3 py-2.5 rounded-md text-sm font-medium transition-all ${
          active
            ? 'bg-lime-400/10 text-lime-400 border border-lime-400/20'
            : 'text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800/60 border border-transparent'
        }`}
      >
        <span
          className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
            active ? 'bg-lime-400' : 'bg-zinc-700'
          }`}
        />
        {label}
      </Link>
    );
  }

  return (
    <aside className="w-60 min-h-screen bg-zinc-950 border-r border-zinc-800 flex flex-col">
      {/* Logo */}
      <div className="px-6 py-5 border-b border-zinc-800">
        <span className="text-xl font-black tracking-tight text-zinc-100">
          Treino<span className="text-lime-400">Zap</span>
        </span>
        <p className="text-xs text-zinc-500 mt-0.5 font-medium uppercase tracking-widest">
          Personal Trainer
        </p>
      </div>

      {/* Nav */}
      <nav className="flex-1 py-4 px-3 space-y-0.5 overflow-y-auto">
        {navItems.map((item) => (
          <NavLink key={item.href} href={item.href} label={item.label} />
        ))}

        {isAdmin && (
          <>
            <div className="pt-4 pb-1 px-3">
              <p className="text-xs font-semibold text-zinc-600 uppercase tracking-widest">
                Admin
              </p>
            </div>
            {adminItems.map((item) => (
              <NavLink key={item.href} href={item.href} label={item.label} />
            ))}
          </>
        )}
      </nav>

      {/* Logout */}
      <div className="p-3 border-t border-zinc-800">
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-md text-sm font-medium text-zinc-400 hover:text-red-400 hover:bg-red-400/10 border border-transparent hover:border-red-400/20 transition-all"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>
          Sair da conta
        </button>
      </div>
    </aside>
  );
}
