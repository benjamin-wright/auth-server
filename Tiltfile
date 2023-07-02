allow_k8s_contexts(['auth-server'])
update_settings ( max_parallel_updates = 3 , k8s_upsert_timeout_secs = 60 , suppress_unused_image_warnings = None ) 

load('ext://helm_resource', 'helm_resource')
load('ext://namespace', 'namespace_yaml')

namespace = "auth-server"

k8s_yaml(namespace_yaml(namespace))

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
build('tokens')
build('verify')
build('forms')

k8s_yaml(helm(
    'deploy/chart',
    'auth-server',
    namespace=namespace,
    set=[
        "cockroach.create=true",
        "redis.create=true",
        "users.image=users",
        "tokens.image=tokens",
        "verify.image=verify",
        "forms.image=forms",
    ],
))
