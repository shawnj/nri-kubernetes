# -*- mode: Python -*-

# Settings and defaults.

project_name = 'nri-kubernetes'

chart_name= 'newrelic-infrastructure-v3'

settings = {
  'kind_cluster_name': 'nri-k8s-dev',
  'default_registry': 'localhost:5005',
  'live_reload': True,
  'chart_path': '../helm-charts-newrelic/charts/%s/' % chart_name,
}

settings.update(read_json('tilt_option.json', default={}))

default_registry(settings.get('default_registry'))

# Only use explicitly allowed kubeconfigs as a safety measure.
allow_k8s_contexts(settings.get("allowed_contexts", "kind-" + settings.get('kind_cluster_name')))


# Building Docker image.
load('ext://restart_process', 'docker_build_with_restart')

if settings.get('live_reload'):
  binary_name = '%s-linux' % project_name

  # Building daemon binary locally.
  local_resource(
    '%s-binary' % project_name, 
    'GOOS=linux BIN_DIR=. make compile',
    deps=[
            "./src",
            "./internal",
            "./cmd"
        ],
    )

  # Use custom Dockerfile for Tilt builds, which only takes locally built daemon binary for live reloading.
  dockerfile = '''
    FROM golang:1.17-alpine AS dlv-builder
    RUN apk add gcc musl-dev && \
        go install github.com/go-delve/delve/cmd/dlv@latest
    COPY %s /usr/local/bin/%s
  ''' % (binary_name, project_name)

  docker_build_with_restart(project_name, '.',
    dockerfile_contents=dockerfile,
    entrypoint=[
      "dlv",
      "--listen=0.0.0.0:50100",
      "--headless=true",
      "--api-version=2",
      "--check-go-version=false",
      "--only-same-user=false",
      "--accept-multiclient",
      "exec",
      "/usr/local/bin/%s" % project_name,
      "--"
    ],
    only=binary_name,
    live_update=[
      # Copy the binary so it gets restarted.
      sync('bin/%s' % binary_name, '/usr/local/bin/%s' % project_name),
    ],
  )
else:
  docker_build(project_name, '.')

# Deploying Kubernetes resources.
k8s_yaml(helm(settings.get('chart_path'), name='nr', values=['values-dev.yaml', 'values-local.yaml']))

k8s_resource(workload='nr-nrk8s-controlplane', port_forwards='50100:50100')
k8s_resource(workload='nr-nrk8s-kubelet', port_forwards='50101:50100')
k8s_resource(workload='nr-nrk8s-ksm', port_forwards='50102:50100')

# e2e resources chart has KSM as a dependency. This resource updates this chart dependency.
local_resource('ksm deps', 'helm dependency update e2e/charts/e2e-resources', deps=['e2e/charts/e2e-resources/Chart.yaml'])

k8s_yaml(helm('e2e/charts/e2e-resources', name='e2e-resources'))
