package restkafka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	. "kafkamongo/internal/model"
)

func JobsGetHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve body from http request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	//  Save Url message into URL struct
	var _Url URL
	err = json.Unmarshal(body, &_Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saveUrlToKafka(_Url)

	// Convert job struct into json
	jsonString, err := json.Marshal(_Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content-type http header
	w.Header().Set("content-type", "application/json")

	// Send back data as response
	w.Write(jsonString)
}

func JobsPostHandler(w http.ResponseWriter, r *http.Request) {

	// Retrieve body from http request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	//  Save data into Job struct
	var _job Job
	err = json.Unmarshal(body, &_job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saveJobToKafka(_job)

	// Convert job struct into json
	jsonString, err := json.Marshal(_job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content-type http header
	w.Header().Set("content-type", "application/json")

	// Send back data as response
	w.Write(jsonString)
}

func saveJobToKafka(job Job) {

	fmt.Println("Save job data to Kafka")

	jsonString, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}

	jobString := string(jsonString)
	fmt.Println(jobString)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	topic := "jobs-topic1"
	for _, word := range []string{string(jobString)} {
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny)},
			Value:          []byte(word),
		}, nil)
	}
}

func saveUrlToKafka(url URL) {

	fmt.Println("Start saving URL message to Kafka")

	jsonString, err := json.Marshal(url)
	if err != nil {
		panic(err)
	}

	urlString := string(jsonString)
	fmt.Println(urlString)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	topic := "jobs-topic2"
	for _, word := range []string{string(urlString)} {
		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny)},
			Value:          []byte(word),
		}, nil)
	}
}
