import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Huddle - Real-time Chat Application",
  description:
    "Connect, Chat, & Collaborate with Huddle - The next generation of real-time communication platform.",
  keywords: [
    "chat",
    "messaging",
    "real-time",
    "communication",
    "team collaboration",
  ],
  authors: [{ name: "Huddle Team" }],
  creator: "Huddle",
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "https://huddle.com",
    title: "Huddle - Real-time Chat Application",
    description:
      "Connect, Chat, & Collaborate with Huddle - The next generation of real-time communication platform.",
    siteName: "Huddle",
  },
  twitter: {
    card: "summary_large_image",
    title: "Huddle - Real-time Chat Application",
    description:
      "Connect, Chat, & Collaborate with Huddle - The next generation of real-time communication platform.",
  },
  robots: {
    index: true,
    follow: true,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="scroll-smooth">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
        <Toaster />
      </body>
    </html>
  );
}
