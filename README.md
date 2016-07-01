# redisman
a wrapper for redigo

**Usage:**

1. Write a main.go, build and run
2. curl http://localhost:6667

``` c++  
func main() {

	var config = &redisman.RedisConfig{"localhost", "6379", "", 0, 60}
	var pool = &redisman.PoolConfig{100, 1000, 180}

	redisman.Startup(config, pool)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		var msg string

		ch := make(chan string)
		defer close(ch)

		go RunQuery(ch)

		// fmt.Printf("Result:%s", string(<-ch))
		msg = <-ch
		return c.String(http.StatusOK, msg)
	})

	e.Run(fasthttp.New(":6667"))

}

func RunQuery(ch chan string) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Redis panicing! %s \n", e)
		}
	}()
	conn := redisman.GetConn()
	defer conn.Close()

	for i := 0; i < 1; i++ {
		conn.Do("SET", "TeamName", "Cavaliers")
		res, err := redis.String(conn.Do("GET", "TeamName"))
		if err != nil {
			ch <- "TeamName Not Found: " + strconv.Itoa(i)
		} else {
			ch <- "Found: " + res + "@" + strconv.Itoa(i)
		}
	}
}

```