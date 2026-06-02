import { useI18n } from 'vue-i18n';
import { persistLocale, supportedLocales, type AppLocale } from '@/i18n';

export function useAppLocale() {
  const { locale, t } = useI18n();

  function setLocale(code: AppLocale) {
    locale.value = code;
    persistLocale(code);
  }

  return {
    locale,
    setLocale,
    supportedLocales,
    t,
  };
}
