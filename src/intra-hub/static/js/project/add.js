function startAddProject() {
    var themesSelect = $('#themes');
    themesSelect.select2({
        placeholder: 'Themes',
        tags: true
    });
    themesSelect.on('select2:select', function (event) {
        var themeSelect = $('#themesId');
        var themeId = event.params.data.id;
        var currentIds = themeSelect.attr('value');
        if (!currentIds) {
            themeSelect.attr('value', themeId);
        } else {
            themeSelect.attr('value', currentIds + ',' + themeId);
        }
    });
    themesSelect.on('select2:unselect', function (event) {
        var themeSelect = $('#themesId');
        var currentIds = themeSelect.attr('value');
        var ids = currentIds.split(',');
        for (var i = 0; i < ids.length; ++i) {
            if (event.params.data.id === ids[i]) {
                ids.splice(i, 1);
                break
            }
        }
        themeSelect.attr('value', ids.join());
    });

    var technoSelect = $('#technos');
    technoSelect.select2({
        placeholder: 'Technos',
        tags: true
    });
    technoSelect.on('select2:select', function (event) {
        var technoSelect = $('#technosId');
        var id = event.params.data.id;
        var currentIds = technoSelect.attr('value');
        if (!currentIds) {
            technoSelect.attr('value', id);
        } else {
            technoSelect.attr('value', currentIds + ',' + id);
        }
    });
    technoSelect.on('select2:unselect', function (event) {
        var technoSelect = $('#technosId');
        var id = technoSelect.attr('value');
        var currentIds = id.split(',');
        for (var i = 0; i < currentIds.length; ++i) {
            if (event.params.data.id === currentIds[i]) {
                currentIds.splice(i, 1);
                break
            }
        }
        technoSelect.attr('value', currentIds.join());
    });

    $('#autocomplete').autocomplete({
        source: function (request, autoCompleteResponse) {
            if (request.term.length >= 2) {
                new Http().Post('/users/search', {login: request.term}).success(function (users) {
                    var currentIds = $('#membersId').attr('value');
                    var arrayIds;
                    if (currentIds) {
                        arrayIds = currentIds.split(',');
                    }
                    for (var i = 0; i < users.length; ++i) {
                        var user = users[i];
                        // If the user is already is the project, removes it.
                        if (arrayIds && ~arrayIds.indexOf(String(user.id)) !== 0) {
                            users.splice(i, 1);
                            continue;
                        }
                        user.label = user.firstName + ' ' + user.lastName;
                        user.value = user.login;
                    }
                    autoCompleteResponse(users);
                });
            }
        },
        select: function (event, object) {
            var membersId = $('#membersId');
            var currentIds = membersId.attr('value');
            if (!currentIds) {
                membersId.attr('value', object.item.id);
            } else {
                membersId.attr('value', currentIds + ',' + object.item.id);
            }
            $('#logins').append(generateDeleteMemberListItem(object.item));
            $('#autocomplete').val('');
            return false;
        }
    });

    $('#form').validate({
        onsubmit: true,
        onfocusout: true,
        rules: {
            name: {
                required: true,
                remote: {
                    url: '/projects/checkname'
                }
            }
        },
        messages: {
            name: {
                remote: 'Le nom est deja pris'
            }
        }
    });
    $("form").validate({
        onfocusout: function (element) {
            $(element).valid()
        }
    });
}

function generateDeleteMemberListItem(user) {
    var deleteButton = '<span ' + 'onclick="deleteLogin(\'' + user.login + '\', \'' + user.id + '\')"' + 'class="pull-right glyphicon glyphicon-remove red clickable"></span>';
    var image = '<img onerror="errorImage(this)" class="img-responsive size" src="' + user.picture + '">';
    return '<li id="' + user.login + '" class="list-group-item"><div>' + image + user.login + deleteButton + '</div></li>';
}

function deleteLogin(login, id) {
    var membersId = $('#membersId');
    var currentIds = membersId.attr('value');
    var ids = currentIds.split(',');
    for (var i = 0; i < ids.length; ++i) {
        if (ids[i] === id) {
            ids.splice(i, 1);
        }
    }
    membersId.attr('value', ids);
    $('#' + login).remove();
}
