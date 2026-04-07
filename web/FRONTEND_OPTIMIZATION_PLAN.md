# Frontend Optimization Plan

## Decisions

- Script language is unified to TypeScript. New JavaScript files are forbidden.
- Element Plus remains the component foundation, but it is treated as a vendor layer.
- Styling is split into three layers:
  - Vendor styles: `element-plus/dist/index.css`
  - Project theme and component override layer: `src/styles/*.css`
  - Page semantic styles: page-local class names only

## Phase 1

- Split theme tokens and Element Plus overrides out of `src/index.css`
- Add migration guard scripts to stop new `.js` files from entering `src/`
- Migrate low-risk shared JavaScript files to TypeScript first
- Keep the existing build working while reducing surface area of the style system

## Phase 2

- Migrate create/detail/list composables from `.js` to `.ts`
- Convert high-change pages to `<script setup lang="ts">`
- Add typed DTOs for products, notebooks, training, inference, storage, and clusters
- Remove `Record<string, any>` and `any` from business paths gradually

## Phase 3

- Introduce shared UI shells:
  - `BaseDialog`
  - `BaseFormDialog`
  - `BaseDrawer`
  - `BaseFormDrawer`
  - `BaseTableToolbar`
- Simplify route registration and keep-alive logic
- Split app shell state from feature-domain state

## Phase 4

- Add linting, formatting, and automated tests
- Keep strict `check:no-js` and `check:no-console` enabled
- Lock the current structure behind repeatable quality gates

## Phase 5

- Continue removing legacy global style spillover from `src/index.css`
- Trim remaining heavy runtime assets and recheck bundle composition
- Run manual regression with real backend menus, permissions, and business data

## Execution Log

- `2026-04-05`: Started Phase 1
  - Added `src/styles/tokens.css`
  - Added `src/styles/element-plus-overrides.css`
  - Added JS migration guard scripts
  - Migrated shared validator/constants files from JS to TS
- `2026-04-05`: Continued Phase 1
  - Added `src/types/consoleResource.ts` for shared console create-page DTOs
  - Migrated `useTrainingCreate`, `useNotebookCreate`, and `useInferenceCreate` from JS to TS
  - Reduced remaining JS allowlist entries from `42` to `39`
  - Verified `npm run build` passes
- `2026-04-05`: Continued Phase 1
  - Fixed repo-wide TypeScript blockers in `src/api/sysDictionary.ts`, `src/api/sysDictionaryDetail.ts`, `src/api/sysParams.ts`, and `src/utils/bus.ts`
  - Added `npm run typecheck`
  - Verified repo-wide `tsc --noEmit` passes
- `2026-04-05`: Continued Phase 2
  - Migrated `useTrainingList`, `useNotebookList`, and `useInferenceList` from JS to TS
  - Added shared list DTOs into `src/types/consoleResource.ts`
  - Reduced remaining JS allowlist entries from `39` to `36`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-new-js` all pass
- `2026-04-05`: Completed JS eradication baseline
  - Renamed all remaining `web/src/**/*.js` modules to `.ts`
  - Reduced frontend source JavaScript file count to `0`
  - Enabled strict `check:no-js` with an empty allowlist
  - Kept `26` legacy modules on a temporary `// @ts-nocheck` bridge to avoid blocking the repository
- `2026-04-05`: Continued Phase 2 debt paydown
  - Added shared detail DTOs into `src/types/consoleResource.ts`
  - Removed `// @ts-nocheck` from training detail/log/terminal composables
  - Removed `// @ts-nocheck` from inference detail/api-key/log/terminal composables
  - Rebuilt `useNotebookDetail.ts` as real TypeScript instead of a bridged rename
  - Reduced remaining `// @ts-nocheck` modules from `36` to `26`
  - Verified repo-wide `npm run typecheck` still passes after the detail-module migration
- `2026-04-06`: Continued Phase 2 page-layer typing
  - Converted compute detail page SFCs to `<script setup lang="ts">` for training, notebook, and inference detail entry/panel/overview components
  - Converted compute create entry SFCs and create summary components to `<script setup lang="ts">`
  - Increased typed SFC count from `0` to `14`
  - Reduced plain `<script setup>` SFC count to `233`
  - Verified `npm run typecheck`, `npm run check:no-js`, and `npm run build` still pass after the SFC migration
- `2026-04-06`: Continued Phase 2 console account and order typing
  - Added `src/types/account.ts` and `src/types/order.ts` for account security, invoice, transaction, and usage DTOs
  - Removed `// @ts-nocheck` from `useAccessLog`, `useAccountProfile`, `useAccountSecurity`, `useInvoicePage`, `useOrderOverview`, `useOrderTransactions`, and `useOrderUsage`
  - Converted account and order entry/components to `<script setup lang="ts">`
  - Reduced remaining `// @ts-nocheck` modules from `26` to `19`
  - Verified `npm run typecheck` passes after the account/order migration
