// File fields defer their upload to record save: the editor only stages the
// File objects, and whoever saves the record uploads them and writes the
// returned attachment IDs into the value. Every saving view needs the same two
// steps, so they live here rather than being copied per view.

/**
 * Upload one staged file and return its attachmentID.
 *
 * An empty recordID is how a not-yet-created record claims its attachments —
 * the upload is accepted unbound and the server ties it to the record on save,
 * which is why callers may upload before create as well as after.
 */
export async function uploadRecordAttachment(
  $ComposeAPI,
  { namespaceID, moduleID, recordID, fieldName, file },
) {
  const url = $ComposeAPI.recordUploadEndpoint({ namespaceID, moduleID })
  const formData = new FormData()
  formData.append('recordID', recordID || '')
  formData.append('fieldName', fieldName)
  formData.append('upload', file, file.name)

  const { data } = await $ComposeAPI
    .api()
    .post(url, formData, { headers: { 'Content-Type': undefined } })

  if (data?.error) throw new Error(data.error)

  const attachment = data?.response ?? data
  if (!attachment?.attachmentID) {
    throw new Error(`Upload failed for "${file.name}": no attachmentID in response`)
  }
  return attachment.attachmentID
}

/**
 * Append new attachment IDs to whatever the field already holds, normalising
 * the single-value and multi-value shapes to one array.
 */
export function mergeAttachmentIDs(existing, ids) {
  const existingIDs = Array.isArray(existing)
    ? existing.filter(Boolean)
    : existing
      ? [existing]
      : []
  return [...existingIDs, ...ids]
}
