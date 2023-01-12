var axios = require('axios');
var data = JSON.stringify({
	"email": "kanivec3@gmail.com"
});

var config = {
	method: 'get',
	url: 'http://localhost:8000/auth/resetpassword?email=kanivec3@gmail.com',
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