- `2026-04-06`: Continued Phase 2 console shell typing
  - Added `src/types/dashboard.ts`, `src/types/image.ts`, `src/types/sshkey.ts`, and `src/types/storage.ts`
  - Removed `// @ts-nocheck` from dashboard, images, sshkeys, and storage composables
  - Converted dashboard, images, sshkeys, and storage pages plus key dialog/list/card components to `<script setup lang="ts">`
  - Fixed the missing `handleSizeChange` binding in `images/index.vue` while moving that page to typed setup
  - Reduced remaining `// @ts-nocheck` modules from `19` to `15`
  - Verified `npm run typecheck`, `npm run check:no-js`, and `npm run build` all pass
- `2026-04-06`: Continued Phase 2 app entry typing
  - Added `src/types/auth.ts`, `src/types/landing.ts`, and `src/types/layout.ts`
  - Corrected `src/api/user.ts` register payload typing to match backend JSON fields
  - Removed `// @ts-nocheck` from `useLandingPage`, `useAppLayout`, and `useLoginPage`
  - Converted landing, layout, and login entry pages to `<script setup lang="ts">`
  - Reduced remaining `// @ts-nocheck` modules from `15` to `12`
  - Increased typed SFC count from `14` to `52`
  - Reduced plain `<script setup>` SFC count from `233` to `195`
  - Verified `npm run typecheck`, `npm run check:no-js`, and `npm run build` all pass after the app-entry migration
- `2026-04-06`: Continued Phase 2 superAdmin typing
  - Added `src/types/superAdmin.ts` for API, authority, menu, operation record, user, cluster, node, and CMS product DTOs
  - Removed `// @ts-nocheck` from all remaining super admin composables, including API, authority, menu, operation record, user, cluster, node, and CMS product management
  - Corrected the user list search payload from `nickname` to backend-compatible `nickName`
  - Relaxed mismatched API signatures in `src/api/api.ts`, `src/api/authority.ts`, `src/api/menu.ts`, `src/api/user.ts`, `src/api/cluster.ts`, and `src/api/cms.ts` so the frontend contracts match real runtime usage
  - Converted super admin entry pages to `<script setup lang="ts">`
  - Reduced remaining `// @ts-nocheck` modules from `12` to `0`
  - Increased typed SFC count from `52` to `61`
  - Reduced plain `<script setup>` SFC count from `195` to `186`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` all pass after the super admin migration
- `2026-04-06`: Started Phase 3 shared dialog shell consolidation
  - Added `src/components/base/BaseDialog.vue` and `src/components/base/BaseFormDialog.vue`
  - Migrated representative dialogs in cluster, user, storage, ssh key, account, and order flows onto the shared shells
  - Reduced repeated `el-dialog` visibility, close, footer, and form-validation boilerplate across `11` business dialog components
  - Increased typed SFC count from `61` to `70`
  - Reduced plain `<script setup>` SFC count from `186` to `179`
  - Established `12` current references to the new shared dialog layer
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the dialog-shell migration
- `2026-04-06`: Continued Phase 3 shared drawer and permission-panel consolidation
  - Added `src/components/base/BaseDrawer.vue` and `src/components/base/BaseFormDrawer.vue`
  - Migrated all remaining business drawers in API and authority flows onto the shared drawer shell, leaving `BaseDrawer.vue` as the only direct `<el-drawer>` entry in `src/`
  - Converted API drawer hosts plus the authority permission panels (`apis.vue`, `datas.vue`, `menus.vue`) to `<script setup lang="ts">`
  - Added shared `drawer-shell` theme constraints in `src/styles/element-plus-overrides.css` so drawer sizing stays in the project override layer instead of per-page hacks
  - Increased typed SFC count from `70` to `83`
  - Reduced plain `<script setup>` SFC count from `179` to `168`
  - Established `17` business-component references to the shared dialog/drawer shell layer
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the drawer-shell migration
- `2026-04-06`: Continued Phase 3 list toolbar and filter-bar consolidation
  - Added `src/components/listPage/BaseTableToolbar.vue` and `src/components/listPage/BaseFilterBar.vue`
  - Converted `src/components/listPage/PageIntro.vue` and `src/components/RefreshButton/index.vue` to `<script setup lang="ts">`
  - Migrated representative page headers and filter bars across API, authority, user, storage, images, ssh keys, product, cluster, node, training, notebook, inference, invoice, and usage management onto the shared toolbar/filter shells
  - Reduced repeated `PageIntro + RefreshButton + action buttons` and `console-filter-card + list-filter-bar` wrappers across `18` business entry/components
  - Increased typed SFC count from `83` to `97`
  - Reduced plain `<script setup>` SFC count from `168` to `156`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the list-shell migration
- `2026-04-06`: Continued Phase 3 list-shell completion
  - Migrated the remaining `9` direct business `PageIntro` entry pages to `BaseTableToolbar`, so business-page direct `PageIntro` imports are now `0`
  - Extended `src/components/listPage/BaseFilterBar.vue` with a `plain` mode and `actionsClass` so the same shell can be reused inside table cards without reintroducing vendor/global style hacks
  - Migrated the remaining `8` table-card inline filter bars in training, notebook, inference, API, authority, menu, operation record, and user management to `BaseFilterBar`
  - Converted `TransactionsHeaderActions.vue` and those `8` table-card components to `<script setup lang="ts">`
  - Increased typed SFC count from `97` to `106`
  - Reduced plain `<script setup>` SFC count from `156` to `147`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the shared list-shell completion pass
- `2026-04-06`: Continued Phase 2 wrapper and shared-shell typing
  - Converted remaining wrapper pages and shared shells in `createPage`, `detailPage`, `listPage`, `orderPage`, account, order, landing, layout, login, storage, and super admin feature slices to `<script setup lang="ts">`
  - Cleared the remaining plain `<script setup>` debt in `account`, `order`, and `superAdmin`
  - Increased typed SFC count from `106` to `192`
  - Reduced plain `<script setup>` SFC count from `147` to `61`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the wrapper/shared-shell migration
- `2026-04-06`: Completed Phase 2 script-language unification
  - Converted the remaining `29` console components in `training`, `notebooks`, and `inference` plus the remaining `32` shared/app/error components to `<script setup lang="ts">`
  - Increased typed SFC count from `192` to `253`
  - Reduced plain `<script setup>` SFC count from `61` to `0`
  - Confirmed `web/src/**/*.js` remains `0`, `BaseTableToolbar` is used in `21` business files, `BaseFilterBar` is used in `14` business files, and business-page direct `PageIntro` imports remain `0`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` all pass after the final SFC migration pass
