allow_k8s_contexts(['auth-server'])
update_settings ( max_parallel_updates = 3 , k8s_upsert_timeout_secs = 60 , suppress_unused_image_warnings = None ) 

load('ext://helm_resource', 'helm_resource')
load('ext://namespace', 'namespace_yaml')

namespace = "auth-server"
tokenCert = str(local('cat .scratch/localhost.crt | base64')).rstrip('\n')
tokenKey = str(local('cat .scratch/localhost.key | base64')).rstrip('\n')

k8s_yaml(namespace_yaml(namespace))

def build(service, build_cmd="build", resource_deps=[], port_forwards=[], live_update=[]):
    custom_build(
        service,
        'just {} cmd/{} $EXPECTED_REF'.format(build_cmd, service),
        [ 'cmd/{}'.format(service) ],
        ignore = [
            'dist/*',
            '**/*_test.go'
        ],
        live_update = live_update,
    )

    k8s_resource(
        'auth-{}'.format(service),
        auto_init = True,
        trigger_mode = TRIGGER_MODE_MANUAL,
        labels=['auth'],
        resource_deps=resource_deps,
        port_forwards=port_forwards,
    )

build('migrations', build_cmd='build-mig')
build('users', resource_deps=['auth-migrations'])
build('init', resource_deps=['auth-users'])
build('tokens')
build('verify')
build('forms', live_update=[sync('./cmd/forms/static', '/www/static')])

k8s_yaml(helm(
    'deploy/chart',
    'auth-server',
    namespace=namespace,
    set=[
        # General setup
        "adminPassword=Password1!",
        # Image overrides
        "users.image=users",
        "users.migrations.image=migrations",
        "users.init.image=init",
        "tokens.image=tokens",
        "verify.image=verify",
        "forms.image=forms",
    ],
))

k8s_yaml(helm(
    'deploy/mock',
    'mock-home',
    namespace=namespace,
    set=[
        "name=mock-app"
    ]
))

k8s_resource(
    'mock-app',
    labels=[
        'mock'
    ]
)
