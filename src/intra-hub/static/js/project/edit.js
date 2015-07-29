function startEditProject() {
    var members = JSON.parse("[[ toJSON .Project.Members ]]");
    var themes = JSON.parse("[[ toJSON .Project.Themes ]]");
    var technos = JSON.parse("[[ toJSON .Project.Technos ]]");
    members.forEach(function (member) {
        addMember(member.id);
    });
    themes.forEach(function (theme) {
        addTheme(theme.id);
    });
    technos.forEach(function (techno) {
        addTechno(techno.id);
    });
    var themesSelect = $('#themes');
    themesSelect.select2({
        placeholder: 'Themes',
        tags: true
    });
    themesSelect.on('select2:select', function (event) {
        addTheme(event.params.data.id);
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
        addTechno(event.params.data.id);
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
            addMember(object.item.id);
            $('#logins').append(generateDeleteMemberListItem(object.item));
            $('#autocomplete').val('');
            return false;
        }
    });
}

function addMember(id) {
    var membersId = $('#membersId');
    var currentIds = membersId.attr('value');
    if (!currentIds) {
        membersId.attr('value', id);
    } else {
        membersId.attr('value', currentIds + ',' + id);
    }
}

function addTheme(id) {
    var themeSelect = $('#themesId');
    var currentIds = themeSelect.attr('value');
    if (!currentIds) {
        themeSelect.attr('value', id);
    } else {
        themeSelect.attr('value', currentIds + ',' + id);
    }
}

function addTechno(id) {
    var technoSelect = $('#technosId');
    var currentIds = technoSelect.attr('value');
    if (!currentIds) {
        technoSelect.attr('value', id);
    } else {
        technoSelect.attr('value', currentIds + ',' + id);
    }
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
    $('#' + id).remove();
}
