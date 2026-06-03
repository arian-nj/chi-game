const ADMIN_SECRET_KEY = 'chigame_admin_secret';

export function readAdminSecret(): string | null {
  try {
    const value = sessionStorage.getItem(ADMIN_SECRET_KEY)?.trim();
    return value || null;
  } catch {
    return null;
  }
}

export function saveAdminSecret(secret: string): void {
  try {
    sessionStorage.setItem(ADMIN_SECRET_KEY, secret);
  } catch {
    // ignore
  }
}

export function clearAdminSecret(): void {
  try {
    sessionStorage.removeItem(ADMIN_SECRET_KEY);
  } catch {
    // ignore
  }
}
