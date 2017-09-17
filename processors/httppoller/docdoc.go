// Code generated by "bitfanDoc "; DO NOT EDIT
package httppoller

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "httppoller",
  ImportPath: "github.com/vjeantet/bitfan/processors/httppoller",
  Doc:        "HTTPPoller allows you to call an HTTP Endpoint, decode the output into an event",
  DocShort:   "",
  Options:    &doc.ProcessorOptions{
    Doc:     "",
    Options: []*doc.ProcessorOption{
      &doc.ProcessorOption{
        Name:           "processors.CommonOptions",
        Alias:          ",squash",
        Doc:            "",
        Required:       false,
        Type:           "processors.CommonOptions",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Codec",
        Alias:          "codec",
        Doc:            "The codec used for input data. Input codecs are a convenient method for decoding\nyour data before it enters the input, without needing a separate filter in your bitfan pipeline",
        Required:       false,
        Type:           "codec",
        DefaultValue:   "\"plain\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Interval",
        Alias:          "interval",
        Doc:            "Use CRON or BITFAN notation",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "interval => \"every_10s\"",
      },
      &doc.ProcessorOption{
        Name:           "Method",
        Alias:          "method",
        Doc:            "Http Method",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"GET\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Url",
        Alias:          "url",
        Doc:            "URL",
        Required:       true,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "url=> \"http://google.fr\"",
      },
      &doc.ProcessorOption{
        Name:           "Target",
        Alias:          "target",
        Doc:            "When data is an array it stores the resulting data into the given target field.",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
  Ports: []*doc.ProcessorPort{},
}
}