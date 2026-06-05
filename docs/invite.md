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
| `CreateInviteRoom` | Host creates a room; returns `invite_code`. Does not affect other rooms the caller is in. |
| `JoinInviteRoom` | Joins by code (max 2 players). Idempotent if already in that room. |
| `GetInviteRoom` | Returns room metadata and `players` (accounts). Caller must be a member. Used by the lobby poll (every 2s). |
| `LeaveInviteRoom` | Removes caller from the room identified by `invite_code`; may delete that room (see lifecycle). |

### Request notes

- **`game_key`** on create is optional and currently sent as `""`. Stored in memory for future use (e.g. picking a game inside the room).
- **`invite_code`** on leave identifies which room to leave (required).
- **Invite codes** are 6 characters, uppercase alphanumeric (ambiguous chars like `0`/`O` excluded). Normalized to uppercase on the server.

### Connect error codes (typical)

| Code | When |
|------|------|
| `Unauthenticated` | Missing or invalid JWT |
| `NotFound` | Unknown or expired code (`GetInviteRoom`, `JoinInviteRoom`) |
| `FailedPrecondition` | Room already has 2 players |
| `PermissionDenied` | `GetInviteRoom` but caller is not in that room |
| `Internal` | Unexpected server failure |

## Room lifecycle

```mermaid
sequenceDiagram
  participant Host
  participant API as InviteService
  participant Store as InviteStore
  participant Guest

  Host->>API: CreateInviteRoom
  API->>Store: create room + host player
  API-->>Host: invite_code
  Host->>Guest: share link /en/room/CODE
  Guest->>API: JoinInviteRoom
  API->>Store: add guest player
  loop every 2s
    Host->>API: GetInviteRoom
    Guest->>API: GetInviteRoom
    API->>Store: read room state
  end
  Host->>API: LeaveInviteRoom
  API->>Store: remove player; delete room if host or empty
```

**Membership rules**

- A person may be in **multiple rooms** at once (each membership is per invite code).
- At most **two players** per room.
- **TTL**: 12 hours from create. Expired rooms are removed on next lookup.

**Deletion (hard delete today)**

On leave, the server removes the player from memory. The room is deleted if:

- no players remain, or
- the leaver is the **host**

This is intentional for the current MVP: no history, no rejoin. See [Retention](#retention) if you add soft-close later.

## In-memory store

Invite room state lives in `InviteStore` on [`APIApplication`](../backend/api/api.go) ([`backend/api/invite.go`](../backend/api/invite.go)):

| Field | Description |
|-------|-------------|
| `Code` | 6-char invite code (sole identifier) |
| `GameKey` | Optional label (often empty) |
| `HostPersonID` | Creator |
| `PlayerIDs` | Join order preserved for the lobby UI |
| `CreatedAt` / `ExpiresAt` | TTL enforcement |

**PostgreSQL** is still used for guest accounts (`persons`); `GetInviteRoom` loads player display names via `GetPersonByID`.

**Operational notes**

- All active invite rooms are **lost on server restart or deploy**.
- Designed for a **single API process** (no shared store across replicas).
- Legacy `invite_rooms` DB tables are unused; they can be dropped in a later migration.

## Code map

### Backend

| Path | Role |
|------|------|
| [`backend/api/invite.go`](../backend/api/invite.go) | In-memory `InviteStore` |
| [`backend/api/invite_rpc.go`](../backend/api/invite_rpc.go) | RPC handlers |
| [`backend/api/router.go`](../backend/api/router.go) | Registers `InviteService` |
| [`backend/api/invite_rpc_test.go`](../backend/api/invite_rpc_test.go) | Tests (needs `DATABASE_URL`) |
| [`backend/internals/random/random.go`](../backend/internals/random/random.go) | `GenerateInviteCode` |

Regenerate RPC stubs: `buf generate` from repo root.

### Frontend

| Path | Role |
|------|------|
| [`frontend/src/views/RoomView.vue`](../frontend/src/views/RoomView.vue) | Create / join UI |
| [`frontend/src/views/RoomLobby.vue`](../frontend/src/views/RoomLobby.vue) | Code, link, player list; redirects to `/play` when full |
| [`frontend/src/views/RoomPlayView.vue`](../frontend/src/views/RoomPlayView.vue) | Game picker when 2 players are present |
| [`frontend/src/libs/invite-room.ts`](../frontend/src/libs/invite-room.ts) | Connect client helpers |
| [`frontend/src/libs/room-url.ts`](../frontend/src/libs/room-url.ts) | Invite path / URL builders |
| [`frontend/src/router/router.ts`](../frontend/src/router/router.ts) | `room` and `room-code` routes |

**Route leave behavior:** `RoomView` only calls `LeaveInviteRoom` when navigating *away* from room routes (e.g. home). Navigating `/room` → `/room/:code` must **not** leave, or the room is deleted immediately after create.

## Adding behavior later

1. **New RPC or fields** — extend [`invite.proto`](../proto/invite/v1/invite.proto), `buf generate`, implement in `invite_rpc.go`.
2. **Game inside room** — set `game_key` on create or a new “start game” RPC; wire to a pluggable game module when online play exists.
3. **Live updates** — replace or complement polling with WebSocket `room.v1` messages when that layer is implemented.
