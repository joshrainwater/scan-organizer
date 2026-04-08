interface ToastProps {
  message: string | null;
  variant?: 'error' | 'success';
  onClose?: () => void;
}

export function Toast({ message, variant = 'error', onClose }: ToastProps) {
  if (!message) return null;

  const isError = variant === 'error';
  const bgClass = isError ? 'bg-red-50' : 'bg-green-50';
  const borderClass = isError ? 'border-red-200' : 'border-green-200';
  const textClass = isError ? 'text-red-800' : 'text-green-800';
  const buttonClass = isError ? 'text-red-600 hover:text-red-800' : 'text-green-600 hover:text-green-800';

  return (
    <div className={`mb-4 p-3 ${bgClass} border ${borderClass} rounded-md`}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {!isError && <span className="text-green-600">✓</span>}
          <span className={`text-sm ${textClass}`}>{message}</span>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className={`ml-2 ${buttonClass} font-bold`}
          >
            &times;
          </button>
        )}
      </div>
    </div>
  );
}
