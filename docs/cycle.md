# Cycle

## Introduction

The file on disk is encrypted with the key in Vault.
The file and key are linked, which means that
- you need both to get to the decrypted data
- a file cannot be decrypted by another (version of) the key.

Should a version of the key, or a version of the file be exposed by any means, one could gain access to the other and get access to the decrypted data.
As such, it is best practice to periodically expire old versions by cycling to a new file and key, invalidating older versions.

**Note** that once a key and file are exposed, there is no option to invalidate the pair, old pairs still work.
Cycling then requires creating a new PostgreSQL cluster with a new key pair and migrating data.
But cycling is great protection against either one being exposed...

## Options

By using the `pgcustodian cycle --help` command all options are revealed.
For a bit more description, please refer to [options](options.md).

```bash
pgcustodian cycle --help
Use this command to setup hashicorp vault so that it is ready to be used by pgcustodian.

Usage:
  pgcustodian stage [flags]

Flags:
  -c, --cfgFile string            config file (default "~/.pgcustodian/config.yaml")
  -C, --customPolicyFile string   file with custom policy to be applied (default "~/.pgcustodian/~/.pgcustodian/custom_policy.hcl")
  -h, --help                      help for stage
  -r, --roleId string             role id for logging into vault.
  -R, --roleIdFile string         path to file with role id for logging into vault. (default "~/.pgcustodian/role-id")
  -n, --roleName string           role id for logging into vault. (default "macbook-pro.home")
  -s, --secretIdFile string       secret id for logging into vault. (default "~/.pgcustodian/secret-id")
  -p, --storePath string          path where private key should be stored.
  -V, --storeVersion uint         version of vault store. (default 2)
  -t, --token string              token for logging into vault.
  -T, --tokenFile string          tokenFile can be set to a path containing the token for logging into vault. (default "~/.pgcustodian/token")
  -v, --verbose count             Be more verbose in the output.
  -w, --wrapped                   wrap replies with the token
```

## How it works

The command `pgcustodian cycle` does the following:
- creates temp backup keys
- retrieves original password from vault
- creates a backup of the original key, encrypted with temp backup keys public key
- generates a new password
- creates a backup of the newly generated key, encrypted with temp backup keys public key
- stream original file, decrypts with original key, encrypts with new key, write to same file with .cycled suffix
- renames cycled file to original file
- update password in vault to newly generated password
- (when public key is available), key will be back'ed up with official public backup key
- shred temp files

