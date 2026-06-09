/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Canonical site origin for og:url (e.g. https://chigame.example) */
  readonly VITE_SITE_URL?: string;
  /** Connect RPC API origin (e.g. http://127.0.0.1:8383). Omit in dev to use Vite proxy. */
  readonly VITE_API_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
