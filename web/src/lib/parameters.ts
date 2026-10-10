export const MAX_SEED = Number.MAX_SAFE_INTEGER;
export const DEFAULTS: Parameters = { seed: 2408, seaLevel: 0.45, roughness: 0.55, contours: true };

export interface Parameters {
  seed: number;
  seaLevel: number;
  roughness: number;
  contours: boolean;
}
export interface Island {
  parameters: Parameters;
  svg: string;
  landPercent: number;
  name: string;
}

export function parseSeed(text: string): number | null {
  if (!/^\d+$/.test(text.trim())) return null;
  const seed = Number(text.trim());
  return Number.isSafeInteger(seed) && seed >= 0 ? seed : null;
}

// 21 + 32 random bits cover the entire JavaScript-safe, nonnegative seed range.
export function randomSeed(): number {
  const bits = crypto.getRandomValues(new Uint32Array(2));
  return (bits[0] & 0x1fffff) * 0x100000000 + bits[1];
}

export function islandQuery(parameters: Parameters): string {
  return new URLSearchParams({ seed: String(parameters.seed), seaLevel: String(parameters.seaLevel),
    roughness: String(parameters.roughness), contours: String(parameters.contours) }).toString();
}

export function saveSVG(island: Island): void {
  const url = URL.createObjectURL(new Blob([island.svg], { type: 'image/svg+xml;charset=utf-8' }));
  const link = document.createElement('a');
  link.href = url;
  link.download = `quiet-atlas-${island.parameters.seed}.svg`;
  document.body.append(link);
  link.click();
  link.remove();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
