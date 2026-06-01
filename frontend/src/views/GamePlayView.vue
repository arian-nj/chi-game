<script setup lang="ts">
import ToggleSwitch from '@/components/ToggleSwitch.vue';
import { useTemplateRef } from 'vue';

const botToggleSwitchRef = useTemplateRef<InstanceType<typeof ToggleSwitch>>('botToggleSwitchRef');

</script>

<template>
  <div class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md
  flex p-3
  ">
    <ToggleSwitch option-one="Bot 🤖" option-two="👥 2 Player" ref="botToggleSwitchRef"></ToggleSwitch>
  </div>
</template>
