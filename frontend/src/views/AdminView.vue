<script setup lang="ts">
import { AdminService, type GetOverviewResponse } from '@/gen/admin/v1/admin_pb';
import { getInitialLocale } from '@/i18n';
import { createApiClient } from '@/libs/api-client';
import {
  clearAdminSecret,
  readAdminSecret,
  saveAdminSecret,
} from '@/libs/admin-secret-storage';
import { Code, ConnectError } from '@connectrpc/connect';
import { computed, onMounted, ref } from 'vue';

const homePath = computed(() => `/${getInitialLocale()}`);

const secretInput = ref('');
const savedSecret = ref<string | null>(readAdminSecret());
const overview = ref<GetOverviewResponse | null>(null);
const errorMessage = ref<string | null>(null);
const isLoading = ref(false);

function connectErrorMessage(err: unknown): string {
  if (err instanceof ConnectError) {
    if (err.code === Code.PermissionDenied) {
      return 'Invalid admin secret.';
    }
    if (err.code === Code.Unavailable) {
      return 'Admin is not enabled on the server (set ADMIN_SECRET).';
    }
    return err.message;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return 'Something went wrong.';
}

async function loadOverview(secret: string) {
  isLoading.value = true;
  errorMessage.value = null;
  overview.value = null;

  try {
    const client = createApiClient(AdminService);
    overview.value = await client.getOverview({ adminSecret: secret });
    savedSecret.value = secret;
    saveAdminSecret(secret);
  } catch (err) {
    errorMessage.value = connectErrorMessage(err);
    if (err instanceof ConnectError && err.code === Code.PermissionDenied) {
      clearAdminSecret();
      savedSecret.value = null;
    }
  } finally {
    isLoading.value = false;
  }
}

function onSubmit() {
  const secret = secretInput.value.trim();
  if (!secret) {
    return;
  }
  loadOverview(secret);
}

function onSignOut() {
  clearAdminSecret();
  savedSecret.value = null;
  secretInput.value = '';
  overview.value = null;
  errorMessage.value = null;
}

onMounted(() => {
  if (savedSecret.value) {
    loadOverview(savedSecret.value);
  }
});
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center p-6 pt-16" lang="en" dir="ltr">
    <section
      class="w-full max-w-3xl mx-auto p-6 bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md"
    >
      <div class="flex items-center justify-between gap-4 mb-6">
        <h1
          class="text-3xl font-extrabold tracking-wide bg-linear-to-r from-white to-blue-200 bg-clip-text text-transparent"
        >
          Admin
        </h1>
        <RouterLink
          :to="homePath"
          class="text-sm font-semibold text-blue-100 hover:text-white underline underline-offset-4 shrink-0"
        >
          Back to site
        </RouterLink>
      </div>

      <p class="text-blue-100 mb-6 leading-relaxed">
        Basic site overview. Requires the server <code class="font-mono text-sm">ADMIN_SECRET</code>.
      </p>

      <form
        v-if="!overview"
        class="space-y-4"
        @submit.prevent="onSubmit"
      >
        <label class="block">
          <span class="text-sm font-semibold text-white/80">Admin secret</span>
          <input
            v-model="secretInput"
            type="password"
            autocomplete="current-password"
            class="mt-2 w-full rounded-xl border border-white/15 bg-custom-deep-blue/40 px-4 py-3 text-white placeholder:text-white/40 focus:outline-none focus:ring-2 focus:ring-blue-400"
            placeholder="Enter ADMIN_SECRET"
          />
        </label>
        <button
          type="submit"
          class="w-full rounded-xl bg-white/10 hover:bg-white/15 border border-white/15 px-4 py-3 text-white font-semibold transition disabled:opacity-50"
          :disabled="isLoading || !secretInput.trim()"
        >
          {{ isLoading ? 'Loading…' : 'View overview' }}
        </button>
      </form>

      <p
        v-if="errorMessage"
        class="rounded-xl bg-rose-950/40 border border-rose-400/30 p-4 text-rose-100 text-sm mb-4"
      >
        {{ errorMessage }}
      </p>

      <div v-if="overview" class="space-y-4">
        <dl class="grid gap-4 sm:grid-cols-2 text-blue-100">
          <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
            <dt class="text-sm font-semibold text-white/80">Total users</dt>
            <dd class="text-2xl font-bold text-white">{{ overview.totalPersons.toString() }}</dd>
          </div>
          <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
            <dt class="text-sm font-semibold text-white/80">Guest users</dt>
            <dd class="text-2xl font-bold text-white">{{ overview.guestPersons.toString() }}</dd>
          </div>
        </dl>

        <div class="rounded-xl bg-custom-deep-blue/30 border border-white/10 p-4">
          <h2 class="text-sm font-semibold text-white/80 mb-3">Recent users</h2>
          <ul
            v-if="overview.recentPersons.length"
            class="space-y-2 text-sm text-blue-100"
          >
            <li
              v-for="person in overview.recentPersons"
              :key="String(person.id)"
              class="flex flex-wrap items-baseline justify-between gap-2 border-b border-white/5 pb-2 last:border-0 last:pb-0"
            >
              <span class="font-mono text-white">@{{ person.username }}</span>
              <span class="text-white/60">
                {{ person.isGuest ? 'Guest' : 'Registered' }}
                <span v-if="person.createdAt"> · {{ person.createdAt }}</span>
              </span>
            </li>
          </ul>
          <p v-else class="text-sm text-white/60">No users yet.</p>
        </div>

        <div class="flex gap-3">
          <button
            type="button"
            class="flex-1 rounded-xl bg-white/10 hover:bg-white/15 border border-white/15 px-4 py-3 text-white font-semibold transition disabled:opacity-50"
            :disabled="isLoading"
            @click="savedSecret && loadOverview(savedSecret)"
          >
            {{ isLoading ? 'Loading…' : 'Refresh' }}
          </button>
          <button
            type="button"
            class="rounded-xl border border-white/15 px-4 py-3 text-blue-100 hover:text-white font-semibold transition"
            @click="onSignOut"
          >
            Sign out
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
