function validateProjects() {
    var form = document.getElementById('form');
    form.action = '/pedago/validate';
    form.submit();
    window.location.reload();
}

function switchValidationTabs() {
    var path = window.location.pathname;
    var indeterminate = $('#indeterminate');
    var refused = $('#refused');
    var validated = $('#validated');
    indeterminate.removeClass('active');
    refused.removeClass('active');
    validated.removeClass('active');
    if (~path.indexOf('indeterminate') !== 0) {
        indeterminate.addClass('active');
    } else if (~path.indexOf('refused') !== 0) {
        refused.addClass('active');
    } else {
        validated.addClass('active');
    }
}

