import { useEffect, useState } from 'react';
import type { Island } from '../lib/parameters';
import { Icon } from './Icon';

export function MapPreview({ island, loading, error, retry }: { island: Island | null; loading: boolean; error: string; retry: () => void }) {
  const [url, setURL] = useState('');
  useEffect(() => {
    if (!island) return;
    const next = URL.createObjectURL(new Blob([island.svg], { type: 'image/svg+xml' }));
    setURL(next);
    return () => URL.revokeObjectURL(next);
  }, [island]);
  return <section className="map-card card" aria-labelledby="map-title" aria-busy={loading}>
    <div className="map-heading"><div><span className="eyebrow">A PLACE ALL YOUR OWN</span><h2 id="map-title">{island?.name ?? 'Your next little escape'}</h2></div><span className={`map-status ${loading ? 'busy' : ''}`}><span />{loading ? 'Dreaming…' : error ? 'Let’s try again' : 'Ready to explore'}</span></div>
    <div className="map-canvas">
      {url && island ? <img className="island-map" src={url} width="800" height="800" alt={`${island.name}, a pastel topographic island map. Seed ${island.parameters.seed}.`} /> : <div className="empty-map"><Icon name="sparkle" /><p>A tiny world is taking shape…</p></div>}
      <div className="map-label" aria-hidden="true"><Icon name="heart" /> A happy little nowhere</div>
      {loading && <div className="loading-badge" role="status"><span className="spinner" />Finding your shoreline…</div>}
      {error && <div className="error-overlay" role="alert"><Icon name="wave" /><strong>A little bump in the tide</strong><p>{error}</p><button className="button" onClick={retry}>Try again<Icon name="reset" /></button></div>}
    </div>
    <div className="map-footer"><div className="legend"><span><i className="sea" />Ocean</span><span><i className="sand" />Shore</span><span><i className="green" />Highlands</span></div><span className="land-stat">{island ? `${island.landPercent}% land` : 'Charting the coast…'}</span></div>
    <div className="map-meta"><span>SEED <b>{island?.parameters.seed ?? '—'}</b></span><span>800 × 800 · VECTOR MAP</span></div>
  </section>;
}
