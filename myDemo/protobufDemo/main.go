package main

import (
	"fmt"
	"zinx/myDemo/protobufDemo/pb"

	"google.golang.org/protobuf/proto"
)

func main() {
	// 定义一个Person结构体对象
	person := &pb.Person{
		Name:  "lijiahao",
		Id:    0,
		Email: "18702748294@163.com",
		Phones: []*pb.Person_PhoneNumber{
			&pb.Person_PhoneNumber{
				Number: "15136588270",
				Type: pb.Person_MOBILE,
			},
			&pb.Person_PhoneNumber{
				Number: "190019393",
				Type: pb.Person_HOME,
			},
			&pb.Person_PhoneNumber{
				Number: "17182920303",
				Type: pb.Person_WORK,
			},
		},
	}

	// 将Person对象进行序列化
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("proto marshal error: ", err)
	}

	// 反序列化
	newData := &pb.Person{}
	err = proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Println("unmarshal error: ", err)
	}
	fmt.Println(newData)
}
