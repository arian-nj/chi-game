<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import ChatInput from '../components/chat/ChatInput.vue';
import { ChatService } from '../gen/chat/v1/chat_pb';
import { authTransport } from '../lib/transport';
import { computed, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useToast } from '../components/Toast.vue';
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { AccountService } from '../gen/account/v1/account_pb';
import ChatBubble from '../components/chat/ChatBubble.vue';


// const allChatMessages = ref(Array<Message>())

const { toast } = useToast()

const route = useRoute()
const chatId = computed(() => Number(route.params.id))
const { data:meData } = useQuery({
  queryKey: ['me'],
  staleTime: 0,
  retryDelay: 10000,
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})

const messagesQueryKey = computed(() => [`chat-${chatId.value}-messages`])

const queryClient = useQueryClient()
async function sendMessage(msgText: string) {
    if (msgText == "") {
        return
    }
    const chatClient = createClient(ChatService, authTransport)
    try {
        const data = await chatClient.sendMessage({
            chatId: BigInt(chatId.value),
            content: msgText,
        })
        console.log(data)
    } catch (error) {
        console.error(error)
        toast.error("خطا در ارسال پیام")
        console.error(error)
    }
    queryClient.invalidateQueries({ queryKey: messagesQueryKey.value })

}

const { data: allChatMessagesData } = useQuery({
    queryKey: messagesQueryKey,
    queryFn: async () => {
        const chatClient = createClient(ChatService, authTransport)
        const data = await chatClient.getMessages({ chatId: BigInt(chatId.value) })
        return data
    }
})

const allChatMessages = computed(() =>
    [...(allChatMessagesData.value?.messages ?? [])].reverse()
)

const pollInterval = setInterval(() => {
    queryClient.invalidateQueries({ queryKey: messagesQueryKey.value })
}, 3_000)
onUnmounted(() => clearInterval(pollInterval))

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
</template>