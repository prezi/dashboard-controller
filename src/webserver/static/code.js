$(document).ready(function() {
    $("#mainform").submit(function(e) {
        var postData = $(this).serializeArray();
        var formURL = $(this).attr("action");
        $.ajax({
            url: formURL,
            type: "POST",
            data: postData,
            success: function(data, textStatus, jqXHR) {
                $(".info").html(data.HTML);
                $(".info").show("slow");
                setTimeout(function() {
                    $(".info").hide("slow");
                }, 5000);
            },
            error: function(jqXHR, textStatus, errorThrown) {}
        });
        e.preventDefault();
    });
});