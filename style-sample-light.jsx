import { useState } from "react";

const mockDocs = [
  {
    id: 1,
    name: "Invoice_2024_Q3.pdf",
    suggestedName: "2024-Q3 Invoice — Acme Corp",
    pages: 3,
    confidence: 0.97,
    tags: ["invoice", "finance"],
    color: "#2563eb",
  },
  {
    id: 2,
    name: "scan_0042.pdf",
    suggestedName: "Medical Record — Jan 2024",
    pages: 7,
    confidence: 0.84,
    tags: ["medical", "record"],
    color: "#0891b2",
  },
  {
    id: 3,
    name: "document_scan.pdf",
    suggestedName: "Lease Agreement — Oak St",
    pages: 12,
    confidence: 0.91,
    tags: ["legal", "housing"],
    color: "#7c3aed",
  },
];

const tagStyles = {
  invoice: { bg: "#fef9c3", text: "#854d0e" },
  finance: { bg: "#dcfce7", text: "#166534" },
  medical: { bg: "#dbeafe", text: "#1e40af" },
  record: { bg: "#e0f2fe", text: "#0c4a6e" },
  legal: { bg: "#ede9fe", text: "#5b21b6" },
  housing: { bg: "#ffedd5", text: "#9a3412" },
};

function ConfidenceDots({ value }) {
  const filled = Math.round(value * 5);
  return (
    <div className="flex items-center gap-1">
      {Array.from({ length: 5 }).map((_, i) => (
        <div
          key={i}
          className="w-1.5 h-1.5 rounded-full"
          style={{
            background: i < filled
              ? value > 0.9 ? "#16a34a" : value > 0.75 ? "#d97706" : "#dc2626"
              : "#e5e7eb",
          }}
        />
      ))}
      <span className="text-xs text-gray-400 ml-1">{Math.round(value * 100)}%</span>
    </div>
  );
}

function PageCard({ num, color, active }) {
  return (
    <div
      className="relative w-9 h-11 rounded-sm cursor-grab select-none transition-all hover:-translate-y-0.5"
      style={{
        background: active ? "#fff" : "#f9fafb",
        border: `1px solid ${active ? color : "#e5e7eb"}`,
        boxShadow: active ? `0 0 0 1px ${color}20, 0 2px 8px ${color}15` : "0 1px 3px rgba(0,0,0,0.06)",
      }}
    >
      <div className="absolute inset-x-1.5 top-1.5 h-0.5 rounded-full" style={{ background: active ? color : "#e5e7eb" }} />
      <div className="absolute inset-x-1.5 top-4 h-px rounded-full" style={{ background: active ? `${color}60` : "#f3f4f6" }} />
      <div className="absolute inset-x-1.5 top-5.5 h-px rounded-full" style={{ background: active ? `${color}40` : "#f3f4f6" }} />
      <div className="absolute bottom-1 right-1 text-[8px] font-mono" style={{ color: active ? color : "#d1d5db" }}>
        {num}
      </div>
    </div>
  );
}

