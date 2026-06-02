<script setup lang="ts">
import DifficultySelector from '@/components/DifficultySelector.vue';
import ToggleSwitch from '@/components/ToggleSwitch.vue';
import { useTextDirection } from '@/composables/use-text-direction';
import type { BotDifficulty } from '@/libs/bot-difficulty';
import type { Connect4Settings } from '@/libs/connect4';
import type { TicTacToeSettings } from '@/libs/tictactoe';
import { computed, ref, unref, useTemplateRef } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

export type PlaySettings = TicTacToeSettings | Connect4Settings;

const route = useRoute();
const { t } = useI18n();
const { textDir } = useTextDirection();
const gameKey = computed(() => route.params.game as string);
const isTicTacToe = computed(() => gameKey.value === 'tictactoe');
const isConnect4 = computed(() => gameKey.value === 'conn4');
const isBotMode = ref(true);
const botDifficulty = ref<BotDifficulty>('medium');

const botToggleSwitchRef = useTemplateRef<InstanceType<typeof ToggleSwitch>>('botToggleSwitchRef');
const boardSizeToggleSwitchRef = useTemplateRef<InstanceType<typeof ToggleSwitch>>('boardSizeToggleSwitchRef');

const difficultyLabels = computed(() => ({
  easy: t('settings.difficultyEasy'),
  medium: t('settings.difficultyMedium'),
  hard: t('settings.difficultyHard'),
}));

function readToggleOption(toggle: InstanceType<typeof ToggleSwitch> | null | undefined): number {
  return unref(toggle?.optionNumber) ?? 1;
}

function onOpponentChange(option: 1 | 2) {
  isBotMode.value = option === 1;
}

function getSettings(): PlaySettings {
  const isBot = readToggleOption(botToggleSwitchRef.value) === 1;
  const settings = { isBot, botDifficulty: botDifficulty.value };

  if (isConnect4.value) {
    return settings;
  }

  return {
    ...settings,
    boardSize: readToggleOption(boardSizeToggleSwitchRef.value) === 1 ? 3 : 5,
  };
}

defineExpose({ getSettings });
</script>

<template>
  <div class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md flex flex-col gap-4 p-4">
    <div class="flex flex-col gap-2">
      <span :dir="textDir" class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('settings.opponent') }}</span>
      <ToggleSwitch
        :option-one="t('settings.bot')"
        :option-two="t('settings.twoPlayer')"
        ref="botToggleSwitchRef"
        @change="onOpponentChange"
      />
    </div>

    <div v-if="isBotMode" class="flex flex-col gap-2">
      <span :dir="textDir" class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('settings.difficulty') }}</span>
      <DifficultySelector v-model="botDifficulty" :labels="difficultyLabels" />
    </div>

    <div v-if="isTicTacToe" class="flex flex-col gap-2">
      <span :dir="textDir" class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('settings.boardSize') }}</span>
      <ToggleSwitch
        :option-one="t('settings.board3x3')"
        :option-two="t('settings.board5x5')"
        ref="boardSizeToggleSwitchRef"
      />
    </div>

    <p v-else-if="isConnect4" :dir="textDir" class="text-sm text-blue-100/80">
      {{ t('settings.connect4Hint') }}
    </p>
  </div>
</template>
