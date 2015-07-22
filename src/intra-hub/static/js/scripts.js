/**
 * Created by Vincent on 27/05/15.
 */


function switchActiveTab() {
    var path = window.location.pathname;
    var projects = $('#projects');
    var hub = $('#hub');
    var staff = $('#staff');
    var students = $('#students');
    projects.removeClass('active');
    hub.removeClass('active');
    staff.removeClass('active');
    students.removeClass('active');
    if (~path.indexOf('projects') !== 0) {
        projects.addClass('active');
    } else if (~path.indexOf('hub')) {
        hub.addClass('active');
    } else if (~path.indexOf('staff')) {
        staff.addClass('active');
    } else if (~path.indexOf('students')) {
        students.addClass('active');
    }
}

function getParameterByName(name) {
    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
        results = regex.exec(location.search);
    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
}

function randomizeLabel() {
    var labels = ['primary', 'info', 'default', 'danger', 'warning', 'success'];
    var random = Math.floor(Math.random() * (labels.length - 1)) + 1;
    return labels[random];
}

if (!String.format) {
    String.format = function (format) {
        var args = Array.prototype.slice.call(arguments, 1);
        return format.replace(/{(\d+)}/g, function (match, number) {
            return typeof args[number] != 'undefined' ? args[number] : match;
        });
    };
}

function errorImage(img) {
    img.src = 'https://cloud.canadastays.com/assets/images/user-placeholder.png';
}

$('document').ready(switchActiveTab);