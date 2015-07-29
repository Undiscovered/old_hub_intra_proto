function startSingleProject() {
    var converter = new showdown.Converter();
    var descSelector = $('#description');
    var text = descSelector.text();
    var html = converter.makeHtml(text);
    descSelector.html(html);
}
