"use client";

import { useState, useRef, useEffect } from 'react';

interface Message {
  role: 'user' | 'agent';
  content: string;
}

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const sendMessage = async () => {
    if (!input.trim() || isLoading) return;

    const userMessage: Message = { role: 'user', content: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    // Simulate AI response
const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const body = { message: input };
      if (sessionId) body.session_id = sessionId;

      try {
        const res = await fetch(`${apiUrl}/api/chat`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        if (!res.ok) throw new Error('API error');
        const data = await res.json();
        setMessages((prev) => [...prev, { role: 'agent', content: data.reply }]);
        setSessionId(data.session_id);
      } catch (err) {
        setMessages((prev) => [...prev, { role: 'agent', content: `Maaf, terjadi kesalahan: ${err.message}` }]);
      }
      setIsLoading(false);
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-zinc-50 dark:from-black dark:to-zinc-950 flex flex-col">
      {/* Header */}
      <div className="bg-background/80 backdrop-blur-md border-b border-foreground/10 p-6 flex items-center gap-4 sticky top-0 z-50">
        <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center">
          <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5 text-white">
            <path stroke="currentColor" strokeWidth="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <div>
          <h1 className="text-2xl font-black text-foreground">Relay Chat</h1>
          <p className="text-sm text-foreground/60">AI Agent siap membantu</p>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4 pb-24">
        {messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-96 text-foreground/50">
            <div className="w-24 h-24 bg-foreground/10 rounded-2xl p-6 mb-4 flex items-center justify-center">
              <svg fill="none" viewBox="0 0 24 24" className="w-12 h-12 text-foreground/30">
                <path stroke="currentColor" strokeWidth="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <h2 className="text-2xl font-bold mb-2">Mulai percakapan</h2>
            <p>Kirim pesan pertama Anda untuk memulai chat dengan AI Agent.</p>
          </div>
        ) : (
          messages.map((message, index) => (
            <div key={index} className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}>
              <div
                className={`max-w-2xl px-6 py-4 rounded-2xl shadow-lg ${
                  message.role === 'user'
                    ? 'bg-gradient-to-r from-blue-500 to-purple-600 text-white'
                    : 'bg-background/80 backdrop-blur-md border border-foreground/20'
                }`}
              >
                <p className="whitespace-pre-wrap">{message.content}</p>
              </div>
            </div>
          ))
        )}
        {isLoading && (
          <div className="flex justify-start">
            <div className="bg-background/80 backdrop-blur-md border border-foreground/20 px-6 py-4 rounded-2xl max-w-2xl">
              <div className="flex items-center gap-2 text-foreground/60">
                <div className="w-6 h-6 border-2 border-foreground/20 border-t-foreground rounded-full animate-spin" />
                <span>Relay AI sedang mengetik...</span>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <div className="bg-background/80 backdrop-blur-md border-t border-foreground/10 p-6 sticky bottom-0">
        <div className="max-w-4xl mx-auto flex items-end gap-3">
          <div className="flex-1 relative">
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyPress}
              placeholder="Ketik pesan Anda di sini..."
              className="w-full p-4 pr-16 resize-none h-16 bg-background/50 backdrop-blur-md border border-foreground/20 rounded-2xl focus:border-foreground/50 focus:outline-none focus:ring-2 focus:ring-foreground/20 text-foreground placeholder-foreground/40"
              disabled={isLoading}
              rows={1}
            />
          </div>
          <button
            onClick={sendMessage}
            disabled={!input.trim() || isLoading}
            className="w-14 h-14 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 disabled:opacity-50 disabled:cursor-not-allowed rounded-2xl flex items-center justify-center transition-all shadow-lg hover:shadow-xl text-white"
            aria-label="Kirim pesan"
          >
            <svg fill="none" viewBox="0 0 24 24" className="w-5 h-5">
              <path
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                d="m4 12 7.75 7.75L20 7"
              />
            </svg>
          </button>
        </div>
        <p className="text-xs text-foreground/50 mt-2 text-center">
          Mock chat. Integrasi API backend nanti via /api/chat endpoint.
        </p>
      </div>
    </div>
  );
}