- `2026-04-06`: Continued Phase 2 type-quality tightening
  - Exported create-page domain types from the training, notebook, and inference composables so child components can reuse real form/filter/update payload contracts instead of `Object` / `Array` / `Function`
  - Tightened the remaining create/detail/dialog child components onto explicit `defineProps<T>()` and typed `defineEmits<T>()`
  - Reduced runtime prop declarations from `29` `type: Object`, `44` `type: Array`, and `43` `type: Function` to `0 / 0 / 0`
  - Expanded `src/types/consoleResource.ts` with shared product fields such as `clusterName`, `nodeType`, `systemDisk`, `cudaVersion`, and `vGpuCores`
  - Verified `npm run typecheck` still passes after the prop-typing cleanup
- `2026-04-06`: Continued Phase 3 style-layer consolidation
  - Reordered frontend style imports to `vendor -> tokens -> element-plus-overrides -> index.css`, so the shared theme override layer now precedes page/global semantic styles
  - Moved the remaining centralized Element Plus form/button/tab/select overrides out of `src/index.css` into `src/styles/element-plus-overrides.css`
  - Reduced `src/index.css` size from `68,817` bytes to `54,648` bytes
  - Left only workspace-semantic `.workspace-shell ... .el-*` bridge rules in `src/index.css`; the direct vendor override block is now centralized in `src/styles/element-plus-overrides.css`
  - Verified `npm run typecheck` and `npm run build` both still pass after the style-layer refactor
- `2026-04-06`: Continued Phase 3 bundle-splitting cleanup
  - Kept `OverviewTrendChart.vue` on tree-shaken `echarts/core` imports and moved the chart itself behind `defineAsyncComponent`, so the overview page now loads the chart chunk on demand instead of baking it into the page chunk
  - Restored a dedicated `vendor-echarts` manual chunk after confirming namespace-style dynamic imports defeated ECharts tree-shaking
  - Added `unplugin-vue-components` with `ElementPlusResolver({ importStyle: false })` and removed full `app.use(ElementPlus)` registration in `src/main.ts`
  - Explicitly registered `v-loading` at the app entry so the project no longer depends on Element Plus plugin-wide directive injection
  - Reduced `vendor-element-plus` from `851.16 kB` to `484.63 kB` and cleared the production chunk-size warning
  - Verified `npm run build`, `npm run typecheck`, and `npm run check:no-js` all pass after the bundle refactor
- `2026-04-06`: Continued Phase 3 icon registration cleanup
  - Removed the global `@element-plus/icons-vue` registration loop from `src/core/global.ts`
  - Added the only missing explicit icon imports (`Sunny`, `Moon`) in `LandingNavbar.vue` after scanning the repo for files that still depended on global icon registration
  - Kept SVG icon registration intact for business icons and the menu-icon catalog flow
  - Verified `npm run build`, `npm run typecheck`, and `npm run check:no-js` still pass after removing the global icon bridge
