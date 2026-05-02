document.addEventListener('DOMContentLoaded', function () {
  var form = document.querySelector('[data-contact-form]');
  if (!form) return;

  form.addEventListener('submit', function (event) {
    event.preventDefault();

    var data = new FormData(form);
    var first = (data.get('first_name') || '').toString().trim();
    var last = (data.get('last_name') || '').toString().trim();
    var email = (data.get('email') || '').toString().trim();
    var telephone = (data.get('telephone') || '').toString().trim();
    var comments = (data.get('comments') || '').toString().trim();

    var subject = encodeURIComponent('Website contact from ' + [first, last].filter(Boolean).join(' '));
    var body = [
      'First Name: ' + first,
      'Last Name: ' + last,
      'Email Address: ' + email,
      'Telephone Number: ' + telephone,
      '',
      comments
    ].join('
');

    window.location.href = 'mailto:writersteve49@gmail.com?subject=' + subject + '&body=' + encodeURIComponent(body);
  });
});
