<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { useI18n } from 'vue-i18n';

type ChangelogEntry = {
  date: string;
  version?: string;
  new?: string[];
  changed?: string[];
  fixed?: string[];
};

const { t } = useI18n();
const { textDir } = useTextDirection();

const entries: ChangelogEntry[] = [
  {
    date: '2026-06-08',
    new: [
      'now we are offline :)',
    ],
  },
  {
    date: '2026-06-07',
    new: [
      'chat ui added.',
      'room member joined and left events added.',
    ],
  },
  {
    date: '2026-06-06',
    new: [
      'merge invite and room.',
      'chat base added.'
    ],
  },
  {
    date: '2026-06-05',
    new: [
      'basic room added.',
      'room lobby page added.'
    ],
  },
  {
    date: '2026-06-04',
    new: [
      'merger game and rules ',
    ],
  },
  {
    date: '2026-06-03',
    new: [
      'We are online 🎉 (now we have backend)',
      'Added guest users.',
      'Added about page.',
    ],
  },
  {
    date: '2026-06-02',
    new: [
      'Changelog page.',
      'Each game has its own rules page.',
      'Added lang and text side in document.',
      'Added robots.txt.',
      'Added sitemap.xml.',
    ],
    changed: ['Offline Tic Tac Toe: medium bot is a little easier (occasionally misses blocks / center).'],
  },
];
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center p-6 pt-16">
    <section
      :dir="textDir"
      class="w-full max-w-3xl mx-auto p-6 bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md"
    >
      <div class="flex items-center justify-between gap-4 mb-6">
        <h2
          class="text-3xl font-extrabold tracking-wide bg-linear-to-r from-white to-blue-200 bg-clip-text text-transparent"
        >
          {{ t('changelog.title') }}
        </h2>

        <RouterLink
          :to="{ name: 'home', params: { locale: $route.params.locale } }"
          class="text-sm font-semibold text-blue-100 hover:text-white underline underline-offset-4"
        >
          {{ t('nav.backToHome') }}
        </RouterLink>
      </div>

      <div class="space-y-6 text-blue-100 leading-relaxed">
        <article
          v-for="entry in entries"
          :key="`${entry.date}:${entry.version ?? ''}`"
          class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4"
        >
          <header class="flex flex-wrap items-baseline justify-between gap-2 mb-3">
            <div class="font-bold text-white">
              {{ entry.date }}
              <span v-if="entry.version" class="opacity-90"> · v{{ entry.version }}</span>
            </div>
          </header>

          <div class="space-y-4">
            <div v-if="entry.new?.length">
              <div class="font-bold text-white mb-2">{{ t('changelog.sections.new') }}</div>
              <ul class="list-disc ps-5 space-y-2">
                <li v-for="(item, idx) in entry.new" :key="`new:${idx}`">{{ item }}</li>
              </ul>
            </div>

            <div v-if="entry.changed?.length">
              <div class="font-bold text-white mb-2">{{ t('changelog.sections.changed') }}</div>
              <ul class="list-disc ps-5 space-y-2">
                <li v-for="(item, idx) in entry.changed" :key="`changed:${idx}`">{{ item }}</li>
              </ul>
            </div>

            <div v-if="entry.fixed?.length">
              <div class="font-bold text-white mb-2">{{ t('changelog.sections.fixed') }}</div>
              <ul class="list-disc ps-5 space-y-2">
                <li v-for="(item, idx) in entry.fixed" :key="`fixed:${idx}`">{{ item }}</li>
              </ul>
            </div>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>
