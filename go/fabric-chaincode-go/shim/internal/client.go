// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"fmt"
	tls "github.com/tjfoc/gmtls"
	"time"

	peerpb "github.com/hyperledger/fabric-protos-go/peer"
	"google.golang.org/grpc"
	credentials "github.com/tjfoc/gmtls/gmcredentials"
	"google.golang.org/grpc/keepalive"
)

const (
	dialTimeout        = 10 * time.Second
	maxRecvMessageSize = 100 * 1024 * 1024 // 100 MiB
	maxSendMessageSize = 100 * 1024 * 1024 // 100 MiB
)

// NewClientConn ...
func NewClientConn(
	address string,
	tlsConf *tls.Config,
	kaOpts keepalive.ClientParameters,
) (*grpc.ClientConn, error) {

	dialOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kaOpts),
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMessageSize),
			grpc.MaxCallSendMsgSize(maxSendMessageSize),
		),
	}

	if tlsConf != nil {
		creds := credentials.NewTLS(tlsConf)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
		fmt.Printf("creds---,%v",creds)
		fmt.Printf("dialOpts---,%v",dialOpts)
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	fmt.Printf("tlsConf.RootCAs---,%v",tlsConf.RootCAs)
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, dialOpts...)
}

// NewRegisterClient ...
func NewRegisterClient(conn *grpc.ClientConn) (peerpb.ChaincodeSupport_RegisterClient, error) {
	return peerpb.NewChaincodeSupportClient(conn).Register(context.Background())
}
