import bundler from './bundler'
import execContext from './exec-ctx'

export default function (): void {
  execContext()
  bundler()
}
