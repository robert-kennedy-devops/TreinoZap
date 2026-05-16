import { ButtonHTMLAttributes } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md';
  loading?: boolean;
}

const variants = {
  primary:
    'bg-lime-400 text-zinc-900 hover:bg-lime-300 font-semibold',
  secondary:
    'bg-zinc-800 text-zinc-100 hover:bg-zinc-700 border border-zinc-700 hover:border-zinc-600',
  danger:
    'bg-red-500/10 text-red-400 hover:bg-red-500/20 border border-red-500/30 hover:border-red-500/50',
  ghost:
    'text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800 border border-transparent',
};

const sizes = {
  sm: 'px-3 py-1.5 text-xs',
  md: 'px-4 py-2 text-sm',
};

export default function Button({
  variant = 'primary',
  size = 'md',
  loading,
  children,
  className = '',
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      {...props}
      disabled={disabled || loading}
      className={`inline-flex items-center justify-center gap-2 rounded-md transition-all focus:outline-none focus:ring-2 focus:ring-lime-400/50 focus:ring-offset-1 focus:ring-offset-zinc-950 disabled:opacity-50 disabled:cursor-not-allowed ${variants[variant]} ${sizes[size]} ${className}`}
    >
      {loading ? (
        <>
          <span className="w-3.5 h-3.5 border-2 border-current border-t-transparent rounded-full animate-spin" />
          Aguarde...
        </>
      ) : children}
    </button>
  );
}
