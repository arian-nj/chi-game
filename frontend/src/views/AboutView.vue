<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { textDir } = useTextDirection();

const sectionKeys = ['mission', 'games', 'roadmap'] as const;
</script>

<template>
  <div class="about-root">
    <section
      :dir="textDir"
      class="about-section"
    >
      <div class="about-header-row">
        <h1 class="about-title">
          {{ t('about.title') }}
        </h1>

        <RouterLink
          :to="{ name: 'home', params: { locale: $route.params.locale } }"
          class="about-back-link"
        >
          {{ t('nav.backToHome') }}
        </RouterLink>
      </div>

      <p class="about-tagline">
        {{ t('about.tagline') }}
      </p>

      <div class="about-article-list">
        <article
          v-for="key in sectionKeys"
          :key="key"
          class="about-article"
        >
          <h2 class="about-article-title">
            {{ t(`about.sections.${key}.title`) }}
          </h2>
          <p>{{ t(`about.sections.${key}.body`) }}</p>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.about-root {
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
.about-section {
  width: 100%;
  max-width: 48rem;
  margin-left: auto;
  margin-right: auto;
  padding: 1.5rem;
  background: color-mix(in srgb, var(--color-custom-lite-blue) 40%, transparent);
  border-radius: 1rem;
  border: 1px solid rgba(255,255,255,0.10);
  box-shadow: 0 2px 8px 0 rgba(16,30,87,0.12);
}
.about-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.about-title {
  font-size: 1.875rem;
  font-weight: 800;
  letter-spacing: 0.04em;
  background: linear-gradient(to right, #fff, #bcd7ff);
  background-clip: text;
  -webkit-background-clip: text;
  color: transparent;
}
.about-back-link {
  font-size: 0.875rem;
  font-weight: 600;
  color: #bcd7ff;
  text-decoration: underline;
  text-underline-offset: 4px;
  transition: color 0.15s;
  flex-shrink: 0;
}
.about-back-link:hover {
  color: #fff;
}
.about-tagline {
  font-size: 1.125rem;
  color: #fff;
  font-weight: 500;
  margin-bottom: 2rem;
  line-height: 1.6;
}
.about-article-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  color: #bcd7ff;
  line-height: 1.6;
}
.about-article {
  border-radius: 0.75rem;
  background: color-mix(in srgb, var(--color-custom-deep-blue) 30%, transparent);
  border: 1px solid rgba(255,255,255,0.10);
  padding: 1rem;
}
.about-article-title {
  font-weight: bold;
  color: #fff;
  font-size: 1.125rem;
  margin-bottom: 0.5rem;
}
</style>