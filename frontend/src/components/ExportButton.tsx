import { useState } from 'react';
import * as App from '@bindings/github.com/joshrainwater/scan-organizer/app';

interface ExportButtonProps {
  onExport: (destination: string) => Promise<void>;
  fileCount: number;
  onSuccess: () => void;
}

export function ExportButton({ onExport, fileCount, onSuccess }: ExportButtonProps) {
  const [loading, setLoading] = useState(false);

  const handleExport = async () => {
    try {
      const destination = await App.SelectExportDirectory();
      if (!destination) return;

      setLoading(true);
      await onExport(destination);
      onSuccess();
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      onClick={handleExport}
      disabled={loading || fileCount === 0}
      className="w-full px-4 py-2 bg-green-600 text-white font-semibold uppercase tracking-wide rounded-md hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2"
    >
      {loading ? (
        'Exporting...'
      ) : (
        <>
          <span>📤</span>
          <span>Export ({fileCount} file{fileCount !== 1 ? 's' : ''})</span>
        </>
      )}
    </button>
  );
}
