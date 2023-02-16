package service

import (
	"context"
	"fmt"
	"github.com/theone-daxia/chat-demo/config"
	"github.com/theone-daxia/chat-demo/model/ws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

func InsertMsg(database string, id string, content string, read uint, expire int64) error {
	collection := config.MongoDBClient.Database(database).Collection(id)
	doc := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err := collection.InsertOne(context.TODO(), doc)
	return err
}

func FindMany(database string, id string, sendId string, pageSize int) []ws.Result {
	listMe, _ := findManyHandler(database, id, pageSize)
	listYou, _ := findManyHandler(database, sendId, pageSize)
	res := AppendAndSort(listMe, listYou)
	return res
}

func findManyHandler(database string, table string, pageSize int) ([]ws.Trainer, error) {
	var list []ws.Trainer
	collection := config.MongoDBClient.Database(database).Collection(table)
	cursor, _ := collection.Find(context.TODO(), options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))
	err := cursor.All(context.TODO(), &list)
	return list, err
}

func AppendAndSort(listMe, listYou []ws.Trainer) (res []ws.Result) {
	// 合并
	res = AppendList(listMe, "me")
	res = append(res, AppendList(listYou, "you")...)

	// 排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].StartTime < res[j].StartTime
	})
	return
}

func AppendList(list []ws.Trainer, from string) (res []ws.Result) {
	for _, v := range list {
		result := ws.Result{
			StartTime: v.StartTime,
			Msg:       fmt.Sprintf("%v", v),
			From:      from,
		}
		res = append(res, result)
	}
	return
}
