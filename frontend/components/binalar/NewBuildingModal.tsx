"use client";

import { useState } from "react";
import { X, Building2 } from "lucide-react";
import { useBuildings } from "@/components/BuildingsProvider";

export default function NewBuildingModal({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const { addBuilding } = useBuildings();
  const [name, setName] = useState("");
  const [district, setDistrict] = useState("");
  const [address, setAddress] = useState("");
  const [owners, setOwners] = useState("");
  const [contractor, setContractor] = useState("");

  if (!open) return null;

  const valid = name.trim() && district.trim() && address.trim();

  function reset() {
    setName("");
    setDistrict("");
    setAddress("");
    setOwners("");
    setContractor("");
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!valid) return;
    addBuilding({
      name,
      district,
      address,
      ownersCount: Number(owners) || 0,
      contractor: contractor || undefined,
    });
    reset();
    onClose();
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="w-full max-w-md overflow-hidden rounded-2xl bg-white shadow-soft"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-slate-100 px-5 py-4">
          <div className="flex items-center gap-2.5">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-indigo-50 text-indigo-600">
              <Building2 className="h-5 w-5" />
            </div>
            <h2 className="text-[15px] font-bold text-ink-900">
              Yeni Süreç Başlat
            </h2>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Kapat"
            className="flex h-8 w-8 items-center justify-center rounded-lg text-ink-400 transition-colors hover:bg-slate-100 hover:text-ink-700"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4 px-5 py-5">
          <Field label="Bina Adı" required>
            <input
              autoFocus
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Örn. Yıldız Apartmanı"
              className={inputCls}
            />
          </Field>

          <div className="grid grid-cols-2 gap-3">
            <Field label="İlçe" required>
              <input
                value={district}
                onChange={(e) => setDistrict(e.target.value)}
                placeholder="Kadıköy"
                className={inputCls}
              />
            </Field>
            <Field label="Hak Sahibi Sayısı">
              <input
                type="number"
                min={0}
                value={owners}
                onChange={(e) => setOwners(e.target.value)}
                placeholder="0"
                className={inputCls}
              />
            </Field>
          </div>

          <Field label="Adres" required>
            <input
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Mahalle, Cadde No"
              className={inputCls}
            />
          </Field>

          <Field label="Müteahhit (opsiyonel)">
            <input
              value={contractor}
              onChange={(e) => setContractor(e.target.value)}
              placeholder="Henüz atanmadıysa boş bırakın"
              className={inputCls}
            />
          </Field>

          <div className="flex items-center justify-end gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-xl border border-slate-200 px-4 py-2.5 text-[13px] font-medium text-ink-700 transition-colors hover:bg-slate-50"
            >
              İptal
            </button>
            <button
              type="submit"
              disabled={!valid}
              className="rounded-xl bg-indigo-600 px-4 py-2.5 text-[13px] font-semibold text-white shadow-soft transition-colors hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              Süreci Başlat
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

const inputCls =
  "w-full rounded-xl border border-slate-200 bg-slate-50/60 px-3 py-2.5 text-[13px] text-ink-900 placeholder:text-ink-400 outline-none transition-colors focus:border-indigo-300 focus:bg-white focus:ring-2 focus:ring-indigo-100";

function Field({
  label,
  required,
  children,
}: {
  label: string;
  required?: boolean;
  children: React.ReactNode;
}) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-[12px] font-medium text-ink-700">
        {label}
        {required && <span className="text-rose-500"> *</span>}
      </span>
      {children}
    </label>
  );
}