function DocRow({ doc, selected, onSelect }) {
  return (
    <div
      onClick={onSelect}
      className="group flex items-start gap-5 px-6 py-4 border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50/70"
      style={{ background: selected ? `${doc.color}05` : undefined }}
    >
      {/* Color accent + thumb */}
      <div
        className="w-1 self-stretch rounded-full shrink-0 transition-opacity"
        style={{ background: doc.color, opacity: selected ? 1 : 0.2 }}
      />

      {/* Main */}
      <div className="flex-1 min-w-0">
        <div className="flex items-baseline gap-2 mb-0.5">
          <span className="text-sm font-semibold text-gray-800 truncate" style={{ fontFamily: "'Instrument Serif', Georgia, serif" }}>
            {doc.suggestedName}
          </span>
          <span className="text-xs text-gray-400 font-mono shrink-0">{doc.pages} pages</span>
        </div>
        <div className="text-xs text-gray-400 mb-2.5 truncate">{doc.name}</div>
        <div className="flex items-center gap-3 flex-wrap">
          <ConfidenceDots value={doc.confidence} />
          <div className="flex gap-1.5">
            {doc.tags.map(t => {
              const s = tagStyles[t] || { bg: "#f3f4f6", text: "#6b7280" };
              return (
                <span
                  key={t}
                  className="text-[11px] px-2 py-0.5 rounded-full font-medium"
                  style={{ background: s.bg, color: s.text }}
                >
                  {t}
                </span>
              );
            })}
          </div>
        </div>
      </div>

      {/* Page cards */}
      <div className="flex gap-1.5 shrink-0 items-end">
        {Array.from({ length: Math.min(doc.pages, 5) }).map((_, i) => (
          <PageCard key={i} num={i + 1} color={doc.color} active={selected} />
        ))}
        {doc.pages > 5 && (
          <div className="w-9 h-11 rounded-sm border border-dashed border-gray-200 flex items-center justify-center text-[10px] text-gray-400 font-mono">
            +{doc.pages - 5}
          </div>
        )}
      </div>

      {/* Row actions */}
      <div className={`flex flex-col gap-1 shrink-0 transition-opacity ${selected ? "opacity-100" : "opacity-0 group-hover:opacity-100"}`}>
        <button className="text-xs px-2.5 py-1 rounded border border-gray-200 text-gray-500 hover:text-gray-800 hover:border-gray-300 transition-colors bg-white">
          rename
        </button>
        <button className="text-xs px-2.5 py-1 rounded border border-gray-200 text-gray-500 hover:text-gray-800 hover:border-gray-300 transition-colors bg-white">
          split
        </button>
      </div>
    </div>
  );
}

export default function LightEditorialSample() {
  const [selected, setSelected] = useState(null);

  return (
    <div className="min-h-screen bg-white flex flex-col" style={{ fontFamily: "'DM Sans', 'Helvetica Neue', sans-serif" }}>
      {/* Titlebar */}
      <div
        className="flex items-center justify-between px-6 py-3 border-b"
        style={{ borderColor: "#f0f0f0" }}
      >
        <div className="flex items-center gap-3">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-gray-200" />
            <div className="w-3 h-3 rounded-full bg-gray-200" />
            <div className="w-3 h-3 rounded-full bg-gray-200" />
          </div>
          <span className="text-sm font-medium text-gray-800 ml-1" style={{ fontFamily: "'Instrument Serif', Georgia, serif" }}>
            Scan Organizer
          </span>
        </div>
        <span className="text-xs text-gray-400">3 documents · 22 pages</span>
      </div>

      {/* Toolbar */}
      <div className="flex items-center gap-2 px-6 py-3 border-b border-gray-100">
        <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-gray-600 hover:text-gray-900 border border-gray-200 hover:border-gray-300 rounded-md bg-white transition-colors">
          <span>+</span> Import
        </button>
        <button className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-gray-600 hover:text-gray-900 border border-gray-200 hover:border-gray-300 rounded-md bg-white transition-colors">
          Auto-name all
        </button>
        <div className="flex-1" />
        <button
          className="px-4 py-1.5 text-sm font-medium text-white rounded-md transition-colors"
          style={{ background: "#1a1a1a" }}
        >
          Export →
        </button>
      </div>

      {/* Column headers */}
      <div className="flex items-center gap-5 px-6 py-2 border-b border-gray-100">
        <div className="w-1 shrink-0" />
        <span className="flex-1 text-[11px] font-semibold text-gray-400 uppercase tracking-wider">Document</span>
        <span className="w-32 text-[11px] font-semibold text-gray-400 uppercase tracking-wider">Confidence</span>
        <span className="w-56 text-[11px] font-semibold text-gray-400 uppercase tracking-wider">Pages</span>
        <span className="w-16" />
      </div>

      {/* Rows */}
      <div className="flex-1">
        {mockDocs.map(doc => (
          <DocRow
            key={doc.id}
            doc={doc}
            selected={selected === doc.id}
            onSelect={() => setSelected(selected === doc.id ? null : doc.id)}
          />
        ))}
      </div>

      {/* Footer */}
      <div className="flex items-center px-6 py-2 border-t border-gray-100">
        <span className="text-xs text-gray-300">drag pages to reorder · drag to merge</span>
      </div>
    </div>
  );
}
