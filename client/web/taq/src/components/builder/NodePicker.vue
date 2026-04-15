<template>
  <div class="flex flex-col h-[600px]">
    <!-- Search Bar -->
    <div class="p-4 border-b border-surface shrink-0">
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

    <!-- Panes -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Left Pane: Groups -->
      <div class="w-1/3 border-r border-surface p-4 overflow-y-auto">
        <div class="space-y-1">
          <button
            v-for="category in filteredCategories"
            :key="category.id"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-border text-left transition-colors"
            :class="
              selectedGroup === category.id
                ? 'bg-primary text-primary-contrast'
                : 'hover:bg-emphasis text-color'
            "
            @click="selectedGroup = category.id"
          >
            <TaqIcon :icon="category.icon" class="text-lg" />
            <span class="flex-1 font-medium truncate">{{ category.label }}</span>
          </button>
        </div>
      </div>

      <!-- Right Pane: Items -->
      <div class="flex-1 flex flex-col p-4 overflow-hidden">
        <!-- Items Grid/List -->
        <div class="flex-1 overflow-y-auto pr-2">
          <div v-if="filteredNodes.length > 0" class="space-y-2">
            <div
              v-for="node in filteredNodes"
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
                  <div class="text-sm text-muted-color truncate" :title="node.description">
                    {{ node.description }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty state for search -->
          <div v-else class="h-full flex flex-col items-center justify-center text-muted-color">
            <i class="pi pi-search text-3xl mb-3 opacity-50" />
            <p>{{ $t('builder.nodePicker.noResults', { query: searchQuery }) }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { normalizeIcon } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import TaqIcon from '../common/TaqIcon.vue'
import { useAutomationStore } from '../../stores/automation'
import {
  DEFAULT_ACTION_ICON,
  DEFAULT_TRIGGER_ICON,
  getTriggerMeta,
} from '../../utils/flow-constants'

const { t } = useI18n()
const store = useAutomationStore()

const props = defineProps({
  filterCategory: {
    type: String,
    default: null, // 'trigger' or 'step'
  },
})

const emit = defineEmits(['select', 'close'])

const searchQuery = ref('')
const selectedGroup = ref(null)

// Helper to reliably sanitize group names for IDs
function slugify(text) {
  return (text || '')
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^\w-]/g, '')
}

// Icon mapping for groups
const GROUP_ICONS = {
  // Step categories
  branches: { type: 'name', value: 'sitemap' },
  loops: { type: 'name', value: 'refresh' },
  records: { type: 'name', value: 'database' },
  attachments: { type: 'name', value: 'paperclip' },
  notifications: { type: 'name', value: 'bell' },
  users: { type: 'name', value: 'user' },
  roles: { type: 'name', value: 'users' },
  reminders: { type: 'name', value: 'clock' },
  templates: { type: 'name', value: 'file' },
  email: { type: 'name', value: 'envelope' },
  http: { type: 'name', value: 'link' },
  agents: { type: 'name', value: 'sparkles' },
  logging: { type: 'name', value: 'list' },

  // Trigger categories
  manual: { type: 'name', value: 'user' },
  schedule: { type: 'name', value: 'calendar' },
  general: { type: 'name', value: 'bolt' },
  system: { type: 'name', value: 'cog' },
}

function getGroupIcon(groupName) {
  const slug = slugify(groupName)
  return GROUP_ICONS[slug] || { type: 'name', value: 'folder' }
}

// Helper to map a trigger to a node format
function mapTriggerNode(trigger) {
  const meta = getTriggerMeta(trigger.eventType, trigger.resourceType)
  const icon = normalizeIcon(trigger.meta?.icon) || meta?.icon || DEFAULT_TRIGGER_ICON
  return {
    id: `${trigger.resourceType}:${trigger.eventType}`,
    type: 'trigger',
    label: trigger.meta?.short || trigger.eventType,
    icon,
    description: trigger.meta?.description || '',
    eventType: trigger.eventType,
    resourceType: trigger.resourceType,
    group: trigger.groups?.[0] || 'General',
    weight: trigger.meta?.weight || 0,
  }
}

