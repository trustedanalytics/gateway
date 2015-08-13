/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
var stdio = require('stdio'),
    os = require('os'),
    uuid = require('node-uuid'),
    WebSocket = require('ws');

// runtime arguments
var env = stdio.getopt({
    'protocal': { key: 's', description: 'WebSocket: [ws | wss]', default: 'ws', args: 1 },
    'endpoint': { key: 'e', description: 'Endpoint', default: '127.0.0.1', args: 1 },
    'port':     { key: 'p', description: 'Port:', default: 4443, args: 1 },
    'path':     { key: 'r', description: 'Route:', default: '/ws', args: 1 },
    'token':    { key: 't', description: 'Token:', default: '', args: 1  },
    'freq':     { key: 'f', description: 'Frequency: [1] (in sec)', default: 1000, args: 1  }
});

console.log("Client started using: ");
console.dir(env);

var getMessage = function(){
  var load = os.loadavg();
  return {
    'source_id': os.hostname(),
    'event_id': uuid.v4(),
    'event_ts': new Date().getTime(),
    'metrics': [
      { 'key': 'cpu_load_5min', 'value': load[0] },
      { 'key': 'cpu_load_10min', 'value': load[1] },
      { 'key': 'cpu_load_15min', 'value': load[2] },
      { 'key': 'free_memory', 'value': os.freemem() }
    ]
  };
}

// WS configuration
var hd = { headers: { 'Authorization': 'Bearer ' + env.token } };
var uri = env.protocal + '://' + env.endpoint + ':' + env.port + env.path;
var ws = new WebSocket(uri, hd);

var i = 0, printOn = 1000
var send = function(){
    var msg = getMessage();
    msg.event_number = i++;
    var str = JSON.stringify(msg);
    setTimeout(function() {
        ws.send(str, function(err) {
            if (err != null) {
                console.log('error: %s', err);
            }else{
                console.log('sent: %s', str);
            }
            if (i % printOn == 0) console.log(i);
        });
        send();
    }, parseInt(env.freq));
}

console.log('connecting: %s', uri);
ws.on('open', function() {
    send();
});


/*
    # this and that

    # loop
    while true; do node client.js; sleep .1; done

    # token
    echo -n 'password' | openssl base64

    # self-signed certs conf, only if needed
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0"

*/