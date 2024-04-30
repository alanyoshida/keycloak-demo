#!/bin/bash

COLOR_RED="31"
COLOR_GREEN="32"
COLOR_YELLOW="33"
COLOR_BLUE="34"

RED="\e[${COLOR_RED}m"
GREEN="\e[${COLOR_GREEN}m"
YELLOW=="\e[${COLOR_YELLOW}m"
BLUE=="\e[${COLOR_BLUE}m"

BOLD="\e[1m"
BOLDGREEN="\e[1;${COLOR_GREEN}m"
BOLDRED="\e[1;${COLOR_RED}m"
BOLDBLUE="\e[1;${COLOR_BLUE}m"
ITALICRED="\e[3;${COLOR_RED}m"
CLEARFORMAT="\e[0m"

check_dependencies(){
    check go
    check yq
    check git
    check gum
    check kind
    check kubectl
    check docker
    check helm
}

check(){
    if [ ! command -v $1 &> /dev/null ]; then
        echo "${BOLDRED}$1 is not installed, please install. Check link in README.md.${CLEARFORMAT}"
        exit 1
    else
        echo -e "${GREEN}$1 - OK${CLEARFORMAT}"
    fi
}

up (){
  CONFIG=$(kubectl config current-context)
  if [ "$CONFIG" == "kind-kind" ]; then
    tilt up
  else
    echo -e "\e[1;31mAtenção: Você não esta no kind, troque o contexto\e[0m"
  fi
}

setup (){
    # DEPENDENCIES
    echo -e "${BOLD}Checking project dependencies ...${CLEARFORMAT}"
    check_dependencies

    # CREATE KIND CLUSTER
    gum confirm "Quer criar um novo cluster kind ?"
    if [ $? -eq 0 ]; then
        echo -e "\n${BOLD}Creating kind cluster${CLEARFORMAT}"
        bash ./kind-with-registry.sh
        kubectl wait -A --for=condition=ready pod --field-selector=status.phase!=Succeeded --timeout=5m
    fi
    # gum spin --show-output --title "Waiting 10s for cluster ..." -- sleep 10

    # METAL LB
    gum confirm "Deseja configurar o metallb?"
    if [ $? -eq 0 ]; then
      echo -e "\n${BOLD}Configuring metallb${CLEARFORMAT}"
      DOCKER_CIDR=$(docker network inspect kind -f '{{(index .IPAM.Config 0).Subnet}}')
      DOCKER_CIDR_2_OCTECTS=$(echo $DOCKER_CIDR | sed -E 's/([0-9]{0,3}\.[0-9]{1,3}).*/\1/')
      cat charts/metallb/values.yaml | yq ".metallb.IPAddressPool.addresses[0]=\"$DOCKER_CIDR_2_OCTECTS.255.200-$DOCKER_CIDR_2_OCTECTS.255.250\"" > charts/metallb/values.yaml
    fi

    # DNSMASQ
    gum confirm "Quer configurar o dnsmasq?"
    if [ $? -eq 0 ]; then
      LB_IP=$(kubectl get svc -n ingress-nginx ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
      # point kind.cluster domain (and subdomains) to our load balancer
      echo "address=/cluster.local/$LB_IP" | sudo tee /etc/dnsmasq.d/kind.k8s.conf
      # restart dnsmasq
      sudo systemctl restart dnsmasq
    fi

    # K8S CONFIG
    up
}

check_folder() {
    if [ -d "$1" ]; then
        # Folder exists
        return 0
    else
        # Folder does not exists
        return 1
    fi
}

install_gum() {
  read -p "Do you want to install gum? (N/y): " -n 1 -r
  echo    # move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    go install github.com/charmbracelet/gum@latest
  fi
}

print_help(){
__usage="
${BOLD}Usage: script.sh [OPTIONS]${CLEARFORMAT}

${BOLD}Options:${CLEARFORMAT}
    up             Runs tilt up if in context kind
    setup, build   Configure everything, create cluster and start tilt
    dependencies   Check if you have the requirements
    -h, --help     Print Help
"
echo -e "$__usage"
}

print_start(){
__usage="${BOLDBLUE}
=================
# Script start  #
=================${CLEARFORMAT}
"

echo -e "$__usage"
}

print_start

case "$1" in

  up)
    up
    ;;
  setup|build)
    setup
    ;;
  dependencies)
    check_dependencies
    ;;
  install_gum)
    install_gum
    ;;
  "--help"|"-h")
    print_help
    ;;
  *)
    echo -e "\e[31merror: Parameter not found.\e[0m"
    print_help

esac

# echo -e "\e[33mScript end\e[0m"