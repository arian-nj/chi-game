<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { textDir } = useTextDirection();
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
        {{ t('rules.tictactoe.title') }}
      </h2>

      <div class="space-y-6 text-blue-100 leading-relaxed">
        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            🎯 {{ t('rules.tictactoe.objective') }}
          </h3>
          <i18n-t keypath="rules.tictactoe.objectiveText" tag="p">
            <template #highlight>
              <strong>{{ t('rules.tictactoe.objectiveHighlight') }}</strong>
            </template>
          </i18n-t>
        </div>

        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            📜 {{ t('rules.tictactoe.gameplay') }}
          </h3>
          <ul class="list-disc ps-5 space-y-2">
            <i18n-t keypath="rules.tictactoe.rule1" tag="li">
              <template #highlight>
                <strong>{{ t('rules.tictactoe.rule1Highlight') }}</strong>
              </template>
            </i18n-t>
            <i18n-t keypath="rules.tictactoe.rule2" tag="li">
              <template #x>
                <strong>X</strong>
              </template>
              <template #o>
                <strong>O</strong>
              </template>
            </i18n-t>
            <li>{{ t('rules.tictactoe.rule3') }}</li>
            <li>{{ t('rules.tictactoe.rule4') }}</li>
          </ul>
        </div>

        <div>
          <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
            🏆 {{ t('rules.tictactoe.winning') }}
          </h3>
          <ul class="list-disc ps-5 space-y-2">
            <li>
              <i18n-t keypath="rules.tictactoe.victory" tag="span">
                <template #label>
                  <strong>{{ t('rules.tictactoe.victoryLabel') }}</strong>
                </template>
              </i18n-t>
            </li>
            <li>
              <i18n-t keypath="rules.tictactoe.drawRule" tag="span">
                <template #label>
                  <strong>{{ t('rules.tictactoe.drawLabel') }}</strong>
                </template>
              </i18n-t>
            </li>
          </ul>
        </div>

        <blockquote class="border-s-4 border-green-400 bg-green-500/10 p-4 rounded-e-xl mt-4">
          <span class="font-bold text-white block mb-1">💡 {{ t('rules.tictactoe.hintTitle') }}</span>
          {{ t('rules.tictactoe.hintText') }}
        </blockquote>
      </div>
    </section>
  </div>
</template>
