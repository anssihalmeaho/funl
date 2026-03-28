package std

import (
	"fmt"
	"time"

	"github.com/anssihalmeaho/funl/funl"
)

func initSTDCal(interpreter *funl.Interpreter) (err error) {
	stdModuleName := "stdcal"
	topFrame := funl.NewTopFrameWithInterpreter(interpreter)
	stdFuncs := []stdFuncInfo{
		{
			Name:   "now",
			Getter: getStdCalNow,
		},
		{
			Name:       "sub",
			Getter:     getStdCalSub,
			IsFunction: true,
		},
		{
			Name:       "add",
			Getter:     getStdCalAdd,
			IsFunction: true,
		},
		{
			Name:       "parts",
			Getter:     getStdCalParts,
			IsFunction: true,
		},
		{
			Name:       "utc",
			Getter:     getStdCalUTC,
			IsFunction: true,
		},
		{
			Name:       "local",
			Getter:     getStdCalLocal,
			IsFunction: true,
		},
		{
			Name:       "to-unix",
			Getter:     getStdCalToUnix,
			IsFunction: true,
		},
		{
			Name:       "from-unix",
			Getter:     getStdCalFromUnix,
			IsFunction: true,
		},
		{
			Name:       "parse",
			Getter:     getStdCalParse,
			IsFunction: true,
		},
		{
			Name:       "parse-in-location",
			Getter:     getStdCalParseInLocation,
			IsFunction: true,
		},
		{
			Name:       "format",
			Getter:     getStdCalFormat,
			IsFunction: true,
		},
		{
			Name:       "fixed-zone",
			Getter:     getStdCalFixedZone,
			IsFunction: true,
		},
		{
			Name:       "load-location",
			Getter:     getStdCalLoadLocation,
			IsFunction: true,
		},
		{
			Name:       "date",
			Getter:     getStdCalDate,
			IsFunction: true,
		},
	}

	err = topFrame.Syms.Add("hour", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Hour)}})
	if err != nil {
		return
	}
	err = topFrame.Syms.Add("minute", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Minute)}})
	if err != nil {
		return
	}
	err = topFrame.Syms.Add("second", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Second)}})
	if err != nil {
		return
	}
	err = topFrame.Syms.Add("millisecond", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Millisecond)}})
	if err != nil {
		return
	}
	err = topFrame.Syms.Add("microsecond", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Microsecond)}})
	if err != nil {
		return
	}
	err = topFrame.Syms.Add("nanosecond", &funl.Item{Type: funl.ValueItem, Data: funl.Value{Kind: funl.IntValue, Data: int(time.Nanosecond)}})
	if err != nil {
		return
	}

	err = setSTDFunctions(topFrame, stdModuleName, stdFuncs, interpreter)
	return
}

var timeFormats = map[string]string{
	"Layout":      time.Layout,
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
	"DateTime":    time.DateTime,
	"DateOnly":    time.DateOnly,
	"TimeOnly":    time.TimeOnly,
}

func makeTripletList(frame *funl.Frame, ok bool, reason string, val funl.Value) funl.Value {
	return funl.MakeListOfValues(frame, []funl.Value{
		{
			Kind: funl.BoolValue,
			Data: ok,
		},
		{
			Kind: funl.StringValue,
			Data: reason,
		},
		val,
	})
}

type OpaqueTime struct {
	tim *time.Time
}

func (*OpaqueTime) TypeName() string {
	return "time"
}

func (ot *OpaqueTime) Str() string {
	return ot.TypeName() + ": " + ot.tim.String()
}

func (ot *OpaqueTime) Equals(with funl.OpaqueAPI) bool {
	other, ok := with.(*OpaqueTime)
	if !ok || other == nil {
		return false
	}
	return ot.tim.Equal(*other.tim)
}

type OpaqueLocation struct {
	loca *time.Location
}

func (ol *OpaqueLocation) Str() string {
	return ol.TypeName() + ": " + ol.loca.String()
}

func (*OpaqueLocation) TypeName() string {
	return "location"
}

func (ol *OpaqueLocation) Equals(with funl.OpaqueAPI) bool {
	other, ok := with.(*OpaqueLocation)
	if !ok || other == nil {
		return false
	}
	return ol.loca.String() == other.loca.String()
}

func getStdCalDate(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 8 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		for i := range 7 {
			if arguments[i].Kind != funl.IntValue {
				funl.RunTimeError2(frame, "%s: requires int value", name)
			}
		}
		year := arguments[0].Data.(int)
		month := time.Month(arguments[1].Data.(int))
		day := arguments[2].Data.(int)
		hour := arguments[3].Data.(int)
		min := arguments[4].Data.(int)
		sec := arguments[5].Data.(int)
		nsec := arguments[6].Data.(int)

		location := arguments[7]
		if location.Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		loc := location.Data.(*OpaqueLocation)
		t := time.Date(year, month, day, hour, min, sec, nsec, loc.loca)
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &t}}
	}
}

func getStdCalLoadLocation(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		loc, err := time.LoadLocation(arguments[0].Data.(string))
		if err != nil {
			return makeTripletList(frame, false, fmt.Sprintf("%s: %v", name, err), funl.Value{Kind: funl.StringValue, Data: ""})
		}
		val := funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueLocation{loca: loc}}
		return makeTripletList(frame, true, "", val)
	}
}

