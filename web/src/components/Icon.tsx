import type { CSSProperties } from 'react';

type IconName = 'sparkle' | 'download' | 'dice' | 'wave' | 'mountain' | 'layers' | 'arrow' | 'check' | 'reset' | 'heart';
const paths: Record<IconName, React.ReactNode> = {
  sparkle: <><path d="m12 3 2.5 6.5L21 12l-6.5 2.5L12 21l-2.5-6.5L3 12l6.5-2.5Z" /></>,
  download: <><path d="M12 3v12m-4-4 4 4 4-4M4 16v4h16v-4" /></>,
  dice: <><rect x="3" y="3" width="18" height="18" rx="5" /><circle cx="8" cy="8" r=".8"/><circle cx="16" cy="16" r=".8"/><circle cx="12" cy="12" r=".8"/></>,
  wave: <><path d="M2 8q2.5-4 5 0t5 0t5 0t5 0M2 16q2.5-4 5 0t5 0t5 0t5 0" /></>,
  mountain: <><path d="m2 20 8-15 5 9 3-5 5 11ZM7 10l3 3 3-3" /></>,
  layers: <><path d="m2 8 10-5 10 5-10 5ZM2 12l10 5 10-5M2 16l10 5 10-5" /></>,
  arrow: <><path d="M4 12h16m-6-6 6 6-6 6" /></>,
  check: <><path d="m5 12 4 4 10-10" /></>,
  reset: <><path d="M4 10a8 8 0 1 1 1 8M4 4v6h6" /></>,
  heart: <><path d="M12 20 4 12a5 5 0 0 1 8-6 5 5 0 0 1 8 6Z" /></>,
};
export function Icon({ name, className = '', style }: { name: IconName; className?: string; style?: CSSProperties }) {
  return <svg className={`icon ${className}`} style={style} viewBox="0 0 24 24" fill="none" stroke="currentColor"
    strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{paths[name]}</svg>;
}
export function IslandMark() {
  return <svg viewBox="0 0 64 64" aria-hidden="true"><path d="M10 43Q19 31 31 35T54 42Q48 54 30 51T10 43Z" fill="#ffe1a3"/><path d="m20 38 12-22 13 24Z" fill="#79b79a"/><path d="m27 25 5-9 5 9-5-3Z" fill="#fff9e6"/><path d="M11 56q4-3 8 0t8 0m13-4q4-3 8 0t8 0" stroke="#78aaa4" fill="none" strokeWidth="2" strokeLinecap="round"/></svg>;
}
