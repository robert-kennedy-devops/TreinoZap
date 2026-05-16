import { TextareaHTMLAttributes } from 'react';

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
}

export default function Textarea({ label, error, className = '', ...props }: TextareaProps) {
  return (
    <div className="flex flex-col gap-1.5">
      {label && (
        <label className="text-xs font-semibold text-zinc-400 uppercase tracking-wider">
          {label}
        </label>
      )}
      <textarea
        {...props}
        rows={props.rows ?? 3}
        className={`block w-full rounded-md bg-zinc-900 border ${
          error
            ? 'border-red-500/60 focus:border-red-500 focus:ring-red-500/30'
            : 'border-zinc-700 focus:border-lime-400 focus:ring-lime-400/20'
        } px-3 py-2 text-sm text-zinc-100 placeholder-zinc-600 focus:outline-none focus:ring-2 resize-none transition-colors ${className}`}
      />
      {error && <span className="text-xs text-red-400">{error}</span>}
    </div>
  );
}
