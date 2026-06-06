import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Sidebar from "@/components/Sidebar";
import Topbar from "@/components/Topbar";
import { BuildingsProvider } from "@/components/BuildingsProvider";
import { ContractorsProvider } from "@/components/ContractorsProvider";

const inter = Inter({
  subsets: ["latin", "latin-ext"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "ParselTakip — Kentsel Dönüşüm Yönetim Sistemi",
  description:
    "Belediyeler için kentsel dönüşüm süreç yönetim paneli — ParselTakip.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="tr" className={inter.variable}>
      <body className="font-sans">
        <BuildingsProvider>
          <ContractorsProvider>
            <div className="flex h-screen w-full overflow-hidden bg-canvas">
              <Sidebar />
              <div className="flex min-w-0 flex-1 flex-col">
                <Topbar />
                <div className="flex min-h-0 flex-1">{children}</div>
              </div>
            </div>
          </ContractorsProvider>
        </BuildingsProvider>
      </body>
    </html>
  );
}
