var axios = require('axios');
var data = JSON.stringify({
	"login": "kanivec3@gmail.com",
	"password": "admin"
});

var config = {
	method: 'post',
	url: 'http://localhost:8000/auth/sing-in',
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
