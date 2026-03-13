import { useState, useCallback, useRef } from 'react';

interface DropEvent {
  preventDefault(): void;
  stopPropagation(): void;
  dataTransfer?: DataTransfer;
}

export function useDropzone(onDrop: (paths: string[]) => void) {
  const [isDragging, setIsDragging] = useState(false);
  const dragCounter = useRef(0);

  const handleDragEnter = useCallback((e: DropEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounter.current++;
    if (e.dataTransfer?.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  }, []);

  const handleDragLeave = useCallback((e: DropEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounter.current--;
    if (dragCounter.current === 0) {
      setIsDragging(false);
    }
  }, []);

  const handleDragOver = useCallback((e: DropEvent) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const processEntry = async (
    entry: FileSystemEntry,
    path: string
  ): Promise<string[]> => {
    if (entry.isFile) {
      const file = entry as FileSystemFileEntry;
      return new Promise((resolve) => {
        file.file((f) => {
          if (f.name.toLowerCase().endsWith('.pdf')) {
            resolve([path + '/' + f.name]);
          } else {
            resolve([]);
          }
        });
      });
    } else if (entry.isDirectory) {
      const dir = entry as FileSystemDirectoryEntry;
      const reader = dir.createReader();
      return new Promise((resolve) => {
        reader.readEntries(async (entries) => {
          const allPaths: string[] = [];
          for (const ent of entries) {
            const paths = await processEntry(ent, path + '/' + dir.name);
            allPaths.push(...paths);
          }
          resolve(allPaths);
        });
      });
    }
    return [];
  };

  const handleDrop = useCallback(async (e: DropEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    dragCounter.current = 0;

    const items = e.dataTransfer?.items;
    if (!items) {
      return;
    }

    const paths: string[] = [];

    for (let i = 0; i < items.length; i++) {
      const item = items[i].webkitGetAsEntry();
      if (item) {
        const entryPaths = await processEntry(item, '');
        paths.push(...entryPaths);
      }
    }

    if (paths.length > 0) {
      onDrop(paths);
    }
  }, [onDrop]);

  return {
    isDragging,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
  };
}
