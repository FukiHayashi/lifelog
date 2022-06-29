$(document).ready(function () {
    $('.btn-logout').click(function (e) {
        var xhr = new XMLHttpRequest();
        xhr.open("DELETE","/logout");
        xhr.send();
        cookieStore.delete('auth-session')
        window.location.href = '/';
    });
});