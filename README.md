## gohttpdummy


- [x] This is a simple golang script that will do a benchmark on a certain URL


- [x] Output the statistics details of the benchmarking



## Compile

```sh

     git clone https://github.com/bayugyug/gohttpdummy .git && cd gohttpdummy

     git pull && make 

```


## Usage

```sh

$ ./gohttpdummy



Version 0.1.0-0

Usage: gohttpdummy [options] [http[s]://]hostname[:port]/path

           Options are:


  -c int
        concurrency  Number of multiple requests to make at a time (default 1)
  -m string
        method       Method to use during the http request (default "GET")
  -r int
        requests     Number of requests to perform (default 1)
  -t int
        timeout      Seconds to max. wait for each response (default 60)

        Example:


                $ ./gohttpdummy -c 5 -r 10 'http://127.0.0.1:7777/parasql/?a=values-a'


                Version 0.1.0-0

                Benchmarking is now in progress ....

                Please be patient!

                Statistics :


                Server Hostname: 127.0.0.1
                Server Port    : 7777
                Document Path  : /parasql/

                SUCCESS : 10
                Elapsed : 22 millisecs
                Requests: 454.54545454545456  (# per seconds)

```

## Docker Binary

- [x] In order to  use it as dockerize binary


``` sh


    sudo  sysctl -w net.ipv4.ip_forward=1

    sudo  docker run --rm  bayugyug/gohttpdummy


```


### License

[MIT](https://bayugyug.mit-license.org/)
