"use client";

import { Building2, BadgeCheck, LifeBuoy, ChevronRight } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { navItems } from "@/lib/data";

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex h-full w-[248px] shrink-0 flex-col border-r border-slate-200 bg-white">
      {/* Logo */}
      <div className="flex items-center gap-2.5 px-5 py-5">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-brand-500 to-brand-700 text-white shadow-soft">
          <Building2 className="h-5 w-5" />
        </div>
        <div className="leading-tight">
          <p className="text-[15px] font-bold tracking-tight text-ink-900">
            ParselTakip
          </p>
          <p className="text-[10px] font-medium text-ink-400">
            Kentsel Dönüşüm Yönetim Sistemi
          </p>
        </div>
      </div>

      {/* Nav */}
      <nav className="scroll-thin mt-1 flex-1 overflow-y-auto px-3">
        <ul className="space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active =
              item.href === "/"
                ? pathname === "/"
                : item.href !== "#" && pathname.startsWith(item.href);
            return (
              <li key={item.label}>
                <Link
                  href={item.href}
                  className={[
                    "group flex items-center gap-3 rounded-xl px-3 py-2.5 text-[13.5px] font-medium transition-colors",
                    active
                      ? "bg-brand-50 text-brand-700"
                      : "text-ink-500 hover:bg-slate-50 hover:text-ink-900",
                  ].join(" ")}
                >
                  <Icon
                    className={[
                      "h-[18px] w-[18px]",
                      active
                        ? "text-brand-600"
                        : "text-ink-400 group-hover:text-ink-700",
                    ].join(" ")}
                  />
                  <span className="flex-1">{item.label}</span>
                  {item.badge ? (
                    <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-brand-500 px-1.5 text-[11px] font-semibold text-white">
                      {item.badge}
                    </span>
                  ) : null}
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* Municipality card */}
      <div className="px-3 pb-3">
        <div className="rounded-2xl border border-slate-200 bg-slate-50/70 p-3">
          <div className="flex items-start gap-2.5">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white text-brand-600 shadow-card">
              <Building2 className="h-[18px] w-[18px]" />
            </div>
            <div className="leading-tight">
              <div className="flex items-center gap-1">
                <p className="text-[12.5px] font-semibold text-ink-900">
                  İstanbul Büyükşehir
                </p>
                <BadgeCheck className="h-3.5 w-3.5 text-brand-500" />
              </div>
              <p className="text-[12.5px] font-semibold text-ink-900">
                Belediyesi
              </p>
              <p className="mt-1 text-[10.5px] text-ink-400">
                Kullanıcı Türü
              </p>
              <p className="text-[11px] font-medium text-ink-500">Belediye</p>
            </div>
          </div>
        </div>
      </div>

      {/* Help */}
      <div className="px-3 pb-4">
        <a
          href="#"
          className="flex items-center gap-2.5 rounded-xl border border-slate-200 bg-white px-3 py-2.5 text-[13px] font-medium text-ink-700 transition-colors hover:bg-slate-50"
        >
          <LifeBuoy className="h-[18px] w-[18px] text-ink-400" />
          <span className="flex-1">Yardım &amp; Destek</span>
          <ChevronRight className="h-4 w-4 text-ink-400" />
        </a>
      </div>
    </aside>
  );
}
