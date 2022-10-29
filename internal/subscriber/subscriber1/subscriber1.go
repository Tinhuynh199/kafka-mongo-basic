package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/mgo.v2"
	. "kafkamongo/internal/model"
	"kafkamongo/internal/app"
)

var mongoStore = MongoStore{}

func main() {
	// Create MongoDB session
	session := initializeMongo()
	mongoStore.Session = session

	receiveFromKafka()
}

func initializeMongo() (session *mgo.Session) {

	// Load Config
	config := app.GetConfig("../../../configs/config.yml")
	fmt.Println(config.Mongo.URI)

	info := &mgo.DialInfo{
		Addrs:    []string{config.Mongo.URI},
		Timeout:  60 * time.Second,
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

	fmt.Println("Start receiving from Kafka")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "group-id-1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{"jobs-topic1"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {
			fmt.Printf("Received from Kafka %s: %s\n", msg.TopicPartition, string(msg.Value))
			job := string(msg.Value)
			saveJobToMongo(job)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	consumer.Close()
}

func saveJobToMongo(jobString string) {

	// Load Config
	config := app.GetConfig("../../../configs/config.yml")

	// Save job data to mongo
	fmt.Println("Save data to MongoDB")
	col := mongoStore.Session.DB(config.Mongo.Database).C(config.Mongo.Collection)

	// Save data into Job struct
	var _job Job
	body := []byte(jobString)
	err := json.Unmarshal(body, &_job)
	if err != nil {
		panic(err)
	}

	// Insert job data into MongoDB
	errMongo := col.Insert(_job)
	if errMongo != nil {
		panic(errMongo)
	}

	fmt.Printf("Save data to MongoDB: %s\n", jobString)
}