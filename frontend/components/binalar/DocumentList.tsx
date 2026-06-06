"use client";

import { FileCheck2, Clock, FileWarning } from "lucide-react";
import { type BuildingDoc } from "@/lib/buildings";

export default function DocumentList({
  documents,
  onUpload,
}: {
  documents: BuildingDoc[];
  onUpload: (docName: string) => void;
}) {
  return (
    <ul>
      {documents.map((doc, i) => {
        const isLast = i === documents.length - 1;
        return (
          <li
            key={doc.name}
            className={[
              "flex items-center justify-between gap-3 py-3.5",
              isLast ? "" : "border-b border-slate-100",
            ].join(" ")}
          >
            <div className="flex min-w-0 items-center gap-3">
              <DocIcon status={doc.status} />
              <div className="min-w-0">
                <p className="truncate text-[13.5px] font-semibold text-ink-900">
                  {doc.name}
                </p>
                {doc.status === "Yüklenmedi" ? (
                  <p className="text-[12px] text-rose-500">
                    Bu belge henüz yüklenmedi.
                  </p>
                ) : (
                  <p className="truncate text-[12px] text-ink-400">
                    {doc.uploadedBy} tarafından yüklendi • {doc.date}
                  </p>
                )}
              </div>
            </div>

            <div className="shrink-0">
              {doc.status === "Onaylandı" && (
                <span className="rounded-full bg-emerald-50 px-3 py-1 text-[11.5px] font-medium text-emerald-600">
                  Onaylandı
                </span>
              )}
              {doc.status === "Onay Bekliyor" && (
                <span className="rounded-full bg-amber-50 px-3 py-1 text-[11.5px] font-medium text-amber-600">
                  Onay Bekliyor
                </span>
              )}
              {doc.status === "Yüklenmedi" && (
                <button
                  type="button"
                  onClick={() => onUpload(doc.name)}
                  className="rounded-lg px-3 py-1 text-[12.5px] font-semibold text-indigo-600 transition-colors hover:bg-indigo-50"
                >
                  Yükle
                </button>
              )}
            </div>
          </li>
        );
      })}
    </ul>
  );
}

function DocIcon({ status }: { status: BuildingDoc["status"] }) {
  if (status === "Onaylandı") {
    return (
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-500">
        <FileCheck2 className="h-[18px] w-[18px]" />
      </span>
    );
  }
  if (status === "Onay Bekliyor") {
    return (
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-amber-50 text-amber-500">
        <Clock className="h-[18px] w-[18px]" />
      </span>
    );
  }
  return (
    <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-slate-400">
      <FileWarning className="h-[18px] w-[18px]" />
    </span>
  );
}
