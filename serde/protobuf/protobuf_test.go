/**
 * Copyright 2022 Confluent Inc.
 * Copyright 2023 Jerome Bidault (jeromedbtdev@gmail.com).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file has been modified by Jerome Bidault (jeromebdtdev@gmail.com) to include additional functionality.
 */

package protobuf

import (
	"errors"
	// "fmt"
	"fmt"
	"testing"

	schemaregistry "github.com/djedjethai/gokfk-regent"
	"github.com/djedjethai/gokfk-regent/serde"
	"github.com/djedjethai/gokfk-regent/test"
	"github.com/djedjethai/gokfk-regent/test/proto/recordname"
	trn "github.com/djedjethai/gokfk-regent/test/proto/topicrecordname"
	"google.golang.org/protobuf/proto"
)

func TestProtobufSerdeWithSimple(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	obj := test.Author{
		Name:  "Kafka",
		Id:    123,
		Works: []string{"The Castle", "The Trial"},
	}
	bytes, err := ser.Serialize("topic1", &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	newobj, err := deser.Deserialize("topic1", bytes)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), obj.ProtoReflect()))
}

func TestProtobufSerdeWithSecondMessage(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	obj := test.Pizza{
		Size:     "Extra extra large",
		Toppings: []string{"anchovies", "mushrooms"},
	}
	bytes, err := ser.Serialize("topic1", &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	newobj, err := deser.Deserialize("topic1", bytes)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), obj.ProtoReflect()))
}

func TestProtobufSerdeWithNestedMessage(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	obj := test.NestedMessage_InnerMessage{
		Id: "inner",
	}
	bytes, err := ser.Serialize("topic1", &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	newobj, err := deser.Deserialize("topic1", bytes)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), obj.ProtoReflect()))
}

func TestProtobufSerdeWithReference(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	msg := test.TestMessage{
		TestString:   "hi",
		TestBool:     true,
		TestBytes:    []byte{1, 2},
		TestDouble:   1.23,
		TestFloat:    3.45,
		TestFixed32:  67,
		TestFixed64:  89,
		TestInt32:    100,
		TestInt64:    200,
		TestSfixed32: 300,
		TestSfixed64: 400,
		TestSint32:   500,
		TestSint64:   600,
		TestUint32:   700,
		TestUint64:   800,
	}
	obj := test.DependencyMessage{
		IsActive:     true,
		TestMesssage: &msg,
	}

	bytes, err := ser.Serialize("topic1", &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	newobj, err := deser.Deserialize("topic1", bytes)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), obj.ProtoReflect()))
}

func TestProtobufSerdeWithCycle(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	inner := test.LinkedList{
		Value: 100,
	}
	obj := test.LinkedList{
		Value: 1,
		Next:  &inner,
	}
	bytes, err := ser.Serialize("topic1", &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	newobj, err := deser.Deserialize("topic1", bytes)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), obj.ProtoReflect()))
}

// Test strategies
func TestProtobufSerdeDeserializeInto(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	obj := test.Pizza{
		Size:     "Extra extra large",
		Toppings: []string{"anchovies", "mushrooms"},
	}

	topic := "topic"

	bytesInner, err := ser.Serialize(topic, &obj)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(obj.ProtoReflect().Type())

	innerReceiver := &test.LinkedList{}

	err = deser.DeserializeInto(topic, bytesInner, innerReceiver)
	serde.MaybeFail("deserializeRecordNameValidSchema", serde.Expect(err.Error(), "recipient proto object differs from incoming events"))
}

// -------------------- RecordName ------------------------------------

const (
	linkedList      = "recordname.LinkedList"
	linkedListValue = "recordname.LinkedList-value"
	pizza           = "recordname.Pizza"
	pizzaValue      = "recordname.Pizza-value"
	invalidSchema   = "invalidSchema"
)

var (
	recLinked = recordname.LinkedList{
		Value: 100,
	}

	recPiz = recordname.Pizza{
		Size:     "Extra extra large",
		Toppings: []string{"anchovies", "mushrooms"},
	}
)

func TestProtobufSerdeDeserializeRecordName(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeRecordName(&recLinked, linkedList)
	serde.MaybeFail("serialization", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	fmt.Println("")

	newobj, err := deser.DeserializeRecordName(bytesInner)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), recLinked.ProtoReflect()))

	newobj, err = deser.DeserializeRecordName(bytesObj)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(proto.Message).ProtoReflect(), recPiz.ProtoReflect()))
}

func RegisterMessageFactory() func([]string, string) (interface{}, error) {
	return func(subject []string, name string) (interface{}, error) {
		switch name {
		case linkedList:
			return &recordname.LinkedList{}, nil
		case pizza:
			return &recordname.Pizza{}, nil
		}
		return nil, errors.New("No matching receiver")
	}
}

