import { useState, useRef } from 'react';

interface TrashButtonProps {
  onTrash: () => Promise<void>;
}

export function TrashButton({ onTrash }: TrashButtonProps) {
  const [submitting, setSubmitting] = useState(false);
  const buttonRef = useRef<HTMLButtonElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (submitting) return;

    setSubmitting(true);
    try {
      await onTrash();
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <button
        ref={buttonRef}
        type="submit"
        disabled={submitting}
        className="w-full px-4 py-2 bg-red-600 text-white font-semibold uppercase tracking-wide rounded-md hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {submitting ? 'Trashing...' : 'Trash'}
      </button>
    </form>
  );
}
