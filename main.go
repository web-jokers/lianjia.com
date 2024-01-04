package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"
	"xietong.me/LianjiaSpider/common"
	"xietong.me/LianjiaSpider/spider"
)

// https://cd.lianjia.com/ershoufang/ 长沙市链家网爬虫
func main() {
	//初始化配置
	InitConfig()
	db := common.InitDB()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(20)
	defer db.Close()
	district := []string{"jinjiang", "qingyang", "wuhou", "gaoxin7", "chenghua", "jinniu", "tianfuxinqu",
		"gaoxinxi1", "shangliu", "wenjiang", "pidou", "longquanyi", "xindou", "tianfuxinqunanqu", "qingbaijiang",
		"doujiangyan", "pengzhou", "jianyang", "xinjin", "chongzhou1", "dayi", "jintang", "pujiang", "qionglai"}
	sellingQueue := make(chan interface{}, 10)
	soldQueue := make(chan interface{}, 10)
	go func() {
		for {
			sellingQueue <- struct{}{}
			go func() {
				defer func() {
					<-sellingQueue
				}()
				districtName := district[rand.Intn(len(district))]
				totalSellingPage := spider.GetSellingPageSpider(db, districtName)
				for page := 1; page < totalSellingPage; page++ {
					time.Sleep(time.Duration(page) * time.Millisecond)
					fmt.Println("start spider", page)
					spider.GetSellingInfoSpider(db, districtName, page)
				}
			}()
		}
	}()
	go func() {
		for {
			soldQueue <- struct{}{}
			go func() {
				defer func() {
					<-soldQueue
				}()
				districtName := district[rand.Intn(len(district))]
				totalSoldPage := spider.GetSoldPageSpider(db, districtName)
				for page := 1; page < totalSoldPage; page++ {
					time.Sleep(time.Duration(page*20) * time.Millisecond)
					fmt.Println("start spider", page)
					spider.GetSoldInfoSpider(db, districtName, page)
				}
			}()
		}
	}()
	select {}

	////通过循环来爬取不同地区，同时获取不同地区的总分页数来爬取不同页面的数据
	//var wgSelling sync.WaitGroup
	//var wgSold sync.WaitGroup
	//for _, districtName := range district {
	//	totalSellingPage := spider.GetSellingPageSpider(db, districtName)
	//	for page := 1; page < totalSellingPage; page++ {
	//		wgSelling.Add(1)
	//		time.Sleep(time.Duration(page) * time.Millisecond)
	//		go func(page int) {
	//			fmt.Println("start spider", page)
	//			defer wgSelling.Done()
	//			spider.GetSellingInfoSpider(db, districtName, page)
	//		}(page)
	//	}
	//}
	//wgSelling.Wait()
	//
	//for _, districtName := range district {
	//	totalSoldPage := spider.GetSoldPageSpider(db, districtName)
	//	for page := 1; page < totalSoldPage; page++ {
	//		wgSold.Add(1)
	//		time.Sleep(time.Duration(page*20) * time.Millisecond)
	//		go func(page int) {
	//			fmt.Println("start spider", page)
	//			defer wgSold.Done()
	//			spider.GetSoldInfoSpider(db, districtName, page)
	//		}(page)
	//	}
	//}
	//wgSold.Wait()
}

// 初始化配置函数
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
