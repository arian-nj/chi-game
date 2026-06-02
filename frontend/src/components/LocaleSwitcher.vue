<script setup lang="ts">
import { useAppLocale } from '@/composables/use-app-locale';
import type { AppLocale } from '@/i18n';

const { locale, setLocale, supportedLocales, t } = useAppLocale();
</script>

<template>
  <div
    dir="ltr"
    class="flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-2 py-1 shadow-md"
    role="group"
    :aria-label="t('locale.label')"
  >
    <button
      v-for="entry in supportedLocales"
      :key="entry.code"
      type="button"
      class="rounded-lg px-3 py-1.5 text-sm font-bold transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :class="
        locale === entry.code
          ? 'bg-white/90 text-custom-blue shadow'
          : 'text-blue-100 hover:bg-white/20 hover:text-white'
      "
      :aria-pressed="locale === entry.code"
      @click="setLocale(entry.code as AppLocale)"
    >
      {{ t(entry.labelKey) }}
    </button>
  </div>
</template>
