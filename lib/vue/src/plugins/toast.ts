import type { App } from 'vue'

export const ToastPlugin = {
  install(app: App) {
    const addToast = (options: any) => {
      app.config.globalProperties.$toast.add(options)
    }

    const toastSuccess = (message: string, title = 'Success') => {
      addToast({
        severity: 'success',
        summary: title,
        detail: message,
        life: 7000,
      })
    }

    const toastWarning = (message: string, title = 'Warning') => {
      addToast({
        severity: 'warn',
        summary: title,
        detail: message,
        life: 7000,
      })
    }

    const toastInfo = (message: string, title = 'Info') => {
      addToast({
        severity: 'info',
        summary: title,
        detail: message,
        life: 7000,
      })
    }

    const toastDanger = (message: string, title = 'Error') => {
      addToast({
        severity: 'error',
        summary: title,
        detail: message,
        life: 7000,
      })
    }

    const toastErrorHandler = (prefix = 'Error', title = 'Error') => {
      return (err: any = {}) => {
        const message = err.message || err.toString() || 'Unknown error'
        const msg = prefix ? `${prefix}: ${message}` : message
        toastDanger(msg, title)
        return message
      }
    }

    const toastService = {
      addToast,
      toastSuccess,
      toastWarning,
      toastInfo,
      toastDanger,
      toastErrorHandler,
    }

    app.provide('$toast', toastService)

    console.log('Toast plugin installed')
  },
}
