import { useEffect, useRef, useState } from 'react';
import { usePreview } from './hooks/usePreview';
import { Preview } from './components/Preview';
import { RenameForm } from './components/RenameForm';
import { AppendForm } from './components/AppendForm';
import { TrashButton } from './components/TrashButton';
import { ExportButton } from './components/ExportButton';
import { Toast } from './components/Toast';
import { Home } from './pages/Home';
import * as AppBindings from '@bindings/github.com/joshrainwater/scan-organizer/app';

interface OrganizerProps {
  outputCount: number;
  onExportSuccess: () => void;
}

function Organizer({ outputCount, onExportSuccess }: OrganizerProps) {
  const { data, loading, error, rename, append, trash, exportFiles, refresh } = usePreview();
  const folderInputRef = useRef<HTMLInputElement>(null);
  const trashButtonRef = useRef<HTMLButtonElement>(null);
  const appendSelectRef = useRef<HTMLSelectElement>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.ctrlKey) {
        switch (e.key) {
          case '1':
            e.preventDefault();
            folderInputRef.current?.focus();
            break;
          case '2':
            e.preventDefault();
            trashButtonRef.current?.focus();
            break;
          case '3':
            e.preventDefault();
            appendSelectRef.current?.focus();
            break;
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const hasFiles = data && (data.preview || loading);

  const handleExportSuccess = () => {
    setSuccessMessage('Export complete! Returning to home...');
    setTimeout(() => {
      onExportSuccess();
    }, 1500);
  };

  return (
    <div className="min-h-screen bg-gray-100 flex">
      <div className="w-1/3 p-8 flex flex-col gap-6">
        <Toast 
          message={error} 
          variant="error"
          onClose={refresh}
        />
        <Toast 
          message={successMessage} 
          variant="success"
        />

        <div>
          <h2 className="text-lg font-semibold text-gray-700 mb-4">Organize PDF</h2>
          {hasFiles ? (
            <>
              <RenameForm
                folders={data?.folders || []}
                onRename={rename}
                inputFilesLength={1}
              />
              <hr className="my-6 border-gray-300" />
              <TrashButton onTrash={trash} />
              <hr className="my-6 border-gray-300" />
              <AppendForm
                previousRenamed={data?.previousRenamed || []}
                onAppend={append}
              />
              <hr className="my-6 border-gray-300" />
              {outputCount > 0 && (
                <ExportButton
                  onExport={exportFiles}
                  fileCount={outputCount}
                  onSuccess={handleExportSuccess}
                />
              )}
            </>
          ) : (
            <div className="text-gray-500 text-center py-8">
              No PDFs found in input folder.
            </div>
          )}
        </div>
      </div>

      <Preview src={data?.preview || null} loading={loading} />
    </div>
  );
}

function App() {
  const [isReady, setIsReady] = useState(false);
  const [checkDone, setCheckDone] = useState(false);
  const [outputCount, setOutputCount] = useState(0);

  const checkStaging = async () => {
    try {
      const status = await AppBindings.GetStatus();
      setOutputCount(status.outputCount);
      if (status.inputCount > 0 || status.outputCount > 0) {
        setIsReady(true);
      }
    } catch (e) {
      console.error('Failed to check staging status:', e);
    } finally {
      setCheckDone(true);
    }
  };

  useEffect(() => {
    checkStaging();
  }, []);

  const handleReady = () => {
    setIsReady(true);
  };

  const handleExportSuccess = async () => {
    const status = await AppBindings.GetStatus();
    if (status.inputCount === 0) {
      setIsReady(false);
      setOutputCount(0);
    } else {
      setOutputCount(status.outputCount);
    }
  };

  if (!checkDone) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (!isReady) {
    return <Home onReady={handleReady} />;
  }

  return <Organizer outputCount={outputCount} onExportSuccess={handleExportSuccess} />;
}

export default App;
