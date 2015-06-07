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

function getParameterByName(name) {
    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
        results = regex.exec(location.search);
    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
}

function randomizeLabel() {
    var labels = ['primary', 'info', 'default', 'danger', 'warning', 'success'];
    var random = Math.floor(Math.random() * labels.length) + 1;
    return labels[random];
}

if (!String.format) {
    String.format = function (format) {
        var args = Array.prototype.slice.call(arguments, 1);
        return format.replace(/{(\d+)}/g, function (match, number) {
            return typeof args[number] != 'undefined'
                ? args[number]
                : match
                ;
        });
    };
}

$('document').ready(switchActiveTab);