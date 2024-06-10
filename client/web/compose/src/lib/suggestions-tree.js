/* eslint-disable no-template-curly-in-string */
export const getRecordBasedSuggestions = ({ moduleFields = [], userProperties = [], recordProperties = undefined, recordValueFields = undefined, isRecordPage = undefined }) => {
  const result = {
    '': ['${userID}', '${user.', 'AND', 'OR'].concat(moduleFields), // eslint-disable-line no-template-curly-in-string
    '$': ['${userID}', '${user.'].map(v => ({ // eslint-disable-line no-template-curly-in-string
      caption: v.trimStart('$'),
      value: v,
    })),
    '${': ['${userID}', '${user.'].map(v => ({ // eslint-disable-line no-template-curly-in-string
      caption: v.trimStart('${'),
      value: v,
    })),
    '${user': userProperties.map(m => {
      return {
        caption: m.trimStart('${user'),
        value: '${user.' + m + '}',
      }
    }),
  }

  if (recordProperties && recordValueFields && isRecordPage) {
    result[''] = ['${record.', '${recordID}', '${ownerID}'].concat(result[''])

    result['$'] = ['${record.', '${recordID}', '${ownerID}'].map(v => ({
      caption: v.trimStart('$'),
      value: v,
    })).concat(result['$'])

    result['${'] = ['${record.', '${recordID}', '${ownerID}'].map(v => ({
      caption: v.trimStart('$'),
      value: v,
    })).concat(result['${'])

    result['${record'] = recordProperties.map(v => ({
      caption: v.trimStart('${record.'),
      value: '${record.' + v + `${v === 'values' ? '.' : '}'}`,
    }))

    result['${record.values'] = recordValueFields.map(v => ({
      caption: v.trimStart('${record.values'),
      value: '${record.values.' + v + '}',
    }))
  }

  return result
}

export const getVisibilityConditionBasedSuggestions = ({ moduleFields = [], userProperties = [], recordProperties = undefined, recordValueFields = undefined, isRecordPage = undefined }) => {
  const result = {
    '': ['user', 'record', 'screen']
      .map(v => ({ caption: v, value: v + '.' }))
      .concat(['oldLayout', 'layout', 'isView', 'isCreate', 'isEdit'].concat(moduleFields)),
    'user': userProperties.map(v => ({
      caption: v,
      value: 'user.' + v,
    })),
    'screen': ['width', 'height', 'userAgent', 'breakpoint'].map(v => ({
      caption: v,
      value: 'screen.' + v,
    })),
  }

  if (recordProperties && recordValueFields && isRecordPage) {
    result['record'] = recordProperties.map(v => ({
      caption: v,
      value: 'record.' + v,
    }))
    result['record.values'] = recordValueFields.map(v => ({
      caption: v,
      value: 'record.values.' + v,
    }))
  }

  return result
}
