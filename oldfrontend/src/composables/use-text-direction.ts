import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { isRtlLocale } from '@/i18n';

export function useTextDirection() {
  const { locale } = useI18n();

  const isRtl = computed(() => isRtlLocale(locale.value));
  const textDir = computed(() => (isRtl.value ? 'rtl' : 'ltr'));

  return { isRtl, textDir };
}
