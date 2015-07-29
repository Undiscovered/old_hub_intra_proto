function resetPassword() {
    var p1 = $('#p').val();
    var p2 = $('#p2').val();
    if (p1 !== p2) {
        alert('Passwords mismatch');
    } else if (p1.length < 5) {
        alert('Password too short.');
    } else {
        $('#myform').submit();
    }
}
