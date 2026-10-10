import type { CSSProperties } from 'react';
import { Icon } from './Icon';
import type { Parameters } from '../lib/parameters';

interface Props {
  seedText: string; seedValid: boolean; seaLevel: number; roughness: number; contours: boolean;
  onSeed: (value: string) => void; onSea: (value: number) => void; onRoughness: (value: number) => void;
  onContours: (value: boolean) => void; onRandom: () => void; onReset: () => void;
}
function Range({ id, title, description, value, min, max, onChange, kind, low, high }:
  { id: string; title: string; description: string; value: number; min: number; max: number;
    onChange: (value: number) => void; kind: 'wave' | 'mountain'; low: string; high: string }) {
  return <div className={`range-group ${kind}`}>
    <div className="control-label"><label htmlFor={id}><span className="label-icon"><Icon name={kind} /></span>{title}</label><output htmlFor={id}>{Math.round(value * 100)}%</output></div>
    <p className="control-description" id={`${id}-hint`}>{description}</p>
    <input id={id} type="range" min={min} max={max} step="0.01" value={value}
      aria-describedby={`${id}-hint`} aria-valuetext={`${Math.round(value * 100)} percent`}
      style={{ '--progress': `${(value-min)/(max-min)*100}%` } as CSSProperties}
      onChange={event => onChange(Number(event.target.value))} />
    <div className="range-captions"><span>{low}</span><span>{high}</span></div>
  </div>;
}
export function Controls(props: Props) {
  return <section className="controls card" aria-labelledby="controls-title">
    <div className="card-heading"><span className="eyebrow">THE LITTLE DETAILS</span><button className="icon-button" onClick={props.onReset} aria-label="Reset settings" title="Reset settings"><Icon name="reset" /></button></div>
    <h2 id="controls-title">Make it yours<span className="pink-dot">.</span></h2>
    <div className="seed-group">
      <label htmlFor="seed">Your seed <span className="tiny-tag">the starting point</span></label>
      <div className={`seed-input ${!props.seedValid ? 'invalid' : ''}`}><span aria-hidden="true">#</span>
        <input id="seed" type="text" inputMode="numeric" value={props.seedText} spellCheck={false} autoComplete="off"
          aria-invalid={!props.seedValid} aria-describedby={props.seedValid ? 'seed-hint' : 'seed-error'}
          onChange={event => props.onSeed(event.target.value)} />
        <button className="icon-button" onClick={props.onRandom} aria-label="Randomize seed" title="Randomize seed"><Icon name="dice" /></button>
      </div>
      {props.seedValid ? <p className="control-description" id="seed-hint">Same seed. Same little world.</p> : <p className="field-error" id="seed-error" role="alert">Use a whole number from 0 to 9007199254740991.</p>}
    </div>
    <div className="divider" />
    <Range id="sea-level" title="Sea level" description="A little more ocean, a little less island." value={props.seaLevel} min={0.25} max={0.7} onChange={props.onSea} kind="wave" low="More land" high="More sea" />
    <Range id="roughness" title="Terrain roughness" description="From mellow hills to wild little ridges." value={props.roughness} min={0} max={1} onChange={props.onRoughness} kind="mountain" low="Smooth" high="Rugged" />
    <div className="divider" />
    <label className="toggle-label" htmlFor="contours"><span className="label-icon lavender"><Icon name="layers" /></span><span><strong>Contour lines</strong><small>A little extra character</small></span>
      <input id="contours" type="checkbox" checked={props.contours} onChange={event => props.onContours(event.target.checked)} /><span className="toggle-track" aria-hidden="true" />
    </label>
    <button className="button surprise" onClick={props.onRandom}><Icon name="sparkle" />Surprise me<Icon name="arrow" /></button>
    <p className="auto-note"><span />Changes grow into a new map automatically.</p>
  </section>;
}
export function settingsMatch(a: Parameters, b: Parameters) {
  return a.seed === b.seed && a.seaLevel === b.seaLevel && a.roughness === b.roughness && a.contours === b.contours;
}
