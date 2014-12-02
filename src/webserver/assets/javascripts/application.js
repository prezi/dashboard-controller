"use strict";

$(document).ready(function() {
    $("#mainform").submit(function(e) {
        var postData = $(this).serializeArray();
        var formURL = $(this).attr("action");
        $.ajax({
            url: formURL,
            type: "POST",
            data: postData,
            timeout: 8000,
            cache: false,
            success: function(data, textStatus, jqXHR) {
                var newInfoBoxContent = data.HTML;
                var isPersistent = data.IsPersistent == "true";
                $(".info").html(data.HTML);
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
});