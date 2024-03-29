load('ext://helm_resource', 'helm_resource', 'helm_repo')

# Install nginx ingress
# k8s_yaml('nginx-ingress/deploy.yaml')
k8s_yaml('nginx-ingress/keycloak-ingress.yaml')

# Install prometheus
# helm_repo('prometheus-community', 'https://prometheus-community.github.io/helm-charts')
# helm_resource('local-prometheus', 'prometheus-community/prometheus', flags=['--version=24.3.0'])

# Install keycloak
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource('keycloak', 'bitnami/keycloak', flags=['--version=19.3.0'])

# Install Grafana
# helm_repo('grafana', 'https://grafana.github.io/helm-charts')
# helm_resource('local-grafana', 'grafana/grafana', flags=['--version=6.59.4'])

# Install app
k8s_yaml(helm('charts/backend-api-1', name="backend-api-1"))
k8s_yaml(helm('charts/backend-api-2', name="backend-api-2"))
k8s_yaml(helm('charts/frontend', name="frontend"))

# Build: tell Tilt what images to build from which directories
docker_build('alanyoshida/backend-api-1', './backend-api-1')
docker_build('alanyoshida/backend-api-2', './backend-api-2')
docker_build('alanyoshida/frontend', './react-app/frontend')

# Watch: tell Tilt how to connect locally (optional)
k8s_resource('frontend', port_forwards=9999)
k8s_resource('backend-api-1', port_forwards=3000)
k8s_resource('backend-api-2', port_forwards=4000)
# k8s_resource('local-grafana', port_forwards=8080)
# k8s_resource('local-prometheus', port_forwards=["9292:80"])
k8s_resource('keycloak', port_forwards=["8282:8080"])