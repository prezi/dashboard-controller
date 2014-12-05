"use strict";

$(document).ready(function() {
    $("#mainform").submit(function(e) {
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
});