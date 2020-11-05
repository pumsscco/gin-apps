package main

import (
    "fmt"
    "database/sql"
    "log"
    "net/http"
    yaml "gopkg.in/yaml.v2"
	"github.com/go-redis/redis"
	"io/ioutil"
    "github.com/julienschmidt/httprouter"
    _ "github.com/go-sql-driver/mysql"
    "os"
)
var Db *sql.DB
var logger *log.Logger
var client *redis.Client
type Conf struct {
    Listen struct {
        Host string `yaml:"host"`
        Port int `yaml:"port"`
    }
    MySQL struct {
        Db string `yaml:"db"`
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        User string `yaml:"user"`
        Pass string `yaml:"pass"`
    }
    Redis struct {
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        Db int `yaml:"db"`
        Pass string `yaml:"pass"`
    }
	Logfile string `yaml:"logfile"`
}
var cnf Conf
func init() {
    //抓全部的配置信息
    yamlBytes, _ := ioutil.ReadFile("config.yml")
    yaml.Unmarshal(yamlBytes,&cnf)
    file, err := os.OpenFile(cnf.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("无法打开日志文件", err)
    }
    logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
    client=redis.NewClient(&redis.Options{
        Addr:       fmt.Sprintf("%s:%d",cnf.Redis.Host,cnf.Redis.Port),
        Password:   cnf.Redis.Pass,
        DB:         cnf.Redis.Db,
    })
    pong, err := client.Ping().Result()
    logger.Println(pong, err)
    dsn:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", cnf.MySQL.User, cnf.MySQL.Pass, cnf.MySQL.Host, cnf.MySQL.Port, cnf.MySQL.Db)
    //fmt.Println("Data Source Name: ",dsn)
    Db,err=sql.Open("mysql",dsn)
    if err!=nil {
        logger.Fatalf("open mysql failed: %v",err)
    }
}
func main() {
    // handle static assets
    router := httprouter.New()
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    //页面
    router.GET("/", index)
    //以时间逆序列出该股票代码的相关所有交易记录，新增记录成功后，也会重定向到这里
    router.GET("/name-list", dealList)
    router.POST("/name-list", dealList)
    router.GET("/hold-last-deal", holdLastDeal)
    router.GET("/clearance", clearance)
    router.GET("/position", position)
    router.GET("/add", newDeal)
    router.POST("/add", newDeal)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d",cnf.Listen.Host,cnf.Listen.Port),router))
}
