import HomeView from '@/views/HomeView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/game/:game',
      name: 'game',
      component: () => import('../views/GameView.vue'),
      children : [
        {
          // Matches exactly "/game/:game"
          path: '', 
          component: () => import('../views/GamePlayView.vue'),
          name: 'game-play'
        },
        {
          // Matches exactly "/game/:game/rules"
          path: 'rules', 
          component: () => import('../views/GameRulesView.vue'),
          name: 'game-rules'
        }
      ]
    },
    // {
    //   path: '/finder',
    //   name: 'finder',
    //   component: () => import('../views/FinderView.vue'),
    //   meta: { showBottomNav: false} 
    // },
    // {
    //   path: '/room',
    //   name: 'room',
    //   component: () => import('../views/RoomView.vue'),
    //   meta: { showBottomNav: false} 
    // },
    // {
    //   path: '/inbox',
    //   name: 'inbox',
    //   component: () => import('../views/InboxView.vue')
    // },
    // {
    //   path: '/me',
    //   name: 'me',
    //   component: () => import('../views/MeView.vue')
    // },
    // {
    //   path: '/profile/:id',
    //   name: 'profile',
    //   component: () => import('../views/ProfileView.vue')
    // },
    // {
    //   path: '/chat/:id',
    //   name: 'chat',
    //   component: () => import('../views/ChatView.vue')
    // },
    // {
    //   path: '/card',
    //   name: 'card',
    //   component: () => import('../components/game/card/CardView.vue')
    // },
  ],
})

export default router
