const DEVICE_ID_KEY = 'chigame_device_id';
const GUEST_TOKEN_KEY = 'chigame_guest_token';

export function readGuestDeviceId(): string | null {
  try {
    const value = localStorage.getItem(DEVICE_ID_KEY)?.trim();
    return value || null;
  } catch {
    return null;
  }
}

export function readGuestToken(): string | null {
  try {
    const value = localStorage.getItem(GUEST_TOKEN_KEY)?.trim();
    return value || null;
  } catch {
    return null;
  }
}

export function saveGuestSession(deviceId: string, token: string): void {
  try {
    localStorage.setItem(DEVICE_ID_KEY, deviceId);
    localStorage.setItem(GUEST_TOKEN_KEY, token);
  } catch {
    // Storage may be unavailable; session stays in memory for this tab only.
  }
}

export function clearGuestSession(): void {
  try {
    localStorage.removeItem(DEVICE_ID_KEY);
    localStorage.removeItem(GUEST_TOKEN_KEY);
  } catch {
    // ignore
  }
}
