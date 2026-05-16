interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export default function Card({ children, className = '' }: CardProps) {
  return (
    <div className={`bg-zinc-900 rounded-lg border border-zinc-800 ${className}`}>
      {children}
    </div>
  );
}
