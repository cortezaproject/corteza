<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('system.users.editor.password.title')"
    class="w-full max-w-md"
    data-testid="user-password-dialog"
    @update:visible="close"
    @after-hide="reset"
  >
    <Form
      :key="formKey"
      :resolver="resolver"
      :initialValues="{ password: '', confirmPassword: '' }"
      class="flex flex-col gap-4"
      @submit="handleSubmit"
    >
      <p class="text-sm text-muted-color">
        {{ $t('system.users.editor.password.description', { name: userLabel }) }}
      </p>

      <CFormGroup
        name="password"
        input-id="user-password-new"
        :label="$t('system.users.editor.password.new')"
        required
      >
        <Password
          input-id="user-password-new"
          name="password"
          v-model="password"
          toggleMask
          :feedback="false"
          :inputProps="{ autocomplete: 'new-password' }"
          inputClass="w-full"
          class="w-full"
          autofocus
        />
      </CFormGroup>

      <CFormGroup
        name="confirmPassword"
        input-id="user-password-confirm"
        :label="$t('system.users.editor.password.confirm')"
        required
      >
        <Password
          input-id="user-password-confirm"
          name="confirmPassword"
          v-model="confirmPassword"
          toggleMask
          :feedback="false"
          :inputProps="{ autocomplete: 'new-password' }"
          inputClass="w-full"
          class="w-full"
        />
      </CFormGroup>

      <div class="flex justify-end gap-2 pt-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          :disabled="saving"
          @click="close(false)"
        />
        <Button
          type="submit"
          :label="$t('system.users.editor.password.set')"
          icon="pi pi-key"
          :loading="saving"
        />
      </div>
    </Form>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  user: { type: Object, required: true },
})

const emit = defineEmits(['update:visible', 'saved'])

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const password = ref('')
const confirmPassword = ref('')
const saving = ref(false)
const formKey = ref(0)

const userLabel = computed(
  () => props.user.name || props.user.handle || props.user.email || props.user.userID,
)

function resolver({ values }) {
  const errors = {}

  if (!values.password) {
    errors.password = [{ message: t('general.label.required') }]
  }

  if (values.password && values.confirmPassword !== values.password) {
    errors.confirmPassword = [{ message: t('system.users.editor.password.missmatch') }]
  }

  return { errors }
}

function close(value = false) {
  if (!value && saving.value) return
  emit('update:visible', value)
}

function reset() {
  password.value = ''
  confirmPassword.value = ''
  formKey.value++
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    await $SystemAPI.userSetPassword({ userID: props.user.userID, password: password.value })
    $toast.toastSuccess(t('notification.user.passwordChange.success'))
    emit('saved')
    saving.value = false
    close(false)
  } catch (e) {
    console.error('Failed to set password:', e)
    $toast.toastErrorHandler(t('notification.user.passwordChange.error'))(e)
  } finally {
    saving.value = false
  }
}
</script>
