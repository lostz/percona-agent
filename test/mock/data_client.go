/*
   Copyright (c) 2014, Percona LLC and/or its affiliates. All rights reserved.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package mock

import (
	"code.google.com/p/go.net/websocket"
	"github.com/percona/cloud-protocol/proto"
	"reflect"
)

type DataClient struct {
	dataChan chan []byte      // agent -> API (test)
	respChan chan interface{} // agent <- API (test)
	// --
	ErrChan   chan error
	RecvError chan error
	// --
	conn            *websocket.Conn
	connectChan     chan bool
	testConnectChan chan bool
}

func NewDataClient(dataChan chan []byte, respChan chan interface{}) *DataClient {
	c := &DataClient{
		dataChan:    dataChan,
		respChan:    respChan,
		RecvError:   make(chan error),
		conn:        new(websocket.Conn),
		connectChan: make(chan bool, 1),
	}
	return c
}

func (c *DataClient) Connect() {
	c.ConnectOnce()
	return
}

func (c *DataClient) ConnectOnce() error {
	if c.testConnectChan != nil {
		// Wait for test to let user/agent connect.
		select {
		case c.testConnectChan <- true:
		default:
		}
		<-c.testConnectChan
	}
	return nil
}

func (c *DataClient) Disconnect() error {
	c.connectChan <- true
	return nil
}

func (c *DataClient) Start() {
}

func (c *DataClient) Stop() {
}

func (c *DataClient) SendChan() chan *proto.Reply {
	return nil
}

func (c *DataClient) RecvChan() chan *proto.Cmd {
	return nil
}

func (c *DataClient) Send(data interface{}, timeout uint) error {
	return nil
}

// First, agent calls this to send encoded proto.Data to API.
func (c *DataClient) SendBytes(data []byte) error {
	c.dataChan <- data
	return nil
}

// Second, agent calls this to recv response from API to previous send.
func (c *DataClient) Recv(resp interface{}, timeout uint) error {
	select {
	case r := <-c.respChan:
		respVal := reflect.ValueOf(resp).Elem()
		respVal.Set(reflect.ValueOf(r).Elem())
	case err := <-c.RecvError:
		return err
	}
	return nil
}

func (c *DataClient) ConnectChan() chan bool {
	return c.connectChan
}

func (c *DataClient) ErrorChan() chan error {
	return c.ErrChan
}

func (c *DataClient) Conn() *websocket.Conn {
	return c.conn
}

func (c *DataClient) SetConnectChan(connectChan chan bool) {
	c.testConnectChan = connectChan
}

func (c *DataClient) Status() map[string]string {
	return map[string]string{
		"data-client": "ok",
	}
}
