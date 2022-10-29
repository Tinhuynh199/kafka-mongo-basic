package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	. "kafkamongo/internal/app"
	"gopkg.in/mgo.v2"
	. "kafkamongo/internal/model"
)

var mongoStore = MongoStore{}

func main() {

	// Create MongoDB session
	session := initializeMongo()
	mongoStore.Session = session

	receiveFromKafka()
}


func  initializeMongo() (session *mgo.Session){

	// Load Config
	config := GetConfig("../../../configs/config.yml")

	info := &mgo.DialInfo{
		Addrs: []string{config.Mongo.URI},
		Timeout: 60 * time.Second,
		Database: config.Mongo.Database,
		Username: config.Mongo.Username,
		Password: config.Mongo.Password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return
}

func receiveFromKafka() {

	fmt.Println("Start receiving data from Kafka")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "group-id-1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{"jobs-topic2"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received data from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
			urlMess := string(msg.Value)
			getListJobFromMongo(urlMess)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	consumer.Close()
}

func getListJobFromMongo(urlMess string) {

	// Load Config
	config := GetConfig("../../../configs/config.yml")

	// Save job data to mongo
	fmt.Println("Get list job data from MongoDB")
	col := mongoStore.Session.DB(config.Mongo.Database).C(config.Mongo.Collection)

	// Save url message into URL struct
	var _UrlString URL
	body := []byte(urlMess)
	err := json.Unmarshal(body, &_UrlString)
	if err != nil {
		panic(err)
	}

	if _UrlString.URL == "localhost:9090/jobs/list" {
		// Starting get list job data from MongoDB
		var jobs []Job
		err := col.Find(nil).All(&jobs)
		if err != nil {
			panic(err)
		}

		fmt.Println("Here's the list jobs in DB")
		for _, job := range jobs {
			fmt.Printf("Title: %s Description: %s Company: %s Salary: %s\n", job.Title, job.Description, job.Company, job.Salary)
		} 
	} else {
		fmt.Println("The request is invalid !!!")
	}
}

