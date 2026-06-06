import Link from "next/link";
import { Users, ChevronRight, Trash2 } from "lucide-react";
import { type Contractor } from "@/lib/contractors";

export default function ContractorCard({
  contractor,
  activeProjects,
  onDelete,
}: {
  contractor: Contractor;
  activeProjects: number;
  onDelete: (id: string) => void;
}) {
  return (
    <div className="group relative flex flex-col rounded-2xl border border-slate-200 bg-white p-5 shadow-card transition-all hover:-translate-y-0.5 hover:shadow-soft">
      <button
        type="button"
        aria-label="Müteahhiti sil"
        onClick={() => onDelete(contractor.id)}
        className="absolute right-3 top-3 flex h-8 w-8 items-center justify-center rounded-lg text-ink-400 opacity-0 transition-all hover:bg-rose-50 hover:text-rose-500 group-hover:opacity-100"
      >
        <Trash2 className="h-4 w-4" />
      </button>

      <Link href={`/muteahhitler/${contractor.id}`} className="flex flex-col">
        <div className="flex items-center gap-3.5">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-indigo-50 text-indigo-600">
            <Users className="h-6 w-6" />
          </div>
          <div className="min-w-0">
            <h3 className="truncate pr-6 text-[16px] font-bold tracking-tight text-ink-900">
              {contractor.name}
            </h3>
            <p className="text-[12.5px] text-ink-500">
              Aktif Proje: {activeProjects}
            </p>
          </div>
        </div>

        <div className="my-4 h-px bg-slate-100" />

        <span className="flex items-center justify-end gap-1 text-[13px] font-semibold text-indigo-600">
          Profili İncele
          <ChevronRight className="h-3.5 w-3.5" />
        </span>
      </Link>
    </div>
  );
}
