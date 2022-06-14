$(function() {
    $('.js_timepic').datetimepicker({
        step:10,
        dateformat: 'yyyy-mm-dd',
        timeformat: 'hh:mm',
    });
});

// タイトルから分類を自動入力する
$(function(){
    $('#js_appointment_title').on('input', function(){
        var input_text = $(this).val();
        var js_class_select = "other";
        switch(input_text){
            case "睡眠":
            case "昼寝":
            case "仮眠":
                js_class_select = "sleep";
                break;
            case "朝食":
            case "昼食":
            case "夕食":
            case "軽食":
            case "食事":
                js_class_select = "meal";
                break;
            case "風呂":
                js_class_select = "bath";
                break;
            default:
                js_class_select = "action";
                break;
        }
        $('#js_class_select').val(js_class_select);
    });
})