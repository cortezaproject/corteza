import { useConfirm } from 'primevue/useconfirm'
import { useI18n } from 'vue-i18n'

interface ConfirmDeleteOptions {
  message: string
  header?: string
  icon?: string
  acceptProps?: Record<string, unknown>
  rejectProps?: Record<string, unknown>
  onConfirm: () => void
}

export function useConfirmDelete() {
  const confirm = useConfirm()
  const { t } = useI18n()

  function confirmDelete(options: ConfirmDeleteOptions) {
    const {
      message,
      header = '',
      icon = 'pi pi-exclamation-triangle',
      acceptProps = {},
      rejectProps = {},
      onConfirm,
    } = options

    confirm.require({
      message,
      header,
      icon,
      acceptProps: {
        label: t('general.label.delete'),
        severity: 'danger',
        size: 'small',
        ...acceptProps,
      },
      rejectProps: {
        label: t('general.label.cancel'),
        severity: 'secondary',
        text: true,
        size: 'small',
        ...rejectProps,
      },
      accept: onConfirm,
    })
  }

  return { confirmDelete }
}
