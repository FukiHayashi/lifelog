$(function() {
    var date = new Date();
    $('#js_timepic_start').datetimepicker({
        step:10,
        dateformat: 'yyyy-mm-dd',
        timeformat: 'hh:mm',
        minDate: date.getFullYear() + '/' + date.getMonth() + '/01',
        maxDate: date.getFullYear() + '/' + (date.getMonth()+3) + '/01'
    });
});

$(function(){
    var date = new Date($('#js_timepic_start').val())
    $('#js_timepic_start').change(function(){
        date = new Date($('#js_timepic_start').val());
        date.setMinutes(date.getMinutes()+30);
        $('#js_timepic_end').val(date.getFullYear() + '/' + (date.getMonth()+1).toString().padStart(2, '0') + '/' + (date.getDate()).toString().padStart(2, '0') + ' ' + date.getHours().toString().padStart(2, '0') + ':' + date.getMinutes().toString().padStart(2, '0'));
        $('#js_timepic_end').datetimepicker({
            step:10,
            dateformat: 'yyyy-mm-dd',
            timeformat: 'hh:mm',
            minDate: $('#js_timepic_start').val(),
            maxDate: date.getFullYear() + '/' + (date.getMonth()+1) + '/' + (date.getDate()+1)
        });
    });
    $('#js_timepic_end').datetimepicker({
        step:10,
        dateformat: 'yyyy-mm-dd',
        timeformat: 'hh:mm',
        minDate: $('#js_timepic_start').val(),
        maxDate: date.getFullYear() + '/' + (date.getMonth()+1) + '/' + (date.getDate()+1)
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