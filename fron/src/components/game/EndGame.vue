<script setup lang="ts">
import type { Account } from '../../gen/account/v1/account_pb';
import gsap from 'gsap';
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from '../Toast.vue';
import { createClient } from '@connectrpc/connect';
import { RoomService } from '../../gen/room/v1/room_pb';
import { authTransport } from '../../lib/transport';
import { useQueryClient } from '@tanstack/vue-query';

const props = defineProps<{
  winner: Account | undefined,
  loser: Account | undefined,
}>()

const endPanelDivRef = useTemplateRef("end-panel-div-ref")
const homeBtnRef = useTemplateRef('home-btn-ref')

onMounted(() => {
  if (endPanelDivRef.value) {
    gsap.from(endPanelDivRef.value, {
      opacity: .5,
      scaleX: .5,
      scaleY: .8,
      y: 100,
      duration: .4,
    })
  }
})

onMounted(() => {
  if (homeBtnRef.value) {
    gsap.fromTo(homeBtnRef.value,
      {
        rotate: 2
      },
      {
        rotate: -2,
        duration: 1,
        yoyo: true,
        repeat: -1,
        ease: 'power1.inOut'
      },
    )
  }
})
const router = useRouter()
function goToHome() {
  handleExistRoom()
  router.push('/')
}
const remainedRoomTime = ref(30)

let interval: number | undefined

onMounted(() => {
  interval = setInterval(() => {
    if (remainedRoomTime.value > 0) {
      remainedRoomTime.value--;
    } else {
      clearInterval(interval)
    }
  }, 1000)
})

onBeforeUnmount(() => {
  clearInterval(interval)
})

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
  <div ref="end-panel-div-ref" class="
    absolute min-w-[80%] max-w-[90%] max-h-[90%] p-8
    bg-slate-800/70 backdrop-blur-lg
    rounded-2xl border border-slate-600 shadow-2xl shadow-black/50
    flex flex-col justify-center items-center gap-6
    text-center text-white
    overflow-y-auto
  ">
    <div v-if="props.winner" class="flex flex-col items-center">
      <h1
        class="text-5xl font-black bg-linear-to-r from-amber-400 to-yellow-300 bg-clip-text text-transparent drop-shadow-lg">
        🏆 {{ props.winner.displayName }}
      </h1>
    </div>

    <div v-if="props.loser" class="flex flex-col items-center opacity-60">
      <h1 class="text-4xl font-medium text-slate-400 line-through">
        💀 {{ props.loser.displayName }}
      </h1>
    </div>

    <hr class="w-1/2 border-slate-600 my-2" />
    <div v-if="remainedRoomTime != 0">
      <h1 class="text-4xl font-bold text-red-200">{{ remainedRoomTime }} </h1>
      <h1 class="text-lg font-bold text-red-200" dir="auto">ثانیه دیگه چت بسته میشه</h1>
    </div>
    <div v-else>
      <h1 class="text-2xl font-bold text-red-200">چت بسته شد </h1>
    </div>
    <button class="
        p-4 rounded-full text-5xl
        bg-slate-700/50 text-white
        border border-slate-600
        active:scale-95
        shadow-xl
      " ref="home-btn-ref" @click="goToHome()">
      🏠
    </button>
  </div>
</template>
