import { Check } from "lucide-react";
import { PROCESS_STEPS } from "@/lib/buildings";

export default function ProcessStepper({ current }: { current: number }) {
  return (
    <div className="flex items-start">
      {PROCESS_STEPS.map((label, i) => {
        const step = i + 1;
        const isLast = i === PROCESS_STEPS.length - 1;
        const done = step < current;
        const active = step === current;

        return (
          <div key={label} className="relative flex flex-1 flex-col items-center">
            {!isLast && (
              <span
                className={[
                  "absolute left-1/2 top-5 h-1 w-full rounded-full",
                  done ? "bg-indigo-500" : "bg-slate-200",
                ].join(" ")}
              />
            )}

            <div
              className={[
                "relative z-10 flex h-10 w-10 items-center justify-center rounded-full text-[14px] font-semibold transition-all",
                done
                  ? "bg-indigo-500 text-white"
                  : active
                    ? "border-2 border-indigo-500 bg-white text-indigo-600 ring-4 ring-indigo-100"
                    : "border-2 border-slate-200 bg-white text-ink-400",
              ].join(" ")}
            >
              {done ? <Check className="h-4 w-4" /> : step}
            </div>

            <p
              className={[
                "mt-2 text-center text-[12px] font-medium leading-tight",
                done
                  ? "text-ink-700"
                  : active
                    ? "text-indigo-700"
                    : "text-ink-400",
              ].join(" ")}
            >
              {label}
            </p>
          </div>
        );
      })}
    </div>
  );
}
