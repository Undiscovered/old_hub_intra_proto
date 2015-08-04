function openSkillModal(id) {
    $('#skill-edit-name').val($('#skill-' + id).html());
    $('#skill-edit-id').val(id)
}

function editSkill() {
    var obj = {
        id: parseInt($('#skill-edit-id').val()),
        name: $('#skill-edit-name').val()
    };
    if (obj.name === '') {
        return;
    }
    new Http().Put('/skills', obj).success(function (skill) {
        $('#skill-' + skill.id).html(skill.name);
        $('#skillModal').modal('hide');
    });
}

function addSkill() {
    var th = $('#skill');
    var skillName = th.val().trim();
    if (skillName === '') {
        return;
    }
    th.val('');
    new Http().Post('/skills', {name: skillName}).done(function (skill) {
        appendSkill(skill);
        $('#empty').hide();
    }).fail(function () {
    });
}

function appendSkill(skill) {
    $('#list-skills').append('<li class="list-group-item" id="' + skill.id + '"><span id="skill-'+skill.id+'" class="label label-' + randomizeLabel() +
        ' s-1x" >' + skill.name + '</span>' +
        '<span class="pull-right"><span class="btn-group btn-group-xs">' +
        '<button class="btn btn-warning" data-toggle="modal" data-target="#skillModal"' +
        'onclick="openSkillModal(' + skill.id + ')">' +
        '<span class="fa fa-pencil color-black"></span>' +
        '</button> ' +
        '<button class="btn btn-danger" onclick=deleteSkill(' + skill.id + ')><span class="fa fa-remove color-white"></span></button>' +
        '<span/><span/></li> ');
}

function deleteSkill(id) {
    new Http().Delete('/skills/' + id).success(function () {
        $('#' + id).remove();
    });
}