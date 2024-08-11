# Package
`github.com/tappoy/backupd`

See [Document.txt](Document.txt)

# Config of Source
- [backupd.config](./test-data/config/backupd.config)

# Config of Destination
- aws: [.aws.config](https://github.com/tappoy/storage/blob/main/.aws.config.sample)
- local: [.local.config](https://github.com/tappoy/storage/blob/main/.local.config.sample)
- openstack: [.openstack.config](https://github.com/tappoy/storage/blob/main/.openstack.config.sample)

# Operation Example
```sh
# set vault dir and password
user@local$ cat .vault-password | ssh user@remote 'backupd vault /opt/vault'
```

# Structure Example
```mermaid
flowchart

%% Declaration
%%% Elements
Memo2[Private Git]
Memo3[Private Git]
%%Memo4[Crypted]

%%% Process
localssh(ssh)
backupd(backupd)

%%% Data
LocalVault[(VaultPassword)]
RemoteVault[(Vault\nDest Config\n*Crypted*)]
SourceConfig[(SourceConfig)]
Cloud2Bucket[(Cloud2Bucket\n*Crypted*)]
Cloud1Bucket[(Cloud1Bucket\n*Crypted*)]
LocalBucket[(LocalBucket\n*Crypted*)]
SourceData[(SourceData)]

%%% Subgraph
subgraph Laptop
    LocalVault
    localssh
end

subgraph Server
    RemoteVault
    backupd
    SourceConfig
    SourceData
    LocalBucket
end

subgraph Cloud1
    Cloud1Bucket
end

subgraph Cloud2
    Cloud2Bucket
end


%% Arrow
localssh -->|vault password| backupd
localssh -.-> LocalVault
backupd -.-> RemoteVault
backupd -.-> SourceConfig
backupd --> Cloud1Bucket
backupd --> Cloud2Bucket
backupd --> LocalBucket
backupd -.-> SourceData

SourceConfig -.- Memo2
RemoteVault -.- Memo3


%% Class
%%% Memo
classDef Memo fill:#9f9,color:#333,stroke:none
class Memo,Memo1,Memo2,Memo3,Memo4,Memo5 Memo

classDef Data fill:#f66,color:#333
class SourceData Data

classDef Buckup fill:#f99,color:#333
class Cloud1Bucket,Cloud2Bucket,LocalBucket Buckup

classDef Config fill:#9ff,color:#333
class SourceConfig,RemoteVault Config

classDef Vault fill:#dd6,color:#333
class LocalVault Vault
```


# Dependencies
- [github.com/tappoy/env](https://github.com/tappoy/env)
- [github.com/tappoy/storage](https://github.com/tappoy/storage)
- [github.com/tappoy/vault-cli](https://github.com/tappoy/vault-cli)
- [github.com/tappoy/vault](https://github.com/tappoy/vault)

# Why it is this way.
See [Philosophy](https://github.com/tappoy/philosophy) for more details.

# License
[GPL-3.0](LICENSE)

# Author
[tappoy](https://github.com/tappoy)