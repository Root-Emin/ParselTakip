import { Search, Bell, ChevronDown } from "lucide-react";

export default function Topbar() {
  return (
    <header className="flex h-[68px] shrink-0 items-center gap-4 border-b border-slate-200 bg-white px-6">
      <div className="min-w-0">
        <h1 className="truncate text-[17px] font-bold tracking-tight text-ink-900">
          Hoş geldiniz, Ahmet Yılmaz <span className="font-normal">👋</span>
        </h1>
      </div>

      <div className="mx-auto hidden w-full max-w-md md:block">
        <div className="relative">
          <Search className="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-ink-400" />
          <input
            type="text"
            placeholder="Bina, proje, müteahhit ara..."
            className="w-full rounded-xl border border-slate-200 bg-slate-50/60 py-2.5 pl-10 pr-4 text-[13px] text-ink-900 placeholder:text-ink-400 outline-none transition-colors focus:border-brand-300 focus:bg-white focus:ring-2 focus:ring-brand-100"
          />
        </div>
      </div>

      <div className="ml-auto flex items-center gap-3">
        <button
          type="button"
          aria-label="Bildirimler"
          className="relative flex h-10 w-10 items-center justify-center rounded-xl border border-slate-200 bg-white text-ink-500 transition-colors hover:bg-slate-50"
        >
          <Bell className="h-[18px] w-[18px]" />
          <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-rose-500 ring-2 ring-white" />
        </button>

        <div className="flex items-center gap-2.5 rounded-xl border border-slate-200 bg-white py-1.5 pl-1.5 pr-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-brand-500 to-brand-700 text-[12px] font-semibold text-white">
            AY
          </div>
          <div className="hidden leading-tight sm:block">
            <p className="text-[12.5px] font-semibold text-ink-900">
              Ahmet Yılmaz
            </p>
            <p className="text-[10.5px] text-ink-400">Proje Yöneticisi</p>
          </div>
          <ChevronDown className="h-4 w-4 text-ink-400" />
        </div>
      </div>
    </header>
  );
}
