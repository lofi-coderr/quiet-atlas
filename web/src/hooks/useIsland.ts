import { useEffect, useState } from 'react';
import { islandQuery } from '../lib/parameters';
import type { Island, Parameters } from '../lib/parameters';

export function useIsland(parameters: Parameters | null, retry: number) {
  const [island, setIsland] = useState<Island | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!parameters) { setLoading(false); return; }
    const controller = new AbortController();
    let active = true;
    setLoading(true);
    setError('');
    const timer = setTimeout(async () => {
      try {
        const response = await fetch(`/api/island?${islandQuery(parameters)}`, { signal: controller.signal });
        if (!response.ok) {
          const detail = await response.json() as { error?: string };
          throw new Error(detail.error ?? 'Something went wrong. Try again.');
        }
        const next = await response.json() as Island;
        if (!next.svg || !next.parameters || typeof next.name !== 'string') throw new Error('The map could not be read. Try again.');
        if (active) { setIsland(next); setLoading(false); }
      } catch (cause) {
        if (!active || controller.signal.aborted) return;
        setError(cause instanceof TypeError ? 'Your island could not be reached. Check the local server and try again.'
          : cause instanceof Error ? cause.message : 'Something went wrong. Try again.');
        setLoading(false);
      }
    }, 160);
    return () => { active = false; clearTimeout(timer); controller.abort(); };
  }, [parameters, retry]);

  return { island, loading, error };
}
