import { createI18n } from 'vue-i18n';
import en from './locales/en';
import fa from './locales/fa';

export const LOCALE_STORAGE_KEY = 'chigame-locale';

export const supportedLocales = [
  { code: 'en', labelKey: 'locale.en' },
  { code: 'fa', labelKey: 'locale.fa' },
] as const;

export type AppLocale = (typeof supportedLocales)[number]['code'];

export const rtlLocales = ['fa'] as const;

export function isRtlLocale(locale: string): boolean {
  return (rtlLocales as readonly string[]).includes(locale);
}

export function getInitialLocale(): AppLocale {
  const stored = localStorage.getItem(LOCALE_STORAGE_KEY);
  if (stored === 'en' || stored === 'fa') {
    return stored;
  }

  const browserLang = navigator.language.split('-')[0];
  if (browserLang === 'fa') {
    return 'fa';
  }

  return 'en';
}

export const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'en',
  messages: { en, fa },
});

export function documentDirection(locale: AppLocale): 'ltr' | 'rtl' {
  return isRtlLocale(locale) ? 'rtl' : 'ltr';
}

export function setDocumentLocale(locale: AppLocale) {
  document.documentElement.lang = locale;
  document.documentElement.dir = documentDirection(locale);
}

export function persistLocale(locale: AppLocale) {
  localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  setDocumentLocale(locale);
}

setDocumentLocale(getInitialLocale());
