
function openThemeModal(id) {
    $('#theme-edit-name').val($('#theme-'+id).html());
    $('#theme-edit-id').val(id)
}

function editTheme() {
    var obj = {
        id: parseInt($('#theme-edit-id').val()),
        name: $('#theme-edit-name').val()
    };
    if (obj.name === '') {
        return;
    }
    new Http().Put('/themes', obj).success(function (theme) {
        $('#theme-' + theme.id).html(theme.name);
        $('#themeModal').modal('hide');
    });
}

function addTheme() {
    var th = $('#theme');
    var themeName = th.val().trim();
    if (themeName === '') {
        return;
    }
    th.val('');
    new Http().Post('/themes', {name: themeName}).success(function (theme) {
        appendTheme(theme);
        $('#empty').hide();
    });
}

function appendTheme(theme) {
    $('#list').append('<li class="list-group-item" id="' + theme.id + '"><span id="theme-'+theme.id+'"class="label label-' + randomizeLabel() +
        ' s-1x" >' + theme.name + '</span>' +
        '<span class="pull-right"><span class="btn-group btn-group-xs">' +
        '<button class="btn btn-warning" data-toggle="modal" data-target="#themeModal"' +
        'onclick="openThemeModal(' + theme.id + ')">' +
        '<span class="fa fa-pencil color-black"></span>' +
        '</button> ' +
        '<button class="btn btn-danger" onclick=deleteTheme(' + theme.id + ')><span class="fa fa-remove color-white"></span></button>' +
        '<span/><span/></li> ');
}

function deleteTheme(id) {
    new Http().Delete('/themes/' + id).success(function () {
        $('#' + id).remove();
    });
}