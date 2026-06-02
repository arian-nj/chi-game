import { supportedLocales, type AppLocale } from '@/i18n';
import type { RouteLocationNormalizedLoaded } from 'vue-router';

export type LocaleAlternateUrls = {
  canonical: string;
  alternates: Record<AppLocale, string>;
  xDefault: string;
};

const DEFAULT_LOCALE: AppLocale = 'en';

/** Path after locale segment, e.g. `` or `/changelog` or `/game/tictactoe/rules`. */
export function resolveLocalePathSuffix(
  route: RouteLocationNormalizedLoaded,
): string | null {
  const routeName = route.name;

  if (routeName === 'home') return '';
  if (routeName === 'changelog') return '/changelog';
  if (routeName === 'not-found') return '/404';

  const gameKey = route.params.game;
  if (typeof gameKey !== 'string') return null;

  if (routeName === 'game-play') return `/game/${gameKey}`;
  if (routeName === 'game-rules') return `/game/${gameKey}/rules`;

  return null;
}

function withBase(path: string): string {
  const base = import.meta.env.BASE_URL;
  if (!base || base === '/') return path;
  const prefix = base.endsWith('/') ? base.slice(0, -1) : base;
  return `${prefix}${path}`;
}

function absoluteUrl(siteOrigin: string, path: string): string {
  return `${siteOrigin}${withBase(path)}`;
}

function currentLocaleFromRoute(route: RouteLocationNormalizedLoaded): AppLocale {
  const locale = route.params.locale;
  if (locale === 'en' || locale === 'fa') return locale;
  return DEFAULT_LOCALE;
}

/** Absolute canonical + hreflang URLs for public locale-paired pages, or null if not applicable. */
export function resolveLocaleAlternateUrls(
  route: RouteLocationNormalizedLoaded,
  siteOrigin: string,
): LocaleAlternateUrls | null {
  const suffix = resolveLocalePathSuffix(route);
  if (suffix === null) return null;

  const alternates = Object.fromEntries(
    supportedLocales.map(({ code }) => [
      code,
      absoluteUrl(siteOrigin, `/${code}${suffix}`),
    ]),
  ) as Record<AppLocale, string>;

  const currentLocale = currentLocaleFromRoute(route);

  return {
    canonical: alternates[currentLocale],
    alternates,
    xDefault: alternates[DEFAULT_LOCALE],
  };
}
