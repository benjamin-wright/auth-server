def load_migrations(path):
    file = migration_name(path)
    [ database, index ] = file.split('-')
    migration = str(local(
        'cat %s' % path,
        echo_off=True,
        quiet=True,
    )).rstrip('\n').split('\n')

    return [
        '--set=userService.cockroach.migrations.%s.database=%s' % (file, database),
        '--set=userService.cockroach.migrations.%s.index=%s' % (file, index),
        '--set-file=userService.cockroach.migrations.%s.migration=%s' % (file, path),
    ]

def migration_name(path):
    return path.split('/')[-1].replace('.sql', '')

def migrations():
    files = str(local(
        'find cmd/users/migrations -name *.sql',
        echo_off=True,
        quiet=True,
    )).rstrip('\n').split('\n')

    return [ flag for file in files for flag in load_migrations(file) ]