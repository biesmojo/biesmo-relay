import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
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
  title: "Relay - AI-native CRM & Workflows",
  description: "AI agents, event-driven rules, and CRM in one platform. Handle customer conversations, automate workflows, and manage tickets seamlessly.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}>
      <body className="min-h-full bg-gradient-to-br from-zinc-50 to-background dark:from-black dark:to-zinc-950 flex">
        {/* Sidebar */}
        <div className="fixed left-0 top-0 h-screen w-64 bg-black/80 backdrop-blur-md border-r border-foreground/20 z-50 p-6 hidden lg:block">
          <div className="flex flex-col h-full">
            <h1 className="text-2xl font-black text-white mb-8">Relay</h1>
            <nav className="flex-1 space-y-2">
              <a href="/chat" className="group flex items-center gap-3 p-3 rounded-xl text-white/70 hover:bg-white/10 hover:text-white transition-all">
                <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5 group-hover:scale-110 transition-transform">
                  <path stroke="currentColor" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.97-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.97 4.03-9 9-9s9 4.03 9 9z"/>
                </svg>
                Chat
              </a>
              <a href="/admin/sessions" className="group flex items-center gap-3 p-3 rounded-xl text-white/70 hover:bg-white/10 hover:text-white transition-all">
                <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5 group-hover:scale-110 transition-transform">
                  <path stroke="currentColor" strokeWidth="2" d="M3 19v-8.93a2 2 0 01.89-1.664l7-3.666a2 2 0 012.22 0l7 3.666A2 2 0 0121 10.07V19M3 19h18M7 19v-6M12 19v-8M17 19v-4"/>
                </svg>
                Sesi
              </a>
              <a href="/admin/tickets" className="group flex items-center gap-3 p-3 rounded-xl text-white/70 hover:bg-white/10 hover:text-white transition-all">
                <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5 group-hover:scale-110 transition-transform">
                  <path stroke="currentColor" strokeWidth="2" d="M19 14c1.49 0 2.49 1.01 2.49 2.4L21.5 19c0 1.1-.84 2-1.9 2H4.4c-1.06 0-1.9-.9-1.9-2v-2.6c0-1.39 1-2.4 2.5-2.4h14.5zM12 10l-6 6h12l-6-6z"/>
                </svg>
                Tiket
              </a>
              <a href="/admin/events" className="group flex items-center gap-3 p-3 rounded-xl text-white/70 hover:bg-white/10 hover:text-white transition-all">
                <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5 group-hover:scale-110 transition-transform">
                  <path stroke="currentColor" strokeWidth="2" d="M12 8v4l3 3M16 3.13a4 4 0 010 7.75M5.05 5.05a7 7 0 019.9 9.9L12 16l3 3"/>
                </svg>
                Events
              </a>
            </nav>
          </div>
        </div>
        
        {/* Main */}
        <main className="flex-1 lg:ml-64 p-8">
          {children}
        </main>
      </body>
    </html>
  );
}

