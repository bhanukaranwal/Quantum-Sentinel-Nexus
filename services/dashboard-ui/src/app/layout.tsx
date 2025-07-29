// services/dashboard-ui/src/app/layout.tsx
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import { Providers } from './providers'; // We will create this component next

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'QSN Dashboard',
  description: 'Quantum Sentinel Nexus - Real-time fraud detection and system monitoring',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        {/* The Providers component will wrap the application with the React Query client */}
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
