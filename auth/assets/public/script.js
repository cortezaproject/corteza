$(function () {
  $('.toast').toast('show');

  $('form').not('do-not-disable-buttons-on-submit').on('submit', function() {
    let form = this
    setTimeout(function() {
      $('button, input[type=submit]', form)
      .not('do-not-disable-on-submit')
      .attr('disabled', true)
    }, 50)
  })

  $('input.mfa-code-mask').mask('000 000')
})
