'use client';

import { useState, useEffect } from "react";
import { Plus, Link2, MousePointer2, Calendar } from "lucide-react";
import { motion } from "framer-motion";

export default function Dashboard() {
  const [url, setUrl] = useState("");
  const [stats, setStats] = useState({ total_urls: 0, total_clicks: 0 });
  const [shortenedUrl, setShortenedUrl] = useState<string | null>(null);

  useEffect(() => {
    fetch("http://localhost:8080/stats")
      .then(res => res.json())
      .then(data => setStats(data))
      .catch(err => console.error("Stats fetch error:", err));
  }, []);

  const handleShorten = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/shorten", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ original_url: url }),
      });
      const data = await res.json();
      if (data.short_url) {
        setShortenedUrl(data.short_url);
        // Refresh stats
        fetch("http://localhost:8080/stats").then(res => res.json()).then(setStats);
      }
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="p-12 max-w-7xl mx-auto space-y-12">
      {/* Hero Section */}
      <section className="space-y-6">
        <h2 className="text-4xl font-bold tracking-tight text-on-surface">Architectural Precision</h2>
        <p className="text-lg text-outline max-w-2xl">
          Transform your long URLs into editorial features. Built for performance, designed for digital curators.
        </p>

        <form onSubmit={handleShorten} className="relative group max-w-3xl">
          <div className="absolute inset-y-0 left-6 flex items-center pointer-events-none text-outline group-focus-within:text-primary transition-colors">
            <Link2 size={24} />
          </div>
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Paste your long URL here..."
            className="w-full pl-16 pr-44 py-6 bg-surface-container-lowest rounded-md shadow-sm focus:ring-2 focus:ring-primary/20 focus:bg-surface-bright transition-all outline-none text-lg border border-transparent"
            required
          />
          <button
            type="submit"
            className="absolute right-3 top-3 bottom-3 px-8 bg-gradient-to-r from-primary to-primary-container text-white font-semibold rounded-md hover:opacity-90 transition-all flex items-center space-x-2"
          >
            <Plus size={20} />
            <span>Shorten</span>
          </button>
        </form>

        {shortenedUrl && (
          <motion.div 
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="p-6 bg-primary/5 border border-primary/10 rounded-md flex items-center justify-between max-w-3xl"
          >
            <div className="space-y-1">
              <span className="text-sm font-bold text-primary uppercase tracking-widest">Ready to use</span>
              <p className="text-xl font-medium text-on-surface">{shortenedUrl}</p>
            </div>
            <button 
              onClick={() => navigator.clipboard.writeText(shortenedUrl)}
              className="px-4 py-2 text-primary font-semibold hover:bg-primary/10 rounded-md transition-colors"
            >
              Copy Link
            </button>
          </motion.div>
        )}
      </section>

      {/* Stats Overview */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <StatCard 
          label="Total Links" 
          value={stats.total_urls.toLocaleString()} 
          icon={<Link2 className="text-primary" />} 
          trend="+12 this week"
        />
        <StatCard 
          label="Total Clicks" 
          value={stats.total_clicks.toLocaleString()} 
          icon={<MousePointer2 className="text-primary" />} 
          trend="+542 this week"
        />
        <StatCard 
          label="Avg. CTR" 
          value="4.2%" 
          icon={<Calendar className="text-primary" />} 
          trend="+0.8% increase"
        />
      </section>
    </div>
  );
}

function StatCard({ label, value, icon, trend }: { label: string; value: string; icon: React.ReactNode; trend: string }) {
  return (
    <div className="p-8 bg-surface-container-lowest rounded-md shadow-sm space-y-4 hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between">
        <div className="p-3 bg-surface-container-low rounded-md">
          {icon}
        </div>
        <span className="text-xs font-bold text-green-600 uppercase tracking-widest bg-green-50 px-2 py-1 rounded-full">{trend}</span>
      </div>
      <div>
        <p className="text-sm font-bold text-outline uppercase tracking-widest">{label}</p>
        <p className="text-5xl font-bold tracking-tight text-on-surface mt-1">{value}</p>
      </div>
    </div>
  );
}
