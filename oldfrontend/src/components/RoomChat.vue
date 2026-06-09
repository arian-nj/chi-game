<script setup lang="ts">
import { useRoomSession } from '@/composables/use-room-session';
import { nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const {
  textDir,
  chatMessages,
  chatConnected,
  chatHasMessages,
  sendChatMessage,
  chatPlayerName,
  isOwnChatMessage,
} = useRoomSession();

const draft = ref('');
const messageListRef = ref<HTMLElement | null>(null);
const messagesExpanded = ref(false);

async function scrollToLatest() {
  await nextTick();
  const list = messageListRef.value;
  if (!list) {
    return;
  }
  list.scrollTop = list.scrollHeight;
}

watch(chatMessages, () => {
  void scrollToLatest();
}, { deep: true });

function toggleMessages() {
  messagesExpanded.value = !messagesExpanded.value;
  if (messagesExpanded.value) {
    void scrollToLatest();
  }
}

function submitMessage() {
  if (!sendChatMessage(draft.value)) {
    return;
  }
  draft.value = '';
  void scrollToLatest();
}
</script>

<template>
  <section
    class="fixed inset-x-0 bottom-0 z-40 flex flex-col border-t border-white/10 bg-custom-blue/95 shadow-[0_-4px_24px_rgba(0,0,0,0.2)] backdrop-blur-md lg:static lg:z-auto lg:h-[40rem] lg:w-full lg:overflow-hidden lg:rounded-2xl lg:border lg:border-white/10 lg:bg-custom-lite-blue/40 lg:shadow-md lg:backdrop-blur-none"
    :dir="textDir"
    :aria-label="t('invite.chatLabel')"
  >
    <header class="hidden shrink-0 items-center justify-between gap-2 border-b border-white/10 px-4 py-3 lg:flex">
      <h2 class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">
        {{ t('invite.chatLabel') }}
      </h2>
      <span
        class="inline-flex items-center gap-1.5 text-xs font-medium"
        :class="chatConnected ? 'text-green-300' : 'text-amber-300'"
      >
        <span
          class="h-2 w-2 rounded-full"
          :class="chatConnected ? 'bg-green-400' : 'bg-amber-400'"
          aria-hidden="true"
        />
        {{ chatConnected ? t('invite.chatConnected') : t('invite.chatConnecting') }}
      </span>
    </header>

    <div
      ref="messageListRef"
      class="min-h-0 flex-1 space-y-3 overflow-y-auto px-3 py-4 lg:px-4"
      :class="messagesExpanded
        ? 'flex max-h-[min(50vh,24rem)] flex-col border-b border-white/10 lg:max-h-none lg:border-b-0'
        : 'hidden lg:block'"
      role="log"
      aria-live="polite"
      aria-relevant="additions"
    >
      <p
        v-if="!chatHasMessages"
        class="px-2 py-6 text-center text-sm text-blue-100/70"
      >
        {{ t('invite.chatEmpty') }}
      </p>

      <article
        v-for="message in chatMessages"
        :key="message.id"
        class="flex flex-col gap-1"
        :class="isOwnChatMessage(message.playerId) ? 'items-end' : 'items-start'"
      >
        <span
          v-if="!isOwnChatMessage(message.playerId)"
          class="px-1 text-xs font-semibold text-blue-100/80"
        >
          {{ chatPlayerName(message.playerId) }}
        </span>
        <p
          class="max-w-[88%] wrap-break-word rounded-2xl px-3 py-2 text-sm leading-relaxed text-white sm:max-w-[80%] sm:text-base"
          :class="isOwnChatMessage(message.playerId)
            ? 'rounded-br-md bg-green-500/90'
            : 'rounded-bl-md bg-custom-deep-blue/80'"
        >
          {{ message.text }}
        </p>
      </article>
    </div>

    <form
      class="flex shrink-0 gap-2 p-3 pb-[max(0.75rem,env(safe-area-inset-bottom))] lg:border-t lg:border-white/10 lg:p-4 lg:pb-4"
      @submit.prevent="submitMessage"
    >
      <button
        type="button"
        class="relative shrink-0 rounded-xl border border-white/20 bg-white/10 px-3 py-2.5 text-sm font-bold text-white transition hover:bg-white/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 lg:hidden"
        :aria-expanded="messagesExpanded"
        :aria-label="messagesExpanded ? t('invite.chatHideMessages') : t('invite.chatShowMessages')"
        @click="toggleMessages"
      >
        <span aria-hidden="true">{{ messagesExpanded ? '▾' : '▴' }}</span>
        <span
          v-if="chatMessages.length > 0 && !messagesExpanded"
          class="absolute -right-1 -top-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-green-500 px-1 text-[10px] font-bold leading-none text-white"
        >
          {{ chatMessages.length > 9 ? '9+' : chatMessages.length }}
        </span>
      </button>
      <label class="sr-only" for="room-chat-input">{{ t('invite.chatPlaceholder') }}</label>
      <input
        id="room-chat-input"
        v-model="draft"
        type="text"
        maxlength="500"
        autocomplete="off"
        enterkeyhint="send"
        :placeholder="t('invite.chatPlaceholder')"
        :disabled="!chatConnected"
        class="min-w-0 flex-1 rounded-xl border border-white/15 bg-custom-deep-blue/70 px-3 py-2.5 text-sm text-white placeholder:text-blue-100/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 disabled:cursor-not-allowed disabled:opacity-60 sm:px-4 sm:py-3 sm:text-base"
      >
      <button
        type="submit"
        :disabled="!chatConnected || !draft.trim()"
        class="shrink-0 rounded-xl border border-white/20 bg-white/15 px-4 py-2.5 text-sm font-bold text-white transition hover:bg-white/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 disabled:cursor-not-allowed disabled:opacity-50 sm:px-5 sm:py-3"
      >
        {{ t('invite.chatSend') }}
      </button>
    </form>
  </section>
</template>
