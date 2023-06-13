allow_k8s_contexts(['auth-server'])

load('ext://namespace', 'namespace_yaml')
load('ext://helm_resource', 'helm_resource')
load('./deploy/tilt/db_operator.Tiltfile', 'operator')
load('./deploy/tilt/migrations.Tiltfile', 'migrations')

k8s_yaml(namespace_yaml('auth-server'))

operator(
    name='db-operator',
    namespace='auth-server',
    version='v1.0.3'
)

k8s_yaml(blob("""
    apiVersion: ponglehub.co.uk/v1alpha1
    kind: CockroachDB
    metadata:
        name: cockroach
        namespace: auth-server
    spec:
        storage: 256Mi
"""))

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

    # k8s_resource(
    #     'auth-{}'.format(service),
    #     auto_init = True,
    #     trigger_mode = TRIGGER_MODE_MANUAL,
    #     labels=['auth'],
    # )

build('users')

helm_resource(
    'auth-server',
    'deploy/chart',
    namespace='auth-server',
    flags=[
        "--set=userService.image=users",
    ] + migrations(),
    image_deps=[
        'users'
    ],
    image_keys=[
        'userService.image'
    ],
    labels=['auth']
)
