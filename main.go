package main

import (
	"common/ini"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

func main() {

	conf := "dedup.conf"

	cfg := ini.DumpAll(conf)

	passwd := cfg["DEFAULT:"+"passwd"]
	scan := cfg["DEFAULT:"+"scan"]
	batch, _ := strconv.Atoi(cfg["DEFAULT:"+"batch"])
	maxCPU, _ := strconv.Atoi(cfg["DEFAULT:"+"maxCPU"])
	urls := cfg["DEFAULT:"+"urls"]
	urlSlice := strings.Split(urls, ";")

	runtime.GOMAXPROCS(maxCPU)

	fmt.Printf("%12s %s\n", "passwd", passwd)
	fmt.Printf("%12s %s\n", "scan", scan)
	fmt.Printf("%12s %d\n", "batch", batch)
	fmt.Printf("%12s %d\n", "maxCPU", maxCPU)

	fmt.Println("========================================")

	// urlSlice 转成 urlMap 去重
	urlMap := map[string]string{}
	for _, u := range urlSlice {
		urlMap[u] = "1"
	}

	var wg sync.WaitGroup

	for url := range urlMap {
		// key is url
		wg.Add(2)

		go func() {
			defer wg.Done()
			dedup(url, passwd, scan, "01:*", batch)
		}()

		go func() {
			defer wg.Done()
			dedup(url, passwd, scan, "04:*", batch)
		}()
	}

	wg.Wait()

}

func dedup(url string, passwd string, scan string, pattern string, batch int) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("%s@%s Failed\n", pattern[:2], url)
		}
	}()

	// 纳秒时间戳
	ts := fmt.Sprint(time.Now().UnixNano())
	// 根据 url, pattern, ts 生成文件名
	filename := strings.Replace(url, ":", "_", -1) + "_" + pattern[:2] + "_" + ts + ".log"
	logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0664)
	defer logfile.Close()
	logger := log.New(logfile, "", log.Ldate|log.Ltime)

	if err != nil {
		logger.Printf("%s\n", err.Error())
		panic(err)
	}

	logger.Println(">>> Start")

	// 连接 redis
	rs, err := redis.Dial("tcp", url)
	if err != nil {
		logger.Printf("create connection to redis fail. url: %s\n", url)
		panic(err)
	}

	// 认证 redis
	if passwd != "" {
		if _, err := rs.Do("AUTH", passwd); err != nil {
			logger.Println("AUTH fail.")
			rs.Close()
			panic(err)
		}
	}

	curInit := "0" //cur 初始值
	curCurt := "0" //cur 当前值
	curNext := "x" //cur 下一值
	for {

		repl, err := redis.Values(rs.Do(scan, curCurt, "MATCH", pattern, "COUNT", batch))
		if err != nil {
			logger.Printf("command scan fail. command: %s\n", "SCAN "+curCurt+" MATCH "+"04:* COUNT "+fmt.Sprintf("%d", batch))
			panic(err)
		}

		for _, val := range repl {

			switch val.(type) {
			case []uint8:

				curNext, _ = redis.String(val, nil)

			case []interface{}:

				keys, err := redis.Strings(val, nil)
				if err != nil {
					logger.Printf("get keys from scan fail. %s quit process. \n", err)
					curNext = curInit
					panic(err)
					// break
				}
				for _, key := range keys {

					logger.Printf(">>> processing key [ %s ]\n", key)

					v, err := redis.Values(rs.Do("HKEYS", key))
					if err != nil {
						logger.Printf("hgetall fail. %s quit process. \n", err)
						curNext = curInit
						panic(err)
						// break
					}

					m, err := redis.Strings(v, err)
					if err != nil {
						logger.Printf("map hash fail. %s quit process. \n", err)
						curNext = curInit
						panic(err)
						// break
					}

					oldFields, newFieldsMap := []string{}, map[string]string{}
					for _, field := range m {
						if (field[:1] == "h") && (field[1:2] == "0") {
							newFieldsMap[field] = "exists"
						} else if (field[:1] == "h") && (field[1:2] != "0") {
							oldFields = append(oldFields, field)
						}
					}

					if len(oldFields) > 0 && len(newFieldsMap) > 0 {
						for _, of := range oldFields {

							f1 := "h00" + strings.Repeat("0", 5-len(of)) + of[1:]
							f2 := "h01" + strings.Repeat("0", 5-len(of)) + of[1:]
							if _, ok := newFieldsMap[f1]; ok {
								logger.Printf("... processing field: [ %s ] with the opposite new field: [ %s ]\n", of, f1)
								_, err = rs.Do("HDEL", key, of)
								if err != nil {
									logger.Printf("... command faild: HDEL %s %s \n", key, of)
									panic(err)
								}
								logger.Printf("... new field [ %s ] exists, old field [ %s ] removed\n", f1, of)
							} else if _, ok := newFieldsMap[f2]; ok {
								logger.Printf("... processing field: [ %s ] with the opposite new field: [ %s ]\n", of, f1)
								_, err = rs.Do("HDEL", key, of)
								if err != nil {
									logger.Printf("... command faild: HDEL %s %s \n", key, of)
									panic(err)
								}
								logger.Printf("... new field [ %s ] exists, old field [ %s ] removed\n", f1, of)

							}
						}
					}

					logger.Println()

				}

			}
		}

		if curNext == curInit {
			break
		} else {
			curCurt = curNext
		}

	}

	logger.Println(">>> Done")

}
