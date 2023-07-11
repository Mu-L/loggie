/*
Copyright 2021 Loggie Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kafka

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/segmentio/kafka-go"

	kafkaSink "github.com/loggie-io/loggie/pkg/sink/kafka"
)

const (
	earliestOffsetReset = "earliest"
	latestOffsetReset   = "latest"
)

type Config struct {
	Brokers            []string       `yaml:"brokers,omitempty" validate:"required"`
	Topic              string         `yaml:"topic,omitempty"` // reserved for compatibility
	Topics             []string       `yaml:"topics,omitempty"`
	GroupId            string         `yaml:"groupId,omitempty" default:"loggie"`
	ClientId           string         `yaml:"clientId,omitempty"`
	Worker             int            `yaml:"worker,omitempty" default:"1"`
	QueueCapacity      int            `yaml:"queueCapacity" default:"100"`
	MinAcceptedBytes   int            `yaml:"minAcceptedBytes" default:"1"`
	MaxAcceptedBytes   int            `yaml:"maxAcceptedBytes" default:"1024000"`
	ReadMaxAttempts    int            `yaml:"readMaxAttempts" default:"3"`
	MaxReadWait        time.Duration  `yaml:"maxPollWait" default:"10s"`
	ReadBackoffMin     time.Duration  `yaml:"readBackoffMin" default:"100ms"`
	ReadBackoffMax     time.Duration  `yaml:"readBackoffMax" default:"1s"`
	EnableAutoCommit   bool           `yaml:"enableAutoCommit"`
	AutoCommitInterval time.Duration  `yaml:"autoCommitInterval" default:"1s"`
	AutoOffsetReset    string         `yaml:"autoOffsetReset" default:"latest" validate:"oneof=earliest latest"`
	SASL               kafkaSink.SASL `yaml:"sasl,omitempty"`
	AddonMeta          *bool          `yaml:"addonMeta,omitempty" default:"true"`
}

func getAutoOffset(autoOffsetReset string) int64 {
	switch autoOffsetReset {
	case earliestOffsetReset:
		return kafka.FirstOffset
	case latestOffsetReset:
		return kafka.LastOffset
	}

	return kafka.LastOffset
}

func (c *Config) SetDefaults() {
	if c.SASL.UserName != "" {
		c.SASL.Username = c.SASL.UserName
	}
}

func (c *Config) Validate() error {
	if c.Topic == "" && len(c.Topics) == 0 {
		return errors.New("topic or topics is required")
	}

	if c.Topic != "" {
		_, err := regexp.Compile(c.Topic)
		if err != nil {
			return errors.WithMessagef(err, "compile kafka topic regex %s error", c.Topic)
		}
	}
	if len(c.Topics) > 0 {
		for _, t := range c.Topics {
			_, err := regexp.Compile(t)
			if err != nil {
				return errors.WithMessagef(err, "compile kafka topic regex %s error", t)
			}
		}
	}

	if err := c.SASL.Validate(); err != nil {
		return err
	}

	return nil
}
