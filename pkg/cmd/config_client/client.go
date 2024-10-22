/*******************************************************************************
*  MIT License
*
*  Copyright (c) 2024 Ji-Young Park(jiyoung.park.dev@gmail.com)
*
*  Permission is hereby granted, free of charge, to any person obtaining a copy
*  of this software and associated documentation files (the "Software"), to deal
*  in the Software without restriction, including without limitation the rights
*  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
*  copies of the Software, and to permit persons to whom the Software is
*  furnished to do so, subject to the following conditions:
*
*      The above copyright notice and this permission notice shall be included in all
*      copies or substantial portions of the Software.
*
*      THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
*      IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
*      FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
*      AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
*      LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
*      OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
*      SOFTWARE.
*******************************************************************************/

package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jparklab/xds-test-server/pkg/generated/config"
	"github.com/jparklab/xds-test-server/pkg/utils/log"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address       = flag.String("address", "localhost:5001", "address to connect to")
	action        = flag.String("action", "add", "action to send")
	numServices   = flag.Int("num-services", 1, "number of services to add")
	serviceOffset = flag.Int("service-offset", 0, "offset to add to service index")
	value         = flag.String("value", "", "value to set")

	actions = []string{"add", "set_window", "use_dynamic", "set_rds_initial_fetch_timeout"}
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial(%s) failed: %v", *address, err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configClient := config.NewConfigServiceClient(conn)
	for i := 0; i < *numServices; i++ {
		serviceIdx := i + *serviceOffset
		serviceName := fmt.Sprintf("service-%d", serviceIdx)

		switch *action {
		case "add":
			_, err = configClient.SendConfigCommand(
				ctx,
				&config.Command{
					Type: config.CommandType_COMMAND_ADD_SERVICE,
					Command: &config.Command_AddService{
						AddService: &config.AddServiceCommand{
							ServiceName: serviceName,
						},
					},
				},
			)
			if err != nil {
				log.With("error", err).Errorf("Failed to send a command")
			}

		case "use_dynamic":
			use_dynamic := true
			if *value == "false" {
				use_dynamic = false
			} else if *value != "true" {
				log.Fatalf("Invalid value: %v(must be true or false)", *value)
			}

			_, err = configClient.SendConfigCommand(
				ctx,
				&config.Command{
					Type: config.CommandType_COMMAND_SET_OPTION,
					Command: &config.Command_SetOption{
						SetOption: &config.SetOptionCommand{
							Key: "use_dynamic",
							Value: &config.SetOptionCommand_BoolValue{
								BoolValue: use_dynamic,
							},
						},
					},
				},
			)
			if err != nil {
				log.With("error", err).Errorf("Failed to send a command")
			}
		case "set_rds_initial_fetch_timeout":
			size, err := strconv.Atoi(*value)
			if err != nil {
				log.Fatalf("Invalid value: %v", *value)
			}

			_, err = configClient.SendConfigCommand(
				ctx,
				&config.Command{
					Type: config.CommandType_COMMAND_SET_OPTION,
					Command: &config.Command_SetOption{
						SetOption: &config.SetOptionCommand{
							Key: "rds_initial_fetch_timeout",
							Value: &config.SetOptionCommand_Int64Value{
								Int64Value: int64(size),
							},
						},
					},
				},
			)
			if err != nil {
				log.With("error", err).Errorf("Failed to send a command")
			}

		case "set_window":
			size, err := strconv.Atoi(*value)
			if err != nil {
				log.Fatalf("Invalid value: %v", *value)
			}

			_, err = configClient.SendConfigCommand(
				ctx,
				&config.Command{
					Type: config.CommandType_COMMAND_SET_OPTION,
					Command: &config.Command_SetOption{
						SetOption: &config.SetOptionCommand{
							Key: "window_size",
							Value: &config.SetOptionCommand_Int64Value{
								Int64Value: int64(size),
							},
						},
					},
				},
			)
			if err != nil {
				log.With("error", err).Errorf("Failed to send a command")
			}

		default:
			log.Fatalf("Unknown action: %s, must be one of %v", *action, actions)
		}
	}
}
