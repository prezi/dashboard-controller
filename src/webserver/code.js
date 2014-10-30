/*
function submitForm() {
	alert("sfgbf");
    var http = new XMLHttpRequest();
    http.open("POST", "/form-submit", true);
    http.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    var params = "search="; // probably use document.getElementById(...).value
    http.send(params);
    http.onload = function() {
        alert(http.responseText);
    }
}
*/

$( document ).ready(function() {
    //alert("Loaded");
    /*	
    $('#mainform').submit(function () {
	 alert("WATWAT");
	 return false;
	});
	*/
	//$("#manual-example a[rel=tipsy]").tipsy("show");
	$("#manual-example a[rel=tipsy]").show();
	$("#mainform").submit(function(e)
	{
	    var postData = $(this).serializeArray();
	    var formURL = $(this).attr("action");
	    $.ajax(
	    {
	        url : formURL,
	        type: "POST",
	        data : postData,
	        success:function(data, textStatus, jqXHR) 
	        {
	            //alert("Success");
	            $( "#rb_id" ).find( "#value" ).html(data.ID);
	            $( "#url" ).find( "#value" ).html(data.URL);
	            $( "#statuscode" ).find( "#value" ).html(data.Code);
	            $( ".info" ).show("slow");
	            setTimeout(function() {
						$( ".info" ).hide("slow");
					}, 5000);
	        },
	        error: function(jqXHR, textStatus, errorThrown) 
	        {
	            //alert("Fail"); 
	        }
	    });
	    e.preventDefault(); //STOP default action
	    //e.unbind(); //unbind. to stop multiple form submit.
	});

});


