import {
  FileText,
  FolderOpen,
  Search,
  CheckCircle2,
  Hammer,
  Building2,
  PackageCheck,
  Check,
  type LucideIcon,
} from "lucide-react";
import { timelineStages } from "@/lib/data";

const icons: LucideIcon[] = [
  FileText,
  FolderOpen,
  Search,
  CheckCircle2,
  Hammer,
  Building2,
  PackageCheck,
];

export default function ProcessTimeline() {
  return (
    <section className="flex h-full flex-col rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
      <h2 className="text-[15px] font-semibold text-ink-900">
        Dönüşüm Süreci Aşamaları
      </h2>

      <div className="mt-8 flex items-start justify-between">
        {timelineStages.map((stage, i) => {
          const Icon = icons[i];
          const isLast = i === timelineStages.length - 1;
          const done = stage.status === "done";
          const active = stage.status === "active";

          return (
            <div
              key={stage.label}
              className="relative flex flex-1 flex-col items-center"
            >
              {/* Connector */}
              {!isLast && (
                <span
                  className={[
                    "absolute left-1/2 top-5 h-0.5 w-full",
                    done ? "bg-emerald-400" : "bg-slate-200",
                  ].join(" ")}
                />
              )}

              {/* Node */}
              <div
                className={[
                  "relative z-10 flex items-center justify-center rounded-full transition-all",
                  active
                    ? "h-11 w-11 bg-brand-500 text-white shadow-soft ring-4 ring-brand-100"
                    : done
                      ? "h-9 w-9 bg-emerald-500 text-white"
                      : "h-9 w-9 border border-slate-200 bg-white text-ink-400",
                ].join(" ")}
              >
                {done ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Icon className={active ? "h-5 w-5" : "h-4 w-4"} />
                )}
              </div>

              {/* Label */}
              <p
                className={[
                  "mt-3 text-center text-[11.5px] font-medium leading-tight",
                  active
                    ? "text-brand-700"
                    : done
                      ? "text-ink-900"
                      : "text-ink-400",
                ].join(" ")}
              >
                {stage.label}
              </p>
              <p
                className={[
                  "mt-1 text-center text-[10px]",
                  active
                    ? "font-medium text-brand-500"
                    : done
                      ? "text-emerald-600"
                      : "text-ink-400",
                ].join(" ")}
              >
                {stage.status === "done"
                  ? stage.hint
                  : stage.status === "active"
                    ? stage.date
                    : stage.hint}
              </p>
              {stage.date && stage.status === "done" && (
                <p className="text-center text-[10px] text-ink-400">
                  {stage.date}
                </p>
              )}
            </div>
          );
        })}
      </div>
    </section>
  );
}