func RegisterMessageFactoryNoReceiver() func([]string, string) (interface{}, error) {
	return func(subject []string, name string) (interface{}, error) {
		return nil, errors.New("No matching receiver")
	}
}

func RegisterMessageFactoryInvalidReceiver() func([]string, string) (interface{}, error) {
	return func(subject []string, name string) (interface{}, error) {
		switch name {
		case pizza:
			return &recordname.LinkedList{}, nil
		case linkedList:
			return "", nil
		}
		return nil, errors.New("No matching receiver")
	}
}

func TestProtobufSerdeDeserializeRecordNameWithHandler(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeRecordName(&recLinked, linkedList)
	serde.MaybeFail("serialization", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client
	deser.MessageFactory = RegisterMessageFactory()

	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	newobj, err := deser.DeserializeRecordName(bytesInner)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.LinkedList).Value, recLinked.Value))

	newobj, err = deser.DeserializeRecordName(bytesObj)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Size, recPiz.Size))
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Toppings[0], recPiz.Toppings[0]))
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Toppings[1], recPiz.Toppings[1]))
}

func TestProtobufSerdeDeserializeRecordNameWithHandlerNoReceiver(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client
	// register invalid receiver
	deser.MessageFactory = RegisterMessageFactoryNoReceiver()

	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	newobj, err := deser.DeserializeRecordName(bytesObj)

	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(err.Error(), "No matching receiver"))
	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(newobj, nil))
}

func TestProtobufSerdeDeserializeRecordNameWithInvalidSchema(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeRecordName(&recLinked)
	serde.MaybeFail("serialization", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client
	// register invalid schema
	deser.MessageFactory = RegisterMessageFactoryInvalidReceiver()

	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	_, err = deser.DeserializeRecordName(bytesObj)
	serde.MaybeFail("deserialization", err)

	_, err = deser.DeserializeRecordName(bytesInner)
	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))
}

func TestProtobufSerdeDeserializeIntoRecordName(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeRecordName(&recLinked, linkedList)
	serde.MaybeFail("serialization", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz)
	serde.MaybeFail("serialization", err)

	var receivers = make(map[string]interface{})
	receivers[linkedListValue] = &recordname.LinkedList{}
	receivers[pizzaValue] = &recordname.Pizza{}

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	err = deser.DeserializeIntoRecordName(receivers, bytesInner)
	serde.MaybeFail("deserialization", err, serde.Expect(int(receivers[linkedListValue].(*recordname.LinkedList).Value), 100))

	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
	serde.MaybeFail("deserialization", err, serde.Expect(receivers[pizzaValue].(*recordname.Pizza).Toppings[0], recPiz.Toppings[0]))
	serde.MaybeFail("deserialization", err, serde.Expect(receivers[pizzaValue].(*recordname.Pizza).Toppings[1], recPiz.Toppings[1]))
}

func TestProtobufSerdeDeserializeIntoRecordNameWithInvalidSchema(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz)
	serde.MaybeFail("serialization", err)

	var receivers = make(map[string]interface{})
	receivers[invalidSchema] = &test.Pizza{}

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "unfound subject declaration"))
	serde.MaybeFail("deserialization", serde.Expect(receivers[invalidSchema].(*test.Pizza).Size, ""))
}

func TestProtobufSerdeDeserializeIntoRecordNameWithInvalidReceiver(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
	serde.MaybeFail("serialization", err)

	bytesInner, err := ser.SerializeRecordName(&recLinked)
	serde.MaybeFail("serialization", err)

	aut := recordname.Author{
		Name: "aut",
	}
	bytesAut, err := ser.SerializeRecordName(&aut, "recordname.Author")
	serde.MaybeFail("serialization", err)

	var receivers = make(map[string]interface{})
	receivers[pizzaValue] = &test.LinkedList{}
	receivers[linkedListValue] = ""

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())

	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
	// serde.MaybeFail("deserialization", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))
	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "recipient proto object differs from incoming events"))

	err = deser.DeserializeIntoRecordName(receivers, bytesInner)
	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))

	err = deser.DeserializeIntoRecordName(receivers, bytesAut)
	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "unfound subject declaration"))
}

func TestProtobufSerdeSubjectMismatchPayload(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	_, err = ser.SerializeRecordName(&recPiz, "test.Pizza")
	fmt.Println(err)
	serde.MaybeFail("serialization", serde.Expect(err.Error(), "the payload's fullyQualifiedName: 'recordname.Pizza' does not match the subject: 'test.Pizza'"))
}