func getStdCalFixedZone(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		loc := time.FixedZone(arguments[0].Data.(string), arguments[1].Data.(int))
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueLocation{loca: loc}}
	}
}

func getStdCalFormat(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		layout := timeFormats[arguments[1].Data.(string)]
		return funl.Value{Kind: funl.StringValue, Data: t.tim.Format(layout)}
	}
}

func getStdCalParseInLocation(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 3 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		if arguments[2].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		layoutStr := arguments[0].Data.(string)
		layout, found := timeFormats[layoutStr]
		if !found {
			return makeTripletList(frame, false, fmt.Sprintf("%s: unknown layout: %s", name, layoutStr), funl.Value{Kind: funl.StringValue, Data: ""})
		}
		tval := arguments[1].Data.(string)
		loc := arguments[2].Data.(*OpaqueLocation)

		tim, err := time.ParseInLocation(layout, tval, loc.loca)
		if err != nil {
			return makeTripletList(frame, false, fmt.Sprintf("%s: parse error: %v", name, err), funl.Value{Kind: funl.StringValue, Data: ""})
		}
		val := funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &tim}}
		return makeTripletList(frame, true, "", val)
	}
}

func getStdCalParse(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		if arguments[1].Kind != funl.StringValue {
			funl.RunTimeError2(frame, "%s: requires string value", name)
		}
		layoutStr := arguments[0].Data.(string)
		layout, found := timeFormats[layoutStr]
		if !found {
			return makeTripletList(frame, false, fmt.Sprintf("%s: unknown layout: %s", name, layoutStr), funl.Value{Kind: funl.StringValue, Data: ""})
		}
		tval := arguments[1].Data.(string)
		tim, err := time.Parse(layout, tval)
		if err != nil {
			return makeTripletList(frame, false, fmt.Sprintf("%s: parse error: %v", name, err), funl.Value{Kind: funl.StringValue, Data: ""})

		}
		val := funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &tim}}
		return makeTripletList(frame, true, "", val)
	}
}

func getStdCalFromUnix(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		sec := int64(arguments[0].Data.(int))
		nsec := int64(arguments[1].Data.(int))
		t := time.Unix(sec, nsec)
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &t}}
	}
}

func getStdCalToUnix(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		return funl.Value{Kind: funl.IntValue, Data: int(t.tim.Unix())}
	}
}

func getStdCalLocal(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		localTime := t.tim.Local()
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &localTime}}
	}
}

func getStdCalUTC(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		utcTime := t.tim.UTC()
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &utcTime}}
	}
}

func getStdCalParts(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 1 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		hour, min, sec := t.tim.Clock()
		year, month, day := t.tim.Date()
		weekday := t.tim.Weekday()

		mapv := funl.HandleMapOP(frame, []*funl.Item{})
		mapv = putToMap(frame, mapv, "hour", funl.Value{Kind: funl.IntValue, Data: hour})
		mapv = putToMap(frame, mapv, "min", funl.Value{Kind: funl.IntValue, Data: min})
		mapv = putToMap(frame, mapv, "sec", funl.Value{Kind: funl.IntValue, Data: sec})
		mapv = putToMap(frame, mapv, "year", funl.Value{Kind: funl.IntValue, Data: year})
		mapv = putToMap(frame, mapv, "month", funl.Value{Kind: funl.IntValue, Data: int(month)})
		mapv = putToMap(frame, mapv, "month-name", funl.Value{Kind: funl.StringValue, Data: month.String()})
		mapv = putToMap(frame, mapv, "day", funl.Value{Kind: funl.IntValue, Data: day})
		mapv = putToMap(frame, mapv, "weekday", funl.Value{Kind: funl.IntValue, Data: int(weekday)})
		mapv = putToMap(frame, mapv, "weekday-name", funl.Value{Kind: funl.StringValue, Data: weekday.String()})
		return mapv
	}
}

func getStdCalAdd(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		if arguments[1].Kind != funl.IntValue {
			funl.RunTimeError2(frame, "%s: requires int value", name)
		}
		t := arguments[0].Data.(*OpaqueTime)
		duration := arguments[1].Data.(int)
		newtime := t.tim.Add(time.Duration(duration))
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &newtime}}
	}
}

func getStdCalSub(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		if l := len(arguments); l != 2 {
			funl.RunTimeError2(frame, "%s: wrong amount of arguments (%d)", name, l)
		}
		if arguments[0].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		if arguments[1].Kind != funl.OpaqueValue {
			funl.RunTimeError2(frame, "%s: requires opaque value", name)
		}
		time1 := arguments[0].Data.(*OpaqueTime)
		time2 := arguments[1].Data.(*OpaqueTime)
		duration := time1.tim.Sub(*time2.tim)
		return funl.Value{Kind: funl.IntValue, Data: int(duration)}
	}
}

func getStdCalNow(name string) stdFuncType {
	return func(frame *funl.Frame, arguments []funl.Value) (retVal funl.Value) {
		t := time.Now()
		return funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueTime{tim: &t}}
	}
}
