<template>
  <div>
    <!-- Search -->
    <div class="mb-4">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="searchQuery"
          :placeholder="
            filterCategory === 'trigger'
              ? $t('builder.nodePicker.searchTriggers')
              : $t('builder.nodePicker.searchSteps')
          "
          class="w-full"
          autofocus
        />
      </IconField>
    </div>

    <!-- Categories -->
    <div class="space-y-6">
      <div v-for="category in filteredCategories" :key="category.id">
        <h4 class="text-sm font-semibold text-muted-color uppercase mb-2 flex items-center gap-2">
          <i :class="category.icon" />
          {{ category.label }}
        </h4>
        <div class="space-y-2">
          <div
            v-for="node in category.nodes"
            :key="node.id"
            class="node-item p-3 rounded-border border border-surface cursor-pointer transition-colors hover:bg-emphasis"
            @click="selectNode(node)"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-border flex items-center justify-center">
                <TaqIcon :icon="node.icon" class="text-lg text-primary" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-medium text-color">{{ node.label }}</div>
                <div class="text-sm text-muted-color truncate">{{ node.description }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-if="filteredCategories.length === 0" class="text-center py-8 text-muted-color">
      <i class="pi pi-search text-2xl mb-2" />
      <p>{{ $t('builder.nodePicker.noResults', { query: searchQuery }) }}</p>
    </div>
  </div>
</template>

<script setup>
import { normalizeIcon } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import TaqIcon from '../common/TaqIcon.vue'
import { useAutomationStore } from '../../stores/automation'
import { DEFAULT_ACTION_ICON, DEFAULT_TRIGGER_ICON, TRIGGER_META } from '../../utils/flow-constants'

const { t } = useI18n()
const store = useAutomationStore()

const props = defineProps({
  filterCategory: {
    type: String,
    default: null,
  },
})

const emit = defineEmits(['select', 'close'])

const searchQuery = ref('')

// Helper to map a trigger to a node format with icons and i18n
function mapTriggerNode(trigger) {
  const meta = TRIGGER_META[trigger.eventType]
  const icon = normalizeIcon(trigger.meta?.icon) || meta?.icon || DEFAULT_TRIGGER_ICON
  const i18nPrefix = meta ? `builder.nodePicker.nodes.triggers.${meta.i18nKey}` : ''
  return {
    id: `${trigger.resourceType}:${trigger.eventType}`,
    type: 'trigger',
    label: trigger.meta?.short || (i18nPrefix ? t(`${i18nPrefix}.label`) : trigger.eventType),
    icon,
    description: trigger.meta?.description || (i18nPrefix ? t(`${i18nPrefix}.description`) : ''),
    eventType: trigger.eventType,
    resourceType: trigger.resourceType,
  }
}

// Node categories - triggers grouped by resourceType, functions from API, logic nodes defined here
const nodeCategories = computed(() => {
  // Split triggers into General (non-record) and Record groups
  const generalTriggers = store.triggers
    .filter(t => t.resourceType !== 'compose:record')
    .map(mapTriggerNode)
  const recordTriggers = store.triggers
    .filter(t => t.resourceType === 'compose:record')
    .map(mapTriggerNode)

  // Map API functions to node format (labels/descriptions from backend, exclude gateway - handled as logic)
  const actionNodes = store.functions
    .filter(fn => fn.kind !== 'gateway')
    .map(fn => {
      const icon = normalizeIcon(fn.meta?.icon) || DEFAULT_ACTION_ICON
      return {
        id: fn.ref,
        type: 'action',
        label: fn.meta?.short || fn.ref,
        icon,
        description: fn.meta?.description || '',
        ref: fn.ref,
        kind: fn.kind,
      }
    })

  const categories = []

  if (generalTriggers.length > 0) {
    categories.push({
      id: 'triggers-general',
      label: t('builder.nodePicker.categories.triggersGeneral'),
      icon: 'pi pi-bolt',
      nodes: generalTriggers,
    })
  }

  if (recordTriggers.length > 0) {
    categories.push({
      id: 'triggers-record',
      label: t('builder.nodePicker.categories.triggersRecord'),
      icon: 'pi pi-database',
      nodes: recordTriggers,
    })
  }

  categories.push(
    {
      id: 'logic',
      label: t('builder.nodePicker.categories.logic'),
      icon: 'pi pi-sitemap',
      nodes: [
        {
          id: 'branch',
          type: 'condition',
          label: t('builder.nodePicker.nodes.logic.branch.label'),
          icon: DEFAULT_ICONS.BRANCH,
          description: t('builder.nodePicker.nodes.logic.branch.description'),
          ref: 'gateway',
        },
      ],
    },
    {
      id: 'actions',
      label: t('builder.nodePicker.categories.actions'),
      icon: 'pi pi-cog',
      nodes: actionNodes,
    },
  )

  return categories
})

const filteredCategories = computed(() => {
  let categories = nodeCategories.value

  // Filter by category: 'trigger' shows only trigger groups, otherwise exclude them
  if (props.filterCategory === 'trigger') {
    categories = categories.filter(c => c.id.startsWith('triggers-'))
  } else {
    categories = categories.filter(c => !c.id.startsWith('triggers-'))
  }

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    categories = categories
      .map(category => ({
        ...category,
        nodes: category.nodes.filter(
          node =>
            node.label.toLowerCase().includes(query) ||
            node.description.toLowerCase().includes(query),
        ),
      }))
      .filter(category => category.nodes.length > 0)
  }
  return categories
})

function selectNode(node) {
  emit('select', node)
  emit('close')
}
</script>
