'use client';

import { useState, useEffect } from "react";
import { BarChart3, TrendingUp, Users, Target } from "lucide-react";
import Link from "next/link";

export default function AnalyticsIndex() {
  const [links, setLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/links")
      .then(res => res.json())
      .then(data => {
        setLinks(data || []);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="p-12 max-w-7xl mx-auto space-y-12">
      <div className="space-y-2">
        <h2 className="text-4xl font-bold tracking-tight text-on-surface">Digital Performance</h2>
        <p className="text-lg text-outline">Real-time engagement metrics across your architectural grid.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        <MetricSummary label="Active Campaigns" value={links.length} icon={<Target size={20} />} trend="+2" />
        <MetricSummary label="Unique Visitors" value="1.2k" icon={<Users size={20} />} trend="+15%" />
        <MetricSummary label="Conversion Rate" value="3.8%" icon={<TrendingUp size={20} />} trend="+0.4%" />
        <MetricSummary label="Avg. Latency" value="12ms" icon={<BarChart3 size={20} />} trend="-2ms" />
      </div>

      <div className="space-y-6">
        <h3 className="text-2xl font-bold text-on-surface">Asset Breakdown</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {links.map((link) => (
            <Link 
              key={link.short_url}
              href={`/analytics/${link.short_url}`}
              className="bg-surface-container-lowest p-6 rounded-md shadow-sm border border-outline/5 hover:border-primary/20 hover:shadow-md transition-all group"
            >
              <div className="flex items-center justify-between mb-4">
                <span className="text-xs font-bold text-outline uppercase tracking-widest">Asset /{link.short_url}</span>
                <TrendingUp size={16} className="text-primary opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
              <p className="text-lg font-bold text-on-surface truncate mb-1">{link.original_url}</p>
              <div className="flex items-center space-x-2 text-sm text-outline">
                <span className="font-bold text-primary">{link.click_count}</span>
                <span>total engagements</span>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

function MetricSummary({ label, value, icon, trend }: { label: string; value: any; icon: React.ReactNode; trend: string }) {
  return (
    <div className="p-6 bg-surface-container-high/30 rounded-md border border-outline/10 space-y-3">
      <div className="flex items-center justify-between">
        <div className="p-2 bg-surface-container-lowest rounded-md text-primary shadow-sm">
          {icon}
        </div>
        <span className="text-xs font-bold text-primary bg-primary/10 px-2 py-1 rounded-full">{trend}</span>
      </div>
      <div>
        <p className="text-xs font-bold text-outline uppercase tracking-widest leading-none">{label}</p>
        <p className="text-2xl font-bold text-on-surface tracking-tight mt-1">{value}</p>
      </div>
    </div>
  );
}
