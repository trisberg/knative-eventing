#!/bin/bash

# Copyright 2018 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This is a collection of useful bash functions and constants, intended
# to be used in test scripts and the like. It doesn't do anything when
# called from command line.

# Default GKE version to be used with eventing
readonly EVENTING_GKE_VERSION=1.9.6-gke.1

# Useful environment variables
[[ -n "${PROW_JOB_ID}" ]] && IS_PROW=1 || IS_PROW=0
readonly IS_PROW
readonly EVENTING_ROOT_DIR="$(dirname $(readlink -f ${BASH_SOURCE}))/.."
readonly OUTPUT_GOBIN="${EVENTING_ROOT_DIR}/_output/bin"

# Copy of *_OVERRIDE variables
readonly OG_DOCKER_REPO="${DOCKER_REPO_OVERRIDE}"
readonly OG_K8S_CLUSTER="${K8S_CLUSTER_OVERRIDE}"
readonly OG_K8S_USER="${K8S_USER_OVERRIDE}"
readonly OG_KO_DOCKER_REPO="${KO_DOCKER_REPO}"

# Returns a UUID
function uuid() {
  # uuidgen is not available in kubekins images
  cat /proc/sys/kernel/random/uuid
}

# Simple header for logging purposes.
function header() {
  echo "================================================="
  echo ${1^^}
  echo "================================================="
}

# Simple subheader for logging purposes.
function subheader() {
  echo "-------------------------------------------------"
  echo $1
  echo "-------------------------------------------------"
}

# Restores the *_OVERRIDE variables to their original value.
function restore_override_vars() {
  export DOCKER_REPO_OVERRIDE="${OG_DOCKER_REPO}"
  export K8S_CLUSTER_OVERRIDE="${OG_K8S_CLUSTER}"
  export K8S_USER_OVERRIDE="${OG_K8S_CLUSTER}"
  export KO_DOCKER_REPO="${OG_KO_DOCKER_REPO}"
}

# Waits until all pods are running in the given namespace.
# Parameters: $1 - namespace.
function wait_until_pods_running() {
  echo -n "Waiting until all pods in namespace $1 are up"
  for i in {1..150}; do  # timeout after 5 minutes
    local pods="$(kubectl get pods -n $1 | grep -v NAME)"
    local not_running=$(echo "${pods}" | grep -v Running | wc -l)
    if [[ -n "${pods}" && ${not_running} == 0 ]]; then
      echo -e "\nAll pods are up:"
      kubectl get pods -n $1
      return 0
    fi
    echo -n "."
    sleep 2
  done
  echo -e "\n\nERROR: timeout waiting for pods to come up"
  kubectl get pods -n $1
  return 1
}

# Sets the given user as cluster admin.
# Parameters: $1 - user
#             $2 - cluster name
#             $3 - cluster zone
function acquire_cluster_admin_role() {
  # Get the password of the admin and use it, as the service account (or the user)
  # might not have the necessary permission.
  local password=$(gcloud --format="value(masterAuth.password)" \
      container clusters describe $2 --zone=$3)
  kubectl --username=admin --password=$password \
      create clusterrolebinding cluster-admin-binding \
      --clusterrole=cluster-admin \
      --user=$1
}

# Authenticates the current user to GCR in the current project.
function gcr_auth() {
  echo "Authenticating to GCR"
  # kubekins-e2e images lack docker-credential-gcr, install it manually.
  # TODO(adrcunha): Remove this step once docker-credential-gcr is available.
  gcloud components install docker-credential-gcr
  docker-credential-gcr configure-docker
  echo "Successfully authenticated"
}

# Installs ko in $OUTPUT_GOBIN
function install_ko() {
  GOBIN="${OUTPUT_GOBIN}" go install ./vendor/github.com/google/go-containerregistry/cmd/ko
}

# Runs ko; prefers using the one installed by install_ko().
# Parameters: $1..$n - arguments to ko
function ko() {
  if [[ -e "${OUTPUT_GOBIN}/ko" ]]; then
    "${OUTPUT_GOBIN}/ko" $@
  else
    ko $@
  fi
}
