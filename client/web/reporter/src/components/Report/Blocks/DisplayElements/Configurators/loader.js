import Text from './Text'
import Chart from './Chart'
import Table from './Table'
import Metric from './Metric'

const Registry = {
  Text,
  Chart,
  Table,
  Metric,
}

export default function (k) {
  return Registry[k]
}
