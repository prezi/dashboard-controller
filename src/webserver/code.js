function checkURLValidity(url){
	/*var req= new AJ(); // XMLHttpRequest object
	try {
		req.open("HEAD", url, false);
		req.send(null);		
		return req.status== 200 ? true : false;
	}
	catch (er) {
		return false;
	}*/
	/*$.ajax( url )
  .done(function() {
    alert( "success" );
    return true;
  })
  .fail(function() {
    alert( "error" );
    return false;
  });*/
return true;
var request = new XMLHttpRequest();  
request.open('GET', url, true);
request.onreadystatechange = function(){
    if (request.readyState === 4){
        if (request.status === 404) {  
            alert("Oh no, it does not exist!");
            return false;
        }  
        else{
        	return true;
        }
    }
    else{
    	return true;
    }
};

}
function websitePreview(url){

	if (checkURLValidity(url)){
		document.getElementById("webpagepreview").src=url;
		$("#webpagepreview").toggle();
	}
	else{
		alert("Webpage does not exist!");
	}
	//alert(url);
	//        <label for="URL" onchange="websitePreview()">URL:</label>
	//<script src="code.js"></script>
}