# General TODO

Tracks cross-cutting work that affects the whole monorepo or shared libraries.

## Checklist — Before Merging Any PR

This is the **single canonical checklist**. CLAUDE.md, INDEX.md, and per-app TODO files all point here.

### Docs
- [ ] Shared docs updated (`docs/`) if shared lib behavior changed
- [ ] App docs updated (`DOCS.md`) if routes, stores, views, or components changed
- [ ] Locale strings added/updated in `locale/en/corteza-webapp-{app}/`
- [ ] New doc file added to `docs/INDEX.md`

### Code quality
- [ ] No hardcoded UI strings — all text uses `$t()` with locale YAML files
- [ ] Tailwind classes preferred over custom `<style>` blocks
- [ ] Naming conventions followed (see [code-style.md](code-style.md#naming-conventions))
- [ ] Lint passes (`make lint` or `pnpm --filter <app> lint`)

### Shared lib changes (if applicable)
- [ ] New components/composables exported from `lib/vue/src/index.ts`
- [ ] New PrimeVue globals added to `lib/vue/src/plugins/primevue-components.ts`

## Shared Lib (`lib/vue`)

- [ ] Add more field editors as needed (e.g. File/Attachment, Url, Email, RichText)
- [ ] Add field editor tests
- [ ] Document `useResourceList` options parameter fully

## Shared Lib (`lib/js`)

- [ ] Audit and document exported API client methods
- [ ] Add TypeScript types for all model classes

## Infrastructure

- [ ] CI pipeline for automated lint + test on PRs
- [ ] Vitest coverage thresholds
