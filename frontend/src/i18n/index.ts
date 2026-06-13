import { createI18n } from 'vue-i18n';
import en from './locales/en';
import { type AppLocale } from './config';

export type { AppLocale } from './config';
export { supportedLocales } from './config';

export const LOCALE_STORAGE_KEY = 'chigame-locale';

export function getInitialLocale(): AppLocale {
  return 'en';
}

export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en },
});

export function setDocumentLocale(locale: AppLocale) {
  document.documentElement.lang = locale;
  document.documentElement.dir = 'ltr';
}

export function persistLocale(locale: AppLocale) {
  localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  setDocumentLocale(locale);
}

setDocumentLocale('en');
