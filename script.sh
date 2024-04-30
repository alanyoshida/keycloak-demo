#!/bin/bash

COLOR_RED="31"
COLOR_GREEN="32"
COLOR_YELLOW="33"
COLOR_BLUE="34"

RED="\e[${COLOR_RED}m;"
GREEN="\e[${COLOR_GREEN}m;"
YELLOW=="\e[${COLOR_YELLOW}m;"
BLUE=="\e[${COLOR_BLUE}m;"

BOLD="\e[1m"
BOLDGREEN="\e[1;${COLOR_GREEN}m"
BOLDRED="\e[1;${COLOR_RED}m"
ITALICRED="\e[3;${COLOR_RED}m"
CLEARFORMAT="\e[0m"

check_dependencies(){
    check go
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
  # echo -e "CONFIG=$CONFIG"
  if [ "$CONFIG" == "kind-kind" ]; then
    tilt up
  else
    echo -e "\e[1;31mAtenção: Você não esta no kind, troque o contexto\e[0m"
  fi
}

setup (){
    echo -e "${BOLD}Checking project dependencies ...${CLEARFORMAT}"
    gum confirm "Quer criar um novo cluster kind ?"
    if [ $? -eq 0 ]; then
        echo -e "\n${BOLD}Creating kind cluster${CLEARFORMAT}"
        bash ./kind-with-registry.sh
    fi
    gum spin --show-output --title "Waiting 10s for cluster ..." -- sleep 10
    CONFIG=$(kubectl config current-context)
    # echo -e "CONFIG=$CONFIG"
    if [ "$CONFIG" == "kind-kind" ]; then
        gum confirm "Quer instalar o nginx no cluster $CONFIG ?"
        if [ $? -eq 0 ]; then
            echo -e "kubectl apply -f ./nginx-ingress/deploy.yaml"
            kubectl apply -f ./nginx-ingress/deploy.yaml
        fi
        echo -e "\n${BOLD}Starting tilt${CLEARFORMAT}"
        tilt up
    fi
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
    setup          Configure everything, create cluster and start tilt
    clone          Clone repos
    link           Create repo symlinks
    dependencies   Check if you have the requirements
    -h, --help     Print Help
"
echo -e "$__usage"
}

print_make_help(){
__usage="
${BOLD}Usage: make [OPTIONS]${CLEARFORMAT}

${BOLD}Options:${CLEARFORMAT}
    up              Runs 'tilt up' if in context kind
    setup           Configure everything, create cluster and start tilt
    clone           Clone repos
    link            Create repo symlinks
    dependencies    Check if you have the requirements
    create-cluster  Create kind cluster
"
echo -e "$__usage"
}

print_start(){
__usage="
####################################
#       PROJECT                    #
####################################
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
  clone)
    clone_repos
    ;;
  link)
    symlinks
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
  "make-help")
    print_make_help
    ;;
  *)
    echo -e "\e[31merror: Parameter not found.\e[0m"
    print_help

esac

# echo -e "\e[33mScript end\e[0m"