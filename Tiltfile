allow_k8s_contexts(['auth-server'])

load('ext://namespace', 'namespace_yaml')
load('ext://helm_resource', 'helm_resource')
load('./deploy/tilt/db_operator.Tiltfile', 'operator')

k8s_yaml(namespace_yaml('auth-server'))

operator(
    name='db-operator',
    namespace='auth-server',
    version='v1.0.3'
)

def build(service):
    custom_build(
        service,
        'just build cmd/{} $EXPECTED_REF'.format(service),
        [ 'cmd/{}'.format(service) ],
        ignore = [
            'dist/*',
            '**/*_test.go'
        ]
    )

    k8s_resource(
        'auth-{}'.format(service),
        auto_init = True,
        trigger_mode = TRIGGER_MODE_MANUAL,
        labels=['auth'],
    )

build('users')

k8s_yaml(helm(
    'deploy/chart',
    'auth-server',
    namespace='auth-server',
    set=[
        "cockroach.create=true",
        "userService.image=users",
    ],
))
