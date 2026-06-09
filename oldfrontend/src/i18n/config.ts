export const supportedLocales = [
  { code: 'en', labelKey: 'locale.en' },
  { code: 'fa', labelKey: 'locale.fa' },
] as const;

export type AppLocale = (typeof supportedLocales)[number]['code'];

export const rtlLocales = ['fa'] as const;

export function isRtlLocale(locale: string): boolean {
  return (rtlLocales as readonly string[]).includes(locale);
}
