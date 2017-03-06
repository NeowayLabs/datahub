# Datahub

Just a POC :-)

# Running

Run:

```
make
```

After build you can run:

```
./datahub
```

The server is running locally, to upload a file to it:

```
./tools/upload.sh 127.0.0.1:8080 /tmp/example
```

To run R code on the server:

```
./tools/execr.sh 127.0.0.1:8080 /tmp/rcode
```

R code must work with the files previously uploaded
to the server. Because of that the script will always
upload examples datasets before uploading your code.


# R Examples

If you want to run an R example directly on your machine just run:

```
cd examples
sudo R -f ./installdeps.r
R -f ./code.r
```

Example:
```
curl -XPOST http://192.168.10.169:8080/api/companies/jobs -d '{"title": "Faturamento Presumido de Empresas do Brasil", "description":"Calcular para a Neoway o faturamento presumido de todas as empresas do Brasil baseado no CNAE", "deadline": "2017-04-01", "proposed": 5000, "accuracyRequired": 90.5}'

curl http://192.168.10.169:8080/api/companies/jobs|json_pp

curl -XPOST http://192.168.10.169:8080/api/scientists/4/jobs/5/apply -d '{"counterproposal": 1000}'
curl -XPOST http://192.168.10.169:8080/api/companies/jobs/5/start -d '{"scientists":[{"id":4}]}'
curl http://192.168.10.169:8080/api/scientists/4/jobs/5/workspace|json_pp
```
