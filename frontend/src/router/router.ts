import { createRouter, createWebHistory } from 'vue-router';
import {
  getInitialLocale,
  i18n,
  persistLocale,
  type AppLocale,
} from '@/i18n';
import { supportedLocales } from '@/i18n/config';

const localeParam = supportedLocales.map((l) => l.code).join('|');
const gameKeys = 'minesweeper|tictactoe|connect4|sudoku|solitaire';

function isAppLocale(value: unknown): value is AppLocale {
  return (
    typeof value === 'string' &&
    supportedLocales.some((locale) => locale.code === value)
  );
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => `/${getInitialLocale()}`,
    },
    {
      path: `/game/:game(${gameKeys})`,
      redirect: (to) => `/${getInitialLocale()}/game/${to.params.game}`,
    },
    {
      path: `/:locale(${localeParam})`,
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../views/HomeView.vue'),
        },
        {
          path: 'game',
          children: [
            {
              path: 'minesweeper',
              name: 'minesweeper',
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
            {
              path: 'solitaire',
              name: 'solitaire',
              component: () => import('@/views/SolitaireView.vue'),
            },
          ],
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: () => `/${getInitialLocale()}`,
    },
  ],
});

router.beforeEach((to) => {
  const locale = to.params.locale;
  if (!isAppLocale(locale)) {
    return;
  }

  i18n.global.locale.value = locale;
  persistLocale(locale);
});

export default router;
