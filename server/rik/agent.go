package rik

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type AgentClient struct {
	c       *http.Client
	l       *logrus.Entry
	baseUrl *url.URL
}

func NewAgentClient(l *logrus.Entry, baseUrl *url.URL) *AgentClient {
	return &AgentClient{
		c: &http.Client{},
		l: l.
			WithField("component", "rik-agent-client").
			WithField("url", baseUrl.String()),
		baseUrl: baseUrl,
	}
}

func (agent *AgentClient) InvokeFunction(functionName string) (FunctionResponse, error) {
	l := agent.l.WithField("functionName", functionName)
	l.Debug("Invoke function")

	res, err := agent.c.Get(agent.baseUrl.String())
	var functionResponse FunctionResponse
	agent.l.WithField("response", res).Debug("Response from function")

	// When invoking a function, we have several cases:
	// - The runtime is not available, so we dont have a response (and statusCode) but we have err != nil
	// - The function is available, but the function itself returns an error (statusCode >= 400 && statusCode < 500)
	// - The function is available, and the function itself returns an OK status code (statusCode >= 200 && statusCode < 300)
	if err != nil && res != nil && res.StatusCode >= 400 && res.StatusCode < 500 || err == nil {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			l.WithError(err).Error("Could not read response body")
			return functionResponse, err
		}

		err = json.Unmarshal(b, &functionResponse)
		if err != nil {
			l.WithError(err).Error("Could not unmarshal response body")
			return functionResponse, err
		}
	}

	// If we don't have an error, we can close the response body
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		l.WithError(err).Warn("Function returned an non-OK error")

		// Error given by the function, so we return it
		return functionResponse, nil
	}

	return functionResponse, nil
}
