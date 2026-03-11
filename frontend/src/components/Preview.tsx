interface PreviewProps {
  src: string | null;
  loading?: boolean;
}

export function Preview({ src, loading }: PreviewProps) {
  if (loading) {
    return (
      <div className="flex-1 bg-white flex items-center justify-center border-l border-gray-200">
        <div className="text-gray-400">Loading preview...</div>
      </div>
    );
  }

  if (!src) {
    return (
      <div className="flex-1 bg-white flex items-center justify-center border-l border-gray-200">
        <div className="text-gray-400">No PDF to display</div>
      </div>
    );
  }

  return (
    <div className="flex-1 bg-white flex items-center justify-center border-l border-gray-200 p-4">
      <img 
        src={src} 
        alt="PDF Preview" 
        className="max-h-[50rem] border border-gray-200 shadow-sm"
      />
    </div>
  );
}
