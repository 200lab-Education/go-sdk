/*
 * @author          Viet Tran <viettranx@gmail.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license         Apache-2.0
 */

package jaeger

import (
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	jg "go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type jaeger struct {
	logger            logger.Logger
	stopChan          chan bool
	processName       string
	sampleTraceRating float64
	agentURI          string
	port              int
	stdTracingEnabled bool
}

func NewJaeger(processName string) *jaeger {
	return &jaeger{
		processName: processName,
		stopChan:    make(chan bool),
	}
}

func (j *jaeger) Name() string {
	return "jaeger"
}

func (j *jaeger) GetPrefix() string {
	return j.Name()
}

// Note: this plugin will not return anything
func (j *jaeger) Get() interface{} {
	return nil
}

func (j *jaeger) InitFlags() {
	flag.Float64Var(
		&j.sampleTraceRating,
		"jaeger-trace-sample-rate",
		1.0,
		"sample rating for remote tracing from OpenSensus: 0.0 -> 1.0 (default is 1.0)",
	)

	flag.StringVar(
		&j.agentURI,
		"jaeger-agent-uri",
		"",
		"jaeger agent URI to receive tracing data directly",
	)

	flag.IntVar(
		&j.port,
		"jaeger-agent-port",
		6831,
		"jaeger agent URI to receive tracing data directly",
	)

	flag.BoolVar(
		&j.stdTracingEnabled,
		"jaeger-std-enabled",
		false,
		"enable tracing export to std (default is false)",
	)
}

func (j *jaeger) Configure() error {
	j.logger = logger.GetCurrent().GetLogger(j.Name())
	return nil
}

func (j *jaeger) Run() error {
	if err := j.Configure(); err != nil {
		return err
	}
	return j.connectToJaegerAgent()
}

func (j *jaeger) Stop() <-chan bool {
	go func() {
		if !j.isEnabled() {
			j.stopChan <- true
			return
		}

		j.stopChan <- true
	}()

	return j.stopChan
}

func (j *jaeger) isEnabled() bool {
	return j.agentURI != ""
}

func (j *jaeger) connectToJaegerAgent() error {
	if !j.isEnabled() {
		return nil
	}

	url := fmt.Sprintf("%s:%d", j.agentURI, j.port)
	j.logger.Infof("connecting to Jaeger Agent on %s...", url)

	je, err := jg.NewExporter(jg.Options{
		AgentEndpoint: url,
		Process:       jg.Process{ServiceName: j.processName},
	})

	if err != nil {
		return err
	}

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)

	// Trace view for console
	if j.stdTracingEnabled {
		// Register stats and trace exporters to export
		// the collected data.
		view.RegisterExporter(&PrintExporter{})

		// Register the views to collect server request count.
		if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
			j.logger.Errorf("jaeger error: %s", err.Error())
		}
	}

	if j.sampleTraceRating >= 1 {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	} else {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(j.sampleTraceRating)})
	}

	return nil
}