- `2026-04-06`: Continued Phase 3 icon runtime hardening
  - Reworked `src/components/AppIcon.vue` to resolve business SVG icons first and lazily load `@element-plus/icons-vue` only when a saved Element Plus icon name is actually requested
  - Added `src/generated/elementPlusIconNames.ts` from the package declaration files so the admin icon catalog no longer eagerly imports the full runtime icon module just to enumerate names
  - Updated `useMenuIconOptions.ts` to consume the generated static icon-name list instead of `import * as ElIconModules`
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the icon fallback fix
- `2026-04-06`: Continued Phase 3 dialog-shell completion
  - Migrated the remaining business dialogs in MFA setup/disable, image management, notebook SSH login, authority button assignment, and CMS product management onto `BaseDialog` / `BaseFormDialog`
  - Reduced direct business `<el-dialog>` usage in `src/` from `8` to `0`, leaving `BaseDialog.vue` as the only Element Plus dialog entry in frontend source
  - Kept the existing product-dialog semantic class styling intact by running those dialogs through `BaseDialog` with `shell=false`
  - Renamed the remaining generic `custom-dialog` business classes into semantic dialog classes so account/profile, account/security, invoice, recharge, and storage dialog styling no longer share an ambiguous global selector
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the final dialog-shell migration pass
- `2026-04-06`: Continued Phase 3 style-layer trimming and workspace cleanup
  - Moved the remaining list-filter, pagination, and workspace-shell Element Plus bridge selectors out of `src/index.css` into `src/styles/element-plus-overrides.css`
  - Reduced `src/index.css` size from `54,648` bytes to `51,198` bytes while keeping the business semantic layer intact
  - Cleaned high-confidence generated clutter: frontend `dist/`, temporary helm package staging under `.tmp/`, local `__pycache__/`, and stale server log output
  - Verified `npm run typecheck`, `npm run build`, and `npm run check:no-js` still pass after the style cleanup, then removed regenerated `dist/` again so the workspace stays clean
- `2026-04-06`: Started Phase 4 quality guard foundation
  - Added `vitest` with `npm run test:unit` and created focused unit tests for layout menu grouping and dynamic hidden-route permission gating
  - Extracted menu grouping and hidden-route registration into pure helper modules so they can be verified without booting the app shell
  - Added `npm run check:no-console` to stop new `console.log` calls from entering `src/`
  - Removed dead startup banner logging from `src/core/config.ts` and changed image compression fallback to return the original file instead of emitting noisy logs
  - Verified the new tests and quality guards pass alongside `npm run typecheck`, `npm run build`, and `npm run check:no-js`
- `2026-04-06`: Completed Phase 4 quality gate baseline
  - Added repository-wide `eslint` and `prettier` enforcement with `npm run lint` and `npm run format:check`
  - Verified `lint`, `format:check`, `typecheck`, `test:unit`, `check:no-console`, `check:no-js`, and `build` all pass together
  - Closed the original Phase 4 goal; remaining work moved to manual regression and ongoing debt paydown
- `2026-04-06`: Continued Phase 5 runtime asset cleanup
  - Replaced all frontend runtime usages of the 1.49 MB `src/assets/logo.png` with the shared SVG `src/components/logo/index.vue`
  - Updated the shared logo component so size is controlled by the caller and multiple logo instances no longer collide on SVG gradient ids
  - Removed the unused `src/assets/logo.png` runtime asset without changing the surrounding landing, layout, or login copy/layout structure

## Current Debt

- `0` `.vue` single-file components still use plain `<script setup>`; the script language is now fully unified on `<script setup lang="ts">`
- `0` frontend source `.js` files remain under `web/src/`
- `0` business components under `web/src/` now use direct `<el-dialog>`; all dialog entry points are centralized behind `src/components/base/BaseDialog.vue`
- `0` remaining `type: Object` / `type: Array` / `type: Function` runtime prop declarations remain in `.vue` components; the next TypeScript quality step is reducing broad index signatures and continuing to replace implicit stringly-typed contracts with domain interfaces
- Phase 4 quality gates are now in place and passing locally: `lint`, `format:check`, `typecheck`, `test:unit`, `check:no-console`, `check:no-js`, and `build`
- The styling architecture is in the target three-layer shape, but there is still follow-up debt to keep shrinking legacy global rules in `src/index.css` (`51,198` bytes) and continue moving non-semantic shared overrides into `src/styles/tokens.css` / `src/styles/element-plus-overrides.css`
- Production build no longer reports chunk-size warnings, but bundle follow-up is still worth tracking around the largest async/vendor chunks after each UI dependency change
- Manual regression with real backend menus, role-based permissions, and business data is still pending; the dynamic permission/menu path has helper-level unit coverage, but it has not yet been revalidated end-to-end against a live backend payload in this optimization pass
