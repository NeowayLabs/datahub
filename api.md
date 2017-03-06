# Data Hub API

## Companies - visão da empresa

### Lista de jobs

GET /api/companies/jobs

response: 200 OK
```json
{
    "pending": [
        {
            "jobId": 8,
            "title": "",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "scientists": [
                {
                    "id": 5,
                    "name": "Juliano Galga",
                    "rating": 5
                },
                {
                    "id": 6,
                    "name": "Caio Silva",
                    "rating": 5
                }
            ],
            "status": "pending"
        }
    ],
    "doing": [
        {
            "jobId": 6,
            "title": "",
            "description": "blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 4,
                "name": "Amanda Rosa",
                "rating": 3,
                "accuracy": 50.5
            },
            "status": "doing"
        },
        {
            "jobId": 7,
            "title": "",
            "description": "blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 4,
                "name": "Amanda Rosa",
                "rating": 3,
                "accuracy": 50.5
            },
            "status": "doing"
        }
    ],
    "done": [
        {
            "jobId": 1,
            "title": "Who Killed Who? Ex-Husband VS Ex-Wife",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 2,
                "name": "Tiago Katcipis",
                "rating": 5,
                "accuracy": 80.5
            },
            "status": "done"
        },
        {
            "jobId": 2,
            "title": "Who Killed Who? Ex-Husband VS Ex-Wife",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 2,
                "name": "Tiago Katcipis",
                "rating": 5
            },
            "accuracy": 80.5,
            "status": "done"
        },
        {
            "jobId": 3,
            "title": "Who Killed Who? Ex-Husband VS Ex-Wife",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 2,
                "name": "Tiago Katcipis",
                "rating": 5
            },
            "accuracy": 80.5,
            "status": "done"
        },
        {
            "jobId": 4,
            "title": "Who Killed Who? Ex-Husband VS Ex-Wife",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 2,
                "name": "Tiago Katcipis",
                "rating": 5,
                "accuracy": 80.5
            },
            "status": "done"
        },
        {
            "jobId": 5,
            "title": "Who Killed Who? Ex-Husband VS Ex-Wife",
            "description": "Homicide Report blah blah blah ....",
            "lastUpdate": "2017-03-03",
            "price": 15450.00,
            "scientist": {
                "id": 2,
                "name": "Tiago Katcipis",
                "rating": 5,
                "accuracy": 80.5
            },
            "status": "done"
        }
    ]
}
```

### Criar um novo job

POST /api/companies/jobs
```json
{
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposed": 15000.00
}
```

response: 200 OK
```json
{
    "id": 1
}
```

Obs: job será criado com o status de "pending"

### Upload de arquivos do job

POST /api/companies/jobs/_ID_/upload

response: 200 OK

### Obter dados de um job

GET /api/companies/jobs/_ID_
```json
{
    "id": 6,
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "lastUpdate": "2017-03-03",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposed": 15000.00,
    "status": "pending",
    "datasets": {
        "training": "/datasets/jobs/6/trainingset.csv",
        "challenge": "/datasets/jobs/6/testset.challenge.csv",
        "result": "/datasets/jobs/6/testset.result.csv"
    },
    "candidates": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "counterproposal": 15450.00
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "rating": 5,
            "counterproposal": 18500.00
        }
    ],
    "scientists": [
        {
            "id": 5,
            "name": "Juliano Galga",
            "rating": 5,
            "price": 15450.00,
            "solution": {
                "accuracy": 80.5,
                "result": "/datasets/job/6/scientists/5/result.csv",
                "code": "/datasets/job/6/scientists/5/code.r",
                "description": "blah blah blah..."
            }
        },
        {
            "id": 6,
            "name": "Caio Silva",
            "price": 18500.00,
            "rating": 5,
            "solution": {
                "accuracy": 70.4,
                "result": "/datasets/job/6/scientists/6/result.csv",
                "code": "/datasets/job/6/scientists/6/code.r",
                "description": "blah blah blah..."
            }
        }
    ]
}
```

### Iniciar um job

POST /api/companies/jobs/_ID_/start
```json
{
    "scientists": [
        {
            "id": 5,
        },
        {
            "id": 6,
        }
    ]
}
```

response: 200 OK

## Scientist

### Lista de scientists

GET /api/scientists

response: 200 OK
```json
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

### Lista de jobs do scientists

GET /api/scientists/_ID_/jobs

respoonse: 200 OK
```json
{
    "new": [
        {
            "id": 20,
            "company": "Neoway Business Solution",
            "title": "blah blah blah blah",
            "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
            "accuracyRequired": 90.5,
            "deadline": "2017-03-20",
            "proposed": 15000.00
        },
        {
            "id": 21,
            "company": "Neoway Business Solution",
            "title": "blah blah blah blah",
            "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
            "accuracyRequired": 90.5,
            "deadline": "2017-03-20",
            "proposed": 15000.00
        },
        {
            "id": 22,
            "company": "Facebook",
            "title": "blah blah blah blah",
            "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
            "accuracyRequired": 90.5,
            "deadline": "2017-03-20",
            "proposed": 15000.00
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
            "proposed": 15000.00
        }
    ],
    "doing": [
        {
            "id": 12,
            "company": "Netflix",
            "title": "blah blah blah blah",
            "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
            "accuracyRequired": 90.5,
            "accuracy": 70.5,
            "deadline": "2017-03-20",
            "proposed": 15000.00
        }
    ],
    "done": [
        {
            "id": 12,
            "company": "Netflix",
            "title": "blah blah blah blah",
            "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
            "accuracyRequired": 90.5,
            "accuracy": 70.5,
            "deadline": "2017-03-20",
            "proposed": 15000.00
        }
    ]
}
```

### Aceitar um job

POST /api/scientists/_ID_/jobs/_ID_/apply
```json
{
    "counterproposal": 20300.00
}
```

### Vizualizar área de trabalho de um job

GET /api/scientists/_ID_/jobs/_ID_/workspace
```json
{
    "id": 6,
    "title": "blah blah blah blah",
    "description": "fsdfsd ff sdf sdf sdfsdf sdf sdf",
    "lastUpdate": "2017-03-03",
    "accuracyRequired": 90.5,
    "deadline": "2017-03-20",
    "proposed": 15000.00,
    "status": "pending",
    "datasets": {
        "training": "/datasets/jobs/6/trainingset.csv",
        "challenge": "/datasets/jobs/6/testset.challenge.csv",
    },
    "workspace": {
        "accuracy": 80.5,
        "result": "/datasets/job/6/scientists/5/result.csv",
        "code": "/datasets/job/6/scientists/5/code.r",
        "description": "blah blah blah..."
    }
}
```

### Upload do codigo R de um scientists para um job

POST /api/scientists/_ID_/jobs/_ID_/upload


## Uploads

Uploads are done using multipart forms just as the browser does.
The name of the form field is the name of the file you are uploading.
Multiple files can be uploaded simultaneously.

Available files to upload are:

* **trainingset.csv** : Training set (has the responses), sent by the company
* **testset.challenge.csv** : Challenge test set (does not have the responses), sent by the company
* **testset.result.csv** : Result test set, has all responses from the challenge, used to calculate accuracy
* **code.r** : R code, sent by the statistician. Must be sent before executing code.

The URI is:

```
/api/upload
```

## Executing R code

Just send a **POST** to the URI is:

```
/api/execr
```

Before this, all datasets should already been uploaded
by the company and the R code should be uploaded by the
statistician.
