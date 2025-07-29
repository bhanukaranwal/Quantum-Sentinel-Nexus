// services/dashboard-ui/src/app/providers.tsx
'use client'; // This must be a client component to use state and context.

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';

export function Providers({ children }: { children: React.ReactNode }) {
  // Create a new instance of QueryClient.
  // We use useState to ensure this instance is only created once per component lifecycle.
  const [queryClient] = React.useState(() => new QueryClient());

  return (
    // Provide the client to the rest of the app.
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
