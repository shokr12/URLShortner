'use client';

import { useState } from "react";
import { Shield, Key, Globe, Bell, Save } from "lucide-react";

export default function SettingsPage() {
  const [apiKey, setApiKey] = useState("sk_live_51P...9zK");
  const [copied, setCopied] = useState(false);

  const copyKey = () => {
    navigator.clipboard.writeText(apiKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="p-12 max-w-4xl mx-auto space-y-12">
      <div className="space-y-2">
        <h2 className="text-4xl font-bold tracking-tight text-on-surface">Settings & Security</h2>
        <p className="text-lg text-outline">Manage your digital identity and API integration.</p>
      </div>

      <div className="space-y-8">
        {/* API Section */}
        <section className="bg-surface-container-lowest p-8 rounded-md shadow-sm border border-outline/5 space-y-6">
          <div className="flex items-center space-x-3">
             <Key className="text-primary" size={24} />
             <h3 className="text-xl font-bold text-on-surface">API Credentials</h3>
          </div>
          <p className="text-sm text-outline">Use this key to authenticate requests from your custom applications.</p>
          
          <div className="flex items-center space-x-4">
            <div className="flex-1 p-4 bg-surface-container-low rounded-md font-mono text-sm border border-outline/10 text-on-surface">
              {apiKey}
            </div>
            <button 
              onClick={copyKey}
              className="px-6 py-4 bg-surface-container-high hover:bg-surface-container-highest rounded-md font-bold text-primary transition-colors"
            >
              {copied ? "Copied!" : "Copy Key"}
            </button>
          </div>
        </section>

        {/* Preferences */}
        <section className="bg-surface-container-lowest p-8 rounded-md shadow-sm border border-outline/5 space-y-6">
          <div className="flex items-center space-x-3">
             <Globe className="text-primary" size={24} />
             <h3 className="text-xl font-bold text-on-surface">Global Preferences</h3>
          </div>
          
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 hover:bg-surface-container-low rounded-md transition-colors">
              <div>
                <p className="font-bold text-on-surface">Default Expiration</p>
                <p className="text-sm text-outline">Set how long links remain active by default.</p>
              </div>
              <select className="bg-surface-container-high px-4 py-2 rounded-md font-medium outline-none">
                <option>Never</option>
                <option>30 Days</option>
                <option>7 Days</option>
                <option>24 Hours</option>
              </select>
            </div>

            <div className="flex items-center justify-between p-4 hover:bg-surface-container-low rounded-md transition-colors">
              <div>
                <p className="font-bold text-on-surface">Analytics Tracking</p>
                <p className="text-sm text-outline">Enable detailed session tracking (IP, Device, Location).</p>
              </div>
              <div className="w-12 h-6 bg-primary rounded-full relative">
                <div className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full"></div>
              </div>
            </div>
          </div>
        </section>

        <div className="flex justify-end">
           <button className="flex items-center space-x-2 px-8 py-4 bg-primary text-white font-bold rounded-md hover:opacity-90 transition-all shadow-lg shadow-primary/20">
             <Save size={18} />
             <span>Save Configuration</span>
           </button>
        </div>
      </div>
    </div>
  );
}
