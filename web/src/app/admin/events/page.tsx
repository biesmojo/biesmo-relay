"use client";

import { useState, useEffect } from 'react';

interface Event {
  id: number;
  type: string;
  source: string;
  rule_matched?: string;
  action: string;
  status: string;
  time: string;
  payload?: any;
}

export default function EventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [expanded, setExpanded] = useState<Set<number>>(new Set());

  useEffect(() => {
    loadEvents();
    const interval = setInterval(loadEvents, 10000);
    return () => clearInterval(interval);
  }, []);

  const loadEvents = async () => {
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/admin/events`, {
        headers: { 'Authorization': 'Bearer ' + process.env.NEXT_PUBLIC_ADMIN_TOKEN || 'test' }
      });
      const data = await res.json();
      setEvents(data);
    } catch (err) {
      console.error(err);
    }
  };

  const toggleExpand = (id: number) => {
    const newExpanded = new Set(expanded);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpanded(newExpanded);
  };

  const statusColor = (status: string) => {
    if (status === 'pending') return 'bg-yellow-100 text-yellow-800';
    if (status === 'sent') return 'bg-green-100 text-green-800';
    return 'bg-red-100 text-red-800';
  };

  return (
    <div className="p-6">
      <div className="mb-8 flex justify-between items-center">
        <h1 className="text-3xl font-bold">Event Explorer</h1>
        <button onClick={loadEvents} className="px-4 py-2 bg-blue-500 text-white rounded-lg">
          Refresh
        </button>
      </div>
      <div className="bg-white rounded-xl shadow-lg overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50">
              <th className="p-4 text-left font-bold">Type</th>
              <th className="p-4 text-left font-bold">Source</th>
              <th className="p-4 text-left font-bold">Rule</th>
              <th className="p-4 text-left font-bold">Action</th>
              <th className="p-4 text-left font-bold">Status</th>
              <th className="p-4 text-left font-bold">Time</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event) => (
              <tr key={event.id} className="border-t hover:bg-gray-50">
                <td className="p-4 font-semibold">{event.type}</td>
                <td className="p-4">{event.source}</td>
                <td className="p-4">{event.rule_matched || '-'}</td>
                <td className="p-4">{event.action}</td>
                <td className="p-4">
                  <span className={`px-3 py-1 rounded-full text-xs font-bold ${statusColor(event.status)}`}>
                    {event.status}
                  </span>
                </td>
                <td className="p-4 text-sm">{new Date(event.time).toLocaleString('id-ID')}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {events.length === 0 && (
          <div className="p-12 text-center text-gray-500">
            No events yet. Try firing test event.
          </div>
        )}
      </div>
    </div>
  );
}

