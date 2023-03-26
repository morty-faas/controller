package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyxia-org/morty-gateway/config"
	"github.com/polyxia-org/morty-gateway/server/rik"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Server struct {
	config config.Config
	rik.ControllerClient
	l *logrus.Entry

	// Temporary map a function name to a workload ID
	functions map[string]string
}

func NewServer(config config.Config) (*Server, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	l := logrus.NewEntry(logger)
	controllerClient, err := rik.NewControllerClient(l, config)
	if err != nil {
		return nil, err
	}

	// When initializing the gateway we want to know existing functions in order to have a minimal state
	// this is not perfect, but it will work as a MVP
	functions, err := controllerClient.GetExistingFunctions()
	if err != nil {
		l.WithError(err).Error("Could not load existing function, will start with an empty list")
	}

	return &Server{
		config:           config,
		ControllerClient: *controllerClient,
		l:                l,
		functions:        functions,
	}, nil
}

func (server *Server) Run() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ginlogrus.Logger(logrus.New()))

	router.GET("/health/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// Handlers to create a FunctionRequest
	router.POST("/functions", server.createFunctionHandler)

	// Handle all methods to invoke a FunctionRequest
	router.Any("/invoke/:functionName", server.invokeFunctionHandler)

	listeningPort := fmt.Sprintf(":%d", server.config.Port)
	err := router.Run(listeningPort)
	if err != nil {
		server.l.WithError(err).Error("Could not start server")
		os.Exit(1)
	}
}

func (server *Server) createFunctionHandler(c *gin.Context) {
	functionBody := rik.FunctionRequest{}
	if err := c.ShouldBindJSON(&functionBody); err != nil {
		server.l.WithError(err).Warn("Could not parse create function body")
		c.JSON(400, gin.H{
			"message": "Invalid request body, please check the documentation",
		})
		return
	}

	workloadId, err := server.ControllerClient.CreateFunction(functionBody)
	if err != nil {
		server.l.WithError(err).Error("Could not create function")
		c.JSON(500, gin.H{
			"message": "Could not create function",
		})
		return
	}

	server.functions[functionBody.Name] = workloadId
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func (server *Server) invokeFunctionHandler(c *gin.Context) {
	functionName := c.Param("functionName")
	l := server.l.WithField("functionName", functionName)

	l.Debug("Invoke function")

	// Check if the mapping between function and workload exist
	workloadId, ok := server.functions[functionName]
	if !ok {
		l.Warn("Function not found")
		c.JSON(404, gin.H{
			"message": "Could not find the request resource",
		})
		return
	}

	instances, err := server.ControllerClient.FetchInstances(workloadId)
	if err != nil {
		l.WithError(err).Error("Could not fetch instances")
		c.JSON(500, gin.H{
			"message": "We cannot serve right now",
		})
		return
	}

	// if no active instance were found, we want to create one and serve the request
	if len(instances) == 0 {
		l.Info("No instance found, creating one")
		err = server.ControllerClient.CreateWorkloadInstance(workloadId)
		if err != nil {
			server.l.WithError(err).Error("Could not create instance")
			c.JSON(500, gin.H{
				"message": "We couldn't create an instance for your function right now, please try again later",
			})
			return
		}

		// Simulate time to create the instance and schedule
		time.Sleep(500 * time.Millisecond)

		instances, err = server.ControllerClient.FetchInstances(workloadId)
		if err != nil {
			l.WithError(err).Error("Could not fetch instances")
			c.JSON(500, gin.H{
				"message": "We couldn't know the state of your function right now, please try again later",
			})
			return
		}
	}

	l.WithField("instances_len", len(instances)).Debug("Fetched instances")
	randomIndex := rand.Intn(len(instances))
	instance := instances[randomIndex]

	// Currently, we consider that the underlying orchestrator, RIK, runs its entire stack on the same
	// machine. Once RIK Controller will be able to dynamically return the worker node IP, we will need to
	// update this line of code.
	// See: https://github.com/polyxia-org/polyxia-org/issues/16
	functionAddr := instance.GetRuntimeUrl(server.config.RIKController.Hostname())

	proxy := httputil.NewSingleHostReverseProxy(functionAddr)

	// Modify the response so we can extract from the Alpha
	// response the payload to return it to the caller.
	proxy.ModifyResponse = func(r *http.Response) error {
		var functionResponse rik.FunctionResponse

		by, err := io.ReadAll(r.Body)
		if err != nil {
			l.WithError(err).Error("Could not read response body")
			return err
		}
		defer r.Body.Close()

		if err := json.Unmarshal(by, &functionResponse); err != nil {
			l.WithError(err).Error("Could not unmarshal response body")
			return err
		}

		var responseBytes []byte

		// If the function payload is a string, return it as text
		if value, ok := functionResponse.Payload.(string); ok {
			responseBytes = []byte(value)
		} else {
			responseBytes, err = json.Marshal(functionResponse.Payload)
			if err != nil {
				return err
			}
		}

		contentLength := len(responseBytes)
		r.Body = io.NopCloser(bytes.NewReader(responseBytes))
		r.ContentLength = int64(contentLength)
		r.Header.Set("Content-Length", strconv.Itoa(contentLength))

		return nil
	}

	// Perform healthcheck against the Alpha agent
	// If alpha doesn't anwser to our requests, it probably that
	// the VM isn't ready yet to receive our requests.
	const maxHealthcheckRetries = 10
	healthcheck := functionAddr.String() + "/_/health"
	for i := 0; i < maxHealthcheckRetries; i++ {
		l.Debugf("Performing healthcheck request on Alpha: %s", healthcheck)
		if _, err := http.Get(healthcheck); err != nil {
			if i == maxHealthcheckRetries-1 {
				l.Errorf("failed to perform healthcheck on Alpha: %v", err)
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status":  "OUT_OF_SERVICE",
					"message": "One or more instances of the function can't be marked as ready",
				})
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}
		l.Infof("Function '%s' is healthy and ready to receive requests", functionName)
		break
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
