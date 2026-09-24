<template>
  <figure class="dg taq-flow">
    <svg
      viewBox="0 0 420 470"
      role="img"
      aria-label="A TAQ: a trigger starts it, a step runs, a branch sends it down one of two paths, one of which loops over records, and each path ends"
    >
      <defs>
        <marker
          id="taq-arrow"
          viewBox="0 0 10 10"
          refX="9"
          refY="5"
          markerWidth="7"
          markerHeight="7"
          orient="auto-start-reverse"
        >
          <path d="M0,0 L10,5 L0,10 z" class="head" />
        </marker>
      </defs>

      <!-- spine -->
      <path d="M210,58 V88" class="edge" marker-end="url(#taq-arrow)" />
      <path d="M210,138 V168" class="edge" marker-end="url(#taq-arrow)" />

      <!-- branch out -->
      <path d="M210,236 V256 H105 V286" class="edge" marker-end="url(#taq-arrow)" />
      <path d="M210,236 V256 H315 V286" class="edge" marker-end="url(#taq-arrow)" />
      <text x="140" y="250" class="edge-label">If</text>
      <text x="262" y="250" class="edge-label">Else</text>

      <!-- left path -->
      <path d="M105,336 V410" class="edge" marker-end="url(#taq-arrow)" />

      <!-- right path: loop -->
      <path d="M315,386 V410" class="edge" marker-end="url(#taq-arrow)" />

      <!-- nodes -->
      <g class="node trigger">
        <rect x="120" y="10" width="180" height="48" rx="24" />
        <text x="210" y="31" class="label">Trigger</text>
        <text x="210" y="47" class="hint">After Record Create</text>
      </g>

      <g class="node">
        <rect x="120" y="90" width="180" height="48" rx="10" />
        <text x="210" y="111" class="label">Step</text>
        <text x="210" y="127" class="hint">Find Record</text>
      </g>

      <g class="node branch">
        <polygon points="210,170 290,203 210,236 130,203" />
        <text x="210" y="207" class="label">First Match</text>
      </g>

      <g class="node">
        <rect x="20" y="288" width="170" height="48" rx="10" />
        <text x="105" y="309" class="label">Step</text>
        <text x="105" y="325" class="hint">Send Notification</text>
      </g>

      <g class="node loop">
        <rect x="230" y="288" width="170" height="98" rx="12" />
        <text x="315" y="306" class="label">Loop · Each Record</text>
        <rect x="250" y="320" width="130" height="48" rx="10" class="inner" />
        <text x="315" y="341" class="label">Step</text>
        <text x="315" y="357" class="hint">Update Record</text>
      </g>

      <g class="node end">
        <rect x="60" y="412" width="90" height="34" rx="17" />
        <text x="105" y="434" class="label">End</text>
      </g>
      <g class="node end">
        <rect x="270" y="412" width="90" height="34" rx="17" />
        <text x="315" y="434" class="label">End</text>
      </g>
    </svg>
  </figure>
</template>

<style scoped>
.dg {
  margin: 24px 0;
}

.taq-flow {
  display: flex;
  justify-content: center;
  padding: 16px;
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  background: var(--vp-c-bg-soft);
}

svg {
  width: 100%;
  max-width: 440px;
  height: auto;
  font-family: var(--vp-font-family-base);
}

.edge {
  fill: none;
  stroke: var(--vp-c-text-3);
  stroke-width: 1.5;
}

.head {
  fill: var(--vp-c-text-3);
}

.edge-label {
  font-size: 13px;
  fill: var(--vp-c-text-2);
  text-anchor: middle;
}

.node rect,
.node polygon {
  fill: var(--vp-c-bg);
  stroke: var(--vp-c-border);
  stroke-width: 1.2;
}

.node.trigger rect {
  fill: var(--vp-c-brand-soft);
  stroke: var(--vp-c-brand-1);
}

.node.branch polygon {
  fill: var(--vp-c-bg);
  stroke: var(--vp-c-accent);
}

.node.loop > rect:first-of-type {
  fill: transparent;
  stroke-dasharray: 5 4;
}

.node .inner {
  fill: var(--vp-c-bg);
}

.node.end rect {
  fill: var(--vp-c-bg-soft);
}

.label {
  font-size: 15px;
  font-weight: 500;
  fill: var(--vp-c-text-1);
  text-anchor: middle;
}

.hint {
  font-size: 13px;
  fill: var(--vp-c-text-2);
  text-anchor: middle;
}
</style>
