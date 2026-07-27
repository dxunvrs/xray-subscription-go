package service

import (
	"context"
	"fmt"
	"time"

	command "github.com/xtls/xray-core/app/proxyman/command"
	stats "github.com/xtls/xray-core/app/stats/command"
	protocol "github.com/xtls/xray-core/common/protocol"
	serial "github.com/xtls/xray-core/common/serial"
	vless "github.com/xtls/xray-core/proxy/vless"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type XrayGrpcClient struct {
	conn          *grpc.ClientConn
	handlerClient command.HandlerServiceClient
	statsClient   stats.StatsServiceClient
	inboundTag    string
}

func NewXrayGrpcClient(grpcHost string, grpcPort int, inboundTag string) (*XrayGrpcClient, error) {
	target := fmt.Sprintf("%s:%d", grpcHost, grpcPort)

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to xray grpc at %s: %w", target, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn.Connect()

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(ctx, state) {
			_ = conn.Close()
			return nil, fmt.Errorf("failed to reach xray grpc at %s within 3s (last state: %s)", target, state)
		}
	}

	return &XrayGrpcClient{
		conn:          conn,
		handlerClient: command.NewHandlerServiceClient(conn),
		statsClient:   stats.NewStatsServiceClient(conn),
		inboundTag:    inboundTag,
	}, nil
}

func (c *XrayGrpcClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *XrayGrpcClient) AddUser(email, uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	account := &vless.Account{
		Id:   uuid,
		Flow: "xtls-rprx-vision",
	}
	accountBytes, err := proto.Marshal(account)
	if err != nil {
		return err
	}

	typedAccount := &serial.TypedMessage{
		Type:  "xray.proxy.vless.Account",
		Value: accountBytes,
	}

	user := &protocol.User{
		Email:   email,
		Level:   0,
		Account: typedAccount,
	}

	addUserOp := &command.AddUserOperation{
		User: user,
	}
	addUserOpBytes, err := proto.Marshal(addUserOp)
	if err != nil {
		return err
	}

	typedOp := &serial.TypedMessage{
		Type:  "xray.app.proxyman.command.AddUserOperation",
		Value: addUserOpBytes,
	}

	req := &command.AlterInboundRequest{
		Tag:       c.inboundTag,
		Operation: typedOp,
	}

	_, err = c.handlerClient.AlterInbound(ctx, req)
	return err
}

func (c *XrayGrpcClient) RemoveUser(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	removeOp := &command.RemoveUserOperation{
		Email: email,
	}
	removeOpBytes, err := proto.Marshal(removeOp)
	if err != nil {
		return err
	}

	typedOp := &serial.TypedMessage{
		Type:  "xray.app.proxyman.command.RemoveUserOperation",
		Value: removeOpBytes,
	}

	req := &command.AlterInboundRequest{
		Tag:       c.inboundTag,
		Operation: typedOp,
	}

	_, err = c.handlerClient.AlterInbound(ctx, req)
	return err
}

func (c *XrayGrpcClient) GetUserUplink(email string) int64 {
	statName := fmt.Sprintf("user>>>%s>>>traffic>>>uplink", email)
	return c.getStatValue(statName)
}

func (c *XrayGrpcClient) GetUserDownlink(email string) int64 {
	statName := fmt.Sprintf("user>>>%s>>>traffic>>>downlink", email)
	return c.getStatValue(statName)
}

func (c *XrayGrpcClient) getStatValue(statName string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &stats.GetStatsRequest{
		Name:   statName,
		Reset_: false,
	}

	resp, err := c.statsClient.GetStats(ctx, req)
	if err != nil {
		return 0
	}

	if resp.GetStat() != nil {
		return resp.GetStat().GetValue()
	}

	return 0
}
