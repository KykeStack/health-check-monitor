![Gopher](https://i.imgur.com/glk34vC.png)

# Status Check API

## What is this?

This program will constantly check for one or more http endpoints and notify you
 whenever something goes down and also when it goes up again!

It will additionally expose an api with the current statuses of every endpoint 
configured.

## Configuration

This health check needs a `config.json` file with this minimal structure

```json
{
  "urlMonitors": [
    {
      "name": "Google",
      "url": "https://www.google.com"
    }
  ]
}
```  

You can also specify the config file path with the `-config` flag. For example:

```bash
health-check-monitor -config my/custom/dir/my_config.json
```

Check more configuration options in the conf.json [distribution file]

## Run locally

```bash
go get gitlab.com/skyvet/health-check-monitor
cd $GOPATH/src/gitlab.com/skyvet/health-check-monitor
go install
curl https://github.com/KykeStack/health-check-monitor.git/raw/master/config.json.dist -o config.json
health-check-monitor
```

## Run with docker

### Download a basic config file

```bash
curl https://github.com/KykeStack/health-check-monitor.git/raw/master/config.json.dist -o config.json
```

### Look up latest version and run it

Or just go [check manually the list] and then run docker

```bash
docker run -d \
-v "$(pwd)/config.json:/config.json" \
-p "8001:8001" \
git git+https://github.com/KykeStack/health-check-monitor.git:<insert here latest tag> -config /config.json
```

## Ping it!

```bash
curl http://localhost:8001/status
```

## Licensing

Check [LICENSE.md]

[distribution file]: https://github.com/KykeStack/health-check-monitor/raw/master/config.json.dist
[jq]: https://stedolan.github.io/jq/
[check manually the list]: https://github.com/KykeStack/health-check-monitor/tags
[gitlab tags api]: https://docs.gitlab.com/ce/api/tags.html
[LICENSE.md]: (https://github.com/KykeStack/health-check-monitor/blob/master/LICENSE)
