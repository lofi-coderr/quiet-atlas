import test from 'node:test';
import assert from 'node:assert/strict';
import { parseSeed, randomSeed, islandQuery, MAX_SEED } from './parameters.ts';

test('seed parsing preserves the safe integer bounds and rejects ambiguous inputs', () => {
  for (const text of ['', '-1', '2.5', '1e3', 'NaN', 'Infinity', '9007199254740992']) assert.equal(parseSeed(text), null, text);
  assert.equal(parseSeed('0'), 0); assert.equal(parseSeed(String(MAX_SEED)), MAX_SEED);
  assert.equal(parseSeed(' 002408 '), 2408);
});
test('random seeds are safe, variable integers', () => {
  const seeds = new Set(Array.from({ length: 100 }, randomSeed));
  assert.ok(seeds.size > 90);
  for (const seed of seeds) assert.ok(Number.isSafeInteger(seed) && seed >= 0 && seed <= MAX_SEED);
});
test('API queries preserve every setting including false and a zero seed', () => {
  const query = new URLSearchParams(islandQuery({ seed: 0, seaLevel: 0.7, roughness: 0, contours: false }));
  assert.equal(query.get('seed'), '0'); assert.equal(query.get('contours'), 'false');
  assert.equal(query.get('seaLevel'), '0.7'); assert.equal(query.get('roughness'), '0');
});
