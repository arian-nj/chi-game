import HomeView from '@/views/HomeView.vue'
import { createRouter, createWebHistory } from 'vue-router'
import { gamesData } from '@/libs/game'
import LocaleLayout from '@/views/LocaleLayout.vue'
import { getInitialLocale, i18n, persistLocale, setDocumentLocale, type AppLocale } from '@/i18n'
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // {
    //   path: '/admin',
    //   name: 'admin',
    //   component: AdminView,
    // },
    {
      path: '/',
      redirect: () => `/${getInitialLocale()}`,
    },
    {
      path: '/:locale(en|fa)',
      component: LocaleLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: '404',
          name: 'not-found',
          component: () => import('@/views/NotFoundView.vue'),
        },
        {
          path: 'changelog',
          name: 'changelog',
          component: () => import('@/views/ChangelogView.vue'),
        },
        {
          path: 'about',
          name: 'about',
          component: () => import('@/views/AboutView.vue'),
        },
        // {
        //   path: 'health',
        //   name: 'health',
        //   component: () => import('@/views/HealthView.vue'),
        // },
        // {
        //   path: 'room',
        //   name: 'room',
        //   component: () => import('@/views/NewRoomView.vue'),
        // },
        // {
        //   path: 'room/:code',
        //   component: () => import('@/views/RoomShell.vue'),
        //   children: [
        //     {
        //       path: '',
        //       name: 'room-code',
        //       component: () => import('@/views/RoomLobbyPanel.vue'),
        //     },
        //     {
        //       path: 'play',
        //       name: 'room-play',
        //       component: () => import('@/views/RoomPlayPanel.vue'),
        //     },
        //   ],
        // },
        {
          path: 'game/:game',
          name: 'game',
          component: () => import('@/views/GameView.vue'),
        },
        {
          path: ':pathMatch(.*)*',
          name: 'not-found-catchall',
          component: NotFoundView,
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: (to) => {
        const rest = Array.isArray(to.params.pathMatch)
          ? to.params.pathMatch.join('/')
          : String(to.params.pathMatch ?? '')
        return `/${getInitialLocale()}/${rest}`
      },
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

router.beforeEach((to) => {
  // if (to.name === 'admin') {
  //   return true;
  // }

  const localeParam = to.params.locale
  if (typeof localeParam !== 'string') return true

  const locale = localeParam as AppLocale
  if (locale !== i18n.global.locale.value) {
    i18n.global.locale.value = locale
    persistLocale(locale)
  } else {
    setDocumentLocale(locale)
  }

  if (to.name === 'game') {
    const inviteQuery = to.query.invite
    if (typeof inviteQuery === 'string' && inviteQuery.trim()) {
      return {
        name: 'room-code',
        params: { locale, code: inviteQuery.trim().toUpperCase() },
      }
    }

    const gameKey = to.params.game
    if (typeof gameKey !== 'string') return { name: 'home', params: { locale } }

    const game = gamesData.find(g => g.key === gameKey)
    if (!game || !game.isEnable) return { name: 'home', params: { locale } }
    return true
  }

  if (to.name === 'room-code' || to.name === 'room-play') {
    const code = to.params.code
    if (typeof code !== 'string' || !code.trim()) {
      return { name: 'room', params: { locale } }
    }
  }

  return true
})

export default router
