import { useState, useRef } from 'react';

interface AppendFormProps {
  previousRenamed: string[];
  onAppend: (target: string) => Promise<void>;
}

export function AppendForm({ previousRenamed, onAppend }: AppendFormProps) {
  const [target, setTarget] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const selectRef = useRef<HTMLSelectElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!target || submitting) return;

    setSubmitting(true);
    try {
      await onAppend(target);
      setTarget('');
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitting(false);
      selectRef.current?.focus();
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <select
        ref={selectRef}
        value={target}
        onChange={(e) => setTarget(e.target.value)}
        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
      >
        <option value="">Select a file...</option>
        {previousRenamed.map((name) => (
          <option key={name} value={name}>
            {name}
          </option>
        ))}
      </select>

      <button
        type="submit"
        disabled={submitting || !target}
        className="w-full px-4 py-2 bg-blue-600 text-white font-semibold uppercase tracking-wide rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {submitting ? 'Appending...' : 'Append to selected'}
      </button>
    </form>
  );
}
