export const supportedLocales = [
  { code: 'en', labelKey: 'locale.en' },
] as const;

export type AppLocale = (typeof supportedLocales)[number]['code'];
