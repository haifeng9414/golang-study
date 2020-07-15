package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func main() {
	now := time.Now()
	v := reflect.ValueOf(now) // 获取reflect.Value
	fmt.Printf("%T\n", v)     // 输出：reflect.Value

	i := v.Interface()    // 返回一个interface{}，其指向reflect.Value表示的值
	a := i.(time.Time)    // 强转
	fmt.Printf("%T\n", a) // 输出：time.Time
	fmt.Println(a)        // 输出：2020-07-14 21:05:19.764024 +0800 CST m=+0.000107937

	b := v.Type()                         // 返回一个reflect.Type
	fmt.Printf("%T\n", b)                 // 输出：*reflect.rtype
	fmt.Println(b)                        // 输出：time.Time
	fmt.Println(b == reflect.TypeOf(now)) // 输出：true

	c := v.Kind()                          // 获取值的类型
	fmt.Printf("%T\n", c)                  // 输出：reflect.Kind
	fmt.Println(c)                         // 输出：struct
	fmt.Println(reflect.ValueOf(1).Kind()) // int

	var x int64 = 1
	d := 1 * time.Nanosecond
	fmt.Println(printAny(x))                  // 输出：1
	fmt.Println(printAny(d))                  // 输出：1
	fmt.Println(printAny([]int64{x}))         // 输出：[]int64 0xc0000b4090
	fmt.Println(printAny([]time.Duration{d})) // 输出：[]time.Duration 0xc0000b4098

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb", Year: 1964,
		Color: false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)", "Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	/*
	输出：
	Display strangelove (main.Movie):
	strangelove.Title = "Dr. Strangelove"
	strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
	strangelove.Year = 1964
	strangelove.Color = false
	strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
	strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
	strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
	strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
	strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
	strangelove.Actor["Pres. Merkin Muffley"] = "Peter Sellers"
	strangelove.Oscars[0] = "Best Actor (Nomin.)"
	strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
	strangelove.Oscars[2] = "Best Director (Nomin.)"
	strangelove.Oscars[3] = "Best Picture (Nomin.)"
	strangelove.Sequel = nil
	*/
	Display("strangelove", strangelove)

	/*
	输出：
	Display os.Stderr (*os.File):
	(*(*os.Stderr).file).pfd.fdmu.state = 0
	(*(*os.Stderr).file).pfd.fdmu.rsema = 0
	(*(*os.Stderr).file).pfd.fdmu.wsema = 0
	(*(*os.Stderr).file).pfd.Sysfd = 2
	(*(*os.Stderr).file).pfd.pd.runtimeCtx = 0
	(*(*os.Stderr).file).pfd.iovecs = nil
	(*(*os.Stderr).file).pfd.csema = 0
	(*(*os.Stderr).file).pfd.isBlocking = 1
	(*(*os.Stderr).file).pfd.IsStream = true
	(*(*os.Stderr).file).pfd.ZeroReadIsEOF = true
	(*(*os.Stderr).file).pfd.isFile = true
	(*(*os.Stderr).file).name = "/dev/stderr"
	(*(*os.Stderr).file).dirinfo = nil
	(*(*os.Stderr).file).nonblock = false
	(*(*os.Stderr).file).stdoutOrErr = true
	(*(*os.Stderr).file).appendMode = false
	 */
	Display("os.Stderr", os.Stderr)
}

func printAny(value interface{}) string {
	return printValue(reflect.ValueOf(value))
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i)) // 通过Index方法获取切片和数组的指定元素
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name) // 结构体通过reflect.Type的Field方法获取属性
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, printValue(key)), v.MapIndex(key)) // map通过MapKeys方法获取所有的key，并通过MapIndex获取指定key的value
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem()) // Elem方法返回一个reflect.Value变量，其持有指针指向的变量
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem()) // Elem方法返回一个reflect.Value变量，其持有接口指向的变量
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, printValue(v))
	}
}

func printValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// 简单起见，省略了一些浮点类型的判断
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface...
		return v.Type().String() + " value"
	}
}
