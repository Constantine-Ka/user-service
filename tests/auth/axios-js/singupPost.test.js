var axios = require('axios');
var data = JSON.stringify({
	"login": "admin",
	"firstName": "Константин",
	"email": "kanivec3@gmail.com",
	"password": "admin",
	"secondName": "Vit",
	"lastName": "Ka",
	"imagePath": "https://pkg.go.dev/static/shared/icon/content_copy_gm_grey_24dp.svg",
	"gender": 0,
	"birthday": 788281200,
	"description": "Обо мне очень хорошо отзываются"
});

var config = {
	method: 'post',
	url: 'http://localhost:8000/auth/sing-up',
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
