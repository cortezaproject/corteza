<template>
  <div class="flex flex-col gap-4">
    <p class="text-sm text-muted-color">
      {{ $t('system.user-groups.editor.members.explainer') }}
    </p>

    <CInputUser
      v-if="canManage"
      class="w-full"
      :placeholder="$t('system.user-groups.editor.members.placeholder')"
      :exclude-users="members.map(m => m.userID)"
      clear-on-select
      @select="addMember"
    >
      <template #option="{ option }">
        <div class="flex flex-col w-full min-w-0 whitespace-normal break-words">
          <span>{{ option.label }}</span>
          <small v-if="groupNames[option.userGroupID]" class="text-muted-color">
            {{ groupNames[option.userGroupID] }}
          </small>
        </div>
      </template>
    </CInputUser>

    <CFormItemList
      :items="members"
      :loading="loading"
      :empty-message="$t('system.user-groups.editor.members.empty')"
      item-key="userID"
      hide-remove
    >
      <template #default="{ item }">
        <span class="font-medium">{{ userName(item) }}</span>
        <span
          v-if="(item.name || item.handle) && item.email"
          class="text-xs text-muted-color block"
        >
          {{ item.email }}
        </span>
      </template>

      <template v-if="canManage" #hover-actions="{ item }">
        <Button
          v-tooltip.left="$t('system.user-groups.editor.members.move.action')"
          icon="pi pi-arrow-right-arrow-left"
          severity="secondary"
          text
          size="small"
          :aria-label="$t('system.user-groups.editor.members.move.action')"
          data-testid="member-move"
          @click.stop="openMove(item)"
        />
      </template>
    </CFormItemList>

    <Dialog
      v-model:visible="moving.visible"
      modal
      :header="$t('system.user-groups.editor.members.move.title')"
      class="w-full max-w-lg"
    >
      <div class="flex flex-col gap-4">
        <CFormGroup :label="$t('system.user-groups.editor.members.move.target')">
          <CInputUserGroup
            v-model="moving.targetID"
            class="w-full"
            :placeholder="$t('system.user-groups.editor.members.move.targetPlaceholder')"
            :exclude-user-groups="[userGroup.userGroupID]"
            @select="g => (moving.target = g)"
          />
        </CFormGroup>

        <p v-if="moving.user && moving.target" data-testid="member-move-confirm">
          {{ moveMessage(moving.user, currentGroupName, groupLabel(moving.target)) }}
        </p>
      </div>

      <template #footer>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          @click="moving.visible = false"
        />
        <Button
          :label="$t('system.user-groups.editor.members.move.submit')"
          :disabled="!moving.targetID"
          :loading="moving.saving"
          data-testid="member-move-submit"
          @click="submitMove"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { system } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'

const { CInputUser, CInputUserGroup } = components

const props = defineProps({
  userGroup: {
    type: Object,
    required: true,
  },
})

const { t } = useI18n()
const confirm = useConfirm()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const loading = ref(false)
const members = ref([])
const groupNames = ref({})

const moving = reactive({
  visible: false,
  saving: false,
  user: null,
  targetID: null,
  target: null,
})

const canManage = computed(() => !!props.userGroup.canManageMembersOnUserGroup)
const currentGroupName = computed(() => groupLabel(props.userGroup))

function groupLabel(g) {
  return g?.meta?.short || g?.handle || g?.userGroupID || ''
}

function userName(u) {
  return u.name || u.email || u.handle || u.userID
}

function moveMessage(user, from, to) {
  return t('system.user-groups.editor.members.move.confirm', { user: userName(user), from, to })
}

async function loadGroupNames() {
  try {
    const result = await $SystemAPI.userGroupList({ limit: 100 })
    groupNames.value = Object.fromEntries(
      (result?.set || []).map(g => [g.userGroupID, groupLabel(g)]),
    )
  } catch {
    // names are a hint; the picker still lists users without them
  }
}

async function groupNameOf(userGroupID) {
  if (groupNames.value[userGroupID]) return groupNames.value[userGroupID]
  try {
    return groupLabel(await $SystemAPI.userGroupRead({ userGroupID }))
  } catch {
    return userGroupID
  }
}

async function loadMembers() {
  loading.value = true
  try {
    const result = await $SystemAPI.userGroupMemberList({
      userGroupID: props.userGroup.userGroupID,
    })
    const ids = result?.set || result || []
    if (ids.length > 0) {
      const users = await Promise.all(
        ids.map(id => $SystemAPI.userRead({ userID: id }).catch(() => null)),
      )
      members.value = users.filter(Boolean).map(u => new system.User(u))
    } else {
      members.value = []
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.fetchError'))(e)
  } finally {
    loading.value = false
  }
}

async function join(user) {
  await $SystemAPI.userGroupMemberAdd({
    userGroupID: props.userGroup.userGroupID,
    userID: user.userID,
  })
  if (!members.value.find(m => m.userID === user.userID)) {
    members.value.push(new system.User({ ...user, userGroupID: props.userGroup.userGroupID }))
  }
}

async function addMember(user) {
  if (!user) return

  const fromID = user.userGroupID && user.userGroupID !== '0' ? user.userGroupID : null
  if (!fromID) {
    try {
      await join(user)
    } catch (e) {
      $toast.toastErrorHandler(t('system.user-groups.editor.members.addError'))(e)
    }
    return
  }

  const from = await groupNameOf(fromID)
  confirm.require({
    header: t('system.user-groups.editor.members.move.title'),
    message: moveMessage(user, from, currentGroupName.value),
    icon: 'pi pi-arrow-right-arrow-left',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: { label: t('system.user-groups.editor.members.move.submit'), size: 'small' },
    accept: async () => {
      try {
        await join(user)
        $toast.toastSuccess(
          t('system.user-groups.editor.members.move.success', {
            user: userName(user),
            to: currentGroupName.value,
          }),
        )
      } catch (e) {
        $toast.toastErrorHandler(t('system.user-groups.editor.members.move.error'))(e)
      }
    },
  })
}

function openMove(user) {
  Object.assign(moving, { visible: true, saving: false, user, targetID: null, target: null })
}

async function submitMove() {
  const { user, target } = moving
  if (!user || !target) return

  moving.saving = true
  try {
    await $SystemAPI.userGroupMemberAdd({ userGroupID: target.userGroupID, userID: user.userID })
    members.value = members.value.filter(m => m.userID !== user.userID)
    moving.visible = false
    $toast.toastSuccess(
      t('system.user-groups.editor.members.move.success', {
        user: userName(user),
        to: groupLabel(target),
      }),
    )
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.move.error'))(e)
  } finally {
    moving.saving = false
  }
}

onMounted(() => {
  loadMembers()
  if (canManage.value) loadGroupNames()
})
</script>
