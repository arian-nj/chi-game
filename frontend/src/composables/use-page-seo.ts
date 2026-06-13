import { useHead } from '@unhead/vue';
import type { Composer } from 'vue-i18n';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

type SeoPayload = {
  title: string;
  description: string;
};

function normalizeSiteOrigin(raw: string | undefined): string | undefined {
  if (!raw) return undefined;

  const trimmed = raw.trim().replace(/\/$/, '');
  if (!trimmed) return undefined;

  if (/^https:\/\//i.test(trimmed)) return trimmed;
  if (/^http:\/\//i.test(trimmed)) return trimmed.replace(/^http:\/\//i, 'https://');

  return `https://${trimmed}`;
}

function resolvePathSuffix(route: RouteLocationNormalizedLoaded): string | null {
  if (route.name === 'home') return '';
  const gamePath = route.path.match(/\/game\/[\w-]+/);
  return gamePath?.[0] ?? null;
}

function resolveCanonicalUrl(
  route: RouteLocationNormalizedLoaded,
  siteOrigin: string,
): string | null {
  const suffix = resolvePathSuffix(route);
  if (suffix === null) return null;

  return `${siteOrigin}/en${suffix}`;
}

export function resolvePageSeo(
  route: RouteLocationNormalizedLoaded,
  t: Composer['t'],
  te: Composer['te'],
): SeoPayload {
  const routeName = route.name;

  if (routeName === 'home') {
    return {
      title: t('seo.home.title'),
      description: t('seo.home.description'),
    };
  }

  if (typeof routeName === 'string' && te(`seo.${routeName}.title`)) {
    return {
      title: t(`seo.${routeName}.title`),
      description: t(`seo.${routeName}.description`),
    };
  }

  return {
    title: t('seo.home.title'),
    description: t('seo.defaultDescription'),
  };
}

export function usePageSeo() {
  const route = useRoute();
  const { t, te } = useI18n();

  const siteOrigin = normalizeSiteOrigin(import.meta.env.VITE_SITE_URL);

  const pageSeo = computed(() => resolvePageSeo(route, t, te));
  const canonicalUrl = computed(() =>
    siteOrigin ? resolveCanonicalUrl(route, siteOrigin) : null,
  );

  useHead({
    title: () => pageSeo.value.title,
    meta: () => {
      const { title, description } = pageSeo.value;
      const meta = [
        { name: 'description', content: description },
        { property: 'og:title', content: title },
        { property: 'og:description', content: description },
        { property: 'og:type', content: 'website' },
        { name: 'twitter:card', content: 'summary' },
        { name: 'twitter:title', content: title },
        { name: 'twitter:description', content: description },
      ];

      const canonical = canonicalUrl.value;
      if (canonical) {
        meta.push({ property: 'og:url', content: canonical });
      }

      return meta;
    },
    link: () => {
      const canonical = canonicalUrl.value;
      if (!canonical) return [];

      return [{ rel: 'canonical', href: canonical }];
    },
  });
}
