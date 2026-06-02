/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** Canonical site origin for og:url (e.g. https://chigame.example) */
  readonly VITE_SITE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
