/** Invite link path: /{locale}/room/{code} */
export function roomInvitePath(locale: string, inviteCode: string): string {
  const code = inviteCode.trim().toUpperCase();
  return `/${locale}/room/${code}`;
}

export function roomInviteUrl(locale: string, inviteCode: string): string {
  if (typeof window === 'undefined') {
    return roomInvitePath(locale, inviteCode);
  }
  return `${window.location.origin}${roomInvitePath(locale, inviteCode)}`;
}

export function roomHomePath(locale: string): string {
  return `/${locale}/room`;
}

/** Lobby with invite code: /{locale}/room/{code} */
export function roomLobbyPath(locale: string, inviteCode: string): string {
  return roomInvitePath(locale, inviteCode);
}

/** Pick a game when both players are in the room: /{locale}/room/{code}/play */
export function roomPlayPath(locale: string, inviteCode: string): string {
  const code = inviteCode.trim().toUpperCase();
  return `/${locale}/room/${code}/play`;
}
