// services/dashboard-ui/src/app/page.tsx
'use client'; // This directive is required for components that use client-side hooks like useState and React Query.

import { useQuery } from '@tanstack/react-query';
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';
import React, { useState, useEffect } from 'react';

// --- Mock Data ---
// In a real application, this data would be fetched from the backend.
const mockKpiData = {
  totalTxns: 125890,
  alertsToday: 73,
  modelAccuracy: 0.985,
  p95Latency: 4.2,
};

const mockChartData = [
  { name: '08:00', high: 12, critical: 2 },
  { name: '09:00', high: 15, critical: 5 },
  { name: '10:00', high: 22, critical: 8 },
  { name: '11:00', high: 18, critical: 4 },
  { name: '12:00', high: 25, critical: 11 },
  { name: '13:00', high: 19, critical: 6 },
];

// --- API Fetching Functions ---

// Fetches KPI data from the backend.
// In a real app, this would hit a dedicated API endpoint.
const fetchKpiData = async () => {
  // const res = await fetch('/api/kpis');
  // if (!res.ok) throw new Error('Network response was not ok');
  // return res.json();
  return new Promise((resolve) => setTimeout(() => resolve(mockKpiData), 500));
};

// Fetches real-time alerts.
// This simulates a WebSocket connection.
const useAlerts = () => {
  const [alerts, setAlerts] = useState([]);

  useEffect(() => {
    // In a real app, you would connect to the alert-router-svc WebSocket here.
    // const ws = new WebSocket('ws://alerts.qsn.local/');
    // ws.onmessage = (event) => {
    //   const newAlert = JSON.parse(event.data);
    //   setAlerts(prev => [newAlert, ...prev].slice(0, 10)); // Keep last 10 alerts
    // };
    // return () => ws.close();

    // Mock WebSocket behavior
    const interval = setInterval(() => {
      const newAlert = {
        id: `txn-${Math.floor(Math.random() * 100000)}`,
        severity: Math.random() > 0.8 ? 'CRITICAL' : 'HIGH',
        score: Math.random() * (1.0 - 0.85) + 0.85,
        timestamp: new Date().toISOString(),
      };
      // @ts-ignore
      setAlerts((prev) => [newAlert, ...prev].slice(0, 10));
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  return alerts;
};

// --- Components ---

const KpiCard = ({ title, value, unit }) => (
  <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
    <h3 className="text-sm font-medium text-gray-400">{title}</h3>
    <p className="mt-2 text-3xl font-bold text-white">
      {value}
      <span className="text-lg font-medium text-gray-300 ml-2">{unit}</span>
    </p>
  </div>
);

const AlertsTable = ({ alerts }) => (
  <div className="bg-gray-800 p-6 rounded-lg shadow-lg mt-8">
    <h3 className="text-lg font-bold text-white mb-4">Real-Time Alerts</h3>
    <div className="overflow-x-auto">
      <table className="min-w-full text-left">
        <thead>
          <tr>
            <th className="p-2 text-sm font-semibold text-gray-400">Severity</th>
            <th className="p-2 text-sm font-semibold text-gray-400">Score</th>
            <th className="p-2 text-sm font-semibold text-gray-400">Timestamp</th>
          </tr>
        </thead>
        <tbody>
          {alerts.map((alert) => (
            <tr key={alert.id} className="border-t border-gray-700">
              <td className="p-2">
                <span
                  className={`px-2 py-1 rounded-full text-xs font-bold ${
                    alert.severity === 'CRITICAL'
                      ? 'bg-red-500 text-white'
                      : 'bg-yellow-500 text-black'
                  }`}
                >
                  {alert.severity}
                </span>
              </td>
              <td className="p-2 font-mono text-indigo-300">
                {alert.score.toFixed(4)}
              </td>
              <td className="p-2 text-gray-400">
                {new Date(alert.timestamp).toLocaleTimeString()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

// --- Main Page Component ---

export default function Home() {
  const { data: kpiData, isLoading: kpiLoading } = useQuery({
    queryKey: ['kpis'],
    queryKey: ['kpis'],
    queryFn: fetchKpiData,
  });
  const alerts = useAlerts();

  return (
    <main className="bg-gray-900 min-h-screen text-white p-8">
      <div className="container mx-auto">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-500">
            QSN Dashboard
          </h1>
          <p className="text-gray-400">
            Real-time fraud detection and system monitoring
          </p>
        </header>

        {/* KPI Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <KpiCard
            title="Total Transactions (24h)"
            value={kpiLoading ? '...' : kpiData.totalTxns.toLocaleString()}
            unit=""
          />
          <KpiCard
            title="Alerts Today"
            value={kpiLoading ? '...' : kpiData.alertsToday}
            unit=""
          />
          <KpiCard
            title="Model Accuracy"
            value={kpiLoading ? '...' : `${(kpiData.modelAccuracy * 100).toFixed(2)}`}
            unit="%"
          />
          <KpiCard
            title="P95 Scoring Latency"
            value={kpiLoading ? '...' : kpiData.p95Latency}
            unit="ms"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
          {/* Chart */}
          <div className="lg:col-span-2 bg-gray-800 p-6 rounded-lg shadow-lg">
            <h3 className="text-lg font-bold text-white mb-4">
              Alerts by Severity (Hourly)
            </h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={mockChartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#4A5568" />
                <XAxis dataKey="name" stroke="#A0AEC0" />
                <YAxis stroke="#A0AEC0" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#1A202C',
                    border: '1px solid #4A5568',
                  }}
                />
                <Legend />
                <Bar dataKey="high" stackId="a" fill="#ECC94B" />
                <Bar dataKey="critical" stackId="a" fill="#F56565" />
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Real-time Alerts */}
          <div className="lg:col-span-1">
            <AlertsTable alerts={alerts} />
          </div>
        </div>
      </div>
    </main>
  );
}
