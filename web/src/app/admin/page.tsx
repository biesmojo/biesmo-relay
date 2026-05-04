"use client";

import { useState, useEffect } from 'react';

interface Customer {
  id: number;
  name: string;
  email: string;
  phone: string;
  created_at: string;
}

export default function AdminPage() {
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState('');

  useEffect(() => {
    fetchCustomers();
  }, []);

  const fetchCustomers = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/customers');
      if (!res.ok) throw new Error('Gagal memuat data');
      const data = await res.json();
      setCustomers(data);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const filteredCustomers = customers.filter((c) =>
    c.name.toLowerCase().includes(filter.toLowerCase()) ||
    c.email.toLowerCase().includes(filter.toLowerCase())
  );

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background to-zinc-50 dark:from-black dark:to-zinc-950 flex items-center justify-center">
        <div className="text-foreground/50">
          <div className="w-8 h-8 border-2 border-foreground/20 border-t-foreground rounded-full animate-spin mx-auto mb-2" />
          <p>Memuat data pelanggan...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-zinc-50 dark:from-black dark:to-zinc-950">
      {/* Header */}
      <div className="bg-background/80 backdrop-blur-md border-b border-foreground/10 p-6 sticky top-0 z-50">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-black text-foreground">Admin Panel</h1>
            <p className="text-foreground/60">Kelola Pelanggan, Sesi, Tiket, dan Event</p>
          </div>
          <div className="flex gap-3">
            <a href="/chat" className="px-6 py-2 rounded-xl bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold hover:shadow-lg transition-all">
              Live Chat
            </a>
            <button onClick={fetchCustomers} className="px-6 py-2 rounded-xl border border-foreground/20 text-foreground hover:border-foreground/40 backdrop-blur hover:bg-foreground/5 transition-all">
              Refresh
            </button>
          </div>
        </div>
      </div>

      <div className="max-w-6xl mx-auto p-6">
        {error && (
          <div className="bg-red-500/10 border border-red-500/30 text-red-400 p-4 rounded-2xl mb-6 backdrop-blur">
            Error: {error} (Backend di localhost:8080?)
            <button onClick={fetchCustomers} className="ml-4 underline hover:no-underline">
              Coba lagi
            </button>
          </div>
        )}

        {/* Filter */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="Cari nama atau email pelanggan..."
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="w-full max-w-md p-4 bg-background/50 backdrop-blur border border-foreground/20 rounded-2xl focus:border-foreground/50 focus:outline-none focus:ring-2 focus:ring-foreground/20 text-foreground placeholder-foreground/40"
          />
        </div>

        {/* Customers Table */}
        <div className="bg-background/80 backdrop-blur-md rounded-3xl border border-foreground/10 overflow-hidden shadow-2xl">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="bg-foreground/5">
                  <th className="p-6 text-left font-bold text-foreground text-sm uppercase tracking-wider">ID</th>
                  <th className="p-6 text-left font-bold text-foreground text-sm uppercase tracking-wider">Nama Pelanggan</th>
                  <th className="p-6 text-left font-bold text-foreground text-sm uppercase tracking-wider">Email</th>
                  <th className="p-6 text-left font-bold text-foreground text-sm uppercase tracking-wider">Telepon</th>
                  <th className="p-6 text-left font-bold text-foreground text-sm uppercase tracking-wider">Dibuat Pada</th>
                </tr>
              </thead>
              <tbody>
                {filteredCustomers.map((customer) => (
                  <tr key={customer.id} className="border-t border-foreground/10 hover:bg-foreground/5 transition-all">
                    <td className="p-6 font-mono text-sm text-foreground/80">{customer.id}</td>
                    <td className="p-6 font-semibold text-foreground">{customer.name}</td>
                    <td className="p-6 text-foreground/90 max-w-md truncate">{customer.email}</td>
                    <td className="p-6 text-foreground/90">{customer.phone || '-'}</td>
                    <td className="p-6 text-sm text-foreground/70">{new Date(customer.created_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' })}</td>
                  </tr>
                ))}
                {filteredCustomers.length === 0 && (
                  <tr>
                    <td colSpan={5} className="p-12 text-center text-foreground/50">
                      <div className="flex flex-col items-center gap-2">
                        <svg fill="none" viewBox="0 0 24 24" className="w-12 h-12 text-foreground/30 mx-auto">
                          <path stroke="currentColor" strokeWidth="1.5" d="M9 12h3.75M14.5 9.75h1.75M9.75 16.25h6.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                        <p className="font-semibold">Tidak ada data pelanggan</p>
                        <p className="text-sm">{filter ? 'Tidak ditemukan' : 'Mulai dengan backend API atau isi data'}</p>
                      </div>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Quick Actions */}
          <div className="p-8 border-t border-foreground/10 bg-foreground/5 grid md:grid-cols-4 gap-4">
            <a href="/admin/sessions" className="group flex flex-col items-center p-6 rounded-2xl hover:bg-background transition-all border border-foreground/20 hover:border-foreground/40 hover:shadow-xl">
              <div className="w-12 h-12 bg-gradient-to-r from-orange-500 to-red-500 rounded-2xl flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                <svg fill="none" viewBox="0 0 24 24" className="w-6 h-6 text-white">
                  <path stroke="currentColor" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.97-4.03 9-9 9a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.97 4.03-9 9-9s9 4.03 9 9z"/>
                </svg>
              </div>
              <h3 className="font-bold text-lg mb-1 text-foreground group-hover:text-foreground/90">Sesi Chat</h3>
              <p className="text-sm text-foreground/60">Lihat semua percakapan</p>
            </a>
            <a href="/admin/tickets" className="group flex flex-col items-center p-6 rounded-2xl hover:bg-background transition-all border border-foreground/20 hover:border-foreground/40 hover:shadow-xl">
              <div className="w-12 h-12 bg-gradient-to-r from-yellow-500 to-orange-500 rounded-2xl flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                <svg fill="none" viewBox="0 0 24 24" className="w-6 h-6 text-white">
                  <path stroke="currentColor" strokeWidth="2" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <h3 className="font-bold text-lg mb-1 text-foreground group-hover:text-foreground/90">Tiket</h3>
              <p className="text-sm text-foreground/60">Kelola tiket support</p>
            </a>
            <a href="/admin/events" className="group flex flex-col items-center p-6 rounded-2xl hover:bg-background transition-all border border-foreground/20 hover:border-foreground/40 hover:shadow-xl">
              <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                <svg fill="none" viewBox="0 0 24 24" className="w-6 h-6 text-white">
                  <path stroke="currentColor" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
              </div>
              <h3 className="font-bold text-lg mb-1 text-foreground group-hover:text-foreground/90">Event & Rules</h3>
              <p className="text-sm text-foreground/60">Lihat event dan rules firing</p>
            </a>
            <a href="/admin/kb" className="group flex flex-col items-center p-6 rounded-2xl hover:bg-background transition-all border border-foreground/20 hover:border-foreground/40 hover:shadow-xl">
              <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-indigo-500 rounded-2xl flex items-center justify-center mb-3 group-hover:scale-110 transition-transform">
                <svg fill="none" viewBox="0 0 24 24" className="w-6 h-6 text-white">
                  <path stroke="currentColor" strokeWidth="2" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <h3 className="font-bold text-lg mb-1 text-foreground group-hover:text-foreground/90">Knowledge Base</h3>
              <p className="text-sm text-foreground/60">Kelola artikel KB</p>
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

