## Filtering JSON from S3

Utility in Go, packaged as an executable Docker image, that filters the data from an object in S3 that is GZIP compressed, and prints the result to `stdout`.

Some examples of data:

```json
{"id":5559641812032348625,"time":"1973-03-26T06:18:27.894733102-08:00","words":["scarfing","chip","nuke"]}
{"id":2700778169311697635,"time":"2000-04-12T06:50:23.866113335-06:00","words":["bouncier"]}
{"id":3622438848630508647,"time":"1994-04-15T11:12:04.921735153-06:00","words":["payoffs","winters"]}
{"id":6022900352107475984,"time":"2008-11-24T05:28:08.519322526-08:00","words":["astringency","entertain"]}
```

### Flags

| Name         | Required | Description                                                                                   |
| ------------ | -------- | --------------------------------------------------------------------------------------------- |
| `-input`     | Yes      | An S3 URI (`s3://{bucket}/{key}`) that refers to the source object to be filtered.            |
| `-with-id`   | No       | An integer that contains the `id` of a JSON object to be selected.                            |
| `-from-time` | No       | An RFC3339 timestamp that represents the earliest `time` of a JSON object to be selected.     |
| `-to-time`   | No       | An RFC3339 timestamp that represents the latest `time` of JSON object to be selected.         |
| `-with-word` | No       | A string containing a word that must be contained in `words` of a JSON object to be selected. |

### Usage

```bash
docker build -t danielboakye/s3filter .

docker run --rm -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY danielboakye/s3filter -input s3://sample-data/data.ndjson.gz -from-time=2000-01-01T00:00:00Z -to-time=2001-01-01T00:00:00Z -with-word=payoffs
```
