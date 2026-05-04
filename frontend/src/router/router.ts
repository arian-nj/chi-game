import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { GetDeviceId, GetJwtToken, SetDeviceId, SetJwtToken } from '../lib/auth'

import type { Router } from "vue-router"
import { AuthService } from '../gen/auth/v1/auth_pb'
import { rawTransport } from '../lib/transport'
import { createClient } from '@connectrpc/connect'

// import { initDataRaw, restoreInitData } from '@telegram-apps/sdk';


export function setupRouterGuards(router: Router) {
  router.beforeEach(async (to, _, next) => {
    // Routes that don't need authentication
    if (to.name === "login") return next()

    let token = GetJwtToken()
    if (token) {
      // Optional: verify token expiration here
      return next()
    }

    try {
      const client = createClient(AuthService, rawTransport)
      const deviceId = GetDeviceId() // read from localStorage, may be empty
      const data = await client.validateGuest({ deviceId })
      
      // Save both token and device ID
      SetJwtToken(data.token)
      SetDeviceId(data.deviceId)  // server may return new one
      
      return next()
    } catch (err) {
      console.error("auth failed:", err)
      return next({ name: "login-fail" })
    }

  })
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/finder',
      name: 'finder',
      component: () => import('../views/FinderView.vue')
    },
    {
      path: '/room',
      name: 'room',
      component: () => import('../views/RoomView.vue')
    },
    // {
    //   path: '/card',
    //   name: 'card',
    //   component: () => import('../components/game/card/CardView.vue')
    // },
  ],
})

export default router

