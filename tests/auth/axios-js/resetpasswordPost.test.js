var axios = require('axios');
var data = JSON.stringify({
	"confirm": "6b6b6b6b6b6b6b6b6b6b6b6b6b39e4592986fb69508e88d3a0517efbdb7446034d",
	"password": "admin1",
	"password2": "admin1"
});

var config = {
	method: 'post',
	url: 'http://localhost:8000/auth/resetpassword',
	headers: {
		'Content-Type': 'application/json'
	},
	data : data
};

axios(config)
	.then(function (response) {
		console.log(JSON.stringify(response.data));
	})
	.catch(function (error) {
		console.log(error);
	});
