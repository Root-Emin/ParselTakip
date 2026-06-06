import { STATUS_BADGE, type BuildingStatus } from "@/lib/buildings";

export default function StatusBadge({
  status,
  className = "",
}: {
  status: BuildingStatus;
  className?: string;
}) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-3 py-1 text-[11.5px] font-semibold text-white shadow-soft ${STATUS_BADGE[status]} ${className}`}
    >
      {status}
    </span>
  );
}
