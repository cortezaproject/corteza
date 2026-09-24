<template>
  <figure class="dg perm-order" aria-label="How Human decides whether an operation is allowed">
    <ol class="steps">
      <li class="step">
        <span class="n">1</span>
        <div>
          <span class="title">Bypass role</span>
          <span class="sub">
            A member of a bypass role is allowed, and nothing else is checked.
          </span>
        </div>
        <span class="verdict allow">Allow</span>
      </li>
      <li class="step">
        <span class="n">2</span>
        <div>
          <span class="title">The user's roles, level by level</span>
          <div class="levels">
            <span v-for="l in levels" :key="l" class="level">{{ l }}</span>
          </div>
          <span class="sub">
            The first level where one of the user's roles has a matching rule decides. Within that
            level, Deny beats Allow.
          </span>
        </div>
        <span class="verdict both">Allow or Deny</span>
      </li>
      <li class="step">
        <span class="n">3</span>
        <div>
          <span class="title">No matching rule anywhere</span>
          <span class="sub">Anything not allowed is denied.</span>
        </div>
        <span class="verdict deny">Deny</span>
      </li>
    </ol>
  </figure>
</template>

<script setup>
const levels = ['Contextual roles', 'Ordinary roles', 'Authenticated', 'Anonymous']
</script>

<style scoped>
.dg {
  margin: 24px 0;
}

.steps {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 8px;
}

.step {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 14px;
  padding: 12px 16px;
  border-radius: 10px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
  margin: 0 !important;
}

.n {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--vp-c-brand-1);
  background: var(--vp-c-brand-soft);
}

.title {
  display: block;
  font-family: var(--vp-font-family-heading);
  font-weight: 500;
  font-size: 0.9375rem;
  color: var(--vp-c-text-1);
}

.sub {
  display: block;
  font-size: 0.8125rem;
  line-height: 1.5;
  color: var(--vp-c-text-2);
}

.levels {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 0;
  margin: 6px 0;
}

.level {
  font-size: 0.8125rem;
  padding: 2px 10px;
  border-radius: 999px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  color: var(--vp-c-text-1);
}

.level + .level::before {
  content: '→';
  margin: 0 8px 0 -2px;
  color: var(--vp-c-text-3);
}

.level + .level {
  border: none;
  background: none;
  padding-left: 0;
}

.verdict {
  font-size: 0.8125rem;
  font-weight: 500;
  padding: 2px 10px;
  border-radius: 999px;
  white-space: nowrap;
}

.allow {
  color: #2e7d64;
  background: rgba(67, 170, 139, 0.14);
}

.deny {
  color: #c2361c;
  background: rgba(229, 65, 34, 0.12);
}

.both {
  color: var(--vp-c-text-2);
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
}

@media (max-width: 560px) {
  .step {
    grid-template-columns: auto 1fr;
  }

  .verdict {
    grid-column: 2;
    justify-self: start;
  }
}
</style>

<style>
.dark .perm-order .allow {
  color: #6fd1b1;
}

.dark .perm-order .deny {
  color: #f08a76;
}
</style>
