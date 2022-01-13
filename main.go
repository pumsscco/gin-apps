package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
    "database/sql"
    "log"
    yaml "gopkg.in/yaml.v2"
	"github.com/go-redis/redis"
	"io/ioutil"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "path/filepath"
)
type Conf struct {
    Listen struct {
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
var (
    Db *sql.DB
    logger *log.Logger
    client *redis.Client
    cnf Conf
)
//所有初始化操作
func init() {
    dir,_:=filepath.Abs(filepath.Dir(os.Args[0]))
    //抓全部的配置信息
    yamlBytes, err := ioutil.ReadFile("config.yml")
    if err!=nil {
        log.Fatalf("无法打开环境配置文件: %v",err)
    }
    yaml.Unmarshal(yamlBytes,&cnf)
    file, err := os.OpenFile(dir+"/"+cnf.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("无法打开日志文件：%v\n", err)
    }
    logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
    client=redis.NewClient(&redis.Options{
        Addr:       fmt.Sprintf("%s:%d",cnf.Redis.Host,cnf.Redis.Port),
        Password:   cnf.Redis.Pass,
        DB:         cnf.Redis.Db,
    })
    _, err = client.Ping().Result()
    if err!=nil {
        logger.Fatalf("redis连接异常：%v\n",err)
    }
    dsn:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", cnf.MySQL.User, cnf.MySQL.Pass, cnf.MySQL.Host, cnf.MySQL.Port, cnf.MySQL.Db)
    //fmt.Println("Data Source Name: ",dsn)
    Db,err=sql.Open("mysql",dsn)
    if err!=nil {
        logger.Fatalf("mysql连接异常：%v\n",err)
    }
}
func main() {
	router := gin.Default()
	//以项目为分组，内部接口全面压缩整合
	//先来个考勤项目
	attendance := router.Group("/attendance")
	{
		//依据所带参数是latest/last-week/last-month，转向不同的处理函数或干脆合并，分情形输出，减少接口
		attendance.POST("/rec",rec)
		//与上面采取类似的思路
		attendance.POST("/stat",stat)
		attendance.POST("/add", add)
	}
	//楚留香项目
	crh:=router.Group("/crh")
	{
		crh.POST("/combo", combo)
		crh.POST("/fighter", fighter)
		crh.POST("/item", item)
        crh.POST("/enemy", cEnemy)
        crh.POST("/mission", cMission)
        crh.POST("/role", role)
	}
	//仙剑四项目
	pal4:=router.Group("/pal4")
	{  //除最后拾取与找寻物品的两个接口外，其余均只用一个type参数即可
		pal4.POST("/equipment", equipment)
		pal4.POST("/property", property)
        pal4.POST("/prescription", prescription)
        pal4.POST("/question", question)
        pal4.POST("/mission", pMission)
        pal4.POST("/magic", magic)
        pal4.POST("/stunt", stunt)
        pal4.POST("/upgrade", upgrade)
        pal4.POST("/enemy", pEnemy)
        pal4.POST("/pickup", pickup)
        pal4.POST("/scene", scene)
        pal4.POST("/thing", thing)
        pal4.POST("/find", find)
	}
	//股票项目
	stock:=router.Group("/stock")
	{
		stock.POST("/create", create)
		stock.POST("/list", list)
        stock.POST("/hold-last-deal", holdLastDeal)
        stock.POST("/statistic", statistic)
		stock.POST("/clearance", clearance)
		stock.POST("/position", position)
	}
	router.Run(fmt.Sprintf(":%d",cnf.Listen.Port))
}
