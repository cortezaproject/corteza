export default {
  label: 'Sink: hello (server)',
  description: 'Answers a signed sink request on /hello',

  triggers({ on }) {
    return on('request').for('system:sink').where('request.path', '/hello')
  },

  exec({ $request, $response }) {
    $response.status = 200
    $response.header = { 'Content-Type': ['text/plain'] }
    $response.body = `hello from corredor via ${$request.method}`

    return $response
  },
}
