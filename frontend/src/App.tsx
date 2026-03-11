import { useEffect, useRef } from 'react';
import { usePreview } from './hooks/usePreview';
import { Preview } from './components/Preview';
import { RenameForm } from './components/RenameForm';
import { AppendForm } from './components/AppendForm';
import { TrashButton } from './components/TrashButton';
import { Toast } from './components/Toast';

function App() {
  const { data, loading, error, rename, append, trash, refresh } = usePreview();
  const folderInputRef = useRef<HTMLInputElement>(null);
  const trashButtonRef = useRef<HTMLButtonElement>(null);
  const appendSelectRef = useRef<HTMLSelectElement>(null);

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

  return (
    <div className="min-h-screen bg-gray-100 flex">
      <div className="w-1/3 p-8 flex flex-col gap-6">
        <Toast 
          message={error} 
          onClose={refresh}
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

export default App;
