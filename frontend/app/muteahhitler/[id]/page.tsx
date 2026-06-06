"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  ArrowLeft,
  Trash2,
  Users,
  User,
  Phone,
  Mail,
  MapPin,
  CalendarClock,
  Building2,
} from "lucide-react";
import { useContractors } from "@/components/ContractorsProvider";
import { useBuildings } from "@/components/BuildingsProvider";
import {
  contractorProjects,
  activeProjectCount,
  completedProjectCount,
} from "@/lib/contractors";
import BuildingCard from "@/components/binalar/BuildingCard";

export default function ContractorDetailPage({
  params,
}: {
  params: { id: string };
}) {
  const router = useRouter();
  const { getContractor, deleteContractor } = useContractors();
  const { buildings } = useBuildings();
  const contractor = getContractor(params.id);

  if (!contractor) {
    return (
      <main className="flex min-w-0 flex-1 items-center justify-center p-6">
        <div className="text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-400">
            <Users className="h-6 w-6" />
          </div>
          <p className="mt-3 text-[15px] font-semibold text-ink-900">
            Müteahhit bulunamadı
          </p>
          <Link
            href="/muteahhitler"
            className="mt-3 inline-flex items-center gap-1.5 text-[13px] font-medium text-indigo-600 hover:text-indigo-700"
          >
            <ArrowLeft className="h-4 w-4" />
            Müteahhitlere dön
          </Link>
        </div>
      </main>
    );
  }

  const projects = contractorProjects(buildings, contractor.name);
  const active = activeProjectCount(buildings, contractor.name);
  const completed = completedProjectCount(buildings, contractor.name);

  function handleDelete() {
    if (!contractor) return;
    if (
      window.confirm(
        `"${contractor.name}" müteahhitini silmek istediğinize emin misiniz?`,
      )
    ) {
      deleteContractor(contractor.id);
      router.push("/muteahhitler");
    }
  }

  return (
    <main className="scroll-thin min-w-0 flex-1 overflow-y-auto p-6">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="flex items-start gap-3">
            <Link
              href="/muteahhitler"
              aria-label="Geri"
              className="mt-0.5 flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-ink-500 transition-colors hover:bg-slate-50"
            >
              <ArrowLeft className="h-4 w-4" />
            </Link>
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-indigo-50 text-indigo-600">
                <Users className="h-6 w-6" />
              </div>
              <div>
                <h1 className="text-[22px] font-bold tracking-tight text-ink-900">
                  {contractor.name}
                </h1>
                <p className="mt-0.5 text-[13px] text-ink-500">
                  {contractor.since} yılından beri • {contractor.city}
                </p>
              </div>
            </div>
          </div>

          <button
            type="button"
            onClick={handleDelete}
            className="flex items-center gap-1.5 rounded-xl border border-rose-200 bg-white px-3.5 py-2 text-[13px] font-medium text-rose-600 transition-colors hover:bg-rose-50"
          >
            <Trash2 className="h-3.5 w-3.5" />
            Sil
          </button>
        </div>

        {/* Stats */}
        <div className="mt-6 grid grid-cols-3 gap-4">
          <StatBox label="Aktif Proje" value={active} tone="indigo" />
          <StatBox label="Tamamlanan" value={completed} tone="emerald" />
          <StatBox label="Toplam Proje" value={projects.length} tone="slate" />
        </div>

        {/* Body */}
        <div className="mt-5 grid grid-cols-1 gap-5 lg:grid-cols-3">
          {/* Info */}
          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-card lg:col-span-1">
            <h2 className="text-[15px] font-semibold text-ink-900">
              Firma Bilgileri
            </h2>
            <div className="mt-4 space-y-3.5 text-[13px]">
              <InfoRow icon={User} label="Yetkili" value={contractor.contactPerson} />
              <InfoRow icon={Phone} label="Telefon" value={contractor.phone} />
              <InfoRow icon={Mail} label="E-posta" value={contractor.email} />
              <InfoRow icon={MapPin} label="Şehir" value={contractor.city} />
              <InfoRow
                icon={CalendarClock}
                label="Kuruluş"
                value={String(contractor.since)}
              />
            </div>
          </section>

          {/* Projects */}
          <section className="lg:col-span-2">
            <div className="mb-4 flex items-center gap-2">
              <h2 className="text-[15px] font-semibold text-ink-900">
                Projeler
              </h2>
              <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-indigo-100 px-1.5 text-[11px] font-semibold text-indigo-600">
                {projects.length}
              </span>
            </div>

            {projects.length > 0 ? (
              <div className="grid grid-cols-1 gap-5 sm:grid-cols-2">
                {projects.map((b) => (
                  <BuildingCard key={b.id} building={b} />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-slate-200 bg-white py-16 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-400">
                  <Building2 className="h-6 w-6" />
                </div>
                <p className="mt-3 text-[14px] font-medium text-ink-700">
                  Henüz proje atanmamış
                </p>
                <p className="mt-1 text-[12.5px] text-ink-400">
                  Bu müteahhit bir binaya atandığında projeler burada görünür.
                </p>
              </div>
            )}
          </section>
        </div>
      </div>
    </main>
  );
}

function StatBox({
  label,
  value,
  tone,
}: {
  label: string;
  value: number;
  tone: "indigo" | "emerald" | "slate";
}) {
  const toneCls = {
    indigo: "text-indigo-600",
    emerald: "text-emerald-600",
    slate: "text-ink-900",
  }[tone];
  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-card">
      <p className="text-[12.5px] font-medium text-ink-500">{label}</p>
      <p className={`mt-1.5 text-[26px] font-bold leading-none ${toneCls}`}>
        {value}
      </p>
    </div>
  );
}

function InfoRow({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="flex items-center gap-2 text-ink-500">
        <Icon className="h-4 w-4 text-ink-400" />
        {label}
      </span>
      <span className="truncate text-right font-semibold text-ink-900">
        {value}
      </span>
    </div>
  );
}
