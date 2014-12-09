"use strict";

function getMaxSlaveLabelWidth () {
    var maxWidth=0;
    $(".slave-selector a").each(function(index){
        if ($( this ).width() > maxWidth) {
        maxWidth = $( this ).width();
        }
    });
    return maxWidth;
}

function setSlaveLabelWidth() {
    var maxWidth = getMaxSlaveLabelWidth();
    $(".slave-selector a").width(maxWidth + 1);
}

$(document).ready(function() {
    setSlaveLabelWidth();
    $("#mainform").submit(function(e) {
        if ($('.slave-selector a.strongSelect').size() === 0) {
            $(".info").show("slow");
            $(".info").html('<div>Slave not selected</br> \
            Please select at least one slave before hit the submit button!</div>');
            setTimeout(function() {
                $(".info").hide("slow");
            }, 5000);
            e.preventDefault();
            return;
        }
        var selectedSlave = $('.slave-selector a').filter('.strongSelect').html();
        var usrToDisplay = $('.form-control').val();
        var postData = {
            'url':usrToDisplay,
            'slave-id': selectedSlave
        };
        var formURL = $(this).attr("action");
        $.ajax({
            url: formURL,
            type: "POST",
            data: postData,
            timeout: 8000,
            cache: false,
            success: function(data, textStatus, jqXHR) {
                var newInfoBoxContent = data.StatusMessage;
                var isPersistent = data.IsPersistent == "true";
                $(".info").html(data.StatusMessage);
                if (!isPersistent) {
                    $(".info").show("slow");
                    setTimeout(function() {
                        $(".info").hide("slow");
                    }, 5000);
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                $(".info").show("slow");
                $(".info").html('<div>Error communicating with web server.</br> \
                Please check the web service, and refresh the page!</div>');
            }
        });
        e.preventDefault();
    });
    $('.slave-selector a').on('click', function (e) {
        if ($(this).hasClass('strongSelect')) {
            $(this).removeClass('strongSelect');
        } else {
            $('.slave-selector a').filter('.strongSelect').removeClass('strongSelect');
            $(this).addClass('strongSelect');
        }
    });
    $('#submit-button').tooltip({
        'show': false,
        'placement': 'right'
    });
    $('#submit-button').mouseover(function(){
        var message = "";
        if (!$('.form-control').val()){
            message += 'URL input field is empty. Please provide a URL to display!';
        }
        if ($('.slave-selector a.strongSelect').size() == 0){
            message += "No slaves are selected. Please select a slave on which you can display a URL";
        }
        if (message !== "") {
            $('#submit-button').attr('data-original-title', message)
                               .tooltip('fixTitle')
                               .tooltip('show');
        } else {
            $('#submit-button').attr('data-original-title', '')
                               .tooltip('fixTitle')
                               .tooltip('hide');
        }
    });
//    $('#submit-button').tooltip('show');
});