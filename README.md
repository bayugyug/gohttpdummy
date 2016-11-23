## gohttpdummy


- [x] This is a simple golang script that will do a benchmark on a certain URL


- [x] Output the statistics details of the benchmarking



## Usage

```sh

$ ./gohttpdummy


Version 0.1.0-0

Usage: gohttpdummy [options] [http[s]://]hostname[:port]/path

           Options are:


  -c int
        concurrency  Number of multiple requests to make at a time (default 1)
  -d string
        form         Form data for POST method
  -m string
        method       Method to use during the http request (default "GET")
  -r int
        requests     Number of requests to perform (default 1)
  -t int
        timeout      Seconds to max. wait for each response (default 60)



        *** Example: (GET)


                $  ./gohttpdummy -c 10 -r 500  'http://192.168.2.121:7777/parasql/?p=GAME_ACTION'

                    Version 0.1.0-0

                    Benchmarking is now in progress ....

                    Please be patient!

                    Statistics :


                    Server Hostname: 192.168.2.121
                    Server Port    : 7777
                    Document Path  : /parasql/

                    Success :  500
                    Elapsed :  1309.107836 ( millisecs )
                    Requests:  381.939506  ( # per sec )
                    App Time:  1.309107836s
                    Sys Time:  1.309200187s


        *** Example: (POST)


                $ ./gohttpdummy -c 10 -r 50  -d "m=aguy&r=dabis&t=hehehe&data=mundo" -m "POST" 'http://192.168.2.121:7777/parasql'


                    Version 0.1.0-0

                    Benchmarking is now in progress ....

                    Please be patient!

                    Statistics :


                    Server Hostname: 192.168.2.121
                    Server Port    : 7777
                    Document Path  : /parasql

                    Success :  50
                    Elapsed :  100.263404 ( millisecs )
                    Requests:  498.686440 ( # per sec )
                    App Time:  100.263404ms
                    Sys Time:  100.998003ms


```

## Docker Binary

- [x] In order to  use it as dockerize binary


``` sh


    sudo  sysctl -w net.ipv4.ip_forward=1

    sudo  docker run --rm  bayugyug/gohttpdummy


```


### License

[MIT](https://bayugyug.mit-license.org/)
