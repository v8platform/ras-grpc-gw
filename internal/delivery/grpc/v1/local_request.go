package v1

import (
	"bytes"
	"context"
	"fmt"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type CallOption func()

type Client struct {

}

func (c *Client) InvokeEndpointRequest(ctx context.Context, req interface{}, reply interface{}, opts ...CallOption) error {

	reqMessage, err := x.NewMessage(message)
	if err != nil {
		return nil, err
	}

	respMessage, err := x.client.EndpointMessage(ctx, reqMessage)
	if err != nil {
		return nil, err
	}

	respProtoMessage, err := anypb.UnmarshalNew(req.GetRespond(), proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := respProtoMessage.(*emptypb.Empty); ok {
		if err := x.UnpackMessage(respMessage, nil); err != nil {
			return nil, err
		}
		return anypb.New(respProtoMessage)
	}

	messageParser, ok := respProtoMessage.(v1.EndpointMessageParser)
	if !ok {
		return nil, fmt.Errorf("not parser interface")
	}
	if err := x.UnpackMessage(respMessage, messageParser); err != nil {
		return nil, err
	}
	return anypb.New(respProtoMessage)


}

type EndpointContext interface {
	GetVersion() int32
	GetId() int32
	GetService() string
	GetFormat() int32
}

func invokeEndpointRequest(ctx context.Context, conn net.Conn, endpointContext EndpointContext, req protocolv1.EndpointMessageFormatter, reply protocolv1.EndpointMessageParser, opts ...CallOption) error {

	message, err := protocolv1.NewEndpointMessage(endpointContext, req)
	if err != nil {
		return err
	}

	var resp  protocolv1.EndpointMessage{}


}


func invokeRequest(ctx context.Context, conn net.Conn, req protocolv1.PacketMessageFormatter, reply protocolv1.PacketMessageParser , opts ...CallOption) error {

	cs, err := newClientStream(ctx, conn, opts...)
	if err != nil {
		return err
	}
	if err := cs.SendMsg(req); err != nil {
		return err
	}

	if reply == nil {
		return nil
	}

	return cs.RecvMsg(reply)
}

var unaryStreamDesc = &streamPacket{}

type streamPacket struct {
	ctx context.Context
	cc net.Conn
}

func (s streamPacket) Context() context.Context {
	return s.ctx
}

func newClientStream(ctx context.Context, conn net.Conn, opts... CallOption) (StreamPacket, error) {
	return streamPacket{
		ctx: ctx,
		cc: conn,
	}, nil
}

type StreamPacket interface {
	Context() context.Context
	SendMsg(m protocolv1.PacketMessageFormatter) error
	RecvMsg(m protocolv1.PacketMessageParser) error
}

func (s streamPacket) SendMsg(m protocolv1.PacketMessageFormatter) error {

	var packet protocolv1.Packet
	buf := &bytes.Buffer{}
	if err := m.Formatter(buf, 0); err != nil {
		return err
	}
	packet.Type = m.GetPacketType()
	packet.Data = buf.Bytes()
	packet.Size = int32(len(packet.Data))

	_, err := packet.WriteTo(s.cc)
	if err != nil {
		return err
	}
	return nil

}

func (s streamPacket) RecvMsg(m protocolv1.PacketMessageParser) error {

	var packet protocolv1.Packet

	if err := packet.Parse(s.cc, 0); err != nil {
		return err
	}

	if err := packet.Unpack(m); err != nil {
		return err
	}
	return nil

}

func localRequest_0(ctx context.Context, req *EndpointRequest) (*anypb.Any, error) {


	message, err := anypb.UnmarshalNew(req.GetRequest(), proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	reqMessage, err := x.NewMessage(message)
	if err != nil {
		return nil, err
	}

	respMessage, err := x.client.EndpointMessage(ctx, reqMessage)
	if err != nil {
		return nil, err
	}

	respProtoMessage, err := anypb.UnmarshalNew(req.GetRespond(), proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := respProtoMessage.(*emptypb.Empty); ok {
		if err := x.UnpackMessage(respMessage, nil); err != nil {
			return nil, err
		}
		return anypb.New(respProtoMessage)
	}

	messageParser, ok := respProtoMessage.(v1.EndpointMessageParser)
	if !ok {
		return nil, fmt.Errorf("not parser interface")
	}
	if err := x.UnpackMessage(respMessage, messageParser); err != nil {
		return nil, err
	}
	return anypb.New(respProtoMessage)
}

