import { useState, useCallback } from 'react';
import { useDropzone } from '../hooks/useDropzone';
import * as App from '../../bindings/github.com/joshrainwater/scan-organizer/app';

interface ModeDialogProps {
  isOpen: boolean;
  onSelect: (mode: 'replace' | 'append' | 'skip') => void;
  existingCount: number;
}

function ModeDialog({ isOpen, onSelect, existingCount }: ModeDialogProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 shadow-xl">
        <h2 className="text-xl font-semibold text-gray-800 mb-2">
          Add Files to Input?
        </h2>
        <p className="text-gray-600 mb-4">
          You already have {existingCount} file{existingCount !== 1 ? 's' : ''} in the input folder.
          What would you like to do?
        </p>
        <div className="space-y-3">
          <button
            onClick={() => onSelect('replace')}
            className="w-full px-4 py-3 bg-red-600 text-white font-semibold rounded-md hover:bg-red-700 transition-colors text-left"
          >
            Replace Input
            <span className="block text-sm font-normal text-red-100">
              Clear existing files and add new ones
            </span>
          </button>
          <button
            onClick={() => onSelect('append')}
            className="w-full px-4 py-3 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-700 transition-colors text-left"
          >
            Append to Input
            <span className="block text-sm font-normal text-blue-100">
              Keep existing files and add new ones
            </span>
          </button>
          <button
            onClick={() => onSelect('skip')}
            className="w-full px-4 py-3 bg-gray-200 text-gray-800 font-semibold rounded-md hover:bg-gray-300 transition-colors text-left"
          >
            Skip
            <span className="block text-sm font-normal text-gray-500">
              Don't add new files, keep existing
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}

interface HomeProps {
  onReady: () => void;
}

export function Home({ onReady }: HomeProps) {
  const [pendingPaths, setPendingPaths] = useState<string[]>([]);
  const [showModeDialog, setShowModeDialog] = useState(false);
  const [existingCount, setExistingCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const checkExistingAndProceed = useCallback(async (paths: string[]) => {
    setLoading(true);
    try {
      const status = await App.GetStatus();
      setExistingCount(status.inputCount);
      setPendingPaths(paths);
      if (status.inputCount > 0) {
        setShowModeDialog(true);
      } else {
        await App.AddFiles(paths, 'replace');
        onReady();
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, [onReady]);

  const { isDragging, handleDragEnter, handleDragLeave, handleDragOver, handleDrop } =
    useDropzone(checkExistingAndProceed);

  const handleModeSelect = async (mode: 'replace' | 'append' | 'skip') => {
    setShowModeDialog(false);
    if (mode === 'skip') {
      onReady();
      return;
    }

    setLoading(true);
    try {
      await App.AddFiles(pendingPaths, mode);
      onReady();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
      setPendingPaths([]);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center p-8">
      <div
        className={`w-full max-w-2xl border-4 border-dashed rounded-xl p-16 text-center transition-colors ${
          isDragging
            ? 'border-blue-500 bg-blue-50'
            : 'border-gray-300 bg-white'
        }`}
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
      >
        {loading ? (
          <div className="text-gray-500">Processing...</div>
        ) : isDragging ? (
          <div className="text-blue-600 text-2xl font-semibold">
            Drop files here
          </div>
        ) : (
          <>
            <div className="text-6xl mb-4">📁</div>
            <h1 className="text-2xl font-semibold text-gray-700 mb-2">
              Drop PDFs Here
            </h1>
            <p className="text-gray-500 mb-4">
              Drag and drop files, a folder, or multiple items
            </p>
            <p className="text-sm text-gray-400">
              Files will be copied to staging area for processing
            </p>
          </>
        )}

        {error && (
          <div className="mt-4 p-3 bg-red-100 text-red-700 rounded-md">
            {error}
          </div>
        )}
      </div>

      <ModeDialog
        isOpen={showModeDialog}
        onSelect={handleModeSelect}
        existingCount={existingCount}
      />
    </div>
  );
}
