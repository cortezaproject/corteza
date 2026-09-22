import logger from '../logger'
import { execContext } from '../config'
import axios from 'axios'

export default function (): void {
  const csCtx = execContext.cortezaServers
  const feCtx = execContext.frontend

  const log = logger.child({ name: 'check' })

  log.debug(csCtx.system, 'configuring cServer system API')
  log.debug(csCtx.compose, 'configuring cServer compose API')

  let versionEndpoint = ''
  if (csCtx.system.apiBaseURL) {
    versionEndpoint = csCtx.system.apiBaseURL
      .replace('api/system', 'version')
      .replace('system', 'version')
  } else {
    log.error(
      'Could not auto configure Corredor. Please, configure your Human API base URL by setting CORREDOR_EXEC_CSERVERS_API_HOST and/or CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE env variables',
    )
    throw new Error('Cannot run Corredor due to misconfiguration')
  }

  axios
    .get(versionEndpoint)
    .then(({ data: { response } }) => {
      log.info(response, `Assuming valid Human API at ${csCtx.system.apiBaseURL}`)
    })
    .catch(r => {
      log.error(
        `check your Human API settings: expecting valid response for ${versionEndpoint} got: ${r}`,
      )
    })

  log.info('checking API endpoints')

  log.debug(feCtx, 'frontend settings')
}