// ----------------------- TopicRecordName ---------------------------

// TODO inplement tests for topicRecordNameStrategy....
const (
	topic                         = "topic"
	second                        = "second"
	topicLinkedList               = "topic-recordname.LinkedList"
	topicLinkedListValue          = "topic-recordname.LinkedList-value"
	secondLinkedList              = "second-recordname.LinkedList"
	secondLinkedListValue         = "second-recordname.LinkedList-value"
	topicPizza                    = "topic-recordname.Pizza"
	topicPizzaValue               = "topic-recordname.Pizza-value"
	secondPizza                   = "second-recordname.Pizza"
	secondPizzaValue              = "second-recordname.Pizza-value"
	topicLinkedListNopackage      = "topic-LinkedList"
	topicLinkedListNopackageValue = "topic-LinkedList-value"
)

var (
	topLinked = trn.LinkedList{
		Value: 100,
	}

	topPiz = trn.Pizza{
		Size:     "Extra extra large",
		Toppings: []string{"anchovies", "mushrooms"},
	}
)

func TestProtobufSerdeDeserializeTopicRecordNameWithHandler(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeTopicRecordName(topic, &recLinked, topicLinkedList)
	serde.MaybeFail("serialization", err)

	bytesInner2, err := ser.SerializeTopicRecordName(second, &recLinked, secondLinkedList)
	serde.MaybeFail("serialization", err)

	bytesObj, err := ser.SerializeTopicRecordName(topic, &recPiz, topicPizza)
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client
	// deser.MessageFactory = RegisterMessageFactory()

	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())

	newobj, err := deser.DeserializeTopicRecordName(topic, bytesInner)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.LinkedList).Value, recLinked.Value))
	newobj, err = deser.DeserializeTopicRecordName(second, bytesInner2)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.LinkedList).Value, recLinked.Value))

	// wrong topic, works anyway
	newobj, err = deser.DeserializeTopicRecordName("unknown", bytesInner2)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.LinkedList).Value, recLinked.Value))

	newobj, err = deser.DeserializeTopicRecordName(topic, bytesObj)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Size, recPiz.Size))
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Toppings[0], recPiz.Toppings[0]))
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*recordname.Pizza).Toppings[1], recPiz.Toppings[1]))
}

// Protobuf have no packageName, gokfk-regent will get the goPackageName: protorecordname
func TestProtobufSerdeDeserializeTopicRecordNameWithHandlerAndNoPackagename(t *testing.T) {
	serde.MaybeFail = serde.InitFailFunc(t)
	var err error
	conf := schemaregistry.NewConfig("mock://")

	client, err := schemaregistry.NewClient(conf)
	serde.MaybeFail("Schema Registry configuration", err)

	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
	serde.MaybeFail("Serializer configuration", err)

	bytesInner, err := ser.SerializeTopicRecordName(topic, &topLinked, "topic-protorecordname.LinkedList")
	serde.MaybeFail("serialization", err)

	bytesInner2, err := ser.SerializeTopicRecordName(second, &topLinked, "second-protorecordname.LinkedList")
	serde.MaybeFail("serialization", err)

	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())

	serde.MaybeFail("Deserializer configuration", err)
	deser.Client = ser.Client

	deser.ProtoRegistry.RegisterMessage(topLinked.ProtoReflect().Type())

	newobj, err := deser.DeserializeTopicRecordName(topic, bytesInner)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*trn.LinkedList).Value, topLinked.Value))
	newobj, err = deser.DeserializeTopicRecordName(second, bytesInner2)
	serde.MaybeFail("deserialization", err, serde.Expect(newobj.(*trn.LinkedList).Value, topLinked.Value))
}

