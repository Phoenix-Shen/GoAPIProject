package main

type TimePoint struct {
	startTime int64 `bson:"startTime"`
	endTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	jobName string `bson:"jobName"`
	command string `bson:"command"`
	err     string `bson:"err"`
	content string `bson:"content"`
	tp      TimePoint
}
