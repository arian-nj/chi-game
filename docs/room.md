# Rooms

Rooms let two guests play together online. A host creates a room, shares a short code or link, and a second player joins. The UI lists player display names; when both are present they can pick games (websocket and live games are planned).

Room UX lives under `/room`.

## URLs (frontend)

| Route | Purpose |
|-------|---------|
| `/{locale}/room` | Create a room or enter a code to join |
| `/{locale}/room/{code}` | Lobby: share code, wait for guest |
| `/{locale}/room/{code}/play` | Both players joined — pick a game |

Examples:

- `https://chigame.site/en/room`
- `https://chigame.site/en/room/TYJQ6X`
- `https://chigame.site/en/room/TYJQ6X/play`

## Auth

All `RoomService` RPCs require a **guest JWT** (`Authorization: Bearer …`). See [`docs/auth.md`](auth.md). The frontend attaches the token via [`frontend/src/libs/api-client.ts`](../frontend/src/libs/api-client.ts).

## API (`room.v1.RoomService`)

Proto: [`proto/room/v1/room.proto`](../proto/room/v1/room.proto)

| RPC | Description |
|-----|-------------|
| `CreateRoom` | Host creates a room; returns `code`. Does not affect other rooms the caller is in. |
| `JoinRoom` | Joins by code (max 2 players). Idempotent if already in that room. |
| `GetRoom` | Returns room metadata and `players`. Caller must be a member. Lobby polls every 2s. |
| `LeaveRoom` | Removes caller from the room; may delete that room (see lifecycle). |

### Request notes

- **`game_key`** on create is optional (often `""`). Hint for a suggested first game.
- **`code`** — 6 characters, uppercase alphanumeric. Normalized to uppercase on the server.

### Connect error codes (typical)

| Code | When |
|------|------|
| `Unauthenticated` | Missing or invalid JWT |
| `NotFound` | Unknown or expired code |
| `FailedPrecondition` | Room already has 2 players |
| `PermissionDenied` | `GetRoom` but caller is not in that room |
| `Internal` | Unexpected server failure |

## Room lifecycle

**Membership rules**

- A person may be in **multiple rooms** at once.
- At most **two players** per room.
- **TTL**: 12 hours from create. Expired rooms are removed on next lookup.

**Deletion**

On leave, the room is deleted if no players remain or the leaver is the **host**.

## In-memory store

Room state lives in `RoomsStore` on [`APIApplication`](../backend/api/api.go) ([`backend/api/room.go`](../backend/api/room.go)):

| Field | Description |
|-------|-------------|
| `Code` | 6-char room code (public identifier) |
| `GameKey` | Optional label |
| `HostPersonID` | Creator |
| `PlayerIDs` | Join order for lobby UI |
| `CreatedAt` / `ExpiresAt` | TTL enforcement |
| `MsgChnl` | Event channel for websocket (future) |

PostgreSQL `persons` is used for display names. The `rooms` / `room_players` tables exist for future persistence but are not used by the in-memory store yet.

**Operational notes**

- Active rooms are **lost on server restart**.
- Designed for a **single API process**.

## Code map

### Backend

| Path | Role |
|------|------|
| [`backend/api/room.go`](../backend/api/room.go) | In-memory `RoomsStore` |
| [`backend/api/room_rpc.go`](../backend/api/room_rpc.go) | RPC handlers |
| [`backend/api/room_socket.go`](../backend/api/room_socket.go) | WebSocket skeleton |
| [`backend/api/router.go`](../backend/api/router.go) | Registers `RoomService` |

### Frontend

| Path | Role |
|------|------|
| [`frontend/src/views/NewRoomView.vue`](../frontend/src/views/NewRoomView.vue) | Create / join UI |
| [`frontend/src/views/RoomLobby.vue`](../frontend/src/views/RoomLobby.vue) | Code, link, player list |
| [`frontend/src/views/RoomView.vue`](../frontend/src/views/RoomView.vue) | Game picker |
| [`frontend/src/libs/room-api.ts`](../frontend/src/libs/room-api.ts) | Connect client helpers |
| [`frontend/src/libs/room-url.ts`](../frontend/src/libs/room-url.ts) | Path / URL builders |
