import { useState, useCallback } from 'react';
import * as App from '../../bindings/github.com/joshrainwater/scan-organizer/app';
import { PreviewData } from '../../bindings/github.com/joshrainwater/scan-organizer/internal/scanorganizer/models';

export function usePreview() {
  const [data, setData] = useState<PreviewData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await App.GetPreview();
      if (result) {
        setData(result);
      } else {
        setData(null);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const rename = useCallback(async (newName: string, folder: string) => {
    setError(null);
    try {
      await App.Rename(newName, folder);
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      throw e;
    }
  }, [refresh]);

  const append = useCallback(async (target: string) => {
    setError(null);
    try {
      await App.Append(target);
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      throw e;
    }
  }, [refresh]);

  const trash = useCallback(async () => {
    setError(null);
    try {
      await App.Trash();
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      throw e;
    }
  }, [refresh]);

  return {
    data,
    loading,
    error,
    rename,
    append,
    trash,
    refresh,
  };
}
