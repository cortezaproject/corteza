import ComposeRecordConfigurator from './ComposeRecordConfigurator'

const Registry = {
  composeRecords: ComposeRecordConfigurator,
}

export default function (k) {
  return Registry[k]
}
