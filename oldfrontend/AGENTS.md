## project
this is chi game a website that lets users play diffrent games
even if server is down offline games should be playable
Wrtie the code simple as possible don't overdo things

## Backend API
- Connect RPC to Go backend (`backend/`, port **8383**)
- Generated clients: `src/gen/` (`buf generate` from repo root)
- **TanStack Vue Query** + `@connectrpc/connect`; use `createApiClient` in `src/libs/api-client.ts` with `useQuery` for reads
- Shared transport: `src/libs/api-client.ts`; QueryClient: `src/libs/vue-query.ts`
- Base URL: `src/libs/api-base-url.ts` and `VITE_API_BASE_URL`
- Dev: Vite proxies `*.v1.*` RPC paths to `http://127.0.0.1:8383` when `VITE_API_BASE_URL` is unset
- Run API: `./mash.sh run-back` from repo root
- Guest auth: device ID + JWT via `useGuestAuth()` — see [`docs/auth.md`](../docs/auth.md)
- Online rooms: `/{locale}/room` — see [`docs/room.md`](../docs/room.md)

## Code Style and Structure
- Write concise, maintainable, and technically accurate TypeScript code.
- Use descriptive variable names with auxiliary verbs (e.g., isLoading, hasError).
- Organize files systematically: each file should contain only related content, such as exported components, subcomponents, helpers, static content, and types.

Naming Conventions
- Use lowercase with dashes for directories (e.g., components/auth-wizard).
- Favor named exports for functions.

TypeScript Usage
- Use TypeScript for all code.

Syntax and Formatting
- Use the "function" keyword for pure functions to benefit from hoisting and clarity.
- Always use the Vue Composition API script setup style.

UI and Styling
- Use Tailwind for components and styling.
- Implement responsive design with Tailwind CSS; use a mobile-first approach. but big monitors are important too.

Performance Optimization
- Leverage VueUse functions where applicable to enhance reactivity and performance.
- Wrap asynchronous components in Suspense with a fallback UI.
- Use dynamic loading for non-critical components.
- Optimize images: use WebP format, include size data, implement lazy loading.
- Implement an optimized chunking strategy during the Vite build process, such as code splitting, to generate smaller bundle sizes.

Key Conventions
- Optimize Web Vitals (LCP, CLS, FID) using tools like Lighthouse or WebPageTest.

- Web Page size is important specially initial loads
- SEO is very important write search engine (crawler) friendly code
- alt text on any images/icons you add later.