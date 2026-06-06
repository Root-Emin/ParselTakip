"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import {
  type Contractor,
  type NewContractorInput,
  createContractor,
  seedContractors,
} from "@/lib/contractors";

const STORAGE_KEY = "parseltakip.contractors.v1";

interface ContractorsContextValue {
  contractors: Contractor[];
  getContractor: (id: string) => Contractor | undefined;
  addContractor: (input: NewContractorInput) => Contractor;
  deleteContractor: (id: string) => void;
}

const ContractorsContext = createContext<ContractorsContextValue | null>(null);

export function ContractorsProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [contractors, setContractors] = useState<Contractor[]>(seedContractors);

  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) setContractors(JSON.parse(raw) as Contractor[]);
    } catch {
      // ignore corrupted storage
    }
  }, []);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(contractors));
    } catch {
      // storage may be unavailable
    }
  }, [contractors]);

  const getContractor = useCallback(
    (id: string) => contractors.find((c) => c.id === id),
    [contractors],
  );

  const addContractor = useCallback((input: NewContractorInput) => {
    const next = createContractor(input);
    setContractors((prev) => [next, ...prev]);
    return next;
  }, []);

  const deleteContractor = useCallback((id: string) => {
    setContractors((prev) => prev.filter((c) => c.id !== id));
  }, []);

  return (
    <ContractorsContext.Provider
      value={{ contractors, getContractor, addContractor, deleteContractor }}
    >
      {children}
    </ContractorsContext.Provider>
  );
}

export function useContractors() {
  const ctx = useContext(ContractorsContext);
  if (!ctx) {
    throw new Error("useContractors must be used within ContractorsProvider");
  }
  return ctx;
}
