package main

import (
	"go.k6.io/croconf"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
)

type ScriptConfig struct {
	cm *croconf.Manager
	*GlobalConfig

	UserAgent string
	VUs       int64

	Duration Duration

	Scenarios lib.ScenarioConfigs

	// TODO: have a sub-config
}

func NewScriptConfig(
	globalConf *GlobalConfig,
	cliSource *croconf.SourceCLI,
	envVarsSource *croconf.SourceEnvVars,
	jsonSource *croconf.SourceJSON,
) (*ScriptConfig, error) {
	cm := croconf.NewManager()
	conf := &ScriptConfig{GlobalConfig: globalConf, cm: cm} // TODO: somehow save the sources in the struct as well?

	cm.AddField(
		croconf.NewStringField(
			&conf.UserAgent,
			croconf.DefaultStringValue("croconf example demo v0.0.1 (https://k6.io/)"),
			jsonSource.From("userAgent"),
			envVarsSource.From("K6_USER_AGENT"),
			cliSource.FromName("user-agent"),
			// TODO: figure this out...
			// croconf.WithDescription("user agent for http requests")
		),
	)

	cm.AddField(croconf.NewInt64Field(
		&conf.VUs,
		croconf.DefaultInt64Value(1),
		jsonSource.From("vus"),
		envVarsSource.From("K6_VUS"),
		cliSource.FromNameAndShorthand("vus", "u"),
		// croconf.WithDescription("number of virtual users") // TODO
	))

	cm.AddField(croconf.NewTextBasedField(
		&conf.Duration,
		croconf.DefaultStringValue("1s"),
		jsonSource.From("duration"),
		envVarsSource.From("K6_DURATION"),
		cliSource.FromNameAndShorthand("duration", "d"),
	))

	cm.AddField(croconf.NewCustomField(
		&conf.Scenarios,
		croconf.DefaultCustomValue(func() {
			conf.Scenarios = lib.ScenarioConfigs{
				lib.DefaultScenarioName: executor.NewPerVUIterationsConfig(lib.DefaultScenarioName),
			}
		}),
		jsonSource.From("scenarios").To(&conf.Scenarios),
	))

	// TODO: add the other options and actually process and consolidate the
	// config values and handle any errors... Here we probably want to error out
	// if we see unknown CLI flags or JSON options

	// TODO: automatically do this on Consolidate()?
	if err := cliSource.Parse(); err != nil {
		return nil, err
	}

	if err := cm.Consolidate(); err != nil {
		return nil, err
	}

	return conf, nil
}