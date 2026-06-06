export type BuildingStatus =
  | "Başvuru Yapıldı"
  | "Belge Toplama"
  | "Yıkım Onayı Bekleniyor"
  | "İnşaat Aşaması"
  | "Teslim Edildi";

export type DocState = "Onaylandı" | "Onay Bekliyor" | "Yüklenmedi";

export interface BuildingDoc {
  name: string;
  status: DocState;
  uploadedBy?: string;
  date?: string;
}

export interface Building {
  id: string;
  name: string;
  district: string;
  address: string;
  status: BuildingStatus;
  contractor: string | null;
  progress: number;
  ownersCount: number;
  startDate: string;
  documents: BuildingDoc[];
}

export interface NewBuildingInput {
  name: string;
  district: string;
  address: string;
  ownersCount: number;
  contractor?: string;
}

export const PROCESS_STEPS = [
  "Başvuru",
  "Belge Toplama",
  "Yıkım Onayı",
  "İnşaat Aşaması",
  "Teslim",
] as const;

export const STATUS_STEP: Record<BuildingStatus, number> = {
  "Başvuru Yapıldı": 1,
  "Belge Toplama": 2,
  "Yıkım Onayı Bekleniyor": 3,
  "İnşaat Aşaması": 4,
  "Teslim Edildi": 5,
};

export const STATUS_BADGE: Record<BuildingStatus, string> = {
  "Başvuru Yapıldı": "bg-blue-500",
  "Belge Toplama": "bg-slate-700",
  "Yıkım Onayı Bekleniyor": "bg-amber-500",
  "İnşaat Aşaması": "bg-emerald-500",
  "Teslim Edildi": "bg-indigo-500",
};

export const ALL_STATUSES: BuildingStatus[] = [
  "Başvuru Yapıldı",
  "Belge Toplama",
  "Yıkım Onayı Bekleniyor",
  "İnşaat Aşaması",
  "Teslim Edildi",
];

export const DEFAULT_DOC_NAMES = [
  "Tapu Fotokopisi",
  "Kimlik Fotokopisi",
  "Deprem Risk Raporu",
  "Muvafakatname",
  "Vekaletname",
  "Yapı Ruhsatı",
  "Zemin Etüdü",
  "Aplikasyon Krokisi",
  "İmar Durumu",
  "Numarataj Belgesi",
];

export function docsDone(b: Building): number {
  return b.documents.filter((d) => d.status !== "Yüklenmedi").length;
}

const approved = (
  name: string,
  by = "Ahmet Yılmaz",
  date = "2024-05-15",
): BuildingDoc => ({ name, status: "Onaylandı", uploadedBy: by, date });

const pending = (
  name: string,
  by = "Av. Mehmet Can",
  date = "2024-06-01",
): BuildingDoc => ({ name, status: "Onay Bekliyor", uploadedBy: by, date });

const missing = (name: string): BuildingDoc => ({
  name,
  status: "Yüklenmedi",
});

export function createBuilding(input: NewBuildingInput): Building {
  const today = new Date().toISOString().slice(0, 10);
  return {
    id: `b-${Date.now()}`,
    name: input.name.trim(),
    district: input.district.trim(),
    address: input.address.trim(),
    status: "Başvuru Yapıldı",
    contractor: input.contractor?.trim() ? input.contractor.trim() : null,
    progress: 5,
    ownersCount: input.ownersCount,
    startDate: today,
    documents: DEFAULT_DOC_NAMES.map((n) => missing(n)),
  };
}

export const seedBuildings: Building[] = [
  {
    id: "gunes-apartmani",
    name: "Güneş Apartmanı",
    district: "Kadıköy",
    address: "Caferağa Mah. No: 12",
    status: "Belge Toplama",
    contractor: null,
    progress: 25,
    ownersCount: 12,
    startDate: "2024-05-12",
    documents: [
      approved("Tapu Fotokopisi"),
      approved("Kimlik Fotokopisi"),
      missing("Deprem Risk Raporu"),
      missing("Muvafakatname"),
      pending("Vekaletname"),
      approved("Yapı Ruhsatı", "Ahmet Yılmaz", "2024-05-18"),
      approved("Zemin Etüdü", "Ahmet Yılmaz", "2024-05-18"),
      approved("Aplikasyon Krokisi", "Ahmet Yılmaz", "2024-05-20"),
      pending("İmar Durumu", "Ahmet Yılmaz", "2024-05-28"),
      missing("Numarataj Belgesi"),
    ],
  },
  {
    id: "yildiz-sitesi-a-blok",
    name: "Yıldız Sitesi A Blok",
    district: "Üsküdar",
    address: "Acıbadem Mah. No: 45",
    status: "Yıkım Onayı Bekleniyor",
    contractor: "ABC Yapı A.Ş.",
    progress: 60,
    ownersCount: 24,
    startDate: "2024-03-10",
    documents: DEFAULT_DOC_NAMES.map((n) => approved(n, "Ahmet Yılmaz", "2024-04-02")),
  },
  {
    id: "huzur-apartmani",
    name: "Huzur Apartmanı",
    district: "Şişli",
    address: "Fulya Mah. No: 8",
    status: "İnşaat Aşaması",
    contractor: "Demir İnşaat",
    progress: 85,
    ownersCount: 18,
    startDate: "2024-01-20",
    documents: [
      ...DEFAULT_DOC_NAMES,
      "İnşaat Ruhsatı",
      "Statik Proje",
      "Mekanik Proje",
      "Elektrik Projesi",
      "İskan Belgesi",
    ].map((n) => approved(n, "Ahmet Yılmaz", "2024-02-15")),
  },
  {
    id: "kardesler-apartmani",
    name: "Kardeşler Apartmanı",
    district: "Beşiktaş",
    address: "Türkali Mah. No: 22",
    status: "Başvuru Yapıldı",
    contractor: null,
    progress: 10,
    ownersCount: 8,
    startDate: "2024-06-01",
    documents: DEFAULT_DOC_NAMES.map((n, i) =>
      i < 2 ? approved(n, "Ahmet Yılmaz", "2024-06-02") : missing(n),
    ),
  },
  {
    id: "bahar-sitesi",
    name: "Bahar Sitesi",
    district: "Maltepe",
    address: "Bağlarbaşı Mah. No: 5",
    status: "Yıkım Onayı Bekleniyor",
    contractor: "Kaya İnşaat",
    progress: 55,
    ownersCount: 30,
    startDate: "2024-02-28",
    documents: [
      ...DEFAULT_DOC_NAMES,
      "İnşaat Ruhsatı",
      "Statik Proje",
    ].map((n) => approved(n, "Ahmet Yılmaz", "2024-03-15")),
  },
  {
    id: "cinar-konutlari",
    name: "Çınar Konutları",
    district: "Ataşehir",
    address: "Barbaros Mah. No: 17",
    status: "Teslim Edildi",
    contractor: "Mega Yapı A.Ş.",
    progress: 100,
    ownersCount: 40,
    startDate: "2023-09-15",
    documents: DEFAULT_DOC_NAMES.map((n) => approved(n, "Ahmet Yılmaz", "2023-10-01")),
  },
];
