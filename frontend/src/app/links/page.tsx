'use client';

import { useState, useEffect } from "react";
import { Search, Filter, MoreHorizontal, Copy, Trash2, ExternalLink } from "lucide-react";
import Link from "next/link";

export default function LinkManagement() {
  const [links, setLinks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("url-shortner-neon-omega.vercel.app/links")
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

  const handleDelete = async (key: string) => {
    if (!confirm("Are you sure you want to delete this link?")) return;
    try {
      await fetch(`http://localhost:8080/links/${key}`, { method: "DELETE" });
      setLinks(links.filter(l => l.short_url !== key));
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="p-12 max-w-7xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-on-surface">Link Management</h2>
        <div className="flex items-center space-x-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-outline" size={18} />
            <input 
              type="text" 
              placeholder="Search links..." 
              className="pl-10 pr-4 py-2 bg-surface-container-lowest border border-transparent focus:border-primary/20 rounded-md outline-none transition-all"
            />
          </div>
          <button className="flex items-center space-x-2 px-4 py-2 bg-surface-container-high rounded-md font-medium hover:bg-surface-container-highest transition-colors">
            <Filter size={18} />
            <span>Filter</span>
          </button>
        </div>
      </div>

      <div className="bg-surface-container-lowest rounded-md shadow-sm overflow-hidden border border-outline/5">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-surface-container-low/50">
              <th className="px-6 py-4 text-xs font-bold text-outline uppercase tracking-widest">Digital Asset</th>
              <th className="px-6 py-4 text-xs font-bold text-outline uppercase tracking-widest">Destination</th>
              <th className="px-6 py-4 text-xs font-bold text-outline uppercase tracking-widest">Performance</th>
              <th className="px-6 py-4 text-xs font-bold text-outline uppercase tracking-widest">Status</th>
              <th className="px-6 py-4 text-xs font-bold text-outline uppercase tracking-widest text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-outline/5">
            {loading ? (
              <tr><td colSpan={5} className="px-6 py-12 text-center text-outline">Optimizing view...</td></tr>
            ) : links.length === 0 ? (
              <tr><td colSpan={5} className="px-6 py-12 text-center text-outline">No digital assets found.</td></tr>
            ) : (
              links.map((link) => (
                <tr key={link.short_url} className="hover:bg-surface-container-low transition-colors group">
                  <td className="px-6 py-4">
                    <p className="font-bold text-primary flex items-center space-x-2">
                       <span>/{link.short_url}</span>
                       <Copy 
                        size={14} 
                        className="opacity-0 group-hover:opacity-50 cursor-pointer hover:opacity-100 transition-opacity" 
                        onClick={() => navigator.clipboard.writeText(`http://localhost:8080/${link.short_url}`)}
                       />
                    </p>
                  </td>
                  <td className="px-6 py-4 max-w-xs truncate">
                    <span className="text-sm text-outline">{link.original_url}</span>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center space-x-2">
                      <span className="font-bold text-on-surface">{link.click_count}</span>
                      <span className="text-xs text-outline uppercase tracking-tight">Clicks</span>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className="inline-block px-3 py-1 bg-green-100 text-green-700 text-xs font-bold rounded-full uppercase tracking-widest">Active</span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end space-x-3 opacity-0 group-hover:opacity-100 transition-opacity">
                      <Link href={`/analytics/${link.short_url}`} className="p-2 hover:bg-surface-container-high rounded-md text-outline hover:text-primary">
                        <ExternalLink size={18} />
                      </Link>
                      <button 
                        onClick={() => handleDelete(link.short_url)}
                        className="p-2 hover:bg-red-50 rounded-md text-outline hover:text-red-500"
                      >
                        <Trash2 size={18} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
