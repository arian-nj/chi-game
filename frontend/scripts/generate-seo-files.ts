import { writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { loadEnv } from 'vite';
import { supportedLocales } from '../src/i18n/config.ts';
import { getPublicLocalePaths } from '../src/libs/public-path-suffixes.ts';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const publicDir = resolve(root, 'public');

function escapeXml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;');
}

function buildRobots(sitemapUrl: string | undefined): string {
  const lines = ['User-agent: *', 'Allow: /', ''];
  if (sitemapUrl) {
    lines.push(`Sitemap: ${sitemapUrl}`);
  }
  return `${lines.join('\n')}\n`;
}

function buildSitemap(entries: Array<{ loc: string; alternates: Record<string, string> }>): string {
  const urlNodes = entries
    .map(({ loc, alternates }) => {
      const alternateLinks = supportedLocales
        .map(
          ({ code }) =>
            `    <xhtml:link rel="alternate" hreflang="${code}" href="${escapeXml(alternates[code])}" />`,
        )
        .concat(
          `    <xhtml:link rel="alternate" hreflang="x-default" href="${escapeXml(alternates.en)}" />`,
        )
        .join('\n');

      return `  <url>
    <loc>${escapeXml(loc)}</loc>
${alternateLinks}
  </url>`;
    })
    .join('\n');

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset
  xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml"
>
${urlNodes}
</urlset>
`;
}

function main() {
  const env = loadEnv(process.env.MODE ?? 'production', root, '');
  const siteOrigin = env.VITE_SITE_URL?.replace(/\/$/, '');
  const base = env.BASE_URL?.replace(/\/$/, '') ?? '';
  const prefix = base && base !== '/' ? base : '';

  const localePaths = getPublicLocalePaths();

  const sitemapUrl = siteOrigin ? `${siteOrigin}${prefix}/sitemap.xml` : undefined;
  writeFileSync(resolve(publicDir, 'robots.txt'), buildRobots(sitemapUrl));

  if (!siteOrigin) {
    console.warn(
      '[generate-seo] VITE_SITE_URL is not set — wrote robots.txt without Sitemap; skipped sitemap.xml',
    );
    return;
  }

  const bySuffix = new Map<string, Record<string, string>>();

  for (const { locale, suffix } of localePaths) {
    const path = `${prefix}/${locale}${suffix}`;
    const loc = `${siteOrigin}${path}`;
    const group = bySuffix.get(suffix) ?? {};
    group[locale] = loc;
    bySuffix.set(suffix, group);
  }

  const entries = [...bySuffix.entries()].map(([, alternates]) => ({
    loc: alternates.en,
    alternates,
  }));

  writeFileSync(resolve(publicDir, 'sitemap.xml'), buildSitemap(entries));
  console.log(`[generate-seo] Wrote robots.txt and sitemap.xml (${entries.length} URLs)`);
}

main();
