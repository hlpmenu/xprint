package xprint_test

// func TestUnsafePointer(t *testing.T) {
// 	var s *string = strptr("hello")
// 	var i *int = intPtr(42)

// 	testCases := []struct {
// 		name   string
// 		arg    any
// 		format string
// 	}{
// 		{"unsafe pointer", unsafe.Pointer(nil), "%v"},
// 		{"unsafe pointer", unsafe.Pointer(s), "%s"},
// 		{"unsafe pointer", unsafe.Pointer(strptr("hello")), "%v"},
// 		{"unsafe pointer", unsafe.Pointer(i), "%d"},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			o := xprint.Printf("Number: ", tc.format, tc.arg)
// 			fo := fmt.Sprintf("Number: ", tc.format, tc.arg)
// 			if o != fo {
// 				t.Errorf("Expected %s, got %s", fo, o)
// 			}
// 		})
// 	}

// }