// Compute all available categories based on `filterCategory`
const availableCategories = computed(() => {
  const groupsMap = {} // { "slug": { id, label, icon, nodes: [] } }

  const addNodeToGroup = node => {
    const label = node.group || 'System'
    const id = slugify(label)

    if (!groupsMap[id]) {
      groupsMap[id] = {
        id,
        label,
        icon: getGroupIcon(label) || node.icon,
        nodes: [],
      }
    }
    groupsMap[id].nodes.push(node)
  }

  // Populate Triggers (if filterCategory is 'trigger')
  if (props.filterCategory === 'trigger') {
    store.triggers.forEach(t => {
      addNodeToGroup(mapTriggerNode(t))
    })
  }

  // Populate Steps/Actions (if filterCategory isn't 'trigger')
  if (props.filterCategory !== 'trigger') {
    // 1. Branches group (exclusive + inclusive)
    groupsMap['branches'] = {
      id: 'branches',
      label: t('builder.nodePicker.categories.branches'),
      icon: getGroupIcon('branches'),
      nodes: [
        {
          id: 'branchExclusive',
          type: 'condition',
          label: t('builder.nodePicker.nodes.branches.exclusive.label'),
          icon: DEFAULT_ICONS.BRANCH,
          description: t('builder.nodePicker.nodes.branches.exclusive.description'),
          ref: 'gatewayExclusive',
          group: 'Branches',
          weight: 0,
        },
        {
          id: 'branchInclusive',
          type: 'condition',
          label: t('builder.nodePicker.nodes.branches.inclusive.label'),
          icon: DEFAULT_ICONS.BRANCH,
          description: t('builder.nodePicker.nodes.branches.inclusive.description'),
          ref: 'gatewayInclusive',
          group: 'Branches',
          weight: 1,
        },
      ],
    }

    // 2. Loop through all functions from construct library (exclude gateways)
    store.functions
      .filter(fn => fn.kind !== 'gateway')
      .forEach(fn => {
        const icon =
          normalizeIcon(fn.meta?.icon) ||
          (fn.kind === 'iterator' ? DEFAULT_ICONS.ITERATOR : DEFAULT_ACTION_ICON)
        addNodeToGroup({
          id: fn.ref,
          type: 'action',
          label: fn.meta?.short || fn.ref,
          icon,
          description: fn.meta?.description || '',
          ref: fn.ref,
          kind: fn.kind,
          group: fn.groups?.[0] || 'System',
          weight: fn.meta?.weight || 0,
        })
      })
  }

  return Object.values(groupsMap).sort((a, b) => a.label.localeCompare(b.label))
})

const filteredCategories = computed(() => {
  let categories = availableCategories.value

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    categories = categories
      .map(category => ({
        ...category,
        nodes: category.nodes.filter(
          node =>
            node.label.toLowerCase().includes(q) ||
            (node.description && node.description.toLowerCase().includes(q)),
        ),
      }))
      .filter(category => category.nodes.length > 0)
  }

  return categories
})

// Auto-select the first category when the modal opens or search changes
watch(
  filteredCategories,
  newVal => {
    if (!selectedGroup.value || !newVal.find(c => c.id === selectedGroup.value)) {
      if (newVal.length > 0) {
        selectedGroup.value = newVal[0].id
      } else {
        selectedGroup.value = null
      }
    }
  },
  { immediate: true },
)

// Active category being viewed
const activeCategory = computed(() => {
  if (!selectedGroup.value) return null
  return filteredCategories.value.find(c => c.id === selectedGroup.value) || null
})

// Filtered nodes of the active category
const filteredNodes = computed(() => {
  if (!activeCategory.value) return []
  // Sort nodes by weight, then alphabetically
  return [...activeCategory.value.nodes].sort((a, b) => {
    if (a.weight !== b.weight) {
      return a.weight - b.weight
    }
    return a.label.localeCompare(b.label)
  })
})

function selectNode(node) {
  emit('select', node)
  emit('close')
}
</script>
