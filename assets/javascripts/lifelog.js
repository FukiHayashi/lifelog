$(document).ready(function(){
    var list = JSON.parse($("#js_schedulerjs_list").val());
    // Customize what time steps are shown in the scheduler
    var steps = [
        '00:00',
        '01:00',
        '02:00',
        '03:00',
        '04:00',
        '05:00',
        '06:00',
        '07:00',
        '08:00',
        '09:00',
        '10:00',
        '11:00',
        '12:00',
        '13:00',
        '14:00',
        '15:00',
        '16:00',
        '17:00',
        '18:00',
        '19:00',
        '20:00',
        '21:00',
        '22:00',
        '23:00',
        '24:00',
        '25:00'
    ];

    // Set the granularity of the time selectors (what nearest time they snap to)
    var snapTo = 10; // 5 minutes
    var pixelsPerHour = 75; // How wide an hour should be, in pixels
    var headName = 'Date'; // Text displayed on top of the list of names
    var defaultStartTime = '00:00';
    var defaultEndTime = '00:00';
    var onClickAppointment = function(payload){
        // Do something with the payload
    };

    var $scheduler = $("#scheduler").schedulerjs({
        'list': list,
        'steps': steps,
        'snapTo': snapTo,
        'pixelsPerHour': pixelsPerHour,
        'start': defaultStartTime,
        'end': defaultEndTime,
        'headName': headName,
        'onClickAppointment': onClickAppointment
    });
});