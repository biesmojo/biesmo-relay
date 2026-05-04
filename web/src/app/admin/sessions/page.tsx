"use client";

import { useState, useEffect } from 'react';

interface Session {
  id: number;
  customer_id: number;
  customer_name?: string;
  status: string;
  sentiment: string;
  created_at: string;
}

export default function SessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filter, setFilter] = useState('');

  useEffect(() => {
    fetchSessions();
  }, []);

  const fetchSessions = async () => {
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/sessions`);
      if (!res.ok) throw new Error('Gagal memuat sesi');
      const data = await res.json();
      setSessions(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const filtered = sessions.filter(s =>
    s.status.includes(filter) || s.sentiment.includes(filter) || s.customer_name?.includes(filter)
  );

  if (loading) return <div className="p-8 text-center">Loading...</div>;

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Sesi Chat</h1>
      <input
        placeholder="Filter..."
        value={filter}
        onChange={e => setFilter(e.target.value)}
        className="w-full max-w-md p-3 border rounded-lg mb-6"
      />
      <table className="w-full border-collapse bg-white rounded-lg shadow-lg">
        <thead>
          <tr className="bg-gray-50">
            <th className="p-4 text-left font-bold">ID</th>
            <th className="p-4 text-left font-bold">Customer</th>
            <th className="p-4 text-left font-bold">Status</th>
            <th className="p-4 text-left font-bold">Sentiment</th>
            <th className="p-4 text-left font-bold">Created At</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((session) => (
            <tr key={session.id} className="border-t hover:bg-gray-50">
              <td className="p-4 font-mono">{session.id}</td>
              <td className="p-4 font-semibold">#{session.customer_id}</td>
              <td><span className={`px-3 py-1 rounded-full text-xs font-bold ${session.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>{session.status}</span></td>
              <td><span className={`px-3 py-1 rounded-full text-xs font-bold ${session.sentiment === 'Positif' ? 'bg-green-100 text-green-800' : session.sentiment === 'Negatif' ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'}`}>{session.sentiment}</span></td>
              <td className="p-4 text-sm">{new Date(session.created_at).toLocaleString('id-ID')}</td>
            </tr>
          ))}
          {filtered.length === 0 &amp;&amp; (
            <tr>
              <td colSpan={5} className="p-12 text-center text-gray-500">
                Tidak ada sesi
              </td>
            </tr>
          )}
        </tbody>
      </table>
      {error &amp;&amp; <p className="text-red-500 mt-4">{error}</p>}
    </div>
  );
}

