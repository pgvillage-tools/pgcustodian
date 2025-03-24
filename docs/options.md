# Options of pgcustodian

The following options can be used.

**Note** that not all options are applicable to all commands.
Issue `pgcustodian [command] --help` to check if an option is applicable to the specific subcommand

## List of options

- `--cfgFile`, `-c`
  Use this option to point to a config file (defaults to "~/.pgcustodian/config.yaml")
  If config file exists, pgcustodian reads all config options from it. COmmandline options take precedence.
- `--customPolicyFile`, `-C`
  Specific to the stage command.
  Set a path to a custom vault policy to be used during the stage process.
  Defaults to "~/.pgcustodian/~/.pgcustodian/custom_policy.hcl".
  If this exists, pgcustodian reads this and sets it as the policy instead of default.
- `--help` , `-h`
  Request help on pgcustodian, or any subcommand
- `--roleId`, `-r`
  You can  directly set a role-id to be used for logging in.
- `--roleIdFile` `-R`
  You can read the role-id from a file.
  Defaults to "~/.pgcustodian/role-id".
  If the file exists, pgcustodian reads its contents and uses it as role-id.
  Note that this is a direct parse. File should only contain the role-id, and nothing else. No extra spaces, \n characters, # [comments], etc.
- `--roleName`, `-n`
  Specific to the stage command.
  When set, the `stage` command created a role-id with this name sineatd of the default (which is set to the hostname)
- `--secretIdFile`, `-s`
  You can read the secret-id from a file
  Defaults to "~/.pgcustodian/secret-id".
  If the file exists, pgcustodian reads its contents and uses it as secret-id.
  Note that this is a direct parse. File should only contain the secret-id, and nothing else. No extra spaces, \n characters, # [comments], etc.
- `--storePath`, `-p`
  You can set a path to the location of the private key in Vault.
- `--storeVersion`, `-V`
  you can use a kv1 store over a kv2 store if you feel that you must.
- `--token` `-t`
  You can directly set a token (and skip the application role).
  Main usage is during the stage command, but setting this takes precedence over setting role-id and secret-id
- `--tokenFile` `-T`
  You can read the token from a file. Defaults to "~/.pgcustodian/token".
  **Note** that the login command creates this file with the token to be used by subsequent commands.
- `--verbose` `-v`
  Every additional -v argument increases verbosity of pgcustodian.
- `--wrapped` `-w`
  You could disable vault [response wrapping](https://developer.hashicorp.com/vault/docs/concepts/response-wrapping).
  There is no test for this option, if you run into issues, please let us know and we will improve unittests to accompany this feature as well.

  - `--backupFile`, `-b`
  For commands that use the backup option, you can set path to backup file with secret encrypted with public key. (default "~/.pgcustodian")
  -G, --generatedPasswordChars string   character list for generating passwords (default "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~0123456789")
  -g, --generatedPasswordLength uint    length for generated passwords (default 16)
  -k, --publicKeyPath string            path where public key should be stored. (default "~/.pgcustodian/public.pem")
  -S, --secretKey string                path in kv1 or kv2 store where secrets are held.
  -P, --secretPath string               path in kv1 or kv2 store where secrets are held. (default "pgcustodian/macbook-pro.home")
  -x, --shred                           shred the tmp files (default true)
