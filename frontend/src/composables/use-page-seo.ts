import type { GameKey } from '@/libs/game';
import { resolveLocaleAlternateUrls } from '@/libs/locale-alternate-urls';
import { normalizeSiteOrigin } from '@/libs/site-origin';
import type { Composer } from 'vue-i18n';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import { computed } from 'vue';
import { useHead } from '@unhead/vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

type SeoPayload = {
  title: string;
  description: string;
};

const PLAYABLE_GAMES = ['tictactoe', 'conn4'] as const satisfies readonly GameKey[];

function isPlayableGameKey(key: string): key is (typeof PLAYABLE_GAMES)[number] {
  return (PLAYABLE_GAMES as readonly string[]).includes(key);
}

function localizedGameName(t: Composer['t'], gameKey: string): string {
  const key = `games.${gameKey}`;
  const translated = t(key);
  return translated === key ? gameKey : translated;
}

export function resolvePageSeo(
  route: RouteLocationNormalizedLoaded,
  t: Composer['t'],
): SeoPayload {
  const routeName = route.name;

  if (routeName === 'home') {
    return {
      title: t('seo.home.title'),
      description: t('seo.home.description'),
    };
  }

  if (routeName === 'changelog') {
    return {
      title: t('seo.changelog.title'),
      description: t('seo.changelog.description'),
    };
  }

  if (routeName === 'about') {
    return {
      title: t('seo.about.title'),
      description: t('seo.about.description'),
    };
  }

  if (routeName === 'health') {
    return {
      title: t('seo.health.title'),
      description: t('seo.health.description'),
    };
  }

  if (routeName === 'not-found' || routeName === 'not-found-catchall') {
    return {
      title: t('seo.notFound.title'),
      description: t('seo.notFound.description'),
    };
  }

  const gameKey = route.params.game;
  if (typeof gameKey === 'string') {
    if (routeName === 'game-play') {
      if (isPlayableGameKey(gameKey)) {
        return {
          title: t(`seo.gamePlay.${gameKey}.title`),
          description: t(`seo.gamePlay.${gameKey}.description`),
        };
      }
      const game = localizedGameName(t, gameKey);
      return {
        title: t('seo.gamePlay.fallback.title', { game }),
        description: t('seo.gamePlay.fallback.description', { game }),
      };
    }

    if (routeName === 'game-rules') {
      if (isPlayableGameKey(gameKey)) {
        return {
          title: t(`seo.gameRules.${gameKey}.title`),
          description: t(`seo.gameRules.${gameKey}.description`),
        };
      }
      const game = localizedGameName(t, gameKey);
      return {
        title: t('seo.gameRules.fallback.title', { game }),
        description: t('seo.gameRules.fallback.description', { game }),
      };
    }
  }

  return {
    title: t('seo.home.title'),
    description: t('seo.defaultDescription'),
  };
}

export function usePageSeo() {
  const route = useRoute();
  const { t } = useI18n();

  const siteUrl = normalizeSiteOrigin(import.meta.env.VITE_SITE_URL) ?? '';

  const pageSeo = computed(() => resolvePageSeo(route, t));
  const localeAlternates = computed(() =>
    siteUrl ? resolveLocaleAlternateUrls(route, siteUrl) : null,
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

      const canonical = localeAlternates.value?.canonical;
      if (canonical) {
        meta.push({ property: 'og:url', content: canonical });
      }

      return meta;
    },
    link: () => {
      const urls = localeAlternates.value;
      if (!urls) return [];

      return [
        { rel: 'canonical', href: urls.canonical },
        { rel: 'alternate', hreflang: 'en', href: urls.alternates.en },
        { rel: 'alternate', hreflang: 'fa', href: urls.alternates.fa },
        {
          rel: 'alternate',
          hreflang: 'x-default',
          href: urls.xDefault,
        },
      ];
    },
  });
}
