import { supportedLocales, type AppLocale } from '../i18n/config';
import { gamesData } from './game';

/** Locale-agnostic public paths (after `/:locale`), shared by sitemap and hreflang helpers. */
export function getPublicPathSuffixes(): string[] {
  const suffixes = ['', '/changelog'];

  for (const game of gamesData) {
    if (!game.isEnable) continue;
    suffixes.push(`/game/${game.key}`, `/game/${game.key}/rules`);
  }

  return suffixes;
}

export function getPublicLocalePaths(): Array<{ locale: AppLocale; suffix: string }> {
  return supportedLocales.flatMap(({ code }) =>
    getPublicPathSuffixes().map(suffix => ({ locale: code, suffix })),
  );
}
