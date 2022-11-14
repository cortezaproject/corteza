/* eslint-disable @typescript-eslint/explicit-function-return-type */
import mime from 'mime-types'

const types = [
  { type: 'application/pdf', component: 'PDF' },
  { type: 'image/', component: 'IMG' },
]

/**
 * Tells what component (if any) can preview the given file
 * @param {String|undefined} type pre defined mime type of object
 * @param {String|undefined} src object's src
 * @param {String|undefined} name object's name
 * @returns {String|undefined} preview component or undefined
 */
export const getComponent = ({ type, src, name }) => {
  const srcType = type || mime.lookup(src) || mime.lookup(name)
  if (!srcType) {
    return
  }

  for (const { type, component } of types) {
    if (srcType.indexOf(type) >= 0) {
      return component
    }
  }
}

/**
 * Tells if we support the given file type preview
 * @param {String|undefined} type pre defined mime type of object
 * @param {String|undefined} src object's src
 * @param {String|undefined} name object's name
 * @returns {Boolean} if file can be previewed
 */
export const canPreview = ({ type, src, name }) => {
  return !!getComponent({ type, src, name })
}
