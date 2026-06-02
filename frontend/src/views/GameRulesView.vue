<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { textDir } = useTextDirection();
const route = useRoute();

type RulesGameKey = 'tictactoe' | 'conn4';

const rulesKey = computed<RulesGameKey>(() =>
  route.params.game === 'conn4' ? 'conn4' : 'tictactoe',
);

function rulesPath(suffix: string): string {
  return `rules.${rulesKey.value}.${suffix}`;
}
</script>

<template>
  <div>
    <section
      :dir="textDir"
      class="max-w-2xl mx-auto p-6 bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md"
    >
      <h2
        class="text-3xl font-extrabold mb-6 tracking-wide text-center bg-linear-to-r from-white to-blue-200 bg-clip-text text-transparent"
      >
        {{ t(rulesPath('title')) }}
      </h2>

      <div class="space-y-6 text-blue-100 leading-relaxed">
        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            🎯 {{ t(rulesPath('objective')) }}
          </h3>
          <i18n-t :keypath="rulesPath('objectiveText')" tag="p">
            <template #highlight>
              <strong>{{ t(rulesPath('objectiveHighlight')) }}</strong>
            </template>
          </i18n-t>
        </div>

        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            📜 {{ t(rulesPath('gameplay')) }}
          </h3>
          <ul class="list-disc ps-5 space-y-2">
            <i18n-t :keypath="rulesPath('rule1')" tag="li">
              <template #highlight>
                <strong>{{ t(rulesPath('rule1Highlight')) }}</strong>
              </template>
            </i18n-t>

            <i18n-t v-if="rulesKey === 'tictactoe'" keypath="rules.tictactoe.rule2" tag="li">
              <template #x>
                <strong>X</strong>
              </template>
              <template #o>
                <strong>O</strong>
              </template>
            </i18n-t>
            <i18n-t v-else keypath="rules.conn4.rule2" tag="li">
              <template #red>
                <strong>{{ t('game.red') }}</strong>
              </template>
              <template #yellow>
                <strong>{{ t('game.yellow') }}</strong>
              </template>
            </i18n-t>

            <li>{{ t(rulesPath('rule3')) }}</li>
            <li>{{ t(rulesPath('rule4')) }}</li>
          </ul>
        </div>

        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            🏆 {{ t(rulesPath('winning')) }}
          </h3>
          <ul class="list-disc ps-5 space-y-2">
            <li>
              <i18n-t :keypath="rulesPath('victory')" tag="span">
                <template #label>
                  <strong>{{ t(rulesPath('victoryLabel')) }}</strong>
                </template>
              </i18n-t>
            </li>
            <li>
              <i18n-t :keypath="rulesPath('drawRule')" tag="span">
                <template #label>
                  <strong>{{ t(rulesPath('drawLabel')) }}</strong>
                </template>
              </i18n-t>
            </li>
          </ul>
        </div>

        <blockquote class="border-s-4 border-green-400 bg-green-500/10 p-4 rounded-e-xl mt-4">
          <span class="font-bold text-white block mb-1">💡 {{ t(rulesPath('hintTitle')) }}</span>
          {{ t(rulesPath('hintText')) }}
        </blockquote>
      </div>
    </section>
  </div>
</template>
