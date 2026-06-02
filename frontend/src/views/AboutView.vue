<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { textDir } = useTextDirection();

const sectionKeys = ['mission', 'games', 'roadmap'] as const;
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center p-6 pt-16">
    <section
      :dir="textDir"
      class="w-full max-w-3xl mx-auto p-6 bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md"
    >
      <div class="flex items-center justify-between gap-4 mb-6">
        <h1
          class="text-3xl font-extrabold tracking-wide bg-linear-to-r from-white to-blue-200 bg-clip-text text-transparent"
        >
          {{ t('about.title') }}
        </h1>

        <RouterLink
          :to="{ name: 'home', params: { locale: $route.params.locale } }"
          class="text-sm font-semibold text-blue-100 hover:text-white underline underline-offset-4 shrink-0"
        >
          {{ t('nav.backToHome') }}
        </RouterLink>
      </div>

      <p class="text-lg text-white font-medium mb-8 leading-relaxed">
        {{ t('about.tagline') }}
      </p>

      <div class="space-y-6 text-blue-100 leading-relaxed">
        <article
          v-for="key in sectionKeys"
          :key="key"
          class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4"
        >
          <h2 class="font-bold text-white text-lg mb-2">
            {{ t(`about.sections.${key}.title`) }}
          </h2>
          <p>{{ t(`about.sections.${key}.body`) }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
