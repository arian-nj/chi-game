/** Absolute https origin for canonical URLs, sitemap, and hreflang. */
export function normalizeSiteOrigin(raw: string | undefined): string | undefined {
  if (!raw) return undefined;

  const trimmed = raw.trim().replace(/\/$/, '');
  if (!trimmed) return undefined;

  if (/^https:\/\//i.test(trimmed)) return trimmed;
  if (/^http:\/\//i.test(trimmed)) return trimmed.replace(/^http:\/\//i, 'https://');

  return `https://${trimmed}`;
}
