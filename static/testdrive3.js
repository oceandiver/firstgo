$(document).ready(function() {

var userid;
var token;

$('#form-signup').submit( function(event) {
    $.post('v1/user/signup', $('form#form-signup').serialize(),
	   function(data) {
               alert(data);
	   },
	   'json' // I expect a JSON response
	  );
    event.preventDefault();
});

$('#form-signin').submit( function(event) {
    $.post( 'v1/user/signin', $('form#form-signin').serialize(), 
	    function(data) {
		//	var response = jQuery.parseJSON(data);
		alert(data);
	
	    },
	    'json' // I expect a JSON response
	  );
    event.preventDefault();
});

$("#form-new-event").submit(function(event){


    function formatTimestring(s) {
	var b = s.split(/[\/:]/);
	return b[2] + b[1] + b[0] + 'T' + b[3] + b[4] + '00' + 'Z'
    }

    alert(
	formatTimestring("#form-new-event input[name=time]") //20100908T120000Z
    );

    $.post("/v1/addevent", $('form#form-signin').serialize(), 
	   function(data, status) {
               alert("Data: " + data + "\nStatus: " + status);
	   },
	   'json' //expect a JSON response
	  );
    event.preventDefault();
});

$("#form-list-users").submit(function(event){
    var a="";
    $.getJSON("/v1/user/all", $('form#form-list-users').serialize(), function(result){
        $.each(result, function(i, field){
            a = a.concat(field + " ");
        });
	alert(a);
    });
    event.preventDefault();
});

$("#formist-events").submit(function(event){
    $.getJSON("/v1/event/all", $('form#form-list-events').serialize(), function(result){
        $.each(result, function(i, field){
            // $("div").append(field + " ");
        });
    });
    event.preventDefault();
});


});
