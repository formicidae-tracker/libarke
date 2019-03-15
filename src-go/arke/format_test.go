package arke

import (
	"time"

	. "gopkg.in/check.v1"
)

type FormatSuite struct{}

var _ = Suite(&FormatSuite{})

func (s *FormatSuite) TestFormatting(c *C) {
	testdata := []struct {
		M ReceivableMessage
		E string
	}{
		{
			&ZeusSetPoint{51.455, 20.001, 127},
			"Zeus.SetPoint{Humidity: 51.46%, Temperature: 20.00°C, Wind: 127}",
		},
		{
			&ZeusStatus{
				Fans: [3]FanStatusAndRPM{
					FanStatusAndRPM(uint16(FanOK)<<12 | uint16(1234)),
					FanStatusAndRPM(uint16(FanStalled) << 14),
					FanStatusAndRPM(uint16(FanAging)<<14 | uint16(100)),
				},
			},
			"Zeus.Status{General: idle, WindFan: {Status: OK, RPM: 1234}, LeftFan: {Status: Aging, RPM: 100}, RightFan: {Status: Stalled, RPM: 0}}",
		},
		{
			&ZeusReport{12.009, [4]float32{12.001, 13.005, 14, 15}},
			"Zeus.Report{Humidity: 12.01%, Ant: 12.00°C, Aux1: 13.01°C, Aux2: 14.00°C, Aux3: 15.00°C}",
		},
		{
			&ZeusConfig{
				PDConfig{
					ProportionnalMultiplier: 10,
					DerivativeMultiplier:    2,
					IntegralMultiplier:      100,
					DividerPower:            4,
					DividerPowerIntegral:    14,
				},
				PDConfig{
					ProportionnalMultiplier: 12,
					DerivativeMultiplier:    3,
					IntegralMultiplier:      103,
					DividerPower:            5,
					DividerPowerIntegral:    11,
				},
			},
			"Zeus.Config{Humidity:PIDConfig{Proportional:10/16, Derivative: 2/16, Integral: 100/16384}, Temperature:PIDConfig{Proportional:12/32, Derivative: 3/32, Integral: 103/2048}}",
		},
		{
			&ZeusControlPoint{1245, -5469},
			"Zeus.ControlPoint{Humidity: 1245, Temperature: -5469}",
		},
		{
			&ZeusDeltaTemperature{[4]float32{-1.25, -0.0625, 0, 0.125}},
			"Zeus.DeltaTemperature{Ants: -1.2500°C, Aux1: -0.0625°C, Aux2: 0.0000°C, Aux3: 0.1250°C}",
		},
		{
			&CelaenoSetPoint{Power: 156},
			"Celaeno.SetPoint{Power: 156}",
		},
		{
			&CelaenoConfig{
				RampUpTime:    500 * time.Millisecond,
				RampDownTime:  3500 * time.Millisecond,
				DebounceTime:  1000 * time.Millisecond,
				MinimumOnTime: 4 * time.Second,
			},
			"Celaeno.Config{RampUp: 500ms, RampDown: 3.5s, MinimumOn: 4s, Debounce: 1s}",
		},
		{
			&CelaenoStatus{},
			"Celaeno.Status{WaterLevel: nominal, Fan:{Status: OK, RPM: 0}}"},
		{
			&HeliosSetPoint{},
			"Helios.SetPoint{Visible: 0, UV: 0}",
		},
		{
			&HeartBeatData{
				Class: HeliosClass,
				ID:    3,
			},
			"arke.HeartBeat{Class: Helios, ID: 3}",
		},
		{
			&HeartBeatData{
				Class:        CelaenoClass,
				ID:           3,
				MajorVersion: 1,
				MinorVersion: 4,
				PatchVersion: 0,
				TweakVersion: 1,
			},
			"arke.HeartBeat{Class: Celaeno, ID: 3, Version: 1.4.0.1}",
		},
	}

	c.Assert(len(testdata) >= len(messageFactory), Equals, true)

	for _, d := range testdata {
		c.Check(d.M.String(), Equals, d.E)
	}

	zeusStatus := []struct {
		V ZeusStatusValue
		E string
	}{
		{ZeusIdle, "idle"},
		{ZeusTemperatureUnreachable | ZeusActive, "temperature-unreachable|active"},
		{ZeusHumidityUnreachable | ZeusClimateNotControlledWatchDog, "humidity-unreachable|climate-uncontrolled|idle"},
		{ZeusClimateNotControlledWatchDog | ZeusActive, "sensor-issue"},
	}

	for _, d := range zeusStatus {
		c.Check(d.V.String(), Equals, d.E)

	}

	celaenoStatus := []struct {
		V WaterLevelStatus
		E string
	}{
		{CelaenoWaterNominal, "nominal"},
		{CelaenoWaterReadError | CelaenoWaterNominal, "readout-error"},
		{CelaenoWaterReadError | CelaenoWaterWarning, "readout-error|warning"},
		{CelaenoWaterReadError | CelaenoWaterCritical | CelaenoWaterWarning, "readout-error|critical"},
		{CelaenoWaterCritical | CelaenoWaterWarning, "critical"},
	}

	for _, d := range celaenoStatus {
		c.Check(d.V.String(), Equals, d.E)

	}

}
