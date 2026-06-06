"use client";

import { useMemo, useState } from "react";
import { Plus, Search, Filter, Check, Building2 } from "lucide-react";
import { useBuildings } from "@/components/BuildingsProvider";
import BuildingCard from "@/components/binalar/BuildingCard";
import NewBuildingModal from "@/components/binalar/NewBuildingModal";
import { ALL_STATUSES, type BuildingStatus } from "@/lib/buildings";

export default function BinalarPage() {
  const { buildings } = useBuildings();
  const [query, setQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<BuildingStatus | null>(null);
  const [filterOpen, setFilterOpen] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    return buildings.filter((b) => {
      const matchesQuery =
        !q ||
        b.name.toLowerCase().includes(q) ||
        b.district.toLowerCase().includes(q) ||
        b.address.toLowerCase().includes(q);
      const matchesStatus = !statusFilter || b.status === statusFilter;
      return matchesQuery && matchesStatus;
    });
  }, [buildings, query, statusFilter]);

  return (
    <main className="scroll-thin min-w-0 flex-1 overflow-y-auto p-6">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h1 className="text-[24px] font-bold tracking-tight text-ink-900">
              Binalar
            </h1>
            <p className="mt-1 text-[13.5px] text-ink-500">
              Kentsel dönüşüm sürecindeki tüm binaların listesi.
            </p>
          </div>
          <button
            type="button"
            onClick={() => setModalOpen(true)}
            className="flex items-center gap-2 rounded-xl bg-indigo-600 px-4 py-2.5 text-[13.5px] font-semibold text-white shadow-soft transition-colors hover:bg-indigo-700"
          >
            <Plus className="h-4 w-4" />
            Yeni Süreç Başlat
          </button>
        </div>

        {/* Search + filter */}
        <div className="mt-6 flex items-center gap-3">
          <div className="relative flex-1">
            <Search className="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-ink-400" />
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Bina adı veya adres ile ara..."
              className="w-full rounded-xl border border-slate-200 bg-white py-3 pl-10 pr-4 text-[13.5px] text-ink-900 placeholder:text-ink-400 outline-none transition-colors focus:border-indigo-300 focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          <div className="relative">
            <button
              type="button"
              onClick={() => setFilterOpen((v) => !v)}
              className={[
                "flex items-center gap-2 rounded-xl border bg-white px-4 py-3 text-[13.5px] font-medium transition-colors",
                statusFilter
                  ? "border-indigo-300 text-indigo-700"
                  : "border-slate-200 text-ink-700 hover:bg-slate-50",
              ].join(" ")}
            >
              <Filter className="h-4 w-4" />
              {statusFilter ?? "Filtrele"}
            </button>

            {filterOpen && (
              <div className="absolute right-0 z-20 mt-2 w-56 overflow-hidden rounded-xl border border-slate-200 bg-white py-1 shadow-soft">
                <FilterOption
                  label="Tümü"
                  active={statusFilter === null}
                  onClick={() => {
                    setStatusFilter(null);
                    setFilterOpen(false);
                  }}
                />
                {ALL_STATUSES.map((s) => (
                  <FilterOption
                    key={s}
                    label={s}
                    active={statusFilter === s}
                    onClick={() => {
                      setStatusFilter(s);
                      setFilterOpen(false);
                    }}
                  />
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Grid */}
        {filtered.length > 0 ? (
          <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-3">
            {filtered.map((b) => (
              <BuildingCard key={b.id} building={b} />
            ))}
          </div>
        ) : (
          <div className="mt-10 flex flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white py-16 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-400">
              <Building2 className="h-6 w-6" />
            </div>
            <p className="mt-3 text-[14px] font-medium text-ink-700">
              Sonuç bulunamadı
            </p>
            <p className="mt-1 text-[12.5px] text-ink-400">
              Arama veya filtre kriterlerinizi değiştirin.
            </p>
          </div>
        )}
      </div>

      <NewBuildingModal open={modalOpen} onClose={() => setModalOpen(false)} />
    </main>
  );
}

function FilterOption({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="flex w-full items-center justify-between px-3 py-2 text-left text-[13px] text-ink-700 transition-colors hover:bg-slate-50"
    >
      {label}
      {active && <Check className="h-4 w-4 text-indigo-600" />}
    </button>
  );
}
