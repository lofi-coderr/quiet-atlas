import { useMemo, useState } from 'react';
import { Controls, settingsMatch } from './components/Controls';
import { Icon, IslandMark } from './components/Icon';
import { MapPreview } from './components/MapPreview';
import { useIsland } from './hooks/useIsland';
import { DEFAULTS, parseSeed, randomSeed, saveSVG } from './lib/parameters';

export function App() {
  const [seedText, setSeedText] = useState(String(DEFAULTS.seed));
  const [seaLevel, setSeaLevel] = useState(DEFAULTS.seaLevel);
  const [roughness, setRoughness] = useState(DEFAULTS.roughness);
  const [contours, setContours] = useState(DEFAULTS.contours);
  const [retry, setRetry] = useState(0);
  const [saved, setSaved] = useState(false);
  const seed = parseSeed(seedText);
  const parameters = useMemo(() => seed === null ? null : { seed, seaLevel, roughness, contours }, [seed, seaLevel, roughness, contours]);
  const { island, loading, error } = useIsland(parameters, retry);
  const pending = loading || (!error && !!parameters && !!island && !settingsMatch(parameters, island.parameters));
  function reset() {
    setSeedText(String(DEFAULTS.seed)); setSeaLevel(DEFAULTS.seaLevel);
    setRoughness(DEFAULTS.roughness); setContours(DEFAULTS.contours); setRetry(value => value + 1);
  }
  function download() {
    if (!island || pending) return;
    saveSVG(island); setSaved(true); setTimeout(() => setSaved(false), 1800);
  }
  return <div className="app-shell">
    <header className="site-header"><a className="brand" href="/" aria-label="Quiet Atlas home"><span className="brand-mark"><IslandMark /></span><span>quiet atlas<span className="brand-dot">✦</span></span></a><span className="header-note">Tiny worlds, happy discoveries.</span><span className="local-badge"><span />Your own little corner</span></header>
    <main>
      <div className="intro"><div><span className="intro-tag"><Icon name="sparkle" />A LITTLE WORLD OF POSSIBILITY</span><h1>Find your <span>happy place<svg viewBox="0 0 400 20" preserveAspectRatio="none" aria-hidden="true"><path d="M5 12Q180 0 393 10" /></svg></span><span className="title-dot">.</span></h1><p>A seed, a few small tweaks, and an island that’s entirely yours.</p></div><button className="button download" disabled={!island || pending} onClick={download}><Icon name={saved ? 'check' : 'download'} />{saved ? 'SVG saved!' : 'Save SVG'}<span className="button-tag">YOURS TO KEEP</span></button></div>
      <div className="workspace"><Controls seedText={seedText} seedValid={seed !== null} seaLevel={seaLevel} roughness={roughness} contours={contours} onSeed={setSeedText} onSea={setSeaLevel} onRoughness={setRoughness} onContours={setContours} onRandom={() => setSeedText(String(randomSeed()))} onReset={reset} /><MapPreview island={island} loading={pending} error={error} retry={() => setRetry(value => value + 1)} /></div>
      <aside className="little-tip"><span className="tip-star"><Icon name="sparkle" /></span><p><strong>A little tip</strong>Keep a seed you love. You can always find your way back.</p><span className="tip-doodle" aria-hidden="true">✦</span></aside>
    </main>
    <footer className="site-footer"><span>Made for wandering. Built for wonder.</span><span><Icon name="heart" />No accounts. Just islands.</span></footer>
  </div>;
}
