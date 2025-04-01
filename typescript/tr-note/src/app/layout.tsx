import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { TrainingProvider } from "../contexts/TrainingContext";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "トレーニング記録アプリ",
  description: "トレーニングの記録を管理するアプリケーション",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <TrainingProvider>{children}</TrainingProvider>
      </body>
    </html>
  );
}
