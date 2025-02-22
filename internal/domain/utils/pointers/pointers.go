package pointers

func String(s string) *string { return &s }
func Int(i int) *int          { return &i }
func Int32(i int32) *int32    { return &i }
