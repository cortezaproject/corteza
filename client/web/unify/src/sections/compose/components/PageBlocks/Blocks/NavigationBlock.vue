<template>
  <PageBlock :block="block">
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
            <div class="relative" :style="itemStyle(navItem)">
              <Button
                :label="displayDropdownText(navItem)"
                text
                size="small"
                icon="pi pi-chevron-down"
                icon-pos="right"
                :style="{ color: navItem.options?.item?.textColor }"
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
              class="nav-link-item flex items-center justify-center px-3 py-2 no-underline"
              :class="{ 'opacity-50 pointer-events-none': !navItem.options?.enabled }"
              :style="itemStyle(navItem)"
            >
              {{ navItem.options?.item?.label || '' }}
            </router-link>
            <a
              v-else-if="getHrefLink(navItem)"
              :href="getHrefLink(navItem)"
              :target="selectTargetOption(navItem.options?.item?.target)"
              class="nav-link-item flex items-center justify-center px-3 py-2 no-underline text-color"
              :class="{ 'opacity-50 pointer-events-none': !navItem.options?.enabled }"
              :style="itemStyle(navItem)"
            >
              {{ navItem.options?.item?.label || '' }}
            </a>
            <!-- Text section: non-clickable label -->
            <span
              v-else-if="navItem.type === 'text-section'"
              class="flex items-center px-3 py-2 text-sm"
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
import { useRoute } from 'vue-router'
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
const menuRefs = ref({})

const options = computed(() => props.block.options || {})
const navigationItems = computed(() => options.value.navigationItems || [])
const display = computed(() => options.value.display || {})

const navContainerClass = computed(() => {
  const classes = ['gap-1', 'items-center']
  const { alignment, appearance, justify } = display.value

  if (justify === 'justify') classes.push('justify-around')
  else if (alignment === 'right') classes.push('justify-end')
  else if (alignment === 'left') classes.push('justify-start')
  else classes.push('justify-center')

  if (appearance === 'pills') classes.push('nav-pills')
  if (appearance === 'tabs') classes.push('nav-tabs')

  classes.push('flex-wrap')
  return classes
})

function setMenuRef(index, el) {
  menuRefs.value[index] = el
}

function toggleDropdown(index, event) {
  menuRefs.value[index]?.toggle(event)
}

function isComposeDropdownPage(navItem) {
  return navItem.type === 'compose' && navItem.options?.item?.displaySubPages
}

function displayDropdownText(navItem) {
  return navItem.options?.item?.label || ''
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
  return []
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

function getRouterLink(navItem) {
  if (['dropdown', 'text-section'].includes(navItem.type) || isComposeDropdownPage(navItem))
    return null

  if (navItem.type === 'compose') {
    const pageID = navItem.options?.item?.pageID
    const pageLayoutID = navItem.options?.item?.pageLayoutID

    if (!pageID) return null

    const isSamePage = pageID === route.params?.pageID
    const query = pageLayoutID ? { layoutID: pageLayoutID } : {}

    return isSamePage ? { ...route, query } : { name: 'page', params: { pageID }, query }
  }

  return null
}

function getHrefLink(navItem) {
  if (['dropdown', 'text-section'].includes(navItem.type) || isComposeDropdownPage(navItem))
    return null
  return navItem.type === 'url' ? evaluateUrl(navItem.options?.item?.url) : null
}
</script>

<style scoped>
.nav-link-item:hover {
  text-decoration: underline !important;
}
</style>
