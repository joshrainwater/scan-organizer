import { useState, useEffect, useRef } from 'react';

interface RenameFormProps {
  folders: string[];
  onRename: (newName: string, folder: string) => Promise<void>;
  inputFilesLength: number;
}

export function RenameForm({ folders, onRename, inputFilesLength }: RenameFormProps) {
  const [folder, setFolder] = useState('');
  const [newName, setNewName] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const folderInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (inputFilesLength > 0) {
      folderInputRef.current?.focus();
    }
  }, [inputFilesLength]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim() || submitting) return;

    setSubmitting(true);
    try {
      await onRename(newName.trim(), folder.trim());
      setNewName('');
      setFolder('');
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitting(false);
      folderInputRef.current?.focus();
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <input
        ref={folderInputRef}
        type="text"
        list="folder-list"
        value={folder}
        onChange={(e) => setFolder(e.target.value)}
        placeholder="Select or type folder name..."
        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        autoComplete="off"
      />
      <datalist id="folder-list">
        {folders.map((f) => (
          <option key={f} value={f} />
        ))}
      </datalist>

      <input
        type="text"
        value={newName}
        onChange={(e) => setNewName(e.target.value)}
        placeholder="Rename to..."
        required
        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <button
        type="submit"
        disabled={submitting || !newName.trim()}
        className="w-full px-4 py-2 bg-blue-600 text-white font-semibold uppercase tracking-wide rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {submitting ? 'Renaming...' : 'Rename'}
      </button>
    </form>
  );
}
