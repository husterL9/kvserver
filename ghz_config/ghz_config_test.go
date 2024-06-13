package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	pb "github.com/husterL9/kvserver/api/protobuf"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
)

func dataGetFunc(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	msg := &pb.GetRequest{
		Key: "key" + (cd.WorkerID),
	}
	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return binData
}
func dataSetFunc(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	meta := &pb.MetaData{Type: pb.DataType_KV_OBJ, Location: "location-{{.RequestNumber}}"}
	msg := &pb.SetRequest{
		Key:   "key" + (cd.WorkerID),
		Value: []byte(cd.UUID),
		Meta:  meta,
	}
	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return binData
}
func dataGetFileFunc(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	msg := &pb.GetRequest{
		Key: "file" + (cd.WorkerID),
	}
	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return binData

}
func dataSetFileFunc(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	meta := &pb.MetaData{
		Type:     pb.DataType_FILE,
		Location: fmt.Sprintf("/home/ljw/SE8/kvserver/internal/kvstore/fakeRoot/ghz_test/%s.txt", cd.WorkerID),
	}
	msg := &pb.SetRequest{
		Key:   "file" + (cd.WorkerID),
		Value: []byte(cd.UUID),
		Meta:  meta,
	}
	binData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return binData

}
func TestGet(t *testing.T) {
	const REQUESTS_COUNT = 10000
	report, err := runner.Run(
		// 基本配置 call host proto文件 data
		"kvstore.KVStoreService/Get",
		"localhost:50051",
		runner.WithProtoFile("/home/ljw/SE8/kvserver/api/protobuf/kvstore.proto", []string{}),
		runner.WithBinaryDataFunc(dataGetFunc),
		runner.WithInsecure(true),
		runner.WithTotalRequests(REQUESTS_COUNT),
		// 并发参数
		runner.WithConcurrencySchedule(runner.ScheduleLine),
		runner.WithConcurrencyStep(10),
		runner.WithConcurrencyStart(5),
		runner.WithConcurrencyEnd(100),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径 test_get_concurrence_summary.txt test_get.html test_get_concurrence.html
	file, err := os.Create("test_get_concurrence3.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}
func TestSet(t *testing.T) {
	const REQUESTS_COUNT = 10000
	report, err := runner.Run(
		// 基本配置 call host proto文件 data
		"kvstore.KVStoreService/Set",
		"localhost:50051",
		runner.WithProtoFile("/home/ljw/SE8/kvserver/api/protobuf/kvstore.proto", []string{}),
		runner.WithBinaryDataFunc(dataSetFunc),
		runner.WithInsecure(true),
		runner.WithTotalRequests(REQUESTS_COUNT),
		// 并发参数
		runner.WithConcurrencySchedule(runner.ScheduleLine),
		runner.WithConcurrencyStep(10),
		runner.WithConcurrencyStart(5),
		runner.WithConcurrencyEnd(100),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径 test_set_concurrence_summary.txt
	file, err := os.Create("test_set_concurrence2.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}
func TestGetFile(t *testing.T) {
	const REQUESTS_COUNT = 10000
	report, err := runner.Run(
		// 基本配置 call host proto文件 data
		"kvstore.KVStoreService/Get",
		"localhost:50051",
		runner.WithProtoFile("/home/ljw/SE8/kvserver/api/protobuf/kvstore.proto", []string{}),
		runner.WithBinaryDataFunc(dataGetFileFunc),
		runner.WithInsecure(true),
		runner.WithTotalRequests(REQUESTS_COUNT),
		// 并发参数
		runner.WithConcurrencySchedule(runner.ScheduleLine),
		runner.WithConcurrencyStep(10),
		runner.WithConcurrencyStart(5),
		runner.WithConcurrencyEnd(100),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	file, err := os.Create("test_get_file.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}
func TestSetFile(t *testing.T) {
	const REQUESTS_COUNT = 10000
	report, err := runner.Run(
		// 基本配置 call host proto文件 data
		"kvstore.KVStoreService/Set",
		"localhost:50051",
		runner.WithProtoFile("/home/ljw/SE8/kvserver/api/protobuf/kvstore.proto", []string{}),
		runner.WithBinaryDataFunc(dataSetFileFunc),
		runner.WithInsecure(true),
		runner.WithTotalRequests(REQUESTS_COUNT),
		// 并发参数
		runner.WithConcurrencySchedule(runner.ScheduleLine),
		runner.WithConcurrencyStep(10),
		runner.WithConcurrencyStart(5),
		runner.WithConcurrencyEnd(100),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 指定输出路径
	file, err := os.Create("test_set_file.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}
	// 指定输出格式
	_ = rp.Print("html")
}

// 官方文档 https://ghz.sh/docs/intro.html
func main() {

}
