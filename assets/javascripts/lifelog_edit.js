$(document).ready(function () {
    $('#js_delete').click(function (e) {
        var url = $("#js_delete").val()
        $.ajax({
            url: url,
            type: "DELETE",
        })
        window.location.href = '/lifelog';
    });
});