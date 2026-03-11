import { useState, useEffect, useCallback } from 'react';
import { GetPreview, Rename, Append, Trash, GetOutputFoldersRecursive, GetInputFiles } from '../../wailsjs/wailsjs/go/main/App';

interface PreviewData {
  preview: string;
  previousRenamed: string[];
  folders: string[];
}

export function usePreview() {
  const [data, setData] = useState<PreviewData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await GetPreview();
      setData(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const rename = async (newName: string, folder: string) => {
    await Rename(newName, folder);
    await refresh();
  };

  const append = async (target: string) => {
    await Append(target);
    await refresh();
  };

  const trash = async () => {
    await Trash();
    await refresh();
  };

  return { data, loading, error, refresh, rename, append, trash };
}

export function useFolders() {
  const [folders, setFolders] = useState<string[]>([]);

  const refresh = useCallback(async () => {
    try {
      const result = await GetOutputFoldersRecursive();
      setFolders(result);
    } catch (e) {
      console.error('Failed to load folders:', e);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { folders, refresh };
}

export function useInputFiles() {
  const [files, setFiles] = useState<string[]>([]);

  const refresh = useCallback(async () => {
    try {
      const result = await GetInputFiles();
      setFiles(result);
    } catch (e) {
      console.error('Failed to load input files:', e);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { files, refresh };
}
