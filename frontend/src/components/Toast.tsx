interface ToastProps {
  message: string | null;
  onClose?: () => void;
}

export function Toast({ message, onClose }: ToastProps) {
  if (!message) return null;

  return (
    <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
      <div className="flex items-center justify-between">
        <span className="text-red-800 text-sm">{message}</span>
        {onClose && (
          <button
            onClick={onClose}
            className="ml-2 text-red-600 hover:text-red-800 font-bold"
          >
            &times;
          </button>
        )}
      </div>
    </div>
  );
}
