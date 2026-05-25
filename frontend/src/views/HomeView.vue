<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { createClient } from "@connectrpc/connect";
import gsap from 'gsap'
import wasmUrl from "@lottiefiles/dotlottie-web/dotlottie-player.wasm?url";
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { RoomService } from '../gen/room/v1/room_pb';
import MeComponent from '../components/MeComponent.vue';
import { authTransport } from '../lib/transport';
import OnlineComponent from '../components/OnlineComponent.vue';
import { useToast } from '../components/Toast.vue';

onMounted(() => {
  // prefetch
  import('../views/FinderView.vue')
  import(/* @vite-ignore */  wasmUrl)
})

const router = useRouter()

const games = [
  { key: "xo3x3", name: "XO 3x3", emoji: "❌⭕" },
  { key: "conn4", name: "Connect 4", emoji: "🟡🔴" }
]
const selectedGame = ref(games[1].key)

const playBtnRef = ref<HTMLButtonElement | null>(null)
watch(playBtnRef,() => {
  if (playBtnRef.value) {
    gsap.from(playBtnRef.value, {
      yPercent: 100,
      opacity: 0,
      duration: 0.8,
      ease: 'power2.out'
    })
  }
})

const exitBtnRef = ref<HTMLButtonElement | null>(null)
watch(exitBtnRef,()=>{
  if (exitBtnRef.value) {
    gsap.from(exitBtnRef.value, {
      yPercent: 100,
      opacity: 0,
      duration: 0.8,
      ease: 'power2.out'
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
<main class="h-screen font-[Rubik] bg-linear-to-br from-[#01041B] via-[#0B122D] to-[#110D24]">
  <div class="relative h-full w-screen text-white overflow-hidden">

    <!-- Burst of color and subtle glowing particles behind main content -->
    <div class="absolute top-0 left-0 w-full h-full pointer-events-none z-0">
      <div class="absolute top-[8%] left-[10%] w-80 h-80 bg-linear-to-br from-emerald-400 via-blue-700/20 to-pink-400/50 opacity-5 rounded-full blur-3xl animate-[pulse_5s_infinite]" />
      <div class="absolute top-[60%] left-[45%] w-56 h-56 bg-linear-to-br from-pink-600 to-violet-600 opacity-10 rounded-full blur-3xl animate-[pulse_10s_infinite]" />
      <div class="absolute bottom-[-10%] right-[10%] w-96 h-80 bg-linear-to-br from-[#ea00ff] to-cyan-500 opacity-10  rounded-full blur-3xl animate-[pulse_10s_infinite_2s]" />
    </div>

    <div class="flex justify-center pt-16 z-10 relative">
      <div class="bg-linear-to-b from-gray-950 via-gray-900/80 to-gray-800/80 p-8 rounded-3xl shadow-2xl backdrop-blur-xl border border-gray-700 transition-all duration-500 hover:scale-[1.02]">
        <div class="flex items-center gap-5 drop-shadow-lg">
          <OnlineComponent />
          <span class="mx-2 w-0.5 h-12 bg-blue-700/30 rounded-xl"></span>
          <MeComponent />
        </div>
      </div>
    </div>

    <!-- Game Selection Buttons -->
    <div class="flex gap-4 justify-center pt-14 flex-wrap z-10 relative">
      <div
        v-for="game in games"
        :key="game.key"
        @click="selectedGame = game.key"
        :class="[
          'px-8 py-4 cursor-pointer rounded-2xl font-bold text-lg shadow-lg ring-1 ring-blue-800/20 border-none transition-all duration-300 select-none flex items-center gap-2',
          'sm:w-40 sm:max-w-40 lg:w-60 lg:max-w-60',
          'flex flex-col items-center justify-center',
          selectedGame === game.key
            ? 'bg-linear-to-tr text-white scale-105 shadow-2xl ring-4 ring-gray-50/80'
            : 'bg-gray-800/80 text-blue-100 hover:bg-gray-700/70 hover:scale-105'
        ]"
      >
        <div class="flex flex-col items-center gap-2 text-center">
          <span class="text-3xl w-full text-center">{{ game.emoji }}</span>
          <span class="text-lg w-full text-center">{{ game.name }}</span>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="absolute bottom-14 left-0 w-full flex flex-col items-center justify-center gap-4 px-4 pb-12 z-20">
      <button
        class="w-full md:w-2/3 bg-linear-to-tr from-violet-500 via-sky-400 to-emerald-400 hover:from-indigo-400 hover:to-fuchsia-400
          rounded-3xl px-8 py-5
          text-2xl font-extrabold text-white tracking-wide
          shadow-xl hover:shadow-fuchsia-400/20 transition-all duration-300 ease-in-out
          border-none flex items-center justify-center gap-3 group"
        @click="onPlayFriendsClick"
      >
        <span class="text-3xl group-hover:animate-bounce">🎮</span> بازی با دوستان
      </button>

      <transition name="fade-slide">
        <button
          v-if="hasRoom"
          ref="exitBtnRef"
          type="button"
          :disabled="selectedGame == ''"
          @click="handleExistRoom"
          :class="[
            `w-full md:w-2/3 py-5 text-2xl font-bold border-none rounded-3xl
              focus:outline-none ring-2 ring-gray-900 shadow-xl
              transition-all duration-400 flex items-center justify-center gap-2`,
            selectedGame
              ? 'bg-linear-to-r from-pink-400 to-red-500 text-white hover:opacity-95'
              : 'bg-gray-700 text-gray-400 cursor-not-allowed'
          ]"
        >
          <span class="text-2xl animate-pulse">❌</span>
          خروج از بازی قبلی
        </button>
      </transition>

      <button
        ref="playBtnRef"
        type="button"
        :disabled="selectedGame == ''"
        @click="handlePlayClick"
        :class="[
          `w-full md:w-2/3 py-6 text-3xl font-bold rounded-3xl border-none
            focus:outline-none focus:ring-4 focus:ring-pink-400/40
            transition-all duration-400 shadow-2xl mt-2 flex items-center justify-center gap-3`,
          selectedGame
            ? 'bg-linear-to-r from-emerald-400 via-cyan-500 to-green-500 text-white hover:opacity-95'
            : 'bg-gray-700 text-gray-400 cursor-not-allowed'
        ]"
      >
        <span v-if="hasRoom" class="material-icons animate-pulse text-white/90">play_circle</span>
        <span v-else class="text-3xl animate-bounce">🚀</span>
        {{ hasRoom ? "ادامه بازی" : "شروع بازی" }}
      </button>
    </div>
  </div>
</main>
</template>

<style scoped>
@keyframes fade-slide-enter-active {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.97);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
@keyframes fade-slide-leave-active {
  0% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
  100% {
    opacity: 0;
    transform: translateY(20px) scale(0.97);
  }
}
.fade-slide-enter-active {
  animation: fade-slide-enter-active .5s cubic-bezier(.4,2,.75,.82);
}
.fade-slide-leave-active {
  animation: fade-slide-leave-active .4s cubic-bezier(.4,2,.75,.82);
}
</style>
