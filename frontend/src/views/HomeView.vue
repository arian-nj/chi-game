<script setup lang="ts">
import { gamesData } from '../libs/game';
import { onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { useRoute } from 'vue-router';
import { useGuestProfile } from '@/composables/use-guest-profile';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { displayUsername } = useGuestProfile();
const route = useRoute();
const locale = route.params.locale as string;

onMounted(() => {
  import('../views/GameView.vue');
  import('@/views/RoomShell.vue');
});

</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center p-6 pt-16">
    <h1 class="text-4xl font-bold text-white mb-2 mt-2 animate-pop select-none drop-shadow-sm">
      {{ t('app.title') }}
    </h1>
    <p
      v-if="displayUsername"
      class="mb-8 text-lg font-medium text-blue-100 select-none"
    >
      {{ displayUsername }}
    </p>
    <div v-else class="mb-8" aria-hidden="true" ></div>

    <div
      class="grid w-full max-w-5xl gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4 place-items-center mx-auto px-4"
    >
      <component
        :is="game.isEnable ? RouterLink : 'div'"
        v-for="game in gamesData"
        :key="game.key"
        v-bind="game.isEnable ? { to: `/${locale}/game/${game.key}` } : {}"
        class="rounded-2xl w-full transition duration-100 shadow outline-none focus:ring-2 focus:ring-blue-500 aspect-4/3 overflow-hidden flex flex-col items-center justify-center"
        :class="
          game.isEnable
            ? 'bg-custom-lite-blue hover:bg-custom-deep-blue cursor-pointer'
            : 'bg-custom-deep-blue/40 cursor-not-allowed opacity-55 saturate-0'
        "
        :aria-disabled="game.isEnable ? undefined : 'true'"
      >
        <div class="relative flex flex-col items-center justify-center w-full h-full p-4">
          <span class="text-lg md:text-xl font-semibold text-white select-none text-center">
            {{ t(`games.${game.key}`) }}
          </span>
        </div>
      </component>
    </div>

    <nav class="mt-8 flex flex-wrap items-center justify-center gap-x-6 gap-y-2" aria-label="Site">
      <RouterLink
        :to="`/${locale}/room`"
        class="text-blue-100 hover:text-white underline underline-offset-4 font-semibold"
      >
        {{ t('nav.room') }}
      </RouterLink>
      <RouterLink
        :to="`/${locale}/about`"
        class="text-blue-100 hover:text-white underline underline-offset-4 font-semibold"
      >
        {{ t('nav.about') }}
      </RouterLink>
      <RouterLink
        :to="`/${locale}/changelog`"
        class="text-blue-100 hover:text-white underline underline-offset-4 font-semibold"
      >
        {{ t('nav.changelog') }}
      </RouterLink>
    </nav>
  </div>
</template>
