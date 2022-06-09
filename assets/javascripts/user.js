$(document).ready(function () {
    $('.btn-logout').click(function (e) {
        $.removeCookie('auth-session');
        window.location.href = '/';
    });
});