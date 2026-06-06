"use client";

import { useState } from "react";
import Link from "next/link";
import {
  ArrowLeft,
  Pencil,
  MoreVertical,
  Upload,
  Users,
  Calendar,
  HardHat,
  Building2,
  MapPin,
} from "lucide-react";
import { useBuildings } from "@/components/BuildingsProvider";
import { STATUS_STEP, docsDone } from "@/lib/buildings";
import ProcessStepper from "@/components/binalar/ProcessStepper";
import DocumentList from "@/components/binalar/DocumentList";
import StatusBadge from "@/components/binalar/StatusBadge";

export default function BuildingDetailPage({
  params,
}: {
  params: { id: string };
}) {
  const { getBuilding, assignContractor, markDocumentUploaded } = useBuildings();
  const building = getBuilding(params.id);

  const [assigning, setAssigning] = useState(false);
  const [contractorName, setContractorName] = useState("");

  if (!building) {
    return (
      <main className="flex min-w-0 flex-1 items-center justify-center p-6">
        <div className="text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-400">
            <Building2 className="h-6 w-6" />
          </div>
          <p className="mt-3 text-[15px] font-semibold text-ink-900">
            Bina bulunamadı
          </p>
          <Link
            href="/binalar"
            className="mt-3 inline-flex items-center gap-1.5 text-[13px] font-medium text-indigo-600 hover:text-indigo-700"
          >
            <ArrowLeft className="h-4 w-4" />
            Binalara dön
          </Link>
        </div>
      </main>
    );
  }

  const done = docsDone(building);
  const total = building.documents.length;
  const firstMissing = building.documents.find((d) => d.status === "Yüklenmedi");

  function handleAssign(e: React.FormEvent) {
    e.preventDefault();
    if (!building || !contractorName.trim()) return;
    assignContractor(building.id, contractorName.trim());
    setContractorName("");
    setAssigning(false);
  }

  return (
    <main className="scroll-thin min-w-0 flex-1 overflow-y-auto p-6">
      <div className="mx-auto max-w-6xl">
        {/* Header */}
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="flex items-start gap-3">
            <Link
              href="/binalar"
              aria-label="Geri"
              className="mt-0.5 flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-ink-500 transition-colors hover:bg-slate-50"
            >
              <ArrowLeft className="h-4 w-4" />
            </Link>
            <div>
              <div className="flex items-center gap-2.5">
                <h1 className="text-[22px] font-bold tracking-tight text-ink-900">
                  {building.name}
                </h1>
                <StatusBadge status={building.status} />
              </div>
              <p className="mt-1 flex items-center gap-1.5 text-[13px] text-ink-500">
                <MapPin className="h-3.5 w-3.5 text-ink-400" />
                {building.district}, {building.address}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              className="flex items-center gap-1.5 rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-[13px] font-medium text-ink-700 transition-colors hover:bg-slate-50"
            >
              <Pencil className="h-3.5 w-3.5" />
              Düzenle
            </button>
            <button
              type="button"
              aria-label="Daha fazla"
              className="flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-ink-500 transition-colors hover:bg-slate-50"
            >
              <MoreVertical className="h-4 w-4" />
            </button>
          </div>
        </div>

        {/* Body grid */}
        <div className="mt-6 grid grid-cols-1 gap-5 lg:grid-cols-3">
          {/* Left column */}
          <div className="space-y-5 lg:col-span-2">
            {/* Process */}
            <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
              <h2 className="text-[15px] font-semibold text-ink-900">
                Süreç İlerlemesi
              </h2>
              <div className="mt-6">
                <ProcessStepper current={STATUS_STEP[building.status]} />
              </div>
            </section>

            {/* Documents */}
            <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-[15px] font-semibold text-ink-900">
                    Belgeler
                  </h2>
                  <p className="mt-0.5 text-[12.5px] text-ink-400">
                    Süreç için gerekli olan belgelerin durumu. ({done}/{total})
                  </p>
                </div>
                <button
                  type="button"
                  disabled={!firstMissing}
                  onClick={() =>
                    firstMissing &&
                    markDocumentUploaded(building.id, firstMissing.name)
                  }
                  className="flex items-center gap-1.5 rounded-xl bg-indigo-50 px-3.5 py-2 text-[13px] font-semibold text-indigo-600 transition-colors hover:bg-indigo-100 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <Upload className="h-4 w-4" />
                  Belge Yükle
                </button>
              </div>

              <div className="mt-2">
                <DocumentList
                  documents={building.documents}
                  onUpload={(docName) =>
                    markDocumentUploaded(building.id, docName)
                  }
                />
              </div>
            </section>
          </div>

          {/* Right column */}
          <div className="space-y-5">
            {/* Image */}
            <div className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-card">
              <div className="flex h-44 items-center justify-center bg-gradient-to-br from-indigo-100 via-slate-100 to-slate-200">
                <Building2 className="h-14 w-14 text-slate-300" />
              </div>
            </div>

            {/* Bina Bilgileri */}
            <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
              <h2 className="text-[15px] font-semibold text-ink-900">
                Bina Bilgileri
              </h2>
              <div className="mt-4 space-y-3 text-[13px]">
                <div className="flex items-center justify-between">
                  <span className="flex items-center gap-2 text-ink-500">
                    <Users className="h-4 w-4 text-ink-400" />
                    Hak Sahipleri
                  </span>
                  <span className="font-semibold text-ink-900">
                    {building.ownersCount} Kişi
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="flex items-center gap-2 text-ink-500">
                    <Calendar className="h-4 w-4 text-ink-400" />
                    Başlangıç
                  </span>
                  <span className="font-semibold text-ink-900">
                    {building.startDate}
                  </span>
                </div>
              </div>
            </section>

            {/* Müteahhit Bilgisi */}
            <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
              <h2 className="flex items-center gap-2 text-[15px] font-semibold text-ink-900">
                <HardHat className="h-[18px] w-[18px] text-indigo-500" />
                Müteahhit Bilgisi
              </h2>
              <div className="mt-3 border-t border-slate-100 pt-4">
                {building.contractor ? (
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-50 text-indigo-600">
                      <HardHat className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="text-[14px] font-semibold text-ink-900">
                        {building.contractor}
                      </p>
                      <p className="text-[12px] text-ink-400">
                        Atanmış müteahhit
                      </p>
                    </div>
                  </div>
                ) : assigning ? (
                  <form onSubmit={handleAssign} className="space-y-2.5">
                    <input
                      autoFocus
                      value={contractorName}
                      onChange={(e) => setContractorName(e.target.value)}
                      placeholder="Müteahhit firma adı"
                      className="w-full rounded-xl border border-slate-200 bg-slate-50/60 px-3 py-2.5 text-[13px] text-ink-900 placeholder:text-ink-400 outline-none focus:border-indigo-300 focus:bg-white focus:ring-2 focus:ring-indigo-100"
                    />
                    <div className="flex gap-2">
                      <button
                        type="button"
                        onClick={() => setAssigning(false)}
                        className="flex-1 rounded-xl border border-slate-200 py-2.5 text-[13px] font-medium text-ink-700 transition-colors hover:bg-slate-50"
                      >
                        İptal
                      </button>
                      <button
                        type="submit"
                        disabled={!contractorName.trim()}
                        className="flex-1 rounded-xl bg-indigo-600 py-2.5 text-[13px] font-semibold text-white transition-colors hover:bg-indigo-700 disabled:opacity-50"
                      >
                        Ata
                      </button>
                    </div>
                  </form>
                ) : (
                  <>
                    <p className="text-center text-[13px] text-ink-400">
                      Henüz bir müteahhit atanmamış.
                    </p>
                    <button
                      type="button"
                      onClick={() => setAssigning(true)}
                      className="mt-3 w-full rounded-xl bg-indigo-600 py-2.5 text-[13.5px] font-semibold text-white shadow-soft transition-colors hover:bg-indigo-700"
                    >
                      Müteahhit Ata
                    </button>
                  </>
                )}
              </div>
            </section>
          </div>
        </div>
      </div>
    </main>
  );
}
