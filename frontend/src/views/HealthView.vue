<script setup lang="ts">
import { useBackendHealth } from '@/composables/use-backend-health';
import { useTextDirection } from '@/composables/use-text-direction';
import { getApiBaseUrl } from '@/libs/api-base-url';
import { HealthType } from '@/gen/healthcheck/v1/healthcheck_pb';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { textDir } = useTextDirection();

const {
  data,
  isLoading,
  isFetching,
  isBackendHealthy,
  backendError,
  refetch,
  dataUpdatedAt,
  failureCount,
} = useBackendHealth();

const apiBaseLabel = computed(() => getApiBaseUrl() || t('health.apiSameOrigin'));

const statusLabel = computed(() => {
  if (isLoading.value) return t('health.status.checking');
  if (backendError.value) return t('health.status.error');
  if (isBackendHealthy.value) return t('health.status.online');
  return t('health.status.degraded');
});

const healthTypeLabel = computed(() => {
  const type = data.value?.healthType;
  if (type === HealthType.OK) return t('health.healthType.ok');
  if (type === HealthType.UNSPECIFIED) return t('health.healthType.unspecified');
  return t('health.healthType.unknown');
});

const lastCheckedLabel = computed(() => {
  if (!dataUpdatedAt.value) return '—';
  return new Date(dataUpdatedAt.value).toLocaleString();
});
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
          {{ t('health.title') }}
        </h1>

        <RouterLink
          :to="{ name: 'home', params: { locale: $route.params.locale } }"
          class="text-sm font-semibold text-blue-100 hover:text-white underline underline-offset-4 shrink-0"
        >
          {{ t('nav.backToHome') }}
        </RouterLink>
      </div>

      <p class="text-blue-100 mb-6 leading-relaxed">
        {{ t('health.description') }}
      </p>

      <dl class="space-y-4 text-blue-100">
        <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
          <dt class="text-sm font-semibold text-white/80 mb-1">{{ t('health.fields.status') }}</dt>
          <dd class="text-lg font-bold text-white">{{ statusLabel }}</dd>
        </div>

        <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
          <dt class="text-sm font-semibold text-white/80 mb-1">{{ t('health.fields.healthType') }}</dt>
          <dd class="text-lg font-medium text-white">{{ healthTypeLabel }}</dd>
        </div>

        <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
          <dt class="text-sm font-semibold text-white/80 mb-1">{{ t('health.fields.apiBase') }}</dt>
          <dd class="font-mono text-sm text-white break-all">{{ apiBaseLabel }}</dd>
        </div>

        <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
          <dt class="text-sm font-semibold text-white/80 mb-1">{{ t('health.fields.lastChecked') }}</dt>
          <dd class="text-white">{{ lastCheckedLabel }}</dd>
        </div>

        <div
          v-if="backendError"
          class="rounded-xl bg-rose-950/40 border border-rose-400/30 p-4"
        >
          <dt class="text-sm font-semibold text-rose-200 mb-1">{{ t('health.fields.error') }}</dt>
          <dd class="font-mono text-sm text-rose-100 break-all">{{ backendError }}</dd>
        </div>

        <div
          v-if="failureCount > 0 && !backendError"
          class="rounded-xl bg-amber-950/30 border border-amber-400/20 p-4 text-amber-100 text-sm"
        >
          {{ t('health.failures', { count: failureCount }) }}
        </div>
      </dl>

      <button
        type="button"
        class="mt-6 w-full rounded-xl bg-white/10 hover:bg-white/15 border border-white/15 px-4 py-3 text-white font-semibold transition disabled:opacity-50"
        :disabled="isFetching"
        @click="() => refetch()"
      >
        {{ isFetching ? t('health.refreshing') : t('health.refresh') }}
      </button>
    </section>
  </div>
</template>