// func TestProtobufSerdeDeserializeRecordNameWithHandlerNoReceiver(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
// 	serde.MaybeFail("serialization", err)
//
// 	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
//
// 	serde.MaybeFail("Deserializer configuration", err)
// 	deser.Client = ser.Client
// 	// register invalid receiver
// 	deser.MessageFactory = RegisterMessageFactoryNoReceiver()
//
// 	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
//
// 	newobj, err := deser.DeserializeRecordName(bytesObj)
//
// 	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(err.Error(), "No matching receiver"))
// 	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(newobj, nil))
// }
//
// func TestProtobufSerdeDeserializeRecordNameWithInvalidSchema(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	bytesInner, err := ser.SerializeRecordName(&recLinked)
// 	serde.MaybeFail("serialization", err)
//
// 	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
// 	serde.MaybeFail("serialization", err)
//
// 	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
//
// 	serde.MaybeFail("Deserializer configuration", err)
// 	deser.Client = ser.Client
// 	// register invalid schema
// 	deser.MessageFactory = RegisterMessageFactoryInvalidReceiver()
//
// 	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
// 	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
//
// 	_, err = deser.DeserializeRecordName(bytesObj)
// 	serde.MaybeFail("deserialization", err)
//
// 	_, err = deser.DeserializeRecordName(bytesInner)
// 	serde.MaybeFail("deserializeInvalidReceiver", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))
// }
//
// func TestProtobufSerdeDeserializeIntoRecordName(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	bytesInner, err := ser.SerializeRecordName(&recLinked, linkedList)
// 	serde.MaybeFail("serialization", err)
//
// 	bytesObj, err := ser.SerializeRecordName(&recPiz)
// 	serde.MaybeFail("serialization", err)
//
// 	var receivers = make(map[string]interface{})
// 	receivers[linkedListValue] = &recordname.LinkedList{}
// 	receivers[pizzaValue] = &recordname.Pizza{}
//
// 	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
//
// 	serde.MaybeFail("Deserializer configuration", err)
// 	deser.Client = ser.Client
//
// 	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
// 	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesInner)
// 	serde.MaybeFail("deserialization", err, serde.Expect(int(receivers[linkedListValue].(*recordname.LinkedList).Value), 100))
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
// 	serde.MaybeFail("deserialization", err, serde.Expect(receivers[pizzaValue].(*recordname.Pizza).Toppings[0], recPiz.Toppings[0]))
// 	serde.MaybeFail("deserialization", err, serde.Expect(receivers[pizzaValue].(*recordname.Pizza).Toppings[1], recPiz.Toppings[1]))
// }
//
// func TestProtobufSerdeDeserializeIntoRecordNameWithInvalidSchema(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	bytesObj, err := ser.SerializeRecordName(&recPiz)
// 	serde.MaybeFail("serialization", err)
//
// 	var receivers = make(map[string]interface{})
// 	receivers[invalidSchema] = &test.Pizza{}
//
// 	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
//
// 	serde.MaybeFail("Deserializer configuration", err)
// 	deser.Client = ser.Client
//
// 	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
// 	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "unfound subject declaration"))
// 	serde.MaybeFail("deserialization", serde.Expect(receivers[invalidSchema].(*test.Pizza).Size, ""))
// }
//
// func TestProtobufSerdeDeserializeIntoRecordNameWithInvalidReceiver(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	bytesObj, err := ser.SerializeRecordName(&recPiz, pizza)
// 	serde.MaybeFail("serialization", err)
//
// 	bytesInner, err := ser.SerializeRecordName(&recLinked)
// 	serde.MaybeFail("serialization", err)
//
// 	aut := recordname.Author{
// 		Name: "aut",
// 	}
// 	bytesAut, err := ser.SerializeRecordName(&aut, "recordname.Author")
// 	serde.MaybeFail("serialization", err)
//
// 	var receivers = make(map[string]interface{})
// 	receivers[pizzaValue] = &test.LinkedList{}
// 	receivers[linkedListValue] = ""
//
// 	deser, err := NewDeserializer(client, serde.ValueSerde, NewDeserializerConfig())
//
// 	serde.MaybeFail("Deserializer configuration", err)
// 	deser.Client = ser.Client
//
// 	deser.ProtoRegistry.RegisterMessage(recPiz.ProtoReflect().Type())
// 	deser.ProtoRegistry.RegisterMessage(recLinked.ProtoReflect().Type())
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesObj)
// 	// serde.MaybeFail("deserialization", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))
// 	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "recipient proto object differs from incoming events"))
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesInner)
// 	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "deserialization target must be a protobuf message"))
//
// 	err = deser.DeserializeIntoRecordName(receivers, bytesAut)
// 	serde.MaybeFail("deserialization", serde.Expect(err.Error(), "unfound subject declaration"))
// }
//
// func TestProtobufSerdeSubjectMismatchPayload(t *testing.T) {
// 	serde.MaybeFail = serde.InitFailFunc(t)
// 	var err error
// 	conf := schemaregistry.NewConfig("mock://")
//
// 	client, err := schemaregistry.NewClient(conf)
// 	serde.MaybeFail("Schema Registry configuration", err)
//
// 	ser, err := NewSerializer(client, serde.ValueSerde, NewSerializerConfig())
// 	serde.MaybeFail("Serializer configuration", err)
//
// 	_, err = ser.SerializeRecordName(&recPiz, "test.Pizza")
// 	fmt.Println(err)
// 	serde.MaybeFail("serialization", serde.Expect(err.Error(), "the payload's fullyQualifiedName: 'recordname.Pizza' does not match the subject: 'test.Pizza'"))
// }
//
