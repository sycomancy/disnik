package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/client"
	"github.com/sycomancy/glasnik/types"
)

type RequestProcessor struct {
	client      *client.Client
	listenAddr  string
	serviceAddr string
	callback    string
}

func NewRequestProcessor(listenAddr string, path string, serviceAddr string) *RequestProcessor {
	return &RequestProcessor{
		client:      client.NewClient(listenAddr, path, serviceAddr),
		listenAddr:  fmt.Sprintf("%s%s", listenAddr, path),
		serviceAddr: serviceAddr,
		// TODO(sycomancy): http://localhost part should be configurable
		callback: fmt.Sprintf("http://localhost%s%s", listenAddr, path),
	}
}

func (r *RequestProcessor) Run() {
	logrus.WithFields(logrus.Fields{
		"listenAddr":  r.listenAddr,
		"serviceAddr": r.serviceAddr,
	}).Info("request processor started")

	go r.client.Run()

	for {
		data := <-r.client.Data
		logrus.WithFields(logrus.Fields{
			"requestId":   data.RequestID,
			"status":      data.Status,
			"resultCount": len(data.Data),
		}).Info("got results from service")
	}
}

func (r *RequestProcessor) SendRequest(request *types.Request) error {
	request.CallbackURL = r.callback
	result, err := r.client.SendRequest(request)
	ctx := logrus.WithFields(logrus.Fields{
		"requestId": result.RequestID,
		"filter":    request.Filter,
		"status":    result.Status,
	})
	if err != nil {
		ctx.Warn("request failed")
		return err
	}
	ctx.Info("request successfully sent")
	return nil
}
