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
        '24:00'
    ];

    // Set the granularity of the time selectors (what nearest time they snap to)
    var snapTo = 10; // 5 minutes
    var pixelsPerHour = 72; // How wide an hour should be, in pixels
    var headName = 'Date'; // Text displayed on top of the list of names
    var defaultStartTime = '00:00';
    var defaultEndTime = '00:00';
    var onClickAppointment = function(payload){
        // Do something with the payload
        outerIndex:
        for(var index = 0; index < list.length; index++){
            for(var i = 0; i < list[index].appointments.length; i++){
                var e = list[index].appointments[i]
                if(e.payload == payload){
                    var url = '/lifelog/edit/' + payload;
                    window.location.href = url;
                    break outerIndex;
                }else{
                    continue;
                }
            };
        };
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

$(document).ready(function () {
    $('#js_monthly_selector').change(function (e) {
        window.location.href = '/lifelog/' + $('#js_monthly_selector').val();
    });
});

function addRemarksHead(width){
    var remarks_head = document.createElement('div');
    remarks_head.textContent = "備考";
    remarks_head.style.width = width + "px";
    remarks_head.className = "sjs-grid-col-head";
    $('.sjs-grid-row-head').append(remarks_head);
}
function addRemarksCol(width){
    var remarks_col = document.createElement('div');
    remarks_col.style.width = width + "px";
    remarks_col.className = "sjs-grid-col";
    $('.sjs-grid-row').append(remarks_col);
}

function addRemarksOverlay(width){
    var data = JSON.parse($("#js_schedulerjs_list").val());
//    var data = [{"name":"2022/06/01","remarks":{"title":"test","payload":"1"}}]
    var sjs_name= $('.sjs-name');
    var sjs_overlay = $('.sjs-grid-overlay-row');
    data.forEach(function(d){
        for( i=0; i<sjs_name.length; i++){
            if(d.remarks.date == sjs_name[i].textContent){
                var new_element = document.createElement('div');
                var data_element = document.createElement('data');
                var childlen = sjs_overlay[i].getElementsByClassName("sjs-grid-overlay-col");
                var margin = 72 * 24;
                for( j=0; j<childlen.length; j++){
                    margin = margin - (parseInt(childlen[j].style.width) + parseInt(childlen[j].style.marginLeft));
                }
                new_element.textContent = d.remarks.title;
                new_element.style.width = width + "px"
                new_element.style.marginLeft = margin + "px";
                new_element.className = "sjs-grid-overlay-col sjs-grid-overlay-col-clickable remarks";
                data_element.value = d.remarks.payload;
                new_element.appendChild(data_element);
                sjs_overlay[i].append(new_element);
            }
        }
    })
}

$(document).ready(function(){
    const width = 90;
    addRemarksHead(width);
    addRemarksCol(width);
    addRemarksOverlay(width);
});

$(document).ready(function(){
    $('.remarks').click(function(){
        window.location.href = '/remarks/edit/' + this.getElementsByTagName('data')[0].value;
    });
});