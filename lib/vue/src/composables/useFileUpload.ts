import { ref } from 'vue'

interface UploadOptions {
  /** The API instance to use (e.g. $SystemAPI, $ComposeAPI) */
  api: any
  /** The endpoint path to POST to (from an endpoint method like pageUploadEndpoint) */
  endpoint: string
  /** Form field name for the file (default: 'upload') */
  fieldName?: string
  /** Additional form fields to include */
  extraFields?: Record<string, string>
}

interface UploadResult {
  attachmentID?: string
  name?: string
  meta?: any
  [key: string]: any
}

/**
 * Composable for handling file uploads with drag-and-drop support.
 *
 * Encapsulates the common upload pattern used across the codebase:
 * FormData construction, authenticated fetch, error handling.
 */
export function useFileUpload() {
  const uploading = ref(false)
  const uploadError = ref('')
  const dragOver = ref(false)

  function handleDragOver() {
    dragOver.value = true
  }

  function handleDragLeave() {
    dragOver.value = false
  }

  function handleDrop(event: DragEvent): File[] {
    dragOver.value = false
    return Array.from(event.dataTransfer?.files || [])
  }

  function handleFileSelect(event: Event): File[] {
    const input = event.target as HTMLInputElement
    const files = Array.from(input?.files || [])
    // Reset input so the same file can be re-selected
    if (input) input.value = ''
    return files
  }

  /**
   * Upload files to the given endpoint via the given API instance.
   *
   * Uses the API's axios instance for auth (bearer token) and base URL.
   */
  async function uploadFiles(files: File[], options: UploadOptions): Promise<UploadResult[]> {
    const { api, endpoint, fieldName = 'upload', extraFields } = options

    uploading.value = true
    uploadError.value = ''

    const results: UploadResult[] = []

    try {
      for (const file of files) {
        const formData = new FormData()
        formData.append(fieldName, file, file.name)

        if (extraFields) {
          for (const [key, value] of Object.entries(extraFields)) {
            formData.append(key, value)
          }
        }

        const { data } = await api
          .api()
          .post(endpoint, formData, { headers: { 'Content-Type': undefined } })

        if (data?.error) {
          throw new Error(data.error?.message || data.error)
        }

        const attachment = data?.response ?? data
        results.push(attachment)
      }
    } catch (err: any) {
      uploadError.value = err.message || 'Upload failed'
      throw err
    } finally {
      uploading.value = false
    }

    return results
  }

  /**
   * Upload a single file via raw fetch (for endpoints that need direct URL + token).
   *
   * This is useful for import endpoints or other non-standard API calls
   * where the API client's axios may not support the multipart format well.
   */
  async function uploadFileRaw(
    file: File,
    options: {
      url: string
      token?: string
      fieldName?: string
      extraFields?: Record<string, string>
    },
  ): Promise<any> {
    const { url, token, fieldName = 'upload', extraFields } = options

    uploading.value = true
    uploadError.value = ''

    try {
      const formData = new FormData()
      formData.append(fieldName, file)

      if (extraFields) {
        for (const [key, value] of Object.entries(extraFields)) {
          formData.append(key, value)
        }
      }

      const response = await fetch(url, {
        method: 'POST',
        body: formData,
        credentials: 'include',
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      })

      if (!response.ok) {
        const errData = await response.json().catch(() => ({}))
        throw new Error(errData?.error?.message || `Upload failed for ${file.name}`)
      }

      const data = await response.json()
      return data.response || data
    } catch (err: any) {
      uploadError.value = err.message || 'Upload failed'
      throw err
    } finally {
      uploading.value = false
    }
  }

  function reset() {
    uploading.value = false
    uploadError.value = ''
    dragOver.value = false
  }

  return {
    uploading,
    uploadError,
    dragOver,
    handleDragOver,
    handleDragLeave,
    handleDrop,
    handleFileSelect,
    uploadFiles,
    uploadFileRaw,
    reset,
  }
}
