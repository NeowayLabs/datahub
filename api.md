# Data Hub API

## Companies - visão da empresa

### Pegar lista de jobs

```
GET /companies/jobs

response:
200 OK
[
	"new": [
		{
			"id": 20,
			"company": "Neoway Business Solution",
			"title": "blah blah blah blah",
			"description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
			"accuracyRequired": 90.5,
			"deadline": "2017-03-20",
			"proposedPrice": 15000.00
		},
		{
			"id": 21,
			"company": "Neoway Business Solution",
			"title": "blah blah blah blah",
			"description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
			"accuracyRequired": 90.5,
			"deadline": "2017-03-20",
			"proposedPrice": 15000.00
		},
		{
			"id": 22,
			"company": "Facebook",
			"title": "blah blah blah blah",
			"description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
			"accuracyRequired": 90.5,
			"deadline": "2017-03-20",
			"proposedPrice": 15000.00
		}
	],
	"myJobs": [
		{
			"id": 12,
			"company": "Netflix",
			"title": "blah blah blah blah",
			"description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
			"accuracyRequired": 90.5,
			"accuracy": 70.5,
			"deadline": "2017-03-20",
			"proposedPrice": 15000.00
		}
	],
	"pending": [
		{
			"id": 12,
			"company": "Twitter",
			"title": "blah blah blah blah",
			"description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
			"accuracyRequired": 90.5,
			"deadline": "2017-03-20",
			"proposedPrice": 15000.00
		}
	]
]
```

### Criar novo job
POST /companies/jobs
```
{
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposedPrice": 15000.00
}

response:
200 Ok
{
    "id": <id>
}
```

Obs: job será criado com o status de "pending"

### Upload de arquivos do job
POST /companies/job/_id_/upload
```
{
    "datasets": {
        "training": "training.csv",
        "test": "test.csv",
        "result": "result.csv"
    }
}

response:
200 OK
```
### Pegar um job
GET /jobs/_id_
```
{
    "id": 6,
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "lastUpdate": "2017-03-03",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposedPrice": 15000.00,
	"status": "pending",
    "datasets": {
        "training": "/files/jobs/6/training.csv",
        "test": "/files/jobs/6/test.csv",
        "result": "/files/jobs/6/result.csv"
    },
    "candidates": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "proposedPrice": 15450.00
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "rating": 5,
            "proposedPrice": 18500.00
        }
    ],
    "scientists": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "price": 15450.00
			"solution": {
            	"accuracy": 80.5,
    			"result": "/files/job/6/scientists/5/result.csv",
    			"code": "/files/job/6/scientists/5/code.r",
				"description": "blah blah blah..."
			}
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "price": 18500.00
            "rating": 5,
			"solution": {
            	"accuracy": 70.4,
    			"result": "/files/job/6/scientists/6/result.csv",
    			"code": "/files/job/6/scientists/6/code.r",
				"description": "blah blah blah..."
			}
        }
    ]
}
```

### Alterar um job
PUT /jobs/_id_
```
{
    "id": 6,
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "lastUpdate": "2017-03-03",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposedPrice": 15000.00,
	"status": "doing",
    "datasets": {
        "training": "/files/jobs/6/training.csv",
        "test": "/files/jobs/6/test.csv",
        "result": "/files/jobs/6/test.csv"
    },
    "candidates": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "proposedPrice": 15450.00
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "rating": 5,
            "proposedPrice": 18500.00
        }
    ],
    "scientists": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "price": 15450.00
			"solution": {
            	"accuracy": 80.5,
    			"result": "/files/job/6/scientists/5/result.csv",
    			"code": "/files/job/6/scientists/5/code.r",
				"description": "blah blah blah..."
			}
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "price": 18500.00
            "rating": 5,
			"solution": {
            	"accuracy": 70.4,
    			"result": "/files/job/6/scientists/6/result.csv",
    			"code": "/files/job/6/scientists/6/code.r",
				"description": "blah blah blah..."
			}
        }
    ]
}
```

## Scientist

### pegar lista de scientists
GET /scientists

response:
200 OK
```
[
    {
        "id": 1,
        "image": "/images/scientists/1.jpg",
        "name": "Paulo Pizarro",
        "stars": 2,
        "likes": 300
    },
    {
        "id": 2,
        "image": "/images/scientists/2.jpg",
        "name": "Tiago Katcipis",
        "stars": 5,
        "likes": 500
    },
    {
        "id": 3,
        "image": "/images/scientists/3.jpg",
        "name": "Juliano Galgaro",
        "stars": 4,
        "likes": 400
    },
    {
        "id": 4,
        "image": "/images/scientists/4.jpg",
        "name": "Amanda Rosa",
        "stars": 3,
        "likes": 100
    },
    {
        "id": 5,
        "image": "/images/scientists/5.jpg",
        "name": "Jefferson Amorin",
        "stars": 5,
        "likes": 700
    },
    {
        "id": 6,
        "image": "/images/scientists/6.jpg",
        "name": "Caio Silva",
        "stars": 4,
        "likes": 100
    }
]
```

