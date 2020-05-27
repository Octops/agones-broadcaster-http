[![Build Status](https://travis-ci.org/Octops/agones-broadcaster-http.svg?branch=master)](https://travis-ci.org/Octops/agones-broadcaster-http)

# Agones Broadcaster HTTP
Expose Agones GameServers information via HTTP

This project leverages the https://github.com/Octops/agones-event-broadcaster and exposes details about GameServers running within the cluster via an HTTP endpoint.

All the information from the GameServers returned from the [Agones Event Broadcaster](https://github.com/Octops/agones-event-broadcaster) is kept in memory only. There is no persistent storage available.

Considerations:
- It is not possible to recover information from the GameServers if the service is not up and running
- Every time the service starts it will re-sync the in-memory cache from scratch
- If the state of a GameServer changes due to any circumstances, the broadcaster will update the cached info in nearly realtime 
- The service can't be used for updating data

**Important**

Only information from GameServers in a `Ready` state will be returned.

The service returns `json` data in a non specific order. An example is shown below.
```json
{
   "gameservers":[
      {
         "name":"simple-udp-agones-1",
         "namespace":"default",
         "labels":{
            "version":"v1"
         },
         "addr":"172.17.0.2",
         "port":7412,
         "state":"Ready",
         "node_name":"node-us-central1-pool-172-17-0-2"
      },
      {
         "name":"simple-udp-agones-2",
         "namespace":"default",
         "labels":{
            "version":"v1"
         },
         "addr":"172.17.0.2",
         "port":7080,
         "state":"Ready",
         "node_name":"node-us-central1-pool-172-17-0-2"
      },
      {
         "name":"simple-udp-agones-3",
         "namespace":"default",
         "labels":{
            "version":"v1"
         },
         "addr":"172.17.0.2",
         "port":7611,
         "state":"Ready",
         "node_name":"node-us-central1-pool-172-17-0-2"
      }
   ]
}
```

## Install

The command below will push the [install.yaml](install/install.yaml) manifest and deploy the required resources. 

```bash
# Everything will be deployed in the `default` namespace.
$ make install
```

Alternatively, you can deploy the service in a difference namespace

```bash
$ kubectl create ns NAMESPACE_NAME
$ kubectl -n [NAMESPACE_NAME] apply -f install/install.yaml
```

## Fetch Data

### Port-Forward

Use Kubernetes port-forward mechanism to access the service's endpoint running withing the cluster from your local environment.

```bash
# Terminal session #1
$ kubectl [-n NAMESPACE_NAME] port-forward port-forward svc/agones-broadcaster-http 8000

# Terminal session #2
$ curl localhost:8000/api/gameservers
```

### In-Cluster

The service's endpoint will be available to other services running within the cluster using the internal DNS name `agones-broadcaster-http.default.svc.cluster.local`.

### External World

The current install manifest does not expose the service to the external world using `Load Balancers` or the [Ingress Controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).

Check the Kubernetes documentation for more details about [Connecting Applications with Services](https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/).  

## Clean up

```bash
$ kubectl [-n NAMESPACE_NAME] delete -f install/install.yaml
```
