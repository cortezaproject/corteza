// placeholder

function buttonDisabler() {
  setTimeout(function() {
    let bb = document.getElementsByTagName('button')
    for (let i = 0; i < bb.length; i++) {
      bb[i].disabled = true
    }
  }, 100)
}
