# binrpt

binrpt is a daemon that reads MySQL binlog and executes replication SQL.

![](https://user-images.githubusercontent.com/117768/96328810-c9f47980-1081-11eb-93f5-c00cad75e474.png)

## Usage

```
Usage of binrpt:
  -config string
    	Config file path
  -version
    	Print version and exit
```

## Config example

```toml
[source]
charset = "utf8mb4"
host = "source.db.example.com"
password = "replpwd"
port = 3306
replicate_server_id = 100
username = "repl"

[replica]
charset = "utf8mb4"
host = "replica.db.example.com"
password = "scott"
port = 3306
replicate_do_db = "test"
replicate_ignore_tables = ["^secure_"]
username = "tiger"
```

## Environment variables

* REPLICATE_MAX_RECONNECT_ATTEMPTS
