package timevalue

import "fmt"

func resolveArgs(args []any) (any, []string, error) {
	if len(args) == 0 {
		return "now", []string{}, nil
	}

	timeValue := args[0]
	modifiers, err := anySliceToTypedSlice[string](args[1:])
	if err != nil {
		return "", []string{}, err
	}

	s, ok := timeValue.(string)
	if ok && (s == "subsec" || s == "subsecond") {
		// https://sqlite.org/lang_datefunc.html#subsec
		// The "subsecond" and "subsec" modifiers have the special property
		// that they can occur as the first argument to date/time functions
		// When this happens, the time-value that is normally in the first
		// argument is understood to be "now"
		timeValue = "now"
		modifiers = append(modifiers, s)
	}

	return timeValue, modifiers, nil
}

func anySliceToTypedSlice[T any](a []any) ([]T, error) {
	s := make([]T, len(a))
	for i, v := range a {
		str, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("element %d is not a %T (type %T)", i, new(T), v)
		}
		s[i] = str
	}
	return s, nil
}

func UnixEpoch(args ...any) any {
	tv, mods, err := resolveArgs(args)
	if err != nil {
		return nil
	}

	var timeValue TimeValue

	switch tvTyped := tv.(type) {
	case string:
		timeValue, err = NewFromString(tvTyped, mods...)
		if err != nil {
			return nil
		}
	case int64:
		timeValue, err = NewFromFloat(float64(tvTyped), mods...)
		if err != nil {
			return nil
		}
	case float64:
		timeValue, err = NewFromFloat(tvTyped, mods...)
		if err != nil {
			return nil
		}
	default:
		return nil
	}

	subsec := Modifiers(mods).SubSec()
	return timeValue.UnixEpoch(subsec)
}
