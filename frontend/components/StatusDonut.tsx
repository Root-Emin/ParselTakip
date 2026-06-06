import { RefreshCw } from "lucide-react";
import { donutSegments, donutTotal } from "@/lib/data";

const RADIUS = 56;
const STROKE = 18;
const CIRC = 2 * Math.PI * RADIUS;
const GAP = 2; // visual gap between segments (in circumference units)

export default function StatusDonut() {
  let offset = 0;

  return (
    <section className="flex h-full flex-col rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
      <h2 className="text-[15px] font-semibold text-ink-900">
        Proje Durum Dağılımı
      </h2>

      <div className="mt-4 flex flex-1 items-center gap-6">
        {/* Donut */}
        <div className="relative h-[160px] w-[160px] shrink-0">
          <svg
            viewBox="0 0 160 160"
            className="h-full w-full -rotate-90"
            aria-hidden
          >
            <circle
              cx="80"
              cy="80"
              r={RADIUS}
              fill="none"
              stroke="#f1f5f9"
              strokeWidth={STROKE}
            />
            {donutSegments.map((seg) => {
              const fraction = seg.value / donutTotal;
              const len = fraction * CIRC - GAP;
              const dash = `${Math.max(len, 0)} ${CIRC - Math.max(len, 0)}`;
              const circle = (
                <circle
                  key={seg.label}
                  cx="80"
                  cy="80"
                  r={RADIUS}
                  fill="none"
                  stroke={seg.color}
                  strokeWidth={STROKE}
                  strokeDasharray={dash}
                  strokeDashoffset={-offset}
                  strokeLinecap="round"
                />
              );
              offset += fraction * CIRC;
              return circle;
            })}
          </svg>
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <span className="text-[26px] font-bold leading-none text-ink-900">
              {donutTotal}
            </span>
            <span className="mt-1 text-[11px] font-medium text-ink-400">
              Toplam
            </span>
          </div>
        </div>

        {/* Legend */}
        <ul className="flex-1 space-y-2">
          {donutSegments.map((seg) => (
            <li
              key={seg.label}
              className="flex items-center gap-2 text-[12.5px]"
            >
              <span
                className="h-2.5 w-2.5 shrink-0 rounded-full"
                style={{ backgroundColor: seg.color }}
              />
              <span className="flex-1 text-ink-700">{seg.label}</span>
              <span className="font-medium text-ink-900">{seg.value}</span>
              <span className="w-10 text-right text-ink-400">
                ({seg.percent})
              </span>
            </li>
          ))}
        </ul>
      </div>

      <div className="mt-4 flex items-center gap-1.5 border-t border-slate-100 pt-3 text-[11px] text-ink-400">
        <RefreshCw className="h-3 w-3" />
        Son güncelleme: 21.05.2024 10:30
      </div>
    </section>
  );
}
