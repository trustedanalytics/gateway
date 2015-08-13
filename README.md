# WebSocket-based Gateway

The `gateway` is an ingestion bridge connecting edge devices with Cloud-based data processing. It exposes WebSocket interface which publishes all inbound messages to the configurable message back-end (Apache Kafka queue). It is compatible with any `Websocket-based` publisher.

Implemented:

* JSON message format, no mapping required
* Token-based client authentication (OAuth 2.0)
* Supports multiple backends (default to Apache Kafka)
* Configurable server (port, path, service etc.)

TODO:

* Client-level authorization
* Dynamic backend configuration 

## Installation

You can install `gateway` by either cloning this repo (below) or by [downloading the latest binary](https://github.com/mchmarny/gateway/releases) distribution.

    git clone git@github.com:mchmarny/gateway.git
    cd ./gateway

To start the server invoke the `gateway` executable

    ./gateway

## Configuration

The `gateway` comes pre-configured with a default (`config.json`). 

```
{
  "id": "g1",
  "trace": true,
  "server": {
    "root": "/ws",
    "host": "127.0.0.1",
    "port": 8080,
    "token": ""
  },
  "publisher": {
    "uri": [
      "127.0.0.1:9092"
    ],
    "topic": "messages",
    "acks": "false"
  }
}

```

* `topic` will be automatically created if one does not exists
* `acks` if set to true will wait for acknowledgment from all brokers (slower)
* `retries` number of times to retry a metadata request when a partition is in the middle of leader election (10+)


> Note, when runtime is [Cloud Foundry](https://github.com/cloudfoundry) the following configuration attributes are going to be overwritten with [CF environment variables](http://docs.cloudfoundry.org/devguide/deploy-apps/environment-variable.html):

    id = VCAP_APPLICATION.instance_id + VCAP_APPLICATION.instance_index
    server.port = VCAP_APPLICATION.port
    server.host = VCAP_APPLICATION.host
    server.token = $GATEWAY_TOKEN
    publisher.uri = VCAP_SERVICES[x].credentials.uri
    publisher.topic = $GATEWAY_TOPIC
    
## Cloud Foundry Push

Make sure the code build locally 

```
go build
go test
```

Pre-package the dependancies 

```
godep save
```

> If you don't have godep already installed download it using `go get github.com/tools/godep`

Export Application Variables

```
export APP_NAME="my-app-name"
export APP_TOKEN=$(echo -n 'your-secret-here' | openssl base64)
```

Push the code to CF

> These are the minimum variables you need to change. See `manifest.yml` for complete list. Also, make sure you provide the `--no-start` argument to prevent the gateway from publishing to the default topic

```
cf push $APP_NAME -n $APP_NAME --no-start
cf set-env $APP_NAME GATEWAY_TOPIC $APP_NAME
cf set-env $APP_NAME GATEWAY_TOKEN $APP_TOKEN
cf start $APP_NAME
```

## Testing

In addition to the integrated `Go` test, the `gateway` application also includes a simple Node.js test client: `etc/client.js`. This client is intended for perform a simple smoke-test of the deployed `gateway` application. 

```
node client -s 'wss' \
            -e '${$APP_NAME}.<platform_domain>' \
            -p 4443 \
            -r '/ws' \
            -t '${APP_TOKEN}' \
            -f 1000
```

The test client will loop through and send to the gateway individual events at the frequency specified `-f` in milliseconds looking like this:

```
{
   "source_id": "ip-172-1-2-3",
   "event_id": "c4e6e427-7e45-41d8-9ba2-2223062398dc",
   "event_ts": 1417763522224,
   "metrics": [
      {
         "key": "cpu_load_5min",
         "value": 0.0029296875
      },
      {
         "key": "cpu_load_10min",
         "value": 0.0146484375
      },
      {
         "key": "cpu_load_15min",
         "value": 0.04541015625
      },
      {
         "key": "free_memory",
         "value": 3402670080
      }
   ],
   "event_number": 30
}
```

If you are not sure of the arguments, execute `node client.js --help` for some help.

### Scaling 

If your throughput on the gateway is not sufficient, you can increase the number of application instances. Following command sets the total number of application instances to `3` 

```
cf scale $APP_NAME -i 3
```

## VERSIONING:
We use [bumpversion](https://github.com/peritus/bumpversion) tool to manage version written in `manifest.yml`. 
Release versions are in standard `Major.Minor.Patch` format. Snapshot version are in `Major.Minor.Patch.build` format.
Each realease version component can be easily upgraded with `bumpvresion` command, but you don't have to do it manually - all things happen at TeamCity. 
If you want to change version manually by yourself (without bumpversion tool) REMEMBER to change `version` in two files: .bumpversion.cfg and manifest.yml!

**Notable things:**
* Configuration is in `.bumpversion.cfg` file.
* Version info is updated in two files: `.bumpversion.cfg`, `manifest.yml`.
* Release version format is {major}.{minor}.{patch}.
* Snapshot version format is {major}.{minor}.{patch}.{build}.
* Git tags format: v{major}.{minor}.{patch} - only created for release.
* Commit format: [{day}-{month}-{year}] TeamCity build: {build version} release.

**Possible actions:**
* Release: `bumpversion dev=1` then `bumpversion patch` and `git push --tags`
* Snapshot release: it is not possible to make snapshot release by bumpversion. TeamCity do it.
* Update minor/major version: `bumpversion dev=1` then `bumpversion minor` or `bumpversion major` followed by `git push --tags`

## Backends

Currently `gateway` supports:

* [Apache Kafka](http://kafka.apache.org/)

#### Message

The `gateway` decorates the inbound messages with following attributes:

```
{
    id: [v4 uuid],
    on: [UTC timestamp]
    body: [inbound message content in UTF-8 encoded string]
}
```

## Deploying to Marketpalce with app-launching-service-broker

First define your broker name:
```
export BROKER_NAME=<your broker name>
```

To deploy gateway to the marketplace you have to use app-launching-service-broker. Clone it and go to the new folder:
```
git clone git@github.com:trustedanalytics/app-launching-service-broker.git $BROKER_NAME
cd $BROKER_NAME
```
and clone gateway to apps/gateway with command
```
git clone git@github.com:trustedanalytics/gateway apps/gateway
```	
Download cf client from [link](https://cli.run.pivotal.io/stable?release=linux64-binary&source=github) or alternative version [link](https://github.com/cloudfoundry/cli#downloads) and click on Stable Binaries > Linux 64 bit. Then unpack and place it in broker's bin directory.

Go to manifest.yml in broker directory and add to the env values:
```
CF_API: <your_cf_api>
CF_SRC: ./apps/gateway
CF_USER: <admin>
CF_PASS: <admin_password>
CF_DEP: kafka|shared
```
also edit name of broker (same as *$BROKER_NAME*) and delete buildpack value. Example manifest.yml after changes:
```
---
applications:
- name: gateway_broker
  memory: 256M
  instances: 1
  path: .
  env:
    CF_API: http://api.platform.com
    CF_CATALOG_PATH: ./catalog.json
    CF_SRC: ./apps/gateway
    CF_USER: admin
    CF_PASS: password
    CF_DEP: kafka|shared
```	
Open catalog.json and edit values:
* services -> id (must be unique in environment)
* services -> name (must be unique eg. gateway*)
* services -> description
* services -> tags
* services -> plans -> id (must be unique in environment)

Example catalog.json:
```
{
  "services": [{
    "id": "548f9a19-a193-4a86-b449-b448350db54d",
    "name": "gateway_",
    "description": "Simple websocket bridge with kafka back-end",
    "bindable": true,
    "tags": ["gateway"],
    "plans": [{
      "id": "3273fc74-8b8d-422b-8217-4a8eb6b6ce5a",
      "name": "simple",
      "description": "Simple",
      "free": true
    }]
  }]
}
```
Open manifest.yml in apps/gateway/ directory and change name to the same as in catalog.json also delete buildpack value as in previous manifest.
Example of apps/gateway/manifest.yml:
```
---
applications:
- name: gateway_
  memory: 512MB
  instances: 1
  env:
    GATEWAY_TRACE: "false"
    GATEWAY_ACKS: "false"
    VERSION: "0.9.5.0"
```

Next go to apps/gateway and copy Godeps folder to the broker directory (there is already one) - when asked about what to do with duplicates click skip all.
Go to broker directory, push app ,create service broker and enable service access, can be done with commends:
```
cf push
export SERVICE_URL=$(cf app $BROKER_NAME | grep urls: | awk '{print $2}')
cf create-service-broker $BROKER_NAME admin admin https://$SERVICE_URL
cf enable-service-access <app name same as in catalog.json>
```
**You have done it, go to the marketplace and enjoy your new Gateway!**
## License

This project is under the MIT License. See the [LICENSE](https://github.com/mchmarny/gateway/blob/master/LICENSE) file for the full license text.
