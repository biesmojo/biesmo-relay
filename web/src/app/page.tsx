export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-zinc-50 dark:from-black dark:to-zinc-950">
      <main className="flex min-h-screen flex-col items-center justify-center py-24 px-8">
        {/* Hero */}
        <div className="max-w-4xl text-center">
          <h1 className="bg-gradient-to-r from-foreground bg-clip-text text-5xl font-black leading-tight tracking-tight text-transparent md:text-7xl lg:text-8xl">
            Relay
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-xl text-foreground/70 dark:text-zinc-400">
            AI-native CRM × Event-driven Workflows. 
            <br />
            <span className="font-semibold text-foreground">
              Agents handle conversations, rules automate actions, CRM tracks everything.
            </span>
          </p>
          
          {/* CTAs */}
          <div className="mt-12 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-center sm:gap-6">
            <a
              href="/chat"
              className="flex h-14 w-full items-center justify-center rounded-2xl bg-foreground px-8 text-xl font-semibold text-background shadow-2xl transition-all hover:shadow-3xl sm:w-64"
            >
              Mulai Chat
            </a>
            <a
              href="/admin"
              className="flex h-14 w-full items-center justify-center rounded-2xl border-2 border-foreground/20 px-8 text-xl font-semibold text-foreground backdrop-blur-sm transition-all hover:border-foreground/40 hover:bg-foreground/5 sm:w-64"
            >
              Admin Panel
            </a>
          </div>

          {/* Features */}
          <div className="mt-32 grid w-full grid-cols-1 gap-8 px-8 md:grid-cols-3">
            <div className="flex flex-col items-center gap-4 rounded-2xl p-8 text-center transition-all hover:bg-foreground/5">
              <div className="h-16 w-16 rounded-2xl bg-foreground/10 p-4">
                <svg fill="none" viewBox="0 0 24 24" className="h-6 w-6 text-foreground"><path stroke="currentColor" stroke-width="2" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
            </div>
              <h3 className="text-2xl font-bold text-foreground">AI Agents</h3>
              <p className="text-foreground/60 dark:text-zinc-400">
                Tool-using LLM agents handle customer conversations, lookup data, create tickets.
              </p>
            </div>

            <div className="flex flex-col items-center gap-4 rounded-2xl p-8 text-center transition-all hover:bg-foreground/5">
              <div className="h-16 w-16 rounded-2xl bg-foreground/10 p-4">
                <svg fill="none" viewBox="0 0 24 24" className="h-6 w-6 text-foreground"><path stroke="currentColor" stroke-width="2" d="M12 2L2 7v10c0 3.31 2.69 6 6 6h8c3.31 0 6-2.69 6-6V7l-10-5z"/></svg>
              </div>
              <h3 className="text-2xl font-bold text-foreground">Smart Rules</h3>
              <p className="text-foreground/60 dark:text-zinc-400">
                Event-driven workflows. Incoming events → rules match → actions fire automatically.
              </p>
            </div>

            <div className="flex flex-col items-center gap-4 rounded-2xl p-8 text-center transition-all hover:bg-foreground/5">
              <div className="h-16 w-16 rounded-2xl bg-foreground/10 p-4">
                <svg fill="none" viewBox="0 0 24 24" className="h-6 w-6 text-foreground"><path stroke="currentColor" stroke-width="2" d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/></svg>
              </div>
              <h3 className="text-2xl font-bold text-foreground">CRM Core</h3>
              <p className="text-foreground/60 dark:text-zinc-400">
                Customers, sessions, tickets, knowledge base. Full audit trail and admin dashboard.
              </p>
            </div>
          </div>

          {/* Footer */}
          <div className="mt-24 flex flex-col items-center gap-2 text-sm text-foreground/50">
            <p>Built for Indonesia • Bahasa Indonesia supported</p>
            <p className="text-xs">AI-native contact center & workflow platform</p>
          </div>
        </div>
      </main>
    </div>
  );
}
