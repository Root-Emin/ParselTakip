import {
  LayoutDashboard,
  Building2,
  HardHat,
  FileText,
  CheckSquare,
  BarChart3,
  Bell,
  Settings,
  Building,
  RefreshCw,
  FileWarning,
  Hourglass,
  type LucideIcon,
} from "lucide-react";

export type StageStatus = "done" | "active" | "pending";
export type DocStatus = "Tamam" | "Eksik Evrak";

export interface NavItem {
  label: string;
  icon: LucideIcon;
  href: string;
  badge?: number;
}

export const navItems: NavItem[] = [
  { label: "Dashboard", icon: LayoutDashboard, href: "/" },
  { label: "Binalar", icon: Building2, href: "/binalar" },
  { label: "Müteahhitler", icon: HardHat, href: "/muteahhitler" },
  { label: "Belgeler", icon: FileText, href: "#" },
  { label: "Onay Süreçleri", icon: CheckSquare, href: "#" },
  { label: "Raporlar", icon: BarChart3, href: "#" },
  { label: "Bildirimler", icon: Bell, href: "#", badge: 8 },
  { label: "Ayarlar", icon: Settings, href: "#" },
];

export interface StatCard {
  label: string;
  value: string;
  delta: string;
  icon: LucideIcon;
  tint: string;
  iconColor: string;
}

export const statCards: StatCard[] = [
  {
    label: "Toplam Proje",
    value: "128",
    delta: "12% geçen aya göre",
    icon: Building,
    tint: "bg-brand-50",
    iconColor: "text-brand-600",
  },
  {
    label: "Aktif Dönüşüm",
    value: "72",
    delta: "8% geçen aya göre",
    icon: RefreshCw,
    tint: "bg-emerald-50",
    iconColor: "text-emerald-600",
  },
  {
    label: "Eksik Evrak",
    value: "24",
    delta: "5% geçen aya göre",
    icon: FileWarning,
    tint: "bg-orange-50",
    iconColor: "text-orange-500",
  },
  {
    label: "Bekleyen Onay",
    value: "15",
    delta: "3% geçen aya göre",
    icon: Hourglass,
    tint: "bg-violet-50",
    iconColor: "text-violet-500",
  },
];

export type StageLabel =
  | "İnceleme"
  | "Belge Toplama"
  | "Onay"
  | "Yıkım"
  | "Yeniden Yapım";

export interface ProjectRow {
  name: string;
  district: string;
  contractor: string;
  stage: StageLabel;
  progress: number;
  progressColor: string;
  docStatus: DocStatus;
  updated: string;
  selected?: boolean;
}

export const projects: ProjectRow[] = [
  {
    name: "Yıldız Park Evleri",
    district: "Beşiktaş / İstanbul",
    contractor: "Mega Yapı A.Ş.",
    stage: "İnceleme",
    progress: 45,
    progressColor: "bg-brand-500",
    docStatus: "Eksik Evrak",
    updated: "21.05.2024",
    selected: true,
  },
  {
    name: "Güneşli Konutları",
    district: "Bağcılar / İstanbul",
    contractor: "Demir İnşaat",
    stage: "Belge Toplama",
    progress: 30,
    progressColor: "bg-emerald-500",
    docStatus: "Tamam",
    updated: "20.05.2024",
  },
  {
    name: "Koru Sitesi",
    district: "Üsküdar / İstanbul",
    contractor: "Yüksel Yapı",
    stage: "Onay",
    progress: 60,
    progressColor: "bg-violet-500",
    docStatus: "Eksik Evrak",
    updated: "19.05.2024",
  },
  {
    name: "Mavişehir Konakları",
    district: "Beylikdüzü / İstanbul",
    contractor: "Yıldırım İnşaat",
    stage: "Yıkım",
    progress: 75,
    progressColor: "bg-orange-500",
    docStatus: "Tamam",
    updated: "18.05.2024",
  },
  {
    name: "Doğa Rezidans",
    district: "Kadıköy / İstanbul",
    contractor: "Doğa Yapı A.Ş.",
    stage: "Yeniden Yapım",
    progress: 40,
    progressColor: "bg-rose-500",
    docStatus: "Eksik Evrak",
    updated: "17.05.2024",
  },
];

export interface TimelineStage {
  label: string;
  status: StageStatus;
  date?: string;
  hint: string;
}

export const timelineStages: TimelineStage[] = [
  { label: "Başvuru", status: "done", date: "01.03.2024", hint: "Tamamlandı" },
  {
    label: "Belge Toplama",
    status: "done",
    date: "05.04.2024",
    hint: "Tamamlandı",
  },
  { label: "İnceleme", status: "active", date: "Devam Ediyor", hint: "Aktif" },
  { label: "Onay", status: "pending", hint: "Beklemede" },
  { label: "Yıkım", status: "pending", hint: "Beklemede" },
  { label: "Yeniden Yapım", status: "pending", hint: "Beklemede" },
  { label: "Teslim", status: "pending", hint: "Beklemede" },
];

export interface DonutSegment {
  label: string;
  value: number;
  percent: string;
  color: string;
}

export const donutSegments: DonutSegment[] = [
  { label: "Başvuru", value: 18, percent: "%14", color: "#2f6bff" },
  { label: "Belge Toplama", value: 24, percent: "%19", color: "#22c55e" },
  { label: "İnceleme", value: 30, percent: "%23", color: "#f97316" },
  { label: "Onay", value: 20, percent: "%16", color: "#8b5cf6" },
  { label: "Yıkım", value: 16, percent: "%12", color: "#f43f5e" },
  { label: "Yeniden Yapım", value: 12, percent: "%9", color: "#06b6d4" },
  { label: "Teslim", value: 8, percent: "%6", color: "#64748b" },
];

export const donutTotal = donutSegments.reduce((s, d) => s + d.value, 0);

export const missingDocs: string[] = [
  "Zemin Etüd Raporu",
  "Ruhsat Projesi (Mimari)",
  "Ruhsat Projesi (Statik)",
  "Hak Sahipliği Belgesi",
  "Emlak Beyan Formu",
  "Asansör Projesi",
  "Otopark Projesi",
];
