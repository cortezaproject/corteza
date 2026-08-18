<template>
  <PageBlock :block="block" :record="record">
    <div
      v-if="!navigationItems.length"
      class="flex items-center justify-center h-full p-3 text-muted-color italic"
    >
      {{ $t('block.navigation.noNavigationItems') }}
    </div>

    <div v-else class="h-full w-full overflow-auto">
      <div class="flex h-full" :class="navContainerClass">
        <template v-for="(navItem, index) in navigationItems" :key="`nav-${index}`">
          <!-- Dropdown type -->
          <template v-if="navItem.type === 'dropdown' || isComposeDropdownPage(navItem)">
            <div
              class="relative flex items-center"
              :class="itemClass(navItem)"
              :style="itemStyle(navItem)"
            >
              <!-- The trigger takes its size, weight and colour from the item
                   around it, so a dropdown reads as one of the links rather
                   than as a button that wandered into the row. -->
              <Button
                :label="displayDropdownText(navItem)"
                text
                icon="pi pi-chevron-down"
                icon-pos="right"
                class="!p-0 !font-normal !text-inherit"
                :style="{ color: navItem.options?.item?.textColor, fontSize: 'inherit' }"
                @click="toggleDropdown(index, $event)"
              />
              <Menu
                :ref="el => setMenuRef(index, el)"
                :model="getDropdownItems(navItem)"
                :popup="true"
              />
            </div>
          </template>

          <!-- Regular link -->
          <template v-else>
            <router-link
              v-if="getRouterLink(navItem)"
              :to="getRouterLink(navItem)"
              :target="selectTargetOption(navItem.options?.item?.target)"
              class="nav-link-item flex items-center justify-center no-underline"
              :class="[
                itemClass(navItem),
                { 'opacity-50 pointer-events-none': !navItem.options?.enabled },
              ]"
              :style="itemStyle(navItem)"
            >
              {{ navItem.options?.item?.label || '' }}
            </router-link>
            <a
              v-else-if="getHrefLink(navItem)"
              :href="getHrefLink(navItem)"
              :target="selectTargetOption(navItem.options?.item?.target)"
              class="nav-link-item flex items-center justify-center no-underline text-color"
              :class="[
                itemClass(navItem),
                { 'opacity-50 pointer-events-none': !navItem.options?.enabled },
              ]"
              :style="itemStyle(navItem)"
            >
              {{ navItem.options?.item?.label || '' }}
            </a>
            <!-- Text section: non-clickable label -->
            <span
              v-else-if="navItem.type === 'text-section'"
              class="flex items-center text-sm"
              :class="paddingClass"
              :style="itemStyle(navItem)"
            >
              {{ navItem.options?.item?.label || '' }}
            </span>
          </template>
        </template>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NoID } from '@planetcrust/human-js'
import { usePageStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $Auth = inject('$Auth', {})
const route = useRoute()
const router = useRouter()
const pageStore = usePageStore()
const menuRefs = ref({})

const options = computed(() => props.block.options || {})
const navigationItems = computed(() => options.value.navigationItems || [])
const display = computed(() => options.value.display || {})

const appearance = computed(() => display.value.appearance || '')
const isSmall = computed(() => appearance.value === 'small')

// Small is a tighter version of the same row, not a third shape.
const paddingClass = computed(() => (isSmall.value ? 'px-2 py-1 text-xs' : 'px-3 py-2'))

const navContainerClass = computed(() => {
  const classes = ['gap-1', 'flex-wrap']
  const { alignment, justify } = display.value

  if (justify === 'justify') classes.push('justify-around')
  else if (alignment === 'right') classes.push('justify-end')
  else if (alignment === 'left') classes.push('justify-start')
  else classes.push('justify-center')

  // Tabs are full-height so each one's underline meets the row's rule; the
  // other appearances centre their items in the block instead.
  if (appearance.value === 'tabs') classes.push('items-stretch', 'border-b', 'border-surface')
  else classes.push('items-center')

  return classes
})

// The appearance decoration for one item, active state included. Author-set
// colours are inline styles, so they still win over whatever this returns.
function itemClass(navItem) {
  const classes = [paddingClass.value]
  const active = isActiveItem(navItem)

  if (appearance.value === 'tabs') {
    classes.push('-mb-px', 'border-b-2', 'hover:text-primary')
    classes.push(active ? 'border-primary text-primary' : 'border-transparent')
  } else if (appearance.value === 'pills') {
    classes.push('rounded-full')
    classes.push(active ? 'bg-primary text-primary-contrast' : 'hover:bg-emphasis')
  } else {
    classes.push('hover:underline')
    if (active) classes.push('text-primary')
  }

  return classes
}

// Only a compose item can be the page you are on — a URL leaves the app and a
// text section is not a destination. A sub-pages item also answers for its
// children, so the strip stays lit while you are inside one.
function isActiveItem(navItem) {
  if (navItem.type !== 'compose') return false

  const pageID = navItem.options?.item?.pageID
  const current = route.params?.pageID
  if (!pageID || !current) return false
  if (pageID === current) return true

  return isComposeDropdownPage(navItem) && subPages(pageID).some(p => p.pageID === current)
}

function setMenuRef(index, el) {
  menuRefs.value[index] = el
}

function toggleDropdown(index, event) {
  menuRefs.value[index]?.toggle(event)
}

// A sub-pages item only becomes a dropdown once there is something to drop:
// the flag can outlive the children it was set for, and a menu that opens on
// nothing is worse than the plain link it replaced.
function isComposeDropdownPage(navItem) {
  if (navItem.type !== 'compose' || !navItem.options?.item?.displaySubPages) return false
  return subPages(navItem.options?.item?.pageID).length > 0
}

function displayDropdownText(navItem) {
  const item = navItem.options?.item || {}
  if (item.label) return item.label
  return isComposeDropdownPage(navItem) ? pageStore.getByID?.(item.pageID)?.title || '' : ''
}

// The pages nested under a compose page. Record pages are excluded: they need
// a record to show, so they are not somewhere a menu entry can send you.
function subPages(pageID) {
  return (pageStore.set || []).filter(p => p.selfID === pageID && p.moduleID === NoID)
}

function getDropdownItems(navItem) {
  if (navItem.type === 'dropdown') {
    return (navItem.options?.item?.dropdown?.items || []).map(item => {
      const url = evaluateUrl(item.url)
      return {
        label: item.label,
        url,
        target: selectTargetOption(item.target),
        command: () => {
          if (url) window.open(url, selectTargetOption(item.target))
        },
        separator: item.delimiter,
      }
    })
  }

  if (!isComposeDropdownPage(navItem)) return []

  // A sub-pages menu leads with the page it was configured for, so the item
  // still reaches its own destination once it has become a dropdown.
  const item = navItem.options?.item || {}

  return [
    {
      label: displayDropdownText(navItem),
      command: () => goToPage(item.pageID, item.pageLayoutID),
    },
    { separator: true },
    ...subPages(item.pageID).map(child => ({
      label: child.title || child.handle || child.pageID,
      command: () => goToPage(child.pageID),
    })),
  ]
}

function goToPage(pageID, pageLayoutID) {
  const to = composeRoute(pageID, pageLayoutID)
  if (to) router.push(to)
}

// A URL that isn't a valid template (or uses record vars without a record) keeps working unchanged.
function evaluateUrl(url) {
  if (!url) return url

  const record = props.record
  if (!record && (url.includes('${record') || url.includes('${ownerID}'))) {
    return url
  }

  try {
    const user = $Auth?.user || {}
    return evaluatePrefilter(url, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  } catch {
    return url
  }
}

function selectTargetOption(target) {
  switch (target) {
    case 'sameTab':
      return '_self'
    case 'newTab':
      return '_blank'
    default:
      return '_self'
  }
}

function itemStyle(navItem) {
  const style = {}
  const item = navItem.options?.item
  if (item?.textColor) style.color = item.textColor
  if (item?.backgroundColor) style.backgroundColor = item.backgroundColor
  return style
}

// Where a compose page is reached. The namespace slug is not passed: it is the
// one the strip is already rendering in, and the router keeps it.
function composeRoute(pageID, pageLayoutID) {
  if (!pageID) return null

  const query = pageLayoutID ? { layoutID: pageLayoutID } : {}
  const page = pageStore.getByID?.(pageID)

  // A record page draws one record. Reached from a navigation item there is
  // none yet, so it opens on a blank one — the route records are created at.
  if (page?.moduleID && page.moduleID !== NoID) {
    return { name: 'page.record', params: { pageID, recordID: NoID }, query }
  }

  return { name: 'page', params: { pageID }, query }
}

function getRouterLink(navItem) {
  if (['dropdown', 'text-section'].includes(navItem.type) || isComposeDropdownPage(navItem))
    return null

  if (navItem.type === 'compose') {
    return composeRoute(navItem.options?.item?.pageID, navItem.options?.item?.pageLayoutID)
  }

  return null
}

function getHrefLink(navItem) {
  if (['dropdown', 'text-section'].includes(navItem.type) || isComposeDropdownPage(navItem))
    return null
  return navItem.type === 'url' ? evaluateUrl(navItem.options?.item?.url) : null
}
</script>
