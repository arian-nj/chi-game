import { useI18n } from 'vue-i18n';
import { persistLocale, supportedLocales, type AppLocale } from '@/i18n';
import { useRoute, useRouter } from 'vue-router';

export function useAppLocale() {
  const { locale, t } = useI18n();
  const route = useRoute();
  const router = useRouter();

  function setLocale(code: AppLocale) {
    locale.value = code;
    persistLocale(code);

    // Keep the URL in sync for SEO: "/<locale>/..."
    const currentLocale = route.params.locale;
    if (typeof currentLocale === 'string') {
      if (route.name) {
        router.push({
          name: route.name,
          params: { ...route.params, locale: code },
          query: route.query,
          hash: route.hash,
        });
      } else {
        const rest = route.fullPath.replace(/^\/(en|fa)/, '');
        router.push(`/${code}${rest}`);
      }
      return;
    }

    router.push(`/${code}${route.fullPath}`);
  }

  return {
    locale,
    setLocale,
    supportedLocales,
    t,
  };
}
