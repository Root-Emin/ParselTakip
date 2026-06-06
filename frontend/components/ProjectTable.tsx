import { ChevronRight, Building } from "lucide-react";
import { projects, type StageLabel, type DocStatus } from "@/lib/data";

const stageStyles: Record<StageLabel, string> = {
  İnceleme: "bg-brand-50 text-brand-700",
  "Belge Toplama": "bg-emerald-50 text-emerald-700",
  Onay: "bg-violet-50 text-violet-700",
  Yıkım: "bg-orange-50 text-orange-700",
  "Yeniden Yapım": "bg-rose-50 text-rose-700",
};

function DocBadge({ status }: { status: DocStatus }) {
  const ok = status === "Tamam";
  return (
    <span className="inline-flex items-center gap-1.5 text-[12px] font-medium">
      <span
        className={`h-2 w-2 rounded-full ${ok ? "bg-emerald-500" : "bg-orange-400"}`}
      />
      <span className={ok ? "text-ink-700" : "text-ink-700"}>{status}</span>
    </span>
  );
}

export default function ProjectTable() {
  return (
    <section className="rounded-2xl border border-slate-200 bg-white shadow-card">
      <div className="flex items-center justify-between px-5 py-4">
        <h2 className="text-[15px] font-semibold text-ink-900">Proje Takibi</h2>
        <a
          href="#"
          className="flex items-center gap-1 text-[12.5px] font-medium text-brand-600 hover:text-brand-700"
        >
          Tüm Projeleri Gör
          <ChevronRight className="h-3.5 w-3.5" />
        </a>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full min-w-[760px] border-collapse text-left">
          <thead>
            <tr className="border-y border-slate-100 text-[11.5px] font-medium uppercase tracking-wide text-ink-400">
              <th className="px-5 py-2.5 font-medium">Bina / Site Adı</th>
              <th className="px-3 py-2.5 font-medium">Müteahhit</th>
              <th className="px-3 py-2.5 font-medium">Mevcut Aşama</th>
              <th className="px-3 py-2.5 font-medium">Tamamlanma</th>
              <th className="px-3 py-2.5 font-medium">Belge Durumu</th>
              <th className="px-5 py-2.5 font-medium">Son Güncelleme</th>
            </tr>
          </thead>
          <tbody className="text-[13px]">
            {projects.map((p) => (
              <tr
                key={p.name}
                className={[
                  "border-b border-slate-50 transition-colors last:border-0",
                  p.selected
                    ? "bg-brand-50/50"
                    : "hover:bg-slate-50/70",
                ].join(" ")}
              >
                <td className="px-5 py-3.5">
                  <div className="flex items-center gap-3">
                    <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-ink-400">
                      <Building className="h-[18px] w-[18px]" />
                    </div>
                    <div className="leading-tight">
                      <p className="font-semibold text-ink-900">{p.name}</p>
                      <p className="text-[11.5px] text-ink-400">{p.district}</p>
                    </div>
                  </div>
                </td>
                <td className="px-3 py-3.5 text-ink-700">{p.contractor}</td>
                <td className="px-3 py-3.5">
                  <span
                    className={`inline-flex rounded-full px-2.5 py-1 text-[11.5px] font-medium ${stageStyles[p.stage]}`}
                  >
                    {p.stage}
                  </span>
                </td>
                <td className="px-3 py-3.5">
                  <div className="flex items-center gap-2">
                    <div className="h-1.5 w-24 overflow-hidden rounded-full bg-slate-100">
                      <div
                        className={`h-full rounded-full ${p.progressColor}`}
                        style={{ width: `${p.progress}%` }}
                      />
                    </div>
                    <span className="text-[11.5px] font-medium text-ink-500">
                      %{p.progress}
                    </span>
                  </div>
                </td>
                <td className="px-3 py-3.5">
                  <DocBadge status={p.docStatus} />
                </td>
                <td className="px-5 py-3.5 text-[12.5px] text-ink-500">
                  {p.updated}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
