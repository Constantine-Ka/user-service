var axios = require('axios');

var config = {
	method: 'get',
	url: 'http://localhost:8000/auth/confirm?code=6b6b6b6b6b6b6b6b6b6b6b6b6b39e4592986fb69508e88d3a0517efbdb7446034d',
	headers: { }
};

axios(config)
	.then(function (response) {
		console.log(JSON.stringify(response.data));
	})
	.catch(function (error) {
		console.log(error);
	});
