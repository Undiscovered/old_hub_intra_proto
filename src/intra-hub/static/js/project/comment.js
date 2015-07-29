function openCommentModal(id) {
    var converter = new showdown.Converter();
    $('#comment-edit-name').html($('#comment-' + id).text());
    $('#comment-edit-id').val(id)
}

function startCommentProject() {
    var converter = new showdown.Converter();
    $('.comment').each(function () {
        var element = $(this);
        var text = element.text();
        var html = converter.makeHtml(text);
        element.html(html);
    });
}
