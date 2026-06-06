import { type Building } from "@/lib/buildings";

export interface Contractor {
  id: string;
  name: string;
  contactPerson: string;
  phone: string;
  email: string;
  city: string;
  since: number;
}

export interface NewContractorInput {
  name: string;
  contactPerson?: string;
  phone?: string;
  email?: string;
  city?: string;
}

export function createContractor(input: NewContractorInput): Contractor {
  return {
    id: `c-${Date.now()}`,
    name: input.name.trim(),
    contactPerson: input.contactPerson?.trim() || "Belirtilmedi",
    phone: input.phone?.trim() || "-",
    email: input.email?.trim() || "-",
    city: input.city?.trim() || "İstanbul",
    since: new Date().getFullYear(),
  };
}

export function contractorProjects(
  buildings: Building[],
  name: string,
): Building[] {
  return buildings.filter((b) => b.contractor === name);
}

export function activeProjectCount(buildings: Building[], name: string): number {
  return buildings.filter(
    (b) => b.contractor === name && b.status !== "Teslim Edildi",
  ).length;
}

export function completedProjectCount(
  buildings: Building[],
  name: string,
): number {
  return buildings.filter(
    (b) => b.contractor === name && b.status === "Teslim Edildi",
  ).length;
}

export const seedContractors: Contractor[] = [
  {
    id: "abc-yapi",
    name: "ABC Yapı A.Ş.",
    contactPerson: "Mehmet Demir",
    phone: "0212 555 0142",
    email: "info@abcyapi.com",
    city: "İstanbul",
    since: 2009,
  },
  {
    id: "demir-insaat",
    name: "Demir İnşaat",
    contactPerson: "Ali Demir",
    phone: "0216 555 0177",
    email: "iletisim@demirinsaat.com",
    city: "İstanbul",
    since: 2014,
  },
  {
    id: "kaya-insaat",
    name: "Kaya İnşaat",
    contactPerson: "Hasan Kaya",
    phone: "0212 555 0190",
    email: "info@kayainsaat.com",
    city: "İstanbul",
    since: 2011,
  },
  {
    id: "mega-yapi",
    name: "Mega Yapı A.Ş.",
    contactPerson: "Ayşe Yıldız",
    phone: "0212 555 0123",
    email: "info@megayapi.com",
    city: "İstanbul",
    since: 2005,
  },
  {
    id: "yuksel-yapi",
    name: "Yüksel Yapı",
    contactPerson: "Cem Yüksel",
    phone: "0216 555 0211",
    email: "info@yukselyapi.com",
    city: "İstanbul",
    since: 2016,
  },
  {
    id: "doga-yapi",
    name: "Doğa Yapı A.Ş.",
    contactPerson: "Selin Doğa",
    phone: "0212 555 0233",
    email: "info@dogayapi.com",
    city: "İstanbul",
    since: 2018,
  },
];
