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
        {
          path: 'connect4',
          name: 'connect4',
          component: () => import('@/views/Connect4View.vue'),
        },
        {
          path: 'sudoku',
          name: 'sudoku',
          component: () => import('@/views/SudokuView.vue'),
        },
      ],
    },

  ],
})

export default router
