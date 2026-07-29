<template>
  <div class="flex flex-col gap-4">
    <ul class="flex flex-col gap-2.5">
      <li class="flex gap-2.5 text-sm">
        <i class="pi pi-cloud-upload text-muted-color mt-0.5" />
        <span>{{ $t('project.publish.goLive.becomesLive', { name: projectName }) }}</span>
      </li>

      <li v-if="hasParent" class="flex gap-2.5 text-sm">
        <i class="pi pi-arrow-right-arrow-left text-muted-color mt-0.5" />
        <span>{{ $t('project.publish.goLive.recordsMove') }}</span>
      </li>

      <!-- Repeated here on purpose: each loss was agreed to a stage earlier,
           and this is the moment it costs the most to have forgotten it. Named
           one by one rather than summed — "1,284 records lose their fax value"
           is checkable in a way that a single total is not. -->
      <li v-for="loss in losses" :key="loss.path" class="flex gap-2.5 text-sm">
        <i class="pi pi-exclamation-triangle text-red-500 mt-0.5" />
        <span>
          {{
            loss.kind === 'field'
              ? $t('project.publish.goLive.fieldDiscarded', {
                  count: $n(loss.records),
                  field: loss.name,
                  module: loss.module,
                })
              : $t('project.publish.goLive.moduleDiscarded', {
                  count: $n(loss.records),
                  module: loss.name,
                })
          }}
        </span>
      </li>

      <li class="flex gap-2.5 text-sm">
        <i class="pi pi-inbox text-muted-color mt-0.5" />
        <span>{{ $t('project.publish.goLive.parentRetired') }}</span>
      </li>

      <!-- Warns, never blocks (locked contract) — the same completeness the
           Manage & Monitor rail shows, from the same shared fetch. -->
      <li v-if="openWorkItems > 0" class="flex gap-2.5 text-sm">
        <i class="pi pi-list-check text-amber-500 mt-0.5" />
        <span>
          {{
            $t('project.publish.confirm.unfinishedWarning', {
              open: openWorkItems,
              total: assignedWorkItems,
            })
          }}
        </span>
      </li>
    </ul>

    <!-- Type-to-confirm replaces the old modal, and only when records are
         actually lost. A dialog that appears every time is dismissed every
         time; typing the handle is a deliberate act, and it keeps the
         consequences above on screen while you make it. -->
    <div
      v-if="requiresTypedConfirmation"
      class="rounded-lg border border-red-500/40 bg-red-500/5 p-3"
    >
      <label :for="confirmInputId" class="block text-sm mb-2">
        {{ $t('project.publish.goLive.typeToConfirm', { handle }) }}
      </label>
      <InputText
        :id="confirmInputId"
        :model-value="typed"
        size="small"
        class="w-full max-w-xs"
        :disabled="disabled"
        autocomplete="off"
        @update:model-value="$emit('update:typed', $event)"
      />
    </div>
  </div>
</template>

<script setup>
// Publish stage 4 — what happens the moment you publish, in plain language,
// plus the confirmation itself. There is no confirm dialog: this stage IS the
// confirmation (ruled 2026-07-29).
import { computed, useId } from 'vue'

const props = defineProps({
  projectName: { type: String, default: '' },
  // The project handle, typed back to confirm a lossy publish.
  handle: { type: String, default: '' },
  typed: { type: String, default: '' },
  // Whether there is a parent revision to migrate records from at all.
  hasParent: { type: Boolean, default: false },
  // Data the user chose to leave behind: [{ path, kind, name, module, records }]
  losses: { type: Array, default: () => [] },
  openWorkItems: { type: Number, default: 0 },
  assignedWorkItems: { type: Number, default: 0 },
  disabled: { type: Boolean, default: false },
})

defineEmits(['update:typed'])

const confirmInputId = useId()

const requiresTypedConfirmation = computed(() => props.losses.length > 0)
</script>
