import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Link from "next/link";
import { LayoutDashboard, Link2, BarChart3, Settings } from "lucide-react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Shortify | Modern URL Shortener",
  description: "High-performance dashboard for digital curators.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="flex h-screen bg-background">
          {/* Sidebar */}
          <aside className="w-64 glass-panel border-r border-outline/15 flex flex-col pt-8">
            <div className="px-8 mb-12">
              <h1 className="text-2xl font-bold text-primary tracking-tight">Shortify</h1>
            </div>
            
            <nav className="flex-1 space-y-2 px-4">
              <NavItem href="/" icon={<LayoutDashboard size={20} />} label="Dashboard" />
              <NavItem href="/links" icon={<Link2 size={20} />} label="Management" />
              <NavItem href="/analytics" icon={<BarChart3 size={20} />} label="Analytics" />
              <NavItem href="/settings" icon={<Settings size={20} />} label="Settings" />
            </nav>
          </aside>

          {/* Main Content */}
          <main className="flex-1 overflow-y-auto bg-surface-container-low">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}

function NavItem({ href, icon, label }: { href: string; icon: React.ReactNode; label: string }) {
  return (
    <Link 
      href={href}
      className="flex items-center space-x-3 px-4 py-3 rounded-md text-on-surface hover:bg-surface-container-high transition-colors group"
    >
      <span className="text-secondary group-hover:text-primary transition-colors">{icon}</span>
      <span className="font-medium">{label}</span>
    </Link>
  );
}
