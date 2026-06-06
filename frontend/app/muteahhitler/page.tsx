"use client";

import { useMemo, useState } from "react";
import { Plus, Search, Users } from "lucide-react";
import { useContractors } from "@/components/ContractorsProvider";
import { useBuildings } from "@/components/BuildingsProvider";
import { activeProjectCount } from "@/lib/contractors";
import ContractorCard from "@/components/muteahhitler/ContractorCard";
import NewContractorModal from "@/components/muteahhitler/NewContractorModal";

export default function MuteahhitlerPage() {
  const { contractors, deleteContractor } = useContractors();
  const { buildings } = useBuildings();
  const [query, setQuery] = useState("");
  const [modalOpen, setModalOpen] = useState(false);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return contractors;
    return contractors.filter(
      (c) =>
        c.name.toLowerCase().includes(q) ||
        c.contactPerson.toLowerCase().includes(q),
    );
  }, [contractors, query]);

  function handleDelete(id: string) {
    const c = contractors.find((x) => x.id === id);
    if (!c) return;
    if (
      window.confirm(`"${c.name}" müteahhitini silmek istediğinize emin misiniz?`)
    ) {
      deleteContractor(id);
    }
  }

  return (
    <main className="scroll-thin min-w-0 flex-1 overflow-y-auto p-6">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h1 className="text-[24px] font-bold tracking-tight text-ink-900">
              Müteahhitler
            </h1>
            <p className="mt-1 text-[13.5px] text-ink-500">
              Sistemdeki müteahhit firmaların listesi.
            </p>
          </div>
          <button
            type="button"
            onClick={() => setModalOpen(true)}
            className="flex items-center gap-2 rounded-xl bg-indigo-600 px-4 py-2.5 text-[13.5px] font-semibold text-white shadow-soft transition-colors hover:bg-indigo-700"
          >
            <Plus className="h-4 w-4" />
            Yeni Müteahhit
          </button>
        </div>

        {/* Search */}
        <div className="relative mt-6">
          <Search className="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-ink-400" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Firma adı ile ara..."
            className="w-full rounded-xl border border-slate-200 bg-white py-3 pl-10 pr-4 text-[13.5px] text-ink-900 placeholder:text-ink-400 outline-none transition-colors focus:border-indigo-300 focus:ring-2 focus:ring-indigo-100"
          />
        </div>

        {/* Grid */}
        {filtered.length > 0 ? (
          <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-3">
            {filtered.map((c) => (
              <ContractorCard
                key={c.id}
                contractor={c}
                activeProjects={activeProjectCount(buildings, c.name)}
                onDelete={handleDelete}
              />
            ))}
          </div>
        ) : (
          <div className="mt-10 flex flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white py-16 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-400">
              <Users className="h-6 w-6" />
            </div>
            <p className="mt-3 text-[14px] font-medium text-ink-700">
              Müteahhit bulunamadı
            </p>
            <p className="mt-1 text-[12.5px] text-ink-400">
              Arama kriterinizi değiştirin veya yeni müteahhit ekleyin.
            </p>
          </div>
        )}
      </div>

      <NewContractorModal open={modalOpen} onClose={() => setModalOpen(false)} />
    </main>
  );
}
