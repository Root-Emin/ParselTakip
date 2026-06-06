import {
  X,
  Building2,
  ChevronRight,
  FileText,
  ExternalLink,
} from "lucide-react";
import { missingDocs } from "@/lib/data";

const DOC_PERCENT = 68;
const ringCirc = 2 * Math.PI * 26;

export default function DetailPanel() {
  return (
    <aside className="scroll-thin flex h-full w-[320px] shrink-0 flex-col gap-4 overflow-y-auto border-l border-slate-200 bg-white px-4 py-5">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div className="leading-tight">
          <h2 className="text-[15px] font-bold text-ink-900">
            Yıldız Park Evleri
          </h2>
          <p className="text-[12px] text-ink-400">Beşiktaş / İstanbul</p>
        </div>
        <button
          type="button"
          aria-label="Kapat"
          className="flex h-7 w-7 items-center justify-center rounded-lg text-ink-400 transition-colors hover:bg-slate-100 hover:text-ink-700"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      {/* Building image + summary */}
      <div className="flex gap-3">
        <div className="flex h-[88px] w-[110px] shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-slate-100 to-slate-200 text-slate-400">
          <Building2 className="h-9 w-9" />
        </div>
        <div className="flex flex-1 flex-col justify-center gap-2 text-[12px]">
          <div>
            <p className="text-ink-400">Mevcut Aşama</p>
            <span className="mt-1 inline-flex rounded-full bg-brand-50 px-2.5 py-0.5 text-[11.5px] font-medium text-brand-700">
              İnceleme
            </span>
          </div>
          <div>
            <div className="flex items-center justify-between">
              <p className="text-ink-400">Tamamlanma Oranı</p>
              <p className="font-semibold text-ink-900">%45</p>
            </div>
            <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-slate-100">
              <div
                className="h-full rounded-full bg-brand-500"
                style={{ width: "45%" }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Belge Durumu */}
      <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-card">
        <h3 className="text-[13.5px] font-semibold text-ink-900">
          Belge Durumu
        </h3>
        <div className="mt-3 flex items-center gap-4">
          <div className="relative h-[72px] w-[72px] shrink-0">
            <svg viewBox="0 0 64 64" className="h-full w-full -rotate-90">
              <circle
                cx="32"
                cy="32"
                r="26"
                fill="none"
                stroke="#f1f5f9"
                strokeWidth="7"
              />
              <circle
                cx="32"
                cy="32"
                r="26"
                fill="none"
                stroke="#22c55e"
                strokeWidth="7"
                strokeLinecap="round"
                strokeDasharray={`${(DOC_PERCENT / 100) * ringCirc} ${ringCirc}`}
              />
            </svg>
            <div className="absolute inset-0 flex items-center justify-center">
              <span className="text-[15px] font-bold text-ink-900">
                {DOC_PERCENT}%
              </span>
            </div>
          </div>
          <div className="leading-tight">
            <p className="text-[12.5px] font-medium text-ink-700">
              Belgeler Tamamlandı
            </p>
            <p className="text-[12px] text-ink-400">15 / 22</p>
          </div>
        </div>
        <a
          href="#"
          className="mt-3 flex items-center justify-center gap-1 rounded-lg border border-slate-200 py-2 text-[12.5px] font-medium text-brand-600 transition-colors hover:bg-brand-50"
        >
          Tüm Belgeler
          <ChevronRight className="h-3.5 w-3.5" />
        </a>
      </div>

      {/* Eksik Belgeler */}
      <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-card">
        <div className="flex items-center gap-2">
          <h3 className="text-[13.5px] font-semibold text-ink-900">
            Eksik Belgeler
          </h3>
          <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-rose-100 px-1.5 text-[11px] font-semibold text-rose-600">
            {missingDocs.length}
          </span>
        </div>

        <ul className="mt-3 space-y-1">
          {missingDocs.map((doc) => (
            <li
              key={doc}
              className="flex items-center justify-between rounded-lg px-1 py-1.5"
            >
              <span className="flex items-center gap-2 text-[12.5px] text-ink-700">
                <FileText className="h-3.5 w-3.5 text-ink-400" />
                {doc}
              </span>
              <span className="rounded-md bg-rose-50 px-2 py-0.5 text-[10.5px] font-medium text-rose-500">
                Eksik
              </span>
            </li>
          ))}
        </ul>
      </div>

      {/* CTA */}
      <button
        type="button"
        className="mt-auto flex items-center justify-center gap-2 rounded-xl bg-brand-600 py-3 text-[13.5px] font-semibold text-white shadow-soft transition-colors hover:bg-brand-700"
      >
        Detaylı İncele
        <ExternalLink className="h-4 w-4" />
      </button>
    </aside>
  );
}
