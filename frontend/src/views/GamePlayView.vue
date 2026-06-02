<script setup lang="ts">
import ToggleSwitch from '@/components/ToggleSwitch.vue';
import type { TicTacToeSettings } from '@/libs/tictactoe';
import { unref, useTemplateRef } from 'vue';

const botToggleSwitchRef = useTemplateRef<InstanceType<typeof ToggleSwitch>>('botToggleSwitchRef');
const boardSizeToggleSwitchRef = useTemplateRef<InstanceType<typeof ToggleSwitch>>('boardSizeToggleSwitchRef');

function readToggleOption(toggle: InstanceType<typeof ToggleSwitch> | null | undefined): number {
  return unref(toggle?.optionNumber) ?? 1;
}

function getSettings(): TicTacToeSettings {
  return {
    isBot: readToggleOption(botToggleSwitchRef.value) === 1,
    boardSize: readToggleOption(boardSizeToggleSwitchRef.value) === 1 ? 3 : 5,
  };
}

defineExpose({ getSettings });
</script>

<template>
  <div class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md flex flex-col gap-4 p-4">
    <div class="flex flex-col gap-2">
      <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">Opponent</span>
      <ToggleSwitch option-one="Bot 🤖" option-two="👥 2 Player" ref="botToggleSwitchRef" />
    </div>

    <div class="flex flex-col gap-2">
      <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">Board size</span>
      <ToggleSwitch option-one="3×3" option-two="5×5" ref="boardSizeToggleSwitchRef" />
    </div>
  </div>
</template>
