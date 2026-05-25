<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import ChatInput from '../components/chat/ChatInput.vue';
import { ChatService } from '../gen/chat/v1/chat_pb';
import { authTransport } from '../lib/transport';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { useToast } from '../components/Toast.vue';

const { toast } = useToast()

const route = useRoute()
const chatId = ref<number>(
    Number(route.params.id)
)

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
}

</script>

<template>
    <div class="z-50 fixed top-0 left-0 right-0 
    flex flex-col items-center justify-center h-screen bg-neutral-400 gap-6
    ">
        <h1>Chat</h1>
        <div class="absolute bottom-0 left-0 right-0">
            <ChatInput @submit="sendMessage"></ChatInput>
        </div>
    </div>
</template>