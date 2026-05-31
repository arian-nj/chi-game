<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { AccountService } from '../gen/account/v1/account_pb';
import { authTransport } from '../lib/transport';

// You can import a default profile picture or provide a fallback
import profilePic from '../assets/profile.jpg'

const { isPending, error, data } = useQuery({
    queryKey: ['me'],
    staleTime: 0,
    retryDelay: 10000,
    queryFn: async () => {
        const client = createClient(AccountService, authTransport)
        const data = await client.getMe({})
        return data
    }
})
</script>

<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-linear-to-br from-custom-blue to-custom-lite-blue">
        <div 
          v-if="isPending" 
          class="flex flex-col items-center justify-center space-y-4 animate-pulse w-full"
        >
            <div class="w-28 h-28 bg-blue-400 rounded-full mb-4"></div>
            <div class="h-7 w-40 bg-blue-300 rounded"></div>
            <div class="h-5 w-32 bg-blue-200 rounded"></div>
        </div>
        <div 
          v-else-if="error" 
          class="flex flex-col items-center justify-center space-y-2"
        >
            <h1 class="text-xl text-red-400 font-bold">خطا در دریافت اطلاعات</h1>
            <p class="text-gray-300">لطفاً دوباره تلاش کنید.</p>
        </div>
        <div 
          v-else-if="data?.account"
          class="bg-white/10 p-8 rounded-2xl shadow-lg flex flex-col items-center gap-3 w-90 max-w-full"
        >
            <img 
              :src="profilePic" 
              alt="Profile Picture"
              class="w-36 h-36 rounded-full object-cover border-4 border-white/40 shadow-lg mb-5"
            >
            <h1 class="text-3xl font-extrabold text-gray-50 mb-2">
                {{ data.account.displayName || "بدون نام" }}
            </h1>
            <div class="flex items-center bg-custom-lite-blue rounded-xl px-4 py-2">
                <p class="text-xl font-bold text-blue-100">@{{ data.account.username }}</p>
            </div>
            <!-- <div v-if="data.account.bio" class="mt-4 w-full bg-white/5 rounded-lg p-3"> -->
                <!-- <p class="text-base text-gray-200">{{ data.account.bio }}</p> -->
            <!-- </div> -->
        </div>
    </div>
</template>