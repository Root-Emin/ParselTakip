import { ArrowUpRight } from "lucide-react";
import { statCards } from "@/lib/data";

export default function StatCards() {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
      {statCards.map((card) => {
        const Icon = card.icon;
        return (
          <div
            key={card.label}
            className="rounded-2xl border border-slate-200 bg-white p-4 shadow-card"
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="text-[12.5px] font-medium text-ink-500">
                  {card.label}
                </p>
                <p className="mt-2 text-[28px] font-bold leading-none tracking-tight text-ink-900">
                  {card.value}
                </p>
              </div>
              <div
                className={`flex h-11 w-11 items-center justify-center rounded-xl ${card.tint}`}
              >
                <Icon className={`h-5 w-5 ${card.iconColor}`} />
              </div>
            </div>
            <div className="mt-3 flex items-center gap-1 text-[11.5px] font-medium text-emerald-600">
              <ArrowUpRight className="h-3.5 w-3.5" />
              <span>{card.delta}</span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
