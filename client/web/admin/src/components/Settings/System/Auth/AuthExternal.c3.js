const props = {
  // OIDC
  title: 'Title',
  value: {
    enabled: true,
    handle: 'Handle',
    issuer: 'Mara',
    key: 'key',
    secret: 'shhhhh',
    cert: 'cert himself',
    idp: {
      url: 'https://kati.com',
      'ident-name': 'kati',
      'ident-handle': 'kati',
      'ident-identifier': 'kati',
    },
  },
}

const controls = [
  {
    label: 'Title',
    type: 'b-form-input',
    value: props.title,
    handle: (props, value) => {
      props.title = value
    },
  },
  {
    label: 'Enabled',
    type: 'b-form-checkbox',
    value: props.value.enabled,
    handle: (props, value) => {
      props.enabled = value
    },
  },
  {
    label: 'Handle',
    type: 'b-form-input',
    value: props.value.handle,
    handle: (props, value) => {
      props.handle = value
    },
  },
  {
    label: 'Issuer',
    type: 'b-form-input',
    value: props.value.issuer,
    handle: (props, value) => {
      props.issuer = value
    },
  },
  {
    label: 'Key',
    type: 'b-form-input',
    value: props.value.key,
    handle: (props, value) => {
      props.key = value
    },
  },
  {
    label: 'Secret',
    type: 'b-form-input',
    value: props.value.secret,
    handle: (props, value) => {
      props.secret = value
    },
  },
  {
    label: 'Cert',
    type: 'b-form-input',
    value: props.value.cert,
    handle: (props, value) => {
      props.cert = value
    },
  },
  {
    label: 'url',
    type: 'b-form-input',
    value: props.value.idp.url,
    handle: (props, value) => {
      props.value.idp.url = value
    },
  },
  {
    label: 'Name',
    type: 'b-form-input',
    value: props.value.idp['ident-name'],
    handle: (props, value) => {
      props.value.idp['ident-name'] = value
    },
  },
  {
    label: 'Handle',
    type: 'b-form-input',
    value: props.value.idp['ident-handle'],
    handle: (props, value) => {
      props.value.idp['ident-handle'] = value
    },
  },
  {
    label: 'Identifier',
    type: 'b-form-input',
    value: props.value.idp['ident-identifier'],
    handle: (props, value) => {
      props.value.idp['ident-handle'] = value
    },
  },
]

const scenarios = [
  { label: 'full form',
    props,
  },
  { label: 'empty form',
    props: {
      title: 'null',
      value: {
        enabled: true,
        handle: 'null',
        issuer: 'null',
        key: 'null',
        secret: 'null',
        cert: 'null',
        idp: {
          url: 'null',
          'ident-name': 'null',
          'ident-handle': 'null',
          'ident-identifier': 'null',
        },
      },
    },
  },
]

export default {
  props,
  controls,
  scenarios,
}
