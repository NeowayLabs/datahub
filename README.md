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

# R Examples

If you want to run an R example directly on your machine just run:

```
cd examples
sudo R -f ./installdeps.r
R -f ./code.r
```
