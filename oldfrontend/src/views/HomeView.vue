<script setup lang="ts">
import { gamesData } from '../libs/game';
import { onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import LocaleSwitcher from '@/components/LocaleSwitcher.vue';

const { t } = useI18n();
const route = useRoute();
const locale = route.params.locale as string;

onMounted(() => {
  import('../views/GameView.vue');
});
</script>

<template>
  <div class="home-root">

    <div class="loc-switcher-row">
      <LocaleSwitcher />
    </div>
    
    <h1 class="home-title">
      {{ t('app.title') }}
    </h1>

    <div class="games-grid">
      <component
        :is="game.isEnable ? RouterLink : 'div'"
        v-for="game in gamesData"
        :key="game.key"
        v-bind="game.isEnable ? { to: `/${locale}/game/${game.key}` } : {}"
        :class="[
          'game-card',
          game.isEnable ? 'game-card--enabled' : 'game-card--disabled'
        ]"
        :aria-disabled="game.isEnable ? undefined : 'true'"
        tabindex="0"
      >
        <div class="game-card-content">
          <span class="game-card-title">
            {{ t(`games.${game.key}`) }}
          </span>
        </div>
      </component>
    </div>

    <nav class="nav-links" aria-label="Site">
      <RouterLink
        :to="`/${locale}/about`"
        class="nav-link"
      >
        {{ t('nav.about') }}
      </RouterLink>
      <RouterLink
        :to="`/${locale}/changelog`"
        class="nav-link"
      >
        {{ t('nav.changelog') }}
      </RouterLink>
    </nav>
  </div>
</template>

<style scoped>
.home-root {
  background: var(--color-custom-blue);
  min-height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1.5rem;
  padding-top: 4rem;
  box-sizing: border-box;
}
.loc-switcher-row {
  width: 100%;
  display: flex;
  justify-content: flex-start;
}
.home-title {
  font-size: 2.5rem;
  font-weight: bold;
  color: #fff;
  margin-bottom: 0.5rem;
  margin-top: 0.5rem;
  animation: pop 0.18s cubic-bezier(0.1, 1.18, 1, 1.01) both;
  user-select: none;
  filter: drop-shadow(0 1px 3px rgba(16,30,87,0.2));
}
@keyframes pop {
  0% {
    transform: scale(0.97);
    opacity: 0.5;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}
.games-grid {
  display: grid;
  width: 100%;
  max-width: 80rem;
  gap: 1.5rem;
  grid-template-columns: repeat(2, 1fr);
  place-items: center;
  margin-left: auto;
  margin-right: auto;
  padding-left: 1rem;
  padding-right: 1rem;
}
@media (min-width: 768px) {
  .games-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (min-width: 1024px) {
  .games-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
.game-card {
  border-radius: 1rem;
  width: 100%;
  transition: background 100ms, filter 100ms, box-shadow 100ms;
  box-shadow: 0 2px 8px 0 rgba(16,30,87,0.12);
  outline: none;
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.game-card:focus {
  box-shadow: 0 0 0 3px #3b82f6, 0 2px 8px 0 rgba(16,30,87,0.12);
  z-index: 1;
}
.game-card--enabled {
  background: var(--color-custom-lite-blue);
  cursor: pointer;
}
.game-card--enabled:hover {
  background: var(--color-custom-deep-blue);
}
.game-card--disabled {
  background: rgba(37, 52, 84, 0.4);
  cursor: not-allowed;
  opacity: 0.55;
  filter: saturate(0);
}
.game-card-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 1rem;
  box-sizing: border-box;
}
.game-card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #fff;
  user-select: none;
  text-align: center;
}
@media (min-width: 768px) {
  .game-card-title {
    font-size: 1.25rem;
  }
}
.nav-links {
  margin-top: 2rem;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 0.5rem 1.5rem;
}
.nav-link {
  color: #c7d5ef;
  font-weight: 600;
  text-decoration: underline;
  text-underline-offset: 4px;
  transition: color 120ms;
}
.nav-link:hover {
  color: #fff;
}
</style>