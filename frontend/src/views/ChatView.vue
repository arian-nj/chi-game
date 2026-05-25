<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import ChatInput from '../components/chat/ChatInput.vue';
import { ChatService } from '../gen/chat/v1/chat_pb';
import { authTransport } from '../lib/transport';
import { computed, onUnmounted, ref, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useToast } from '../components/Toast.vue';
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { AccountService } from '../gen/account/v1/account_pb';
import ChatBubble from '../components/chat/ChatBubble.vue';

const { toast } = useToast();

const route = useRoute();
const router = useRouter();
const chatId = computed(() => Number(route.params.id));
const { data: meData } = useQuery({
  queryKey: ['me'],
  staleTime: 0,
  retryDelay: 10000,
  queryFn: async () => {
    const client = createClient(AccountService, authTransport);
    const data = await client.getMe({});
    return data;
  }
});

const messagesQueryKey = computed(() => [`chat-${chatId.value}-messages`]);

const queryClient = useQueryClient();

async function sendMessage(msgText: string) {
  if (msgText.trim() === "") {
    toast.error("لطفاً پیامی وارد کنید");
    return;
  }
  const chatClient = createClient(ChatService, authTransport);
  try {
    const data = await chatClient.sendMessage({
      chatId: BigInt(chatId.value),
      content: msgText,
    });
    // Success feedback (e.g., optional animation or sound)
  } catch (error) {
    console.error(error);
    toast.error("خطا در ارسال پیام");
  }
  queryClient.invalidateQueries({ queryKey: messagesQueryKey.value });
}

const { data: allChatMessagesData, isLoading: isMessagesLoading } = useQuery({
  queryKey: messagesQueryKey,
  queryFn: async () => {
    const chatClient = createClient(ChatService, authTransport);
    const data = await chatClient.getMessages({ chatId: BigInt(chatId.value) });
    return data;
  }
});

const allChatMessages = computed(() =>
  [...(allChatMessagesData.value?.messages ?? [])].reverse()
);

// Auto-scroll to bottom functionality
const messagesEndRef = ref<HTMLElement | null>(null);

watch(allChatMessages, async () => {
  await nextTick();
  if (messagesEndRef.value) {
    messagesEndRef.value.scrollIntoView({ behavior: 'smooth' });
  }
});

const pollInterval = setInterval(() => {
  queryClient.invalidateQueries({ queryKey: messagesQueryKey.value });
}, 3_000);
onUnmounted(() => clearInterval(pollInterval));

function goBack() {
  router.back();
}
</script>



<template>
<div class="z-50 absolute bottom-0 bg-neutral-600 h-screen w-screen flex flex-col overflow-y-hidden">
    <!-- messages -->
    <div class="flex flex-col w-full h-full justify-end font-[Rubik] text-lg overflow-y-auto lg:max-w-1/2 m-auto">
        <div class="flex flex-col w-full h-full px-1.5">
            <ChatBubble class="py-1" 
                v-for="msg in allChatMessages" 
                :text="msg.content" 
                :is-me="meData!.account?.id == msg.senderPersonId" />
        </div>

    </div>
    <!-- input -->
    <div class="relative">
        <ChatInput @submit="sendMessage"></ChatInput>
    </div>
</div>


<div class="z-50 absolute bottom-0 bg-linear-to-tr from-custom-lite-blue  to-custom-blue h-screen w-screen flex flex-col overflow-hidden font-[Rubik]">

    <!-- Chat Header -->
    <div class="w-full flex items-center justify-between px-6 py-4 bg-neutral-900 shadow-lg">
      <div class="flex items-center space-x-3">
        <!-- Back Button -->
        <button
          @click="goBack"
          class="text-white hover:bg-neutral-800 rounded-full p-2 transition-colors duration-150 flex items-center"
          title="برگشت"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 rtl:scale-x-[-1]" fill="none" viewBox="0 0 20 20" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7 7-7M3 12h14" />
          </svg>
        </button>
        <span class="text-white text-xl font-bold">گفتگو</span>
      </div>
      <div>
        <span class="text-neutral-400 text-sm">ID: {{ chatId }}</span>
      </div>
    </div>

    <!-- Messages -->
    <div
      class="flex-1 w-full flex flex-col items-center justify-end overflow-y-auto px-0 md:px-0 py-3"
      style="background: transparent;">
      <div class="w-full md:max-w-2xl lg:max-w-xl flex flex-col gap-1 px-2 pb-2 overflow-y-auto">
        <transition-group name="chat-message-fade" tag="div">
          <div v-if="isMessagesLoading" class="w-full flex justify-center items-center py-6">
            <span class="text-neutral-300 animate-pulse">در حال بارگذاری پیام‌ها...</span>
          </div>
          <ChatBubble
            v-for="msg in allChatMessages"
            :key="msg.id.toString()"
            class="py-1"
            :text="msg.content"
            :is-me="meData?.account?.id == msg.senderPersonId"
          />
        </transition-group>
        <div ref="messagesEndRef"></div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="relative w-full bg-neutral-900 py-3 shadow-lg border-t border-neutral-800 flex items-end px-3 md:px-0">
      <div class="w-full md:max-w-2xl lg:max-w-xl m-auto flex items-center">
        <ChatInput @submit="sendMessage" />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Subtle scrollbar */
::-webkit-scrollbar {
  width: 8px;
  background: #23272e;
}
::-webkit-scrollbar-thumb {
  background: #353941;
  border-radius: 4px;
}

/* Animations for message appearance */
.chat-message-fade-enter-active,
.chat-message-fade-leave-active {
  transition: opacity 0.25s;
}
.chat-message-fade-enter-from, .chat-message-fade-leave-to {
  opacity: 0;
}
</style>