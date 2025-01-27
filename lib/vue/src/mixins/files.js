import mime from 'mime'

export default {
  methods: {
    /**
     * Checks if given file
     * @param {String} fileName File in question
     * @param {Array<String>} accepted Array of accepted mime types
     * @returns {Boolean} If this file is acceptable
     */
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    validateFileType (fileName = '', accepted = ['*/*']) {
      const t = mime.getType(fileName)
      return !!accepted.find(at => {
        at = at.split('/')
        at[0] = at[0] === '*' ? '.*?' : at[0]
        at[1] = at[1] === '*' ? '.*?' : at[1]
        return (new RegExp(at.join('\\/'))).exec(t)
      })
    },
  },
}
