# Keycloak Demonstration

This repository is a demonstration of how to configure keycloak and use in a frontend react app, and multiple backends.

The frontend client will access backend-api-1 and backend-api 2 with aud.

The backend-api-1 client will access backend-api-2 with aud.

## Requirements to run in cluster
- Go
- Tilt
- microk8s, kind, or another with local registry

## Requirements to run locally
- node >12
- npm
- npx
- go

Tilt will install helm repos, install the keycloak chart, build the containers and push to local registry, then port forward to access with localhost.

## Steps to run the demo

### 1 - Create the kind cluster with registry using the script

`./kind-with-registry.sh`

### 2 - Use Tilt to set everything up

`tilt up`

### 3 - Open browser in `http://localhost:10350/`, this is Tilt dashboard

### 3.5 - Install Nginx Ingress

Sometimes the ingress definition inside the charts folder, fail to bring nginx controller up, in that case, try to reinstall using the following command.

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

### 4 - Keycloak

For login use `user` as the username, and get the password from the secret with:

```bash
kubectl get secret keycloak -ojsonpath="{.data.admin-password}" | base64 -d
```

You will need to map the following domain inside /etc/hosts, them you can access the keycloak interface using the domain: keycloak.default.svc.cluster.local

```
127.0.0.1 keycloak.default.svc.cluster.local
```

This is because the backend-api will validate the provider domain, must be the same that generated the access token.

### 5 - Keycloak clients configuration

Create the frontend client, as public
```
client-id: frontend
Client authentication: Disabled
Standard flow: Enabled
Valid redirect URIs: *
Web origins: *
Direct access grants: Enabled
Implicit flow: Disabled
Service accounts roles: Disabled
OAuth 2.0 Device Authorization Grant: Disabled
OIDC CIBA Grant: Disabled
Authorization: Disabled
```

Create the backend-api-1 client as confidential
```
client-id: backend-api-1
Client authentication: Enabled
Standard flow: Disabled
Direct access grants: Disabled
Implicit flow: Disabled
Service accounts roles: Disabled
OAuth 2.0 Device Authorization Grant: Disabled
OIDC CIBA Grant: Disabled
Authorization: Disabled
```

Create a new client scope named `backend-api-1` and add a mapper of type `audience`, then select the `backend-api-1` as `Included Client Audience`

Put the new client scope named `backend-api-1` inside the frontend client. This will enable the tokens from frontend to be used in the backend-api-1

Create the backend-api-2 client as confidential
```
client-id: backend-api-2
Client authentication: Enabled
Standard flow: Disabled
Direct access grants: Disabled
Implicit flow: Disabled
Service accounts roles: Disabled
OAuth 2.0 Device Authorization Grant: Disabled
OIDC CIBA Grant: Disabled
Authorization: Disabled
```

Create a new client scope named `backend-api-2` and add a mapper of type `audience`, then select the `backend-api-2` as `Included Client Audience`

Put the new client scope `backend-api-2` inside the frontend client. This will enable the tokens from frontend to be used in backend-api-2.

Them put the client scope `backend-api-1` inside the backend-api-2 client. This will enable tokens from backend-api-1 to be used in backend-api-2.

### 6 - Configure backends charts

Edit the file `./charts/backend-api-1/values.yaml` and set the `client secret` from the `credentials` tab from keycloak

```yaml
- name: CLIENT_SECRET
  value: ""
```

Edit the file `./charts/backend-api-2/values.yaml` and set the `client secret` from the `credentials` tab from keycloak

```yaml
- name: CLIENT_SECRET
  value: ""
```
### How to access the demo

You can configure the `/etc/hosts` to access using the following domains, but you can also use the port-forwarding like:

| app | description |
|-|-|
| frontend  | port_forwards to 9999 |
| backend-api-1 | port_forwards to 3000 |
| backend-api-2 | port_forwards to 4000 |

But for keycloak always use the `keycloak.default.svc.cluster.local` because it matters to validate the provider in token generation.

---

`/etc/hosts` file:
```
127.0.0.1 frontend.default.svc.cluster.local
127.0.0.1 backend-api-1.default.svc.cluster.local
127.0.0.1 backend-api-2.default.svc.cluster.local
```

### How to test

Request a access token for the frontend client

```bash
curl  -X POST \
  'http://keycloak.default.svc.cluster.local/realms/master/protocol/openid-connect/token' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'grant_type=password' \
  --data-urlencode 'username=user' \
  --data-urlencode 'password=<admin-password>' \
  --data-urlencode 'client_id=frontend'
```

You can create you own user in keycloak and use that authentication instead.

With the access token, make a request to the backend-api-1 like the following

```bash
curl  -X GET \
  'http://localhost:3000/profile/name' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer <ACCESS_TOKEN>'
```

You should see the user profile data.

Backend 2 test:

```bash
curl  -X GET \
  'http://localhost:4000/protected/pets/list' \
  --header 'Accept: */*' \
  --header 'Authorization: Bearer <ACCESS_TOKEN>'
```

You should get a json with pets.