'use client';

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { ArrowLeft, MousePointer2, Clock, Globe, Laptop } from "lucide-react";
import Link from "next/link";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export default function AnalyticsDetail() {
  const { key } = useParams();
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/analytics/${key}`)
      .then(res => res.json())
      .then(data => {
        setData(data);
        setLoading(false);
      })
      .catch(err => {
        console.error(err);
        setLoading(false);
      });
  }, [key]);

  if (loading) return <div className="p-12 text-center text-outline">Decrypting analytics...</div>;
  if (!data) return <div className="p-12 text-center text-outline">No asset history found for /{key}</div>;

  return (
    <div className="p-12 max-w-7xl mx-auto space-y-12">
      <Link href="/links" className="inline-flex items-center space-x-2 text-outline hover:text-primary transition-colors font-medium">
        <ArrowLeft size={18} />
        <span>Back to Assets</span>
      </Link>

      <div className="space-y-2">
        <h2 className="text-4xl font-bold tracking-tight text-on-surface">/{key} Analytics</h2>
        <p className="text-lg text-outline">Tracking performance for your editorial asset.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
        <MetricCard label="Total Reach" value={data.click_count} icon={<MousePointer2 size={20} />} />
        <MetricCard label="Morning Clicks" value={data.clicks_today} icon={<Clock size={20} />} />
        <MetricCard label="Top Region" value="North America" icon={<Globe size={20} />} />
        <MetricCard label="Primary Device" value="Desktop" icon={<Laptop size={20} />} />
      </div>

      <div className="bg-surface-container-lowest p-8 rounded-md shadow-sm border border-outline/5 space-y-8">
        <h3 className="text-xl font-bold text-on-surface">Reach Over Time</h3>
        <div className="h-80 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data.history}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#E2E8F0" />
              <XAxis dataKey="date" axisLine={false} tickLine={false} tick={{fill: '#727785', fontSize: 12}} dy={10} />
              <YAxis axisLine={false} tickLine={false} tick={{fill: '#727785', fontSize: 12}} />
              <Tooltip 
                contentStyle={{ backgroundColor: '#ffffff', borderRadius: '12px', border: 'none', boxShadow: '0 10px 15px -3px rgb(0 0 0 / 0.1)' }}
              />
              <Line 
                type="monotone" 
                dataKey="clicks" 
                stroke="#0058be" 
                strokeWidth={3} 
                dot={{ fill: '#0058be', strokeWidth: 2, r: 4 }}
                activeDot={{ r: 6, strokeWidth: 0 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

function MetricCard({ label, value, icon }: { label: string; value: any; icon: React.ReactNode }) {
  return (
    <div className="p-6 bg-surface-container-lowest rounded-md shadow-sm border border-outline/5 space-y-2 hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between text-outline">
        <p className="text-xs font-bold uppercase tracking-widest">{label}</p>
        {icon}
      </div>
      <p className="text-3xl font-bold text-on-surface tracking-tight">{value}</p>
    </div>
  );
}
