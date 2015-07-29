function startListProject() {
    $('#name').val(getParameterByName('name'));
    $('#student').val(getParameterByName('student'));

    var promoSelector = $('#promotions').select2({
        placeholder: 'Promotion',
        tags: true
    });
    var citiesSelector = $('#cities').select2({
        placeholder: 'Ville',
        tags: true
    });
    var managersSelector = $('#managers').select2({
        placeholder: 'Manageur',
        tags: true
    });
    var statusSelector = $('#status').select2({
        placeholder: 'Statut',
        tags: true
    });
    var themesSelector = $('#themes').select2({
        placeholder: 'Theme',
        tags: true
    });
    var technosSelector = $('#technos').select2({
        placeholder: 'Techno',
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
    getParameterByName('technos').split(',').forEach(function (techno) {
        if (techno) {
            technosSelector.append(new Option(techno, techno, true, true));
        }
    });
    getParameterByName('themes').split(',').forEach(function (theme) {
        if (theme) {
            themesSelector.append(new Option(theme, theme, true, true));
        }
    });
    var statuses = {};
    // generation of status and managers via the backend :(
    // [[range $element := .Status]]
    statuses['[[ $element ]]'] = '[[ i18n $.Lang $element ]]';
    // [[end]]
    var managers = {};
    // [[range $element := .Managers]]
    managers['[[ $element.Login ]]'] = '[[ $element.FirstName ]] [[ $element.LastName ]]';
    // [[end]]
    getParameterByName('status').split(',').forEach(function (status) {
        if (status) {
            statusSelector.append(new Option(statuses[status], status, true, true));
        }
    });
    getParameterByName('managers').split(',').forEach(function (manager) {
        if (manager) {
            managersSelector.append(new Option(managers[manager], manager, true, true));
        }
    });
    themesSelector.trigger('change');
    technosSelector.trigger('change');
    promoSelector.trigger('change');
    citiesSelector.trigger('change');
    statusSelector.trigger('change');
    managersSelector.trigger('change');
}

function performProjectSearch(limit, page, force) {
    limit = limit || getParameterByName('limit');
    page = page || getParameterByName('page');
    var promotions = $('#promotions').val();
    var cities = $('#cities').val();
    var name = $('#name').val() ? $('#name').val() : null;
    var managers = $('#managers').val();
    var status = $('#status').val();
    var technos = $('#technos').val();
    var themes = $('#themes').val();
    var student = $('#student').val() ? $('#student').val() : null;

    if (cities === null && promotions === null && name === null && managers === null &&
        status === null && student === null && !force && technos === null && themes === null) {
        return;
    }
    var baseUrl = String.format('{0}{1}?page={2}&limit={3}', window.location.origin, window.location.pathname, page, limit);

    if (cities !== null) {
        baseUrl += String.format('&cities={0}', cities.join());
    }
    if (promotions !== null) {
        baseUrl += String.format('&promotions={0}', promotions.join());
    }
    if (name !== null) {
        baseUrl += String.format('&name={0}', name);
    }
    if (managers !== null) {
        baseUrl += String.format('&managers={0}', managers.join());
    }
    if (status !== null) {
        baseUrl += String.format('&status={0}', status.join());
    }
    if (student !== null) {
        baseUrl += String.format('&student={0}', student);
    }
    if (themes !== null) {
        baseUrl += String.format('&themes={0}', encodeURIComponent(themes.join()));
    }
    if (technos !== null) {
        baseUrl += String.format('&technos={0}', encodeURIComponent(technos.join()));
    }
    window.location.href = baseUrl;
}

function nextPage() {
    performProjectSearch(getParameterByName('limit'), parseInt(getParameterByName('page')) + 1, true);
    return false;
}

function previousPage() {
    performProjectSearch(getParameterByName('limit'), parseInt(getParameterByName('page')) - 1, true);
    return false;
}
