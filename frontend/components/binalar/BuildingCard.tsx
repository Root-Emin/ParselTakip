import Link from "next/link";
import {
  MapPin,
  HardHat,
  Building2,
  CheckCircle2,
  FileWarning,
} from "lucide-react";
import { type Building, docsDone } from "@/lib/buildings";
import StatusBadge from "./StatusBadge";

export default function BuildingCard({ building }: { building: Building }) {
  const done = docsDone(building);
  const total = building.documents.length;
  const complete = done === total;

  return (
    <Link
      href={`/binalar/${building.id}`}
      className="group flex flex-col overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-card transition-all hover:-translate-y-0.5 hover:shadow-soft"
    >
      {/* Image header */}
      <div className="relative flex h-40 items-center justify-center bg-gradient-to-br from-slate-100 to-slate-200">
        <Building2 className="h-12 w-12 text-slate-300" />
        <StatusBadge status={building.status} className="absolute right-3 top-3" />
      </div>

      {/* Body */}
      <div className="flex flex-1 flex-col p-5">
        <h3 className="text-[16px] font-bold tracking-tight text-ink-900">
          {building.name}
        </h3>
        <p className="mt-1 flex items-start gap-1.5 text-[12.5px] text-ink-500">
          <MapPin className="mt-0.5 h-3.5 w-3.5 shrink-0 text-ink-400" />
          <span>
            {building.district}, {building.address}
          </span>
        </p>

        <div className="my-4 h-px bg-slate-100" />

        <div className="space-y-2.5 text-[13px]">
          <div className="flex items-center justify-between">
            <span className="flex items-center gap-2 text-ink-500">
              <HardHat className="h-4 w-4 text-ink-400" />
              Müteahhit:
            </span>
            <span className="font-semibold text-ink-900">
              {building.contractor ?? "Atanmadı"}
            </span>
          </div>

          <div className="flex items-center justify-between">
            <span className="flex items-center gap-2 text-ink-500">
              {complete ? (
                <CheckCircle2 className="h-4 w-4 text-emerald-500" />
              ) : (
                <FileWarning className="h-4 w-4 text-amber-500" />
              )}
              Belgeler:
            </span>
            <span
              className={`font-semibold ${complete ? "text-emerald-600" : "text-amber-600"}`}
            >
              {done} / {total}
            </span>
          </div>
        </div>

        <div className="mt-4">
          <div className="flex items-center justify-between text-[12px]">
            <span className="text-ink-500">Süreç İlerlemesi</span>
            <span className="font-semibold text-ink-900">%{building.progress}</span>
          </div>
          <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-slate-100">
            <div
              className="h-full rounded-full bg-indigo-500"
              style={{ width: `${building.progress}%` }}
            />
          </div>
        </div>
      </div>
    </Link>
  );
}
