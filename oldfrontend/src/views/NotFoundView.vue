<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

const route = useRoute();
const { t } = useI18n();

const locale = computed(() => route.params.locale as string);
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex items-center justify-center p-6 text-white">
    <div class="w-full max-w-xl rounded-2xl border border-white/10 bg-custom-lite-blue/30 p-8 shadow-lg">
      <p class="text-center text-6xl font-extrabold tracking-tight">404</p>
      <h1 class="mt-4 text-center text-2xl font-extrabold">{{ t('notFound.title') }}</h1>
      <p class="mt-3 text-center text-blue-100/85">{{ t('notFound.description') }}</p>

      <div class="mt-8 flex items-center justify-center gap-3">
        <RouterLink
          :to="`/${locale}`"
          class="rounded-xl bg-white/90 px-5 py-2.5 text-sm font-extrabold text-custom-blue shadow transition hover:bg-white"
        >
          {{ t('notFound.backHome') }}
        </RouterLink>
      </div>
    </div>
  </div>
</template>

