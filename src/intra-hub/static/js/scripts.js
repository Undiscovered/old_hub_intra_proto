/**
 * Created by Vincent on 27/05/15.
 */


function switchActiveTab() {
    var path = window.location.pathname;
    var projects = $('#projects');
    var hub = $('#hub');
    projects.removeClass('active');
    hub.removeClass('active');
    if (~path.indexOf('projects') !== 0) {
        projects.addClass('active');
    } else {
        hub.addClass('active');
    }
}

$('document').ready(switchActiveTab);