"use client";

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';

interface Message {
  role: string;
  content: string;
}

export default function SessionDetail() {
  const params = useParams();
  const id = params.id as string;
  const [session, setSession] = useState(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSession();
  }, [id]);

  const fetchSession = async () => {
    try {
      const [sessionRes, msgRes] = await Promise.all([
        fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/sessions/${id}`),
        fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/sessions/${id}/messages`)
      ]);
      const sessionData = await sessionRes.json();
      const msgData = await msgRes.json();
      setSession(sessionData);
      setMessages(msgData);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading detail...</div>;
  if (!session) return <div>Sesi tidak ditemukan</div>;

  return (
    <div className="p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Sesi #{session.id}</h1>
        <div className="bg-gray-50 p-6 rounded-xl">
          <h2 className="text-xl font-semibold mb-4">Info Customer</h2>
          <p><strong>Status:</strong> {session.status}</p>
          <p><strong>Sentiment:</strong> {session.sentiment}</p>
          <p><strong>Summary:</strong> {session.summary || ' - '}</p>
          <p><strong>CSAT:</strong> {session.csat_score || ' - '}</p>
        </div>
      </div>

      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-4">Transkrip Pesan</h2>
        <div className="bg-gray-50 rounded-xl p-6 space-y-4 max-h-96 overflow-y-auto">
          {messages.map((msg, i) => (
            <div key={i} className={`p-4 rounded-lg ${msg.role === 'user' ? 'bg-blue-100 ml-auto max-w-2xl text-right' : 'bg-white mr-auto max-w-2xl'}`}>
              <p className="font-semibold text-sm text-gray-500 mb-1">{msg.role}</p>
              <p>{msg.content}</p>
            </div>
          ))}
        </div>
      </div>

      <div>
        <h2 className="text-2xl font-bold mb-4">Tiket Terkait</h2>
        <p>Tiket (jika ada)</p>
      </div>
    </div>
  );
}

