# Invite rooms

Invite rooms let two guests meet in a shared lobby before any online game runs. A host creates a room, shares a short code or link, and a second player joins. The UI lists player display names only (no moves, chat, or matchmaking).

Rooms are **not** tied to `/game/:game` routes. All invite UX lives under `/room`.

## URLs (frontend)

| Route | Purpose |
|-------|---------|
| `/{locale}/room` | Create a room or enter a code to join |
| `/{locale}/room/{code}` | Invite link (lobby: share code, wait for guest) |
| `/{locale}/room/{code}/play` | Both players joined — pick a game (opens offline game page; no online match yet) |

Examples:

- `https://chigame.site/en/room`
- `https://chigame.site/en/room/TYJQ6X`
- `https://chigame.site/en/room/TYJQ6X/play` (auto-navigated when 2 players are in the room)

## Auth

All `InviteService` RPCs require a **guest JWT** (`Authorization: Bearer …`). See [`docs/auth.md`](auth.md) for how device IDs, tokens, and `ValidateGuest` work. The frontend attaches the token automatically via [`frontend/src/libs/api-client.ts`](../frontend/src/libs/api-client.ts).

## API (`invite.v1.InviteService`)

Proto: [`proto/invite/v1/invite.proto`](../proto/invite/v1/invite.proto)

| RPC | Description |
|-----|-------------|
| `CreateInviteRoom` | Host creates a room; returns `invite_code` and `room_id`. Clears any prior room membership for the caller first. |
| `JoinInviteRoom` | Guest joins by code (max 2 players). Idempotent if already in that room. Leaving a *different* room happens before join. |
| `GetInviteRoom` | Returns room metadata and `players` (accounts). Caller must be a member. Used by the lobby poll (every 2s). |
| `LeaveInviteRoom` | Removes caller from the room; may delete the room (see lifecycle). |

### Request notes

- **`game_key`** on create is optional and currently sent as `""`. It is stored on the row for future use (e.g. picking a game inside the room).
- **Invite codes** are 6 characters, uppercase alphanumeric (ambiguous chars like `0`/`O` excluded). Normalized to uppercase on the server.

### Connect error codes (typical)

| Code | When |
|------|------|
| `Unauthenticated` | Missing or invalid JWT |
| `NotFound` | Unknown or expired code (`GetInviteRoom`, `JoinInviteRoom`) |
| `FailedPrecondition` | Room already has 2 players |
| `PermissionDenied` | `GetInviteRoom` but caller is not in that room |
| `Internal` | Database or unexpected failure |

## Room lifecycle

```mermaid
sequenceDiagram
  participant Host
  participant API as InviteService
  participant DB as PostgreSQL
  participant Guest

  Host->>API: CreateInviteRoom
  API->>DB: insert invite_rooms + host player
  API-->>Host: invite_code
  Host->>Guest: share link /en/room/CODE
  Guest->>API: JoinInviteRoom
  API->>DB: insert guest player
  loop every 2s
    Host->>API: GetInviteRoom
    Guest->>API: GetInviteRoom
  end
  Host->>API: LeaveInviteRoom
  API->>DB: delete player; delete room if host or empty
```

**Membership rules**

- At most **one active room per person** (non-expired row in `invite_room_players`).
- At most **two players** per room.
- **TTL**: `expires_at` is 24 hours from create. Expired rooms are invisible to queries (`expires_at > NOW()`).

**Deletion (hard delete today)**

On leave, the server removes the player row. The room row is deleted if:

- no players remain, or
- the leaver is the **host**

This is intentional for the current MVP: no history, no rejoin. See [Retention](#retention) if you add soft-close later.

**Create / join and prior rooms**

`CreateInviteRoom` and `JoinInviteRoom` (when switching codes) call `clearPersonFromInviteRoom` so callers are not blocked by a stale membership.

## Database

Migration: [`backend/db/migrations/20260604120000_create_invite_rooms.sql`](../backend/db/migrations/20260604120000_create_invite_rooms.sql)

Queries: [`backend/db/query/invite_rooms.sql`](../backend/db/query/invite_rooms.sql)

### `invite_rooms`

| Column | Description |
|--------|-------------|
| `id` | Primary key |
| `invite_code` | Unique, 6-char code |
| `game_key` | Optional label (often empty) |
| `host_person_id` | Creator |
| `created_at` | Set on insert |
| `expires_at` | Room hidden after this time |

### `invite_room_players`

| Column | Description |
|--------|-------------|
| `room_id`, `person_id` | Composite primary key |
| `joined_at` | Order used when listing players |

`ON DELETE CASCADE` from `invite_rooms` removes players when a room is deleted.

## Code map

### Backend

| Path | Role |
|------|------|
| [`backend/api/invite_rpc.go`](../backend/api/invite_rpc.go) | RPC handlers + `clearPersonFromInviteRoom` |
| [`backend/api/router.go`](../backend/api/router.go) | Registers `InviteService` |
| [`backend/api/invite_rpc_test.go`](../backend/api/invite_rpc_test.go) | Tests (needs `DATABASE_URL`) |
| [`backend/internals/random/random.go`](../backend/internals/random/random.go) | `GenerateInviteCode` |

Regenerate RPC stubs: `buf generate` from repo root.

### Frontend

| Path | Role |
|------|------|
| [`frontend/src/views/RoomView.vue`](../frontend/src/views/RoomView.vue) | Create / join UI and lobby shell |
| [`frontend/src/views/RoomPlayView.vue`](../frontend/src/views/RoomPlayView.vue) | Game picker when 2 players are present |
| [`frontend/src/components/invite/InviteRoomLobby.vue`](../frontend/src/components/invite/InviteRoomLobby.vue) | Code, link, player list; redirects to `/play` when full |
| [`frontend/src/libs/invite-room.ts`](../frontend/src/libs/invite-room.ts) | Connect client helpers |
| [`frontend/src/libs/room-url.ts`](../frontend/src/libs/room-url.ts) | Invite path / URL builders |
| [`frontend/src/router/router.ts`](../frontend/src/router/router.ts) | `room` and `room-code` routes |

**Route leave behavior:** `RoomView` only calls `LeaveInviteRoom` when navigating *away* from room routes (e.g. home). Navigating `/room` → `/room/:code` must **not** leave, or the room is deleted immediately after create.

## Adding behavior later

1. **New RPC or fields** — extend [`invite.proto`](../proto/invite/v1/invite.proto), `buf generate`, implement in `invite_rpc.go`.
2. **Game inside room** — set `game_key` on create or a new “start game” RPC; wire to a pluggable game module when online play exists.
3. **Live updates** — replace or complement polling with WebSocket `room.v1` messages when that layer is implemented.
