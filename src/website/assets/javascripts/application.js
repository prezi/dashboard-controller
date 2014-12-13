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
            Please select at least one slave before submitting your URL!</div>');
            setTimeout(function() {
                $(".info").hide("slow");
            }, 5000);
            e.preventDefault();
            return;
        }
        var selectedSlaves = [];
        $('.slave-selector a').filter('.strongSelect').each(function() {
            selectedSlaves.push($( this ).html());
        });
        var usrToDisplay = $('.form-control').val();
        var postData = {
            'URLToDisplay':usrToDisplay,
            'SlaveNames': selectedSlaves
        };
        var formURL = $(this).attr("action");
        $.ajax({
            url: formURL,
            type: "POST",
            data: JSON.stringify(postData),
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
        $(this).toggleClass('strongSelect');
    });
    $('#submit-button').tooltip({
        'show': false,
        'placement': 'right'
    });
    $('#submit-button').mouseover(function(){
        var message = "";
        if (!$('.form-control').val()){
            message += 'Please provide a URL to display.';
        } else if ($('.slave-selector a.strongSelect').size() == 0){
            message += "No destination selected. Please select where you would like to load this URL.";
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
});