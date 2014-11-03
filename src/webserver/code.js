$( document ).ready(function() {
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
	        	//TO DO: Change it to universal info box!
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
	            //TO DO: As above, change it to universal info box
	        }
	    });
	    e.preventDefault(); 
	    //e.unbind(); //unbind. to stop multiple form submit.
	});

});


