import { computed } from 'vue';
import { useRoute } from 'vue-router';
import type { AppLocale } from '@/i18n';
import { getInitialLocale } from '@/i18n';

export function useLocalePath() {
  const route = useRoute();

  const locale = computed(
    () => (route.params.locale as AppLocale | undefined) ?? getInitialLocale(),
  );

  function localePath(path: string): string {
    const loc = locale.value;
    if (path === '/' || path === '') {
      return `/${loc}`;
    }
    const sub = path.startsWith('/') ? path : `/${path}`;
    return `/${loc}${sub}`;
  }

  return { locale, localePath };
}
