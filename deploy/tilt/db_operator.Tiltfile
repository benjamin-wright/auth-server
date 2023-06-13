v1alpha1.extension_repo(name='db-operator', url='https://github.com/benjamin-wright/db-operator')
v1alpha1.extension(name='db_operator', repo_name='db-operator', repo_path='tilt/db_operator')
load('ext://db_operator', 'db_operator')

def operator(namespace, name, version, labels=[]):
   db_operator(namespace, name, version, labels)
