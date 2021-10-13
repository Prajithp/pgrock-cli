# pgrock-cli 
pgrock cli is the commandline client for the pgrock server written in go lang.

### Install
```
go get -v github.com/Prajithp/pgrep-cli
```

### Usage
```
$> pgrock 
Usage: pgrock [OPTIONS]
Options:
  -local-port int
    	Port number of the app server (default 8080)
  -port int
    	Port number of Pgrock Server (default 1080)
  -remote string
    	Pgrock Server Address (default "127.0.0.1")

Example:
		pgrock -r 192.168.22.1 -p 1080 -l 5001

```

