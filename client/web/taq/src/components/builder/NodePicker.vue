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
    <div class="space-y-4">
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
                <i :class="node.icon" class="text-lg text-primary" />
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
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAutomationStore } from '../../stores/automation'

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

// Trigger config: maps eventType to translation key and icon
const triggerConfig = {
  onManual: { key: 'manual', icon: 'pi pi-play' },
  onInterval: { key: 'interval', icon: 'pi pi-clock' },
  onTimestamp: { key: 'timestamp', icon: 'pi pi-calendar' },
  afterCreate: { key: 'recordCreate', icon: 'pi pi-plus-circle' },
  beforeCreate: { key: 'recordCreate', icon: 'pi pi-plus-circle' },
  afterUpdate: { key: 'recordUpdate', icon: 'pi pi-pencil' },
  beforeUpdate: { key: 'recordUpdate', icon: 'pi pi-pencil' },
  afterDelete: { key: 'recordDelete', icon: 'pi pi-trash' },
  beforeDelete: { key: 'recordDelete', icon: 'pi pi-trash' },
  onHTTPRequest: { key: 'httpRequest', icon: 'pi pi-globe' },
}

// Function config: maps ref to translation key and icon
const functionConfig = {
  usersLookup: { key: 'usersLookup', icon: 'pi pi-user' },
}

// Node categories - triggers and functions from API, logic nodes defined here
const nodeCategories = computed(() => {
  // Map API triggers to node format
  const triggerNodes = store.triggers.map(trigger => {
    const config = triggerConfig[trigger.eventType]
    const key = config?.key || trigger.eventType
    return {
      id: `${trigger.resourceType}:${trigger.eventType}`,
      type: 'trigger',
      label: t(`builder.nodePicker.nodes.triggers.${key}.label`, trigger.eventType),
      icon: config?.icon || 'pi pi-bolt',
      description: t(`builder.nodePicker.nodes.triggers.${key}.description`, ''),
      eventType: trigger.eventType,
      resourceType: trigger.resourceType,
    }
  })

  // Map API functions to node format (exclude gateway/branch - handled as logic)
  const actionNodes = store.functions
    .filter(fn => fn.kind !== 'gateway')
    .map(fn => {
      const config = functionConfig[fn.ref]
      const key = config?.key || fn.ref
      return {
        id: fn.ref,
        type: 'action',
        label: t(`builder.nodePicker.nodes.actions.${key}.label`, fn.meta?.short || fn.ref),
        icon: config?.icon || 'pi pi-cog',
        description: t(`builder.nodePicker.nodes.actions.${key}.description`, fn.meta?.description || ''),
        ref: fn.ref,
        kind: fn.kind,
      }
    })

  return [
    {
      id: 'triggers',
      label: t('builder.nodePicker.categories.triggers'),
      icon: 'pi pi-bolt',
      nodes: triggerNodes,
    },
    {
      id: 'logic',
      label: t('builder.nodePicker.categories.logic'),
      icon: 'pi pi-sitemap',
      nodes: [
        {
          id: 'branch',
          type: 'condition',
          label: t('builder.nodePicker.nodes.logic.branch.label'),
          icon: 'pi pi-sitemap',
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
  ]
})

const filteredCategories = computed(() => {
  let categories = nodeCategories.value

  // Filter by category: 'trigger' shows only triggers, otherwise exclude triggers
  if (props.filterCategory === 'trigger') {
    categories = categories.filter(c => c.id === 'triggers')
  } else {
    categories = categories.filter(c => c.id !== 'triggers')
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
