document.addEventListener('DOMContentLoaded', function () {
  // --- Toast auto-show & auto-dismiss ---
  document.querySelectorAll('.toast').forEach(function (toast) {
    var delay = parseInt(toast.getAttribute('data-autohide') || '5000', 10);
    setTimeout(function () {
      toast.style.transition = 'opacity 0.3s ease';
      toast.style.opacity = '0';
      setTimeout(function () { toast.remove(); }, 300);
    }, delay);
  });

  // --- Toast close button ---
  document.addEventListener('click', function (e) {
    var btn = e.target.closest('[data-action="close-toast"]');
    if (btn) {
      var toast = btn.closest('.toast');
      if (toast) {
        toast.style.transition = 'opacity 0.2s ease';
        toast.style.opacity = '0';
        setTimeout(function () { toast.remove(); }, 200);
      }
    }
  });

  // --- Disable buttons on form submit ---
  document.querySelectorAll('form').forEach(function (form) {
    if (form.classList.contains('do-not-disable-buttons-on-submit')) return;
    form.addEventListener('submit', function () {
      setTimeout(function () {
        form.querySelectorAll('button, input[type=submit]').forEach(function (el) {
          if (!el.classList.contains('do-not-disable-on-submit')) {
            el.disabled = true;
          }
        });
      }, 50);
    });
  });

  // --- MFA code mask (6-digit with space: "000 000") ---
  document.querySelectorAll('input.mfa-code-mask').forEach(function (input) {
    input.addEventListener('input', function () {
      var raw = input.value.replace(/[^0-9]/g, '').slice(0, 6);
      if (raw.length > 3) {
        input.value = raw.slice(0, 3) + ' ' + raw.slice(3);
      } else {
        input.value = raw;
      }
    });

    // Support pasting
    input.addEventListener('paste', function (e) {
      e.preventDefault();
      var pasted = (e.clipboardData || window.clipboardData).getData('text');
      var raw = pasted.replace(/[^0-9]/g, '').slice(0, 6);
      if (raw.length > 3) {
        input.value = raw.slice(0, 3) + ' ' + raw.slice(3);
      } else {
        input.value = raw;
      }
    });
  });

  // --- Avatar file input auto-submit ---
  var avatarInput = document.getElementById('avatar');
  if (avatarInput) {
    avatarInput.addEventListener('change', function () {
      avatarInput.closest('form').submit();
    });
  }

  // --- Handle mask (from jQuery .handle-mask) ---
  document.querySelectorAll('input.handle-mask').forEach(function (input) {
    input.addEventListener('input', function () {
      input.value = input.value.replace(/[^a-zA-Z0-9_.-]/g, '').toLowerCase();
    });
  });

  // --- Color picker hex value sync + avatar initials preview ---
  var avatarEl = document.getElementById('avatarInitials');

  // Compute initials from the name field (same logic as topbar)
  if (avatarEl) {
    var nameInput = document.querySelector('input[name="name"]');
    var emailInput = document.querySelector('input[name="email"]');
    var name = nameInput ? nameInput.value : '';
    var initials = '';

    if (name) {
      initials = name.split(' ').map(function (n) { return n[0]; }).join('').toUpperCase().slice(0, 2);
    } else if (emailInput && emailInput.value) {
      initials = emailInput.value[0].toUpperCase();
    }

    avatarEl.textContent = initials || '?';

    // Update initials when name changes
    if (nameInput) {
      nameInput.addEventListener('input', function () {
        var val = nameInput.value;
        if (val) {
          avatarEl.textContent = val.split(' ').map(function (n) { return n[0]; }).join('').toUpperCase().slice(0, 2);
        } else {
          avatarEl.textContent = '?';
        }
      });
    }
  }

  document.querySelectorAll('.color-picker-input').forEach(function (input) {
    input.addEventListener('input', function () {
      // Update hex label
      var label = document.querySelector('[data-color-for="' + input.id + '"]');
      if (label) {
        label.textContent = input.value;
      }

      // Live-update avatar initials preview colors
      if (avatarEl) {
        if (input.id === 'initialColor') {
          avatarEl.style.color = input.value;
        } else if (input.id === 'customColor') {
          avatarEl.style.backgroundColor = input.value;
        }
      }
    });
  });
});
