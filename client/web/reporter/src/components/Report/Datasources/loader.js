import { default as Load } from './Load'
import { default as Join } from './Join'
import { default as Group } from './Group'

const Registry = {
  load: Load,
  join: Join,
  group: Group,
}

export default function (k) {
  return Registry[k]
}
