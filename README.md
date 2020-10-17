# binrpt

binrpt is a daemon that reads MySQL binlog and executes replication SQL.

![](https://user-images.githubusercontent.com/117768/96328810-c9f47980-1081-11eb-93f5-c00cad75e474.png)

## Usage

```
Usage of binrpt:
  -config string
    	Config file path
  -dryrun
    	Dry-run mode
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

## Getting Started

```sh
docker-compose build
docker-compose up -d
cp config.toml.example config.toml
make
./binrpt -config config.toml
 ```

 ```
~$ mysql -u root -h 127.0.0.1 -P 13307 test -e 'select * from test'
+----+-----+-----+
| id | num | str |
+----+-----+-----+
|  1 | 100 | foo |
|  2 | 200 | bar |
|  3 | 300 | zoo |
+----+-----+-----+
~$ mysql -u root -h 127.0.0.1 -P 13306 test -e 'insert into test (num, str) values (1, "abc")'
~$ mysql -u root -h 127.0.0.1 -P 13306 test -e 'update test set num = id + 1000 where id = 2'
~$ mysql -u root -h 127.0.0.1 -P 13306 test -e 'delete from test where id = 1'
~$ mysql -u root -h 127.0.0.1 -P 13307 test -e 'select * from test'
+----+------+-----+
| id | num  | str |
+----+------+-----+
|  2 | 1002 | bar |
|  3 |  300 | zoo |
|  4 |    1 | abc |
+----+------+-----+
```
