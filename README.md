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
[master]
host = "master.db.example.com"
password = "replpwd"
port = 3306
server_id = 100
username = "repl"

[replica]
host = "replica.db.example.com"
password = "scott"
port = 3306
replicate_do_db = "test"
replicate_ignore_tables = ["^secure_"]
username = "tiger"
```
