function startListStudent() {
    $('#name').val(getParameterByName('name'));
    $('#login').val(getParameterByName('login'));
    $('#email').val(getParameterByName('email'));

    var promoSelector = $('#promotions').select2({
        placeholder: 'Promotions',
        tags: true
    });
    var citiesSelector = $('#cities').select2({
        placeholder: 'Villes',
        tags: true
    });
    var skillSelector = $('#skills').select2({
        placeholder: 'Skills',
        tags: true
    });
    var themeSelector = $('#themes').select2({
        placeholder: 'Themes',
        tags: true
    });

    getParameterByName('promotions').split(',').forEach(function (promo) {
        if (promo) {
            promoSelector.append(new Option(promo, promo, true, true));
        }
    });
    getParameterByName('cities').split(',').forEach(function (city) {
        if (city) {
            citiesSelector.append(new Option(city, city, true, true));
        }
    });
    getParameterByName('skills').split(',').forEach(function (skill) {
        if (skill) {
            skillSelector.append(new Option(skill, skill, true, true));
        }
    });
    getParameterByName('themes').split(',').forEach(function (theme) {
        if (theme) {
            themeSelector.append(new Option(theme, theme, true, true));
        }
    });
    promoSelector.trigger('change');
    citiesSelector.trigger('change');
    themeSelector.trigger('change');
    skillSelector.trigger('change');
}

function performStudentSearch(limit, page, force) {
    limit = limit || getParameterByName('limit');
    page = page || getParameterByName('page');
    var promotions = $('#promotions').val();
    var cities = $('#cities').val();
    var email = $('#email').val() ? $('#email').val() : null;
    var login = $('#login').val() ? $('#login').val() : null;
    var name = $('#name').val() ? $('#name').val() : null;
    var skills = $('#skills').val() ? $('#skills').val() : null;
    var themes = $('#themes').val() ? $('#themes').val() : null;

    if (cities === null && promotions === null && name === null && themes === null &&
        skills === null && name === null && login === null && email === null && !force) {
        return;
    }

    var baseUrl = String.format('{0}{1}?page={2}&limit={3}', window.location.origin, window.location.pathname, page, limit);

    if (cities !== null) {
        baseUrl += String.format('&cities={0}', cities.join());
    }
    if (email !== null) {
        baseUrl += String.format('&email={0}', email);
    }
    if (promotions !== null) {
        baseUrl += String.format('&promotions={0}', promotions.join());
    }
    if (name !== null) {
        baseUrl += String.format('&name={0}', name);
    }
    if (skills !== null) {
        baseUrl += String.format('&skills={0}', encodeURIComponent(skills.join()));
    }
    if (themes !== null) {
        baseUrl += String.format('&themes={0}', encodeURIComponent(themes.join()));
    }
    if (login !== null) {
        baseUrl += String.format('&login={0}', login);
    }
    window.location.href = baseUrl;
}

function nextPage() {
    performStudentSearch(getParameterByName('limit'), parseInt(getParameterByName('page')) + 1, true);
    return false;
}

function previousPage() {
    performStudentSearch(getParameterByName('limit'), parseInt(getParameterByName('page')) - 1, true);
    return false;
}
