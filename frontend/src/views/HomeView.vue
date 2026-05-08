<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { createClient } from "@connectrpc/connect";
import gsap from 'gsap'

// import wasmUrl from "@lottiefiles/dotlottie-web/dist/dotlottie-player.wasm?url";
import wasmUrl from "@lottiefiles/dotlottie-web/dotlottie-player.wasm?url";

// import { switchInlineQuery } from '@telegram-apps/sdk';
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { RoomService } from '../gen/room/v1/room_pb';
import MeComponent from '../components/MeComponent.vue';
import GameSelectorBtn from '../components/home/GameSelectorBtn.vue';
import { authTransport } from '../lib/transport';
import OnlineComponent from '../components/OnlineComponent.vue';
import { useToast } from '../components/Toast.vue';

onMounted(() => {
  // prefetch
  import('../views/FinderView.vue')
  import(/* @vite-ignore */  wasmUrl)
})

const router = useRouter()

const games = ["xo3x3", "conn4"];
// const games = ["xo3x3", "conn4","snake","apple","tart","mart","gart","sdar"];

const selectedGame = ref(games[1])

const playBtnRef = ref<HTMLButtonElement | null>(null)

watch(playBtnRef,() => {
  if (playBtnRef.value) {
    gsap.from(playBtnRef.value, {
      yPercent: 100,
      opacity: 0,
      duration: 1,
    })
  }
})

const exitBtnRef = ref<HTMLButtonElement | null>(null)
watch(exitBtnRef,()=>{
    if (exitBtnRef.value) {
    gsap.from(exitBtnRef.value, {
      yPercent: 100,
      opacity: 0,
      duration: 1,
    })
  }
})

function onPlayFriendsClick() {
}


let coolDownTime = 1_000
const maxCoolDownTime = 32_000

const { data: hasRoomData } = useQuery({
  queryKey: ['hasRoom'],
  queryFn: async () => {
    const client = createClient(RoomService, authTransport)
    const data = await client.hasRoom({})
    return data
  },
  refetchInterval: (query) => {
    const hasSes = query.state.data?.hasRoom
    if (hasSes) {
      if (coolDownTime * 2 <= maxCoolDownTime){
        coolDownTime *= 2
      }
      return coolDownTime
    }
    coolDownTime = 1_000
    return false
  }
})


const hasRoom = computed(() => {
  if (hasRoomData.value && hasRoomData.value.hasRoom) {
    return hasRoomData.value.hasRoom
  }
  return false
})

function handlePlayClick() {
  if (hasRoom.value) {
    router.push(`/room`)
    return
  }
  router.push(`/finder?game=${selectedGame.value}`)
}

async function handleExistRoom() {
  const client = createClient(RoomService,authTransport)
  const data = await client.closeRoom({})
  const { toast } = useToast()
  if (!data.isOk) {
    toast.info("یه مشکلی پیش اومده")
  }else{
    toast.success("حله خارج شدی")
    const queryClient = useQueryClient()
    queryClient.invalidateQueries({queryKey:["hasRoom"]})
  }
}

</script>

<template>
  <main class="h-screen">
    <div
      class="relative h-full w-screen bg-linear-to-br from-gray-950 via-gray-900 to-black text-white overflow-hidden">
      
        <div class="flex justify-center pt-12">
    <div class="bg-gray-900/60 p-6 rounded-2xl shadow-xl backdrop-blur-md border border-gray-700">
      <div class="flex items-center gap-3">
        <OnlineComponent />
        <MeComponent />
      </div>
    </div>
  </div>
      <div class="flex justify-center gap-2 pt-16 flex-wrap">
        <GameSelectorBtn v-for="game in games" :key="game" :game-name="game" :selected="game == selectedGame"
          @choosed="selectedGame = game" />
      </div>

      <div class="absolute bottom-0 left-0 w-full flex flex-col gap-1 items-center justify-center">
        <button class="bg-linear-to-r to-cyan-500 from-pink-300
                  rounded-2xl px-10 py-4
                 text-2xl font-extrabold text-white tracking-wide
                 shadow-lg hover:shadow-cyan-400/30
                 hover:scale-105 active:scale-95
                 transition-all duration-300 ease-in-out" @click="onPlayFriendsClick()">
          🎮 بازی با دوستان
        </button>
        
        <button v-if="hasRoom" ref="exitBtnRef" type="button" :disabled="selectedGame == ''" @click="handleExistRoom" :class="[
          `w-1/3 py-6 text-3xl font-bold
              focus:outline-none ring-2 ring-gray-900
              transition-colors duration-400
              rounded-3xl
              `,
          selectedGame
            ? 'bg-linear-to-r to-red-400 from-pink-300 text-white hover:opacity-95 shadow-xl'
            : 'bg-gray-700 text-gray-400 cursor-not-allowed'
        ]">
       ❌ خروج از بازی قبلی
        </button>

        <button ref="playBtnRef" type="button" :disabled="selectedGame == ''" @click="handlePlayClick" :class="[
          `w-full py-6 text-3xl font-bold
              focus:outline-none focus:ring-4 focus:ring-pink-400/40
              transition-colors duration-400`,
          selectedGame
            ? 'bg-linear-to-r from-emerald-400 to-green-500 text-white hover:opacity-95 shadow-xl'
            : 'bg-gray-700 text-gray-400 cursor-not-allowed'
        ]">
          {{ hasRoom ?
            "ادامه بازی"
            : "🚀 شروع بازی" }}
        </button>

      </div>

    </div>
  </main>
</template>
