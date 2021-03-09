package eeprom

type Interval struct {
	Start uint8 `json:"start"`
	End   uint8 `json:"end"`
}
type Address struct {
	Fixed                        Interval `json:"fixed"`
	StretchGain                  Interval `json:"stretch-gain"`
	Speed                        Interval `json:"speed"`
	Punch                        Interval `json:"punch"`
	DeadBand                     Interval `json:"dead-band"`
	Damping                      Interval `json:"damping"`
	SafeTimer                    Interval `json:"safe-timer"`
	Flag                         Interval `json:"flag"`
	MaximumPulseLimit            Interval `json:"maximum-pulse-limit"`
	MinimumPulseLimit            Interval `json:"minimum-pulse-limit"`
	SignalSpeed                  Interval `json:"signal-speed"`
	TemperatureLimit             Interval `json:"temperature-limit"`
	CurrentLimit                 Interval `json:"current-limit"`
	Response                     Interval `json:"response"`
	UserOffset                   Interval `json:"user-offset"`
	ID                           Interval `json:"id"`
	CharacteristicChangeStretch1 Interval `json:"characteristic-change-stretch1"`
	CharacteristicChangeStretch2 Interval `json:"characteristic-change-stretch2"`
	CharacteristicChangeStretch3 Interval `json:"characteristic-change-stretch3"`
}

func NewInterval(start, end uint8) Interval {
	return Interval{Start: start, End: end}
}
