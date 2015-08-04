function startUserAdmin() {
    $('#autocomplete').autocomplete({
        source: function (request, autoCompleteResponse) {
            if (request.term.length >= 2) {
                new Http().Post('/users/search', {login: request.term}).success(function (users) {
                    for (var i = 0; i < users.length; ++i) {
                        var user = users[i];
                        user.label = user.firstName + ' ' + user.lastName;
                        user.value = user.login;
                    }
                    autoCompleteResponse(users);
                });
            }
        },
        select: function (event, object) {
            $('#autocomplete').val('');
            console.log(object);
            window.location.href = '/users/edit/' + object.item.login;
            return false;
        }
    });
}
