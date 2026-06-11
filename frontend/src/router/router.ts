import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/game',
      name: 'game',
      children :[
        {
          path : 'minesweeper',
          name : 'minesweeper',
          component: () => import('@/views/MinesweeperView.vue'),
        },
        {
          path: 'tictactoe',
          name: 'tictactoe',
          component: () => import('@/views/TicTacToeView.vue'),
        },
      ],
    },

  ],
})

export default router
