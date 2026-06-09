/** Connect RPC base URL. Empty string = same origin (Vite dev proxy in development). */
export function getApiBaseUrl(): string {
  const configured = import.meta.env.VITE_API_BASE_URL?.trim();
  if (!configured) {
    return 'http://127.0.0.1:8383';
  }
  return configured.replace(/\/$/, '');
}
