<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 px-4 pt-4">
      <p class="text-sm text-muted-color">{{ $t('project.permissions.hint') }}</p>
    </div>

    <div v-if="!roles.length" class="flex-1 grid place-items-center text-center px-6 text-muted-color">
      <div>
        <i class="pi pi-id-card text-4xl mb-2" />
        <p class="text-sm">{{ $t('project.permissions.noRoles') }}</p>
      </div>
    </div>
    <div
      v-else-if="!sections.length"
      class="flex-1 grid place-items-center text-center px-6 text-muted-color"
    >
      <div>
        <i class="pi pi-sitemap text-4xl mb-2" />
        <p class="text-sm">{{ $t('project.permissions.noResources') }}</p>
      </div>
    </div>

    <div v-else class="flex-1 min-h-0 overflow-auto p-4">
      <div class="inline-block min-w-full rounded-lg border border-surface overflow-hidden">
        <table class="text-sm border-collapse">
          <thead>
            <tr class="bg-emphasis">
              <th class="text-left font-medium px-3 py-2 sticky left-0 z-10 bg-emphasis min-w-64" />
              <th
                v-for="role in roles"
                :key="role.id"
                class="font-medium px-3 py-2 border-l border-surface text-center whitespace-nowrap"
              >
                {{ role.name }}
              </th>
            </tr>
          </thead>
          <tbody>
            <template v-for="section in sections" :key="section.kind">
              <!-- Kind wildcard row: grants on all resources of the kind -->
              <tr class="border-t border-surface bg-surface-100 dark:bg-surface-800">
                <td class="px-3 py-2 font-semibold sticky left-0 z-10 bg-inherit">
                  <span class="inline-flex items-center gap-2 whitespace-nowrap">
                    <span
                      class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
                      :class="[section.cfg.bg, section.cfg.ring]"
                    >
                      <i :class="[section.cfg.icon, section.cfg.text, 'text-xs']" />
                    </span>
                    {{ $t('project.permissions.allOf', { kind: $t(section.cfg.labelKey) }) }}
                  </span>
                </td>
                <td
                  v-for="role in roles"
                  :key="role.id"
                  class="border-l border-surface text-center px-2 py-1.5"
                >
                  <Button
                    :icon="'pi pi-lock'"
                    severity="secondary"
                    text
                    rounded
                    size="small"
                    :disabled="disabled"
                    :title="$t('project.permissions.editFor', { role: role.name })"
                    @click="openDialog(section.allRes, allLabel(section), role.id)"
                  />
                </td>
              </tr>

              <!-- Specific resources -->
              <tr
                v-for="item in section.items"
                :key="item.id"
                class="border-t border-surface odd:bg-surface even:bg-emphasis/20 hover:bg-emphasis/50"
              >
                <td class="pl-10 pr-3 py-2 sticky left-0 z-10 bg-inherit">
                  <span class="truncate">{{ item.name }}</span>
                </td>
                <td
                  v-for="role in roles"
                  :key="role.id"
                  class="border-l border-surface text-center px-2 py-1.5"
                >
                  <Button
                    icon="pi pi-lock"
                    severity="secondary"
                    text
                    rounded
                    size="small"
                    :disabled="disabled"
                    :title="$t('project.permissions.editFor', { role: role.name })"
                    @click="openDialog(item.resource, item.name, role.id)"
                  />
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <ProjectPermissionsDialog
      v-model="dlg.open"
      :resource="dlg.resource"
      :title="dlg.title"
      :preselect-id="dlg.roleId"
      :roles="roles"
    />
  </div>
</template>

<script setup>
import ProjectPermissionsDialog from '@/sections/project/components/permissions/ProjectPermissionsDialog.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, onMounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')

const roles = computed(() => store.rolesFor(props.project.id))
const ns = computed(() => props.project.namespaceID)

// Each kind: how to list its resources and build its RBAC resource string.
// Compose resources are namespaced (type/ns/id); system & automation are not.
const KINDS = [
  { kind: 'module', getter: 'resourcesFor', type: 'corteza::compose:module', ns: true },
  { kind: 'page', getter: 'pagesFor', type: 'corteza::compose:page', ns: true },
  { kind: 'automation', getter: 'automationsFor', type: 'corteza::automation:ng-automation', ns: false },
  { kind: 'agent', getter: 'agentsFor', type: 'corteza::system:agent', ns: false },
  { kind: 'chatbot', getter: 'chatbotsFor', type: 'corteza::system:chatbot', ns: false },
  { kind: 'connection', getter: 'connectionsFor', type: 'corteza::system:configured-connection', ns: false },
]

function resStr(def, id) {
  return def.ns ? `${def.type}/${ns.value}/${id}` : `${def.type}/${id}`
}

const sections = computed(() =>
  KINDS.map(def => {
    const items = store[def.getter](props.project.id) || []
    return {
      kind: def.kind,
      cfg: kindConfig(def.kind),
      allRes: resStr(def, '*'),
      items: items.map(it => ({ id: it.id, name: it.name, resource: resStr(def, it.id) })),
    }
  }).filter(s => s.items.length),
)

function allLabel(section) {
  return t('project.permissions.allOf', { kind: t(section.cfg.labelKey) })
}

const dlg = reactive({ open: false, resource: '', title: '', roleId: '' })
function openDialog(resource, title, roleId) {
  if (props.disabled) return
  dlg.resource = resource
  dlg.title = title
  dlg.roleId = roleId
  dlg.open = true
}

async function refresh(id) {
  if (!id) return
  try {
    await Promise.all([
      store.loadRoles(id),
      store.loadResources(id),
      store.loadPages(id),
      store.loadAutomations(id),
      store.loadAgents(id),
      store.loadChatbots(id),
      store.loadConnections(id),
    ])
  } catch (err) {
    $toast.toastErrorHandler(t('project.permissions.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(() => props.project.id, refresh)
</script>
